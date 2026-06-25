import * as React from 'react';
import {
  Form,
  FormGroup,
  Checkbox,
  Stack,
  StackItem,
  Alert,
  PageSection,
  ActionList,
  ActionListGroup,
  ActionListItem,
  Button,
  TextInput,
  FormHelperText,
  HelperText,
  HelperTextItem,
  FileUpload,
  Sidebar,
  SidebarPanel,
  SidebarContent,
} from '@patternfly/react-core';
import { OpenDrawerRightIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { agentCatalogSettingsUrl } from '~/app/routes/agentCatalogSettings/agentCatalogSettings';
import { AgentCatalogSettingsContext } from '~/app/context/agentCatalogSettings/AgentCatalogSettingsContext';
import { CatalogSourceConfig, CatalogSourceType } from '~/app/modelCatalogTypes';
import { generateSourceIdFromName } from '~/app/pages/modelCatalogSettings/utils/modelCatalogSettingsUtils';
import {
  SOURCE_NAME_CHARACTER_LIMIT,
  EXPECTED_YAML_FORMAT_LABEL,
} from '~/app/pages/modelCatalogSettings/constants';
import {
  ManageAgentSourceFormData,
  useManageAgentSourceData,
} from '~/app/pages/agentCatalogSettings/useManageAgentSourceData';
import AgentVisibilitySection from './AgentVisibilitySection';
import AgentPreviewPanel from './AgentPreviewPanel';
import { useAgentSourcePreview } from '~/app/pages/agentCatalogSettings/useAgentSourcePreview';

type ManageAgentSourceFormProps = {
  existingSourceConfig?: CatalogSourceConfig;
  isEditMode: boolean;
  onToggleExpectedFormatDrawer?: () => void;
};

const agentSourceConfigToFormData = (
  sourceConfig: CatalogSourceConfig,
): Partial<ManageAgentSourceFormData> => ({
  name: sourceConfig.name,
  enabled: sourceConfig.enabled ?? true,
  includedAgents: ((sourceConfig as Record<string, unknown>).includedAgents as string[] || []).join(', '),
  excludedAgents: ((sourceConfig as Record<string, unknown>).excludedAgents as string[] || []).join(', '),
  isDefault: sourceConfig.isDefault,
  id: sourceConfig.id,
  yamlContent: sourceConfig.type === CatalogSourceType.YAML ? (sourceConfig.yaml ?? '') : '',
});

const parseCommaSeparatedList = (value: string): string[] =>
  value
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0);

const transformAgentFormDataToConfig = (
  formData: ManageAgentSourceFormData,
  existingSourceConfig?: CatalogSourceConfig,
): Record<string, unknown> => {
  const config: Record<string, unknown> = {
    id: formData.id || generateSourceIdFromName(formData.name),
    name: formData.name,
    type: CatalogSourceType.YAML,
    enabled: formData.enabled,
    isDefault: formData.isDefault,
    yaml: formData.yamlContent,
  };

  if (existingSourceConfig?.type === CatalogSourceType.YAML) {
    config.yamlCatalogPath = existingSourceConfig.yamlCatalogPath;
  }

  if (formData.includedAgents) {
    config.includedAgents = parseCommaSeparatedList(formData.includedAgents);
  }
  if (formData.excludedAgents) {
    config.excludedAgents = parseCommaSeparatedList(formData.excludedAgents);
  }

  return config;
};

const getAgentPayloadForConfig = (
  config: Record<string, unknown>,
  isEditMode: boolean,
  isDefault: boolean,
): Record<string, unknown> => {
  if (isDefault) {
    return {
      enabled: config.enabled,
      includedAgents: config.includedAgents,
      excludedAgents: config.excludedAgents,
    };
  }

  if (isEditMode) {
    return {
      name: config.name,
      type: config.type,
      enabled: config.enabled,
      isDefault: config.isDefault,
      yaml: config.yaml,
      includedAgents: config.includedAgents,
      excludedAgents: config.excludedAgents,
    };
  }

  return config;
};

const isAgentFormValid = (formData: ManageAgentSourceFormData): boolean => {
  const nameValid =
    formData.name.trim().length > 0 && formData.name.length <= SOURCE_NAME_CHARACTER_LIMIT;
  const yamlValid = formData.isDefault || formData.yamlContent.trim().length > 0;
  return nameValid && yamlValid;
};

