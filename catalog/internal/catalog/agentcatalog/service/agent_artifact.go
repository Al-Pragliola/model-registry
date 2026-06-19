package service

import (
	"errors"
	"fmt"

	"github.com/kubeflow/hub/catalog/internal/catalog/agentcatalog/models"
	"github.com/kubeflow/hub/catalog/internal/db/pagination"
	dbmodels "github.com/kubeflow/hub/internal/platform/db/entity"
	"github.com/kubeflow/hub/internal/platform/db/dbutil"
	service "github.com/kubeflow/hub/internal/platform/db/repository"
	"github.com/kubeflow/hub/internal/platform/db/schema"
	"github.com/kubeflow/hub/internal/platform/db/scopes"
	"github.com/kubeflow/hub/internal/platform/db/utils"
	"gorm.io/gorm"
)

var ErrAgentArtifactNotFound = errors.New("agent_artifact not found")

// AgentArtifactRepositoryImpl implements AgentArtifactRepository using GORM.
type AgentArtifactRepositoryImpl struct {
	*service.GenericRepository[models.AgentArtifact, schema.Artifact, schema.ArtifactProperty, *models.AgentArtifactListOptions]
}

// NewAgentArtifactRepository creates a new AgentArtifactRepository.
func NewAgentArtifactRepository(db *gorm.DB, typeID int32) models.AgentArtifactRepository {
	r := &AgentArtifactRepositoryImpl{}

	r.GenericRepository = service.NewGenericRepository(service.GenericRepositoryConfig[models.AgentArtifact, schema.Artifact, schema.ArtifactProperty, *models.AgentArtifactListOptions]{
		DB:                      db,
		TypeID:                  typeID,
		EntityToSchema:          mapAgentArtifactToSchema,
		SchemaToEntity:          mapSchemaToAgentArtifact,
		EntityToProperties:      mapAgentArtifactToProperties,
		NotFoundError:           ErrAgentArtifactNotFound,
		EntityName:              "agent_artifact",
		PropertyFieldName:       "artifact_id",
		ApplyListFilters:        applyAgentArtifactListFilters,
		CreatePaginationToken:   r.createAgentArtifactPaginationToken,
		ApplyCustomOrdering:     r.applyAgentArtifactCustomOrdering,
		IsNewEntity:             func(entity models.AgentArtifact) bool { return entity.GetID() == nil },
		HasCustomProperties:     func(entity models.AgentArtifact) bool { return entity.GetCustomProperties() != nil },
		EntityMappingFuncs:      newAgentArtifactEntityMappings(),
		PreserveHistoricalTimes: true,
	})

	return r
}

func (r *AgentArtifactRepositoryImpl) Save(entity models.AgentArtifact, parentResourceID *int32) (models.AgentArtifact, error) {
	config := r.GetConfig()
	if entity.GetTypeID() == nil && config.TypeID > 0 {
		entity.SetTypeID(config.TypeID)
	}
	return r.GenericRepository.Save(entity, parentResourceID)
}

func mapAgentArtifactToSchema(artifact models.AgentArtifact) schema.Artifact {
	attrs := artifact.GetAttributes()
	art := schema.Artifact{}
	if typeID := artifact.GetTypeID(); typeID != nil {
		art.TypeID = *typeID
	}
	if artifact.GetID() != nil {
		art.ID = *artifact.GetID()
	}
	if attrs != nil {
		art.Name = attrs.Name
		art.ExternalID = attrs.ExternalID
		if attrs.CreateTimeSinceEpoch != nil {
			art.CreateTimeSinceEpoch = *attrs.CreateTimeSinceEpoch
		}
		if attrs.LastUpdateTimeSinceEpoch != nil {
			art.LastUpdateTimeSinceEpoch = *attrs.LastUpdateTimeSinceEpoch
		}
	}
	return art
}

func mapSchemaToAgentArtifact(art schema.Artifact, props []schema.ArtifactProperty) models.AgentArtifact {
	entity := &models.AgentArtifactImpl{
		ID:     &art.ID,
		TypeID: &art.TypeID,
		Attributes: &models.AgentArtifactAttributes{
			Name:                     art.Name,
			ExternalID:               art.ExternalID,
			CreateTimeSinceEpoch:     &art.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &art.LastUpdateTimeSinceEpoch,
		},
	}

	properties := []dbmodels.Properties{}
	customProperties := []dbmodels.Properties{}
	for _, prop := range props {
		mapped := service.MapArtifactPropertyToProperties(prop)
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapped)
		} else {
			properties = append(properties, mapped)
		}
	}
	entity.Properties = &properties
	entity.CustomProperties = &customProperties
	return entity
}

