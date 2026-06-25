package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kubeflow/hub/ui/bff/internal/constants"
	k8s "github.com/kubeflow/hub/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/hub/ui/bff/internal/models"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type AgentCatalogSettingsRepository struct{}

func NewAgentCatalogSettingsRepository() *AgentCatalogSettingsRepository {
	return &AgentCatalogSettingsRepository{}
}

func (r *AgentCatalogSettingsRepository) GetAllAgentCatalogSourceConfigs(ctx context.Context, client k8s.KubernetesClientInterface, namespace string) (*models.CatalogSourceConfigList, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	catalogMap := make(map[string]models.CatalogSourceConfig)

	if raw, ok := defaultCM.Data[k8s.CatalogSourceKey]; ok {
		defaultSources, err := ParseCatalogYamlSection(raw, true, SectionKeyAgentCatalogs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default agent catalogs: %w", err)
		}
		for _, catalog := range defaultSources {
			catalogMap[catalog.Id] = catalog
		}
	}

	if raw, ok := userCM.Data[k8s.CatalogSourceKey]; ok {
		userSources, err := ParseCatalogYamlSection(raw, false, SectionKeyAgentCatalogs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user agent catalogs: %w", err)
		}
		for _, userSource := range userSources {
			if existing, exists := catalogMap[userSource.Id]; exists {
				catalogMap[userSource.Id] = mergeCatalogSourceConfigs(existing, userSource)
			} else {
				catalogMap[userSource.Id] = userSource
			}
		}
	}

	result := &models.CatalogSourceConfigList{
		Catalogs: make([]models.CatalogSourceConfig, 0),
	}
	for _, c := range catalogMap {
		result.Catalogs = append(result.Catalogs, c)
	}

	return result, nil
}

func (r *AgentCatalogSettingsRepository) GetAgentCatalogSourceConfig(ctx context.Context, client k8s.KubernetesClientInterface, namespace string, catalogSourceId string) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	defaultSource := FindCatalogSourceByIdSection(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId, true, SectionKeyAgentCatalogs)
	userSource := FindCatalogSourceByIdSection(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, false, SectionKeyAgentCatalogs)

	var result *models.CatalogSourceConfig
	if userSource != nil {
		if defaultSource != nil {
			merged := mergeCatalogSourceConfigs(*defaultSource, *userSource)
			result = &merged
		} else {
			result = userSource
		}
	} else if defaultSource != nil {
		result = defaultSource
	} else {
		return nil, fmt.Errorf("%w, %s", ErrCatalogSourceNotFound, catalogSourceId)
	}

	_, yamlFilePath := FindCatalogSourcePropertiesSection(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, SectionKeyAgentCatalogs)
	if yamlFilePath == "" {
		_, yamlFilePath = FindCatalogSourcePropertiesSection(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId, SectionKeyAgentCatalogs)
	}

	if result.Type == CatalogTypeYaml {
		if yamlFilePath != "" {
			result.YamlCatalogPath = &yamlFilePath
			if yamlContent, ok := userCM.Data[yamlFilePath]; ok {
				result.Yaml = &yamlContent
			} else if yamlContent, ok := defaultCM.Data[yamlFilePath]; ok {
				result.Yaml = &yamlContent
			} else if result.IsDefault == nil || !*result.IsDefault {
				sessionLogger := ctx.Value(constants.TraceLoggerKey).(*slog.Logger)
				sessionLogger.Warn("yaml agent catalog content missing from configmap",
					"catalogId", catalogSourceId,
					"expectedPath", yamlFilePath,
				)
			}
		}
	}

	return result, nil
}

func (r *AgentCatalogSettingsRepository) CreateAgentCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	if err := validateAgentSourcePayload(payload); err != nil {
		return nil, err
	}

	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	if FindCatalogSourceByIdSection(defaultCM.Data[k8s.CatalogSourceKey], payload.Id, true, SectionKeyAgentCatalogs) != nil {
		return nil, fmt.Errorf("%w: '%s' already exists in default sources", ErrCatalogSourceAlreadyExist, payload.Id)
	}
	if FindCatalogSourceByIdSection(userCM.Data[k8s.CatalogSourceKey], payload.Id, false, SectionKeyAgentCatalogs) != nil {
		return nil, fmt.Errorf("%w: '%s' already exists in user managed sources", ErrCatalogSourceAlreadyExist, payload.Id)
	}

	if payload.Type != CatalogTypeYaml {
		return nil, fmt.Errorf("%w: agent sources only support yaml type", ErrUnsupportedCatalogType)
	}

	yamlFileName := fmt.Sprintf("%s.yaml", payload.Id)
	yamlContent := make(map[string]string)
	yamlContent[yamlFileName] = *payload.Yaml

	newEntry := ConvertSourceConfigToYamlEntry(payload, yamlFileName, "")

	existingConfigMapEntry := userCM.Data[k8s.CatalogSourceKey]
	updatedConfigMapEntry, err := AppendCatalogSourceToYamlSection(existingConfigMapEntry, newEntry, SectionKeyAgentCatalogs)
	if err != nil {
		return nil, fmt.Errorf("failed to append agent catalog to yaml: %w", err)
	}

	if userCM.Data == nil {
		userCM.Data = make(map[string]string)
	}
	userCM.Data[k8s.CatalogSourceKey] = updatedConfigMapEntry

	for key, value := range yamlContent {
		userCM.Data[key] = value
	}

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		if apierrors.IsConflict(err) {
			return nil, fmt.Errorf("%w: %v", ErrCatalogSourceConflict, err)
		}
		return nil, fmt.Errorf("failed to update user configmap: %w", err)
	}

	return r.GetAgentCatalogSourceConfig(ctx, client, namespace, payload.Id)
}

