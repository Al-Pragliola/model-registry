import { kebabTableColumn, SortableData } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';

export const agentCatalogSourceConfigsColumns: SortableData<CatalogSourceConfig>[] = [
  {
    field: 'name',
    label: 'Source name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 25,
  },
  {
    field: 'type',
    label: 'Source type',
    sortable: (a, b) => a.type.localeCompare(b.type),
    width: 15,
  },
  {
    field: 'enabled',
    label: 'Enabled',
    sortable: (a, b) => Number(a.enabled ?? true) - Number(b.enabled ?? true),
    width: 10,
  },
  {
    field: 'status',
    label: 'Status',
    sortable: false,
    width: 15,
  },
  {
    field: 'actions',
    label: '',
    sortable: false,
    width: 20,
  },
  kebabTableColumn(),
];