const ManageAgentSourceForm: React.FC<ManageAgentSourceFormProps> = ({
  existingSourceConfig,
  isEditMode,
  onToggleExpectedFormatDrawer,
}) => {
  const navigate = useNavigate();
  const { isMUITheme } = useThemeContext();
  const existingData = existingSourceConfig
    ? agentSourceConfigToFormData(existingSourceConfig)
    : undefined;
  const [formData, setData] = useManageAgentSourceData(existingData);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [submitError, setSubmitError] = React.useState<Error | undefined>(undefined);
  const { apiState, refreshCatalogSourceConfigs } = React.useContext(AgentCatalogSettingsContext);

  const preview = useAgentSourcePreview({
    formData,
    existingSourceConfig,
    apiState,
    isEditMode,
  });

  // Name validation state
  const [isNameTouched, setIsNameTouched] = React.useState(false);
  const isNameValid =
    formData.name.trim().length > 0 && formData.name.length <= SOURCE_NAME_CHARACTER_LIMIT;
  const hasNameError = isNameTouched && !isNameValid;

  // YAML validation state
  const [isYamlTouched, setIsYamlTouched] = React.useState(false);
  const [filename, setFilename] = React.useState('');
  const [fileUploadError, setFileUploadError] = React.useState<string | undefined>(undefined);
  const isYamlContentValid = formData.yamlContent.trim().length > 0;

  const isFormComplete = isAgentFormValid(formData);

  const handleFileChange = (
    _event: React.DragEvent<HTMLElement> | React.ChangeEvent<HTMLInputElement> | Event,
    file: File,
  ) => {
    setFilename(file.name);
    setFileUploadError(undefined);
    const reader = new FileReader();
    reader.onload = () => {
      const text = typeof reader.result === 'string' ? reader.result : '';
      setData('yamlContent', text);
      setIsYamlTouched(true);
    };
    reader.onerror = () => {
      setFileUploadError("The YAML file couldn't be uploaded. Check its syntax and structure, then try again.");
    };
    reader.readAsText(file);
  };

  const handleTextChange = (_event: React.ChangeEvent<HTMLTextAreaElement>, value: string) => {
    setData('yamlContent', value);
  };

  const handleClear = () => {
    setFilename('');
    setData('yamlContent', '');
    setIsYamlTouched(true);
    setFileUploadError(undefined);
  };

  const handleSubmit = async () => {
    if (!apiState.apiAvailable) {
      setSubmitError(new Error('API is not available'));
      return;
    }
    setIsSubmitting(true);
    setSubmitError(undefined);

    try {
      const sourceConfig = transformAgentFormDataToConfig(formData, existingSourceConfig);
      const payload = getAgentPayloadForConfig(sourceConfig, isEditMode, formData.isDefault);

      if (isEditMode) {
        await apiState.api.updateCatalogSourceConfig({}, formData.id, payload);
      } else {
        await apiState.api.createCatalogSourceConfig({}, payload);
      }

      refreshCatalogSourceConfigs();
      navigate(agentCatalogSettingsUrl());
    } catch (error) {
      setSubmitError(error instanceof Error ? error : new Error('Failed to save source'));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    navigate(agentCatalogSettingsUrl());
  };

  const expectedFormatButton = onToggleExpectedFormatDrawer ? (
    <Button
      variant="link"
      isInline
      onClick={onToggleExpectedFormatDrawer}
      data-testid="view-expected-agent-yaml-format-link"
      icon={<OpenDrawerRightIcon />}
      iconPosition="end"
    >
      {EXPECTED_YAML_FORMAT_LABEL}
    </Button>
  ) : null;

  const yamlInput = (
    <div data-testid="agent-yaml-content-input">
      <FileUpload
        id="agent-yaml-content"
        type="text"
        value={formData.yamlContent}
        filename={filename}
        filenamePlaceholder="Drag and drop a YAML file or upload one"
        onFileInputChange={handleFileChange}
        onTextChange={handleTextChange}
        onClearClick={handleClear}
        onBlur={() => setIsYamlTouched(true)}
        validated={isYamlTouched && !isYamlContentValid ? 'error' : 'default'}
        browseButtonText="Upload"
        allowEditingUploadedText
        dropzoneProps={{
          accept: { 'text/yaml': ['.yaml', '.yml'] },
        }}
      />
    </div>
  );

  const yamlHelperTxtNode =
    isYamlTouched && !isYamlContentValid ? (
      <FormHelperText>
        <HelperText>
          <HelperTextItem variant="error" data-testid="agent-yaml-content-error">
            YAML content is required
          </HelperTextItem>
        </HelperText>
      </FormHelperText>
    ) : (
      <FormHelperText>
        <HelperText>
          <HelperTextItem>Upload or paste a YAML string.</HelperTextItem>
        </HelperText>
      </FormHelperText>
    );

  return (
    <>
      <Sidebar hasBorder isPanelRight hasGutter>
        <SidebarContent>
      <Form isWidthLimited>
        <Stack hasGutter>
          {/* Source name section */}
          <StackItem>
            <FormSection>
              <FormGroup label="Name" isRequired fieldId="agent-source-name">
                <FormFieldset
                  component={
                    <TextInput
                      isRequired
                      readOnlyVariant={formData.isDefault ? 'plain' : undefined}
                      type="text"
                      id="agent-source-name"
                      name="agent-source-name"
                      data-testid="agent-source-name-input"
                      value={formData.name}
                      onChange={(_event, value) => setData('name', value)}
                      onBlur={() => setIsNameTouched(true)}
                      validated={hasNameError ? 'error' : 'default'}
                    />
                  }
                  field="Name"
                  data-testid="agent-source-name-readonly"
                />
                {hasNameError && (
                  <FormHelperText>
                    <HelperText>
                      <HelperTextItem variant="error" data-testid="agent-source-name-error">
                        {formData.name.trim().length === 0
                          ? 'Name is required'
                          : formData.name.length > SOURCE_NAME_CHARACTER_LIMIT
                            ? `Cannot exceed ${SOURCE_NAME_CHARACTER_LIMIT} characters`
                            : null}
                      </HelperTextItem>
                    </HelperText>
                  </FormHelperText>
                )}
              </FormGroup>
            </FormSection>
          </StackItem>

          {/* YAML section (not shown for default sources) */}
          {!formData.isDefault && (
            <StackItem>
              <FormSection data-testid="agent-yaml-section">
                {fileUploadError && (
                  <Alert
                    variant="danger"
                    isInline
                    title="File upload failed"
                    className="pf-v6-u-mb-md"
                    data-testid="agent-yaml-file-upload-error"
                  >
                    {fileUploadError}
                  </Alert>
                )}
                {isMUITheme && (
                  <div className="pf-v6-u-mb-sm" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span>Upload a YAML file</span>
                    <span>{expectedFormatButton}</span>
                  </div>
                )}
                <FormGroup
                  label={!isMUITheme ? 'Upload a YAML file' : undefined}
                  labelInfo={!isMUITheme ? expectedFormatButton : undefined}
                  isRequired
                  fieldId="agent-yaml-content"
                >
                  <FormFieldset component={yamlInput} field="YAML" />
                  {yamlHelperTxtNode}
                </FormGroup>
              </FormSection>
            </StackItem>
          )}

          {/* Agent visibility section */}
          <StackItem>
            <AgentVisibilitySection
              formData={formData}
              setData={setData as UpdateObjectAtPropAndValue<ManageAgentSourceFormData>}
              isDefaultExpanded={
                existingData?.isDefault ||
                !!existingData?.includedAgents ||
                !!existingData?.excludedAgents
              }
            />
          </StackItem>

          {/* Enable source checkbox */}
          <StackItem>
            <FormSection>
              <FormGroup fieldId="agent-enable-source">
                <Checkbox
                  label={
                    <span className="pf-v6-c-form__label-text">Enable source</span>
                  }
                  id="agent-enable-source"
                  name="agent-enable-source"
                  data-testid="agent-enable-source-checkbox"
                  description="Enable users in your organization to view agents from this source in the agent catalog."
                  isChecked={formData.enabled}
                  onChange={(_event, checked) => setData('enabled', checked)}
                />
              </FormGroup>
            </FormSection>
          </StackItem>
        </Stack>
      </Form>
        </SidebarContent>
        <SidebarPanel width={{ default: 'width_50' }}>
          <AgentPreviewPanel preview={preview} />
        </SidebarPanel>
      </Sidebar>

      {/* Footer */}
      <PageSection hasBodyWrapper={false} stickyOnBreakpoint={{ default: 'bottom' }}>
        <Stack hasGutter>
          {submitError && (
            <StackItem>
              <Alert variant="danger" isInline title="Failed to save source">
                {submitError.message}
              </Alert>
            </StackItem>
          )}
          <StackItem>
            <ActionList>
              <ActionListGroup>
                <ActionListItem>
                  <Button
                    isDisabled={!isFormComplete || isSubmitting}
                    variant="primary"
                    id="agent-submit-button"
                    data-testid="agent-submit-button"
                    isLoading={isSubmitting}
                    onClick={handleSubmit}
                  >
                    {isEditMode ? 'Save' : 'Add'}
                  </Button>
                </ActionListItem>
                <ActionListItem>
                  <Button
                    isDisabled={isSubmitting}
                    variant="link"
                    id="agent-cancel-button"
                    data-testid="agent-cancel-button"
                    onClick={handleCancel}
                  >
                    Cancel
                  </Button>
                </ActionListItem>
              </ActionListGroup>
            </ActionList>
          </StackItem>
        </Stack>
      </PageSection>
    </>
  );
};

export default ManageAgentSourceForm;