func mapAgentArtifactToProperties(artifact models.AgentArtifact, artifactID int32) []schema.ArtifactProperty {
	var properties []schema.ArtifactProperty
	if artifact.GetProperties() != nil {
		for _, prop := range *artifact.GetProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}
	if artifact.GetCustomProperties() != nil {
		for _, prop := range *artifact.GetCustomProperties() {
			properties = append(properties, service.MapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}
	return properties
}

func applyAgentArtifactListFilters(query *gorm.DB, _ *models.AgentArtifactListOptions) *gorm.DB {
	return query
}

func (r *AgentArtifactRepositoryImpl) createAgentArtifactPaginationToken(lastItem schema.Artifact, listOptions *models.AgentArtifactListOptions) string {
	if listOptions.GetOrderBy() == "NAME" {
		return pagination.CreateNamePaginationToken(lastItem.ID, lastItem.Name)
	}
	return r.CreateDefaultPaginationToken(lastItem, listOptions)
}

var AgentArtifactOrderByColumns = map[string]string{
	"ID":               "id",
	"CREATE_TIME":      "create_time_since_epoch",
	"LAST_UPDATE_TIME": "last_update_time_since_epoch",
	"NAME":             "name",
	"id":               "id",
}

func (r *AgentArtifactRepositoryImpl) applyAgentArtifactCustomOrdering(query *gorm.DB, listOptions *models.AgentArtifactListOptions) *gorm.DB {
	db := r.GetConfig().DB
	artifactTable := utils.GetTableName(db, &schema.Artifact{})
	orderBy := listOptions.GetOrderBy()

	if orderBy == "NAME" {
		return pagination.ApplyNameOrdering(query, artifactTable, listOptions.GetSortOrder(), listOptions.GetNextPageToken(), listOptions.GetPageSize(), false)
	}

	return r.ApplyStandardPagination(query, listOptions, []models.AgentArtifact{})
}

func (r *AgentArtifactRepositoryImpl) ApplyStandardPagination(query *gorm.DB, listOptions *models.AgentArtifactListOptions, entities any) *gorm.DB {
	pageSize := listOptions.GetPageSize()
	orderBy := listOptions.GetOrderBy()
	sortOrder := listOptions.GetSortOrder()
	nextPageToken := listOptions.GetNextPageToken()

	pag := &dbmodels.Pagination{
		PageSize:      &pageSize,
		OrderBy:       &orderBy,
		SortOrder:     &sortOrder,
		NextPageToken: &nextPageToken,
	}

	return query.Scopes(scopes.PaginateWithOptions(entities, pag, r.GetConfig().DB, "Artifact", AgentArtifactOrderByColumns))
}

func (r *AgentArtifactRepositoryImpl) DeleteBySource(sourceID string) error {
	config := r.GetConfig()
	tableName := utils.GetTableName(config.DB, &schema.Artifact{})
	propTableName := utils.GetTableName(config.DB, &schema.ArtifactProperty{})

	subQuery := config.DB.Table(tableName).
		Select(tableName + ".id").
		Joins("INNER JOIN " + propTableName + " ON " +
			tableName + ".id = " + propTableName + ".artifact_id").
		Where(propTableName+".name = ? AND "+
			propTableName+".string_value = ? AND "+
			tableName+".type_id = ?",
			"source_id", sourceID, config.TypeID)

	return config.DB.Where("id IN (?)", subQuery).Delete(&schema.Artifact{}).Error
}

func (r *AgentArtifactRepositoryImpl) DeleteByID(id int32) error {
	config := r.GetConfig()
	result := config.DB.Where("id = ? AND type_id = ?", id, config.TypeID).Delete(&schema.Artifact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: id %d", config.NotFoundError, id)
	}
	return nil
}

func (r *AgentArtifactRepositoryImpl) GetDistinctSourceIDs() ([]string, error) {
	config := r.GetConfig()
	var sourceIDs []string

	propTableName := utils.GetTableName(config.DB, &schema.ArtifactProperty{})
	tableName := utils.GetTableName(config.DB, &schema.Artifact{})

	err := config.DB.Table(propTableName+" cp").
		Select("DISTINCT cp.string_value").
		Joins("INNER JOIN "+tableName+" c ON cp.artifact_id = c.id").
		Where("cp.name = ? AND c.type_id = ?", "source_id", config.TypeID).
		Pluck("string_value", &sourceIDs).Error

	if err != nil {
		err = dbutil.SanitizeDatabaseError(err)
		return nil, fmt.Errorf("error querying distinct source IDs: %w", err)
	}
	return sourceIDs, nil
}

func (r *AgentArtifactRepositoryImpl) GetTypeID() int32 {
	return r.GetConfig().TypeID
}

