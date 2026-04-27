import { CatalogModel } from '~/app/modelCatalogTypes';
import {
  ModelRegistryCustomProperties,
  ModelRegistryCustomProperty,
  ModelRegistryMetadataType,
} from '~/app/types';

const FIXED_COLUMNS = [
  'Name',
  'Source ID',
  'Provider',
  'Description',
  'Maturity',
  'License',
  'License Link',
  'Library',
  'Language',
  'Tasks',
  'Created',
  'Last Updated',
] as const;

export const escapeCSVField = (value: string): string => {
  if (value.includes('"') || value.includes(',') || value.includes('\n') || value.includes('\r')) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
};

const epochMsToISO = (epoch: string | undefined): string => {
  if (!epoch) {
    return '';
  }
  const ms = parseInt(epoch, 10);
  if (isNaN(ms)) {
    return epoch;
  }
  return new Date(ms).toISOString();
};

const extractCustomPropertyValue = (prop: ModelRegistryCustomProperty): string => {
  switch (prop.metadataType) {
    case ModelRegistryMetadataType.STRING:
      return prop.string_value;
    case ModelRegistryMetadataType.INT:
      return prop.int_value;
    case ModelRegistryMetadataType.DOUBLE:
      return String(prop.double_value);
    case ModelRegistryMetadataType.BOOL:
      return String(prop.bool_value);
    default:
      return '';
  }
};

const collectCustomPropertyKeys = (models: CatalogModel[]): string[] => {
  const keySet = new Set<string>();
  for (const model of models) {
    if (model.customProperties) {
      for (const key of Object.keys(model.customProperties)) {
        keySet.add(key);
      }
    }
  }
  return Array.from(keySet).sort();
};

const getCustomPropertyValue = (
  props: ModelRegistryCustomProperties | undefined,
  key: string,
): string => {
  if (!props || !(key in props)) {
    return '';
  }
  return extractCustomPropertyValue(props[key]);
};

const modelToRow = (model: CatalogModel, customPropKeys: string[]): string[] => [
  model.name,
  model.source_id ?? '',
  model.provider ?? '',
  model.description ?? '',
  model.maturity ?? '',
  model.license ?? '',
  model.licenseLink ?? '',
  model.libraryName ?? '',
  (model.language ?? []).join('; '),
  (model.tasks ?? []).join('; '),
  epochMsToISO(model.createTimeSinceEpoch),
  epochMsToISO(model.lastUpdateTimeSinceEpoch),
  ...customPropKeys.map((key) => getCustomPropertyValue(model.customProperties, key)),
];

export const generateCatalogModelCSV = (models: CatalogModel[]): string => {
  const customPropKeys = collectCustomPropertyKeys(models);
  const headers = [...FIXED_COLUMNS, ...customPropKeys];

  const lines = [
    headers.map(escapeCSVField).join(','),
    ...models.map((model) => modelToRow(model, customPropKeys).map(escapeCSVField).join(',')),
  ];

  // UTF-8 BOM for Excel compatibility
  return '﻿' + lines.join('\r\n') + '\r\n';
};

export const downloadCSV = (content: string, filename: string): void => {
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

export const getExportFilename = (): string => {
  const date = new Date().toISOString().split('T')[0];
  return `model-catalog-export-${date}.csv`;
};