func (r *AgentCatalogSettingsRepository) UpdateAgentCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	sourceId string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	existingUserSource := FindCatalogSourceByIdSection(userCM.Data[k8s.CatalogSourceKey], sourceId, false, SectionKeyAgentCatalogs)
	existingDefaultSource := FindCatalogSourceByIdSection(defaultCM.Data[k8s.CatalogSourceKey], sourceId, true, SectionKeyAgentCatalogs)

	var isOverridingDefault bool

	if existingUserSource != nil {
		isOverridingDefault = existingDefaultSource != nil
	} else if existingDefaultSource != nil {
		isOverridingDefault = true
	} else {
		return nil, fmt.Errorf("%w: '%s'", ErrCatalogSourceNotFound, sourceId)
	}

	if isOverridingDefault {
		if err := validateUpdatePayloadForDefaultOverride(payload); err != nil {
			return nil, err
		}
	}

	var yamlFilePath string
	if !isOverridingDefault || existingUserSource != nil {
		_, yamlFilePath = FindCatalogSourcePropertiesSection(userCM.Data[k8s.CatalogSourceKey], sourceId, SectionKeyAgentCatalogs)
		if yamlFilePath == "" {
			_, yamlFilePath = FindCatalogSourcePropertiesSection(defaultCM.Data[k8s.CatalogSourceKey], sourceId, SectionKeyAgentCatalogs)
		}
	}

	if userCM.Data == nil {
		userCM.Data = make(map[string]string)
	}

	if !isOverridingDefault || existingUserSource != nil {
		if payload.Yaml != nil && *payload.Yaml != "" {
			if yamlFilePath == "" {
				yamlFilePath = fmt.Sprintf("%s.yaml", sourceId)
			}
			userCM.Data[yamlFilePath] = *payload.Yaml
		}
	}

	if existingUserSource != nil {
		updatedYAML, err := UpdateCatalogSourceInYAMLSection(
			userCM.Data[k8s.CatalogSourceKey],
			sourceId,
			payload,
			"",
			yamlFilePath,
			SectionKeyAgentCatalogs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update agent catalog in yaml: %w", err)
		}
		userCM.Data[k8s.CatalogSourceKey] = updatedYAML
	} else {
		overrideEntry := BuildOverrideEntryForDefaultSource(sourceId, payload)
		updatedYAML, err := AppendCatalogSourceToYamlSection(userCM.Data[k8s.CatalogSourceKey], overrideEntry, SectionKeyAgentCatalogs)
		if err != nil {
			return nil, fmt.Errorf("failed to append override entry: %w", err)
		}
		userCM.Data[k8s.CatalogSourceKey] = updatedYAML
	}

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		if apierrors.IsConflict(err) {
			return nil, fmt.Errorf("%w: %v", ErrCatalogSourceConflict, err)
		}
		return nil, fmt.Errorf("failed to update user configmap: %w", err)
	}

	return r.GetAgentCatalogSourceConfig(ctx, client, namespace, sourceId)
}

func (r *AgentCatalogSettingsRepository) DeleteAgentCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	defaultCM, userCM, err := client.GetAllCatalogSourceConfigs(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog source configmaps: %w", err)
	}

	if FindCatalogSourceByIdSection(defaultCM.Data[k8s.CatalogSourceKey], catalogSourceId, true, SectionKeyAgentCatalogs) != nil {
		return nil, fmt.Errorf("%w: '%s' is a default source", ErrCannotDeleteDefaultSource, catalogSourceId)
	}

	catalogSourceToDelete := FindCatalogSourceByIdSection(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, false, SectionKeyAgentCatalogs)
	if catalogSourceToDelete == nil {
		return nil, fmt.Errorf("%w: '%s' not found in user sources", ErrCatalogSourceNotFound, catalogSourceId)
	}

	_, yamlFilePath := FindCatalogSourcePropertiesSection(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, SectionKeyAgentCatalogs)

	if catalogSourceToDelete.Type == CatalogTypeYaml && yamlFilePath != "" {
		delete(userCM.Data, yamlFilePath)
	}

	updatedYAML, err := RemoveCatalogSourceFromYAMLSection(userCM.Data[k8s.CatalogSourceKey], catalogSourceId, SectionKeyAgentCatalogs)
	if err != nil {
		return nil, fmt.Errorf("failed to remove agent catalog from sources.yaml: %w", err)
	}
	userCM.Data[k8s.CatalogSourceKey] = updatedYAML

	err = client.UpdateCatalogSourceConfig(ctx, namespace, &userCM)
	if err != nil {
		if apierrors.IsConflict(err) {
			return nil, fmt.Errorf("%w: %v", ErrCatalogSourceConflict, err)
		}
		return nil, fmt.Errorf("failed to update configmap after deletion: %w", err)
	}

	return catalogSourceToDelete, nil
}

func validateAgentSourcePayload(payload models.CatalogSourceConfigPayload) error {
	if payload.Id == "" {
		return fmt.Errorf("%w", ErrCatalogSourceIdRequired)
	}
	if err := validateCatalogId(payload.Id); err != nil {
		return err
	}
	if payload.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidationFailed)
	}
	if payload.Type == "" {
		return fmt.Errorf("%w: type is required", ErrValidationFailed)
	}
	if payload.Type != CatalogTypeYaml {
		return fmt.Errorf("%w: agent sources only support yaml type", ErrValidationFailed)
	}
	if payload.Yaml == nil || *payload.Yaml == "" {
		return fmt.Errorf("%w: yaml field is required for agent sources", ErrValidationFailed)
	}
	return nil
}

// Ensure AgentCatalogSettingsRepository uses same error variables
var _ = errors.New
