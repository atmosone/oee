<template>
  <v-treeview
    :items="items"
    item-title="title"
    item-value="value"
    hide-actions
    indent-lines
    separate-roots
    open-on-click
    open-all
    color="primary"
  >
    <template #header="{ props }">
      <v-treeview-item v-bind="props">
        <template #prepend>
          <v-icon>mdi-table</v-icon>
        </template>
        <template #append>
          <create-equipment @created="refreshList"/>
        </template>
      </v-treeview-item>
    </template>
    <template #item="{ item, props }">
      <v-treeview-item v-bind="props" :to="item.value" :subtitle="item.subtitle">
        <template #prepend>
          <v-icon>mdi-table-row</v-icon>
        </template>
        <template #append>
          <v-chip append-icon="mdi-link-variant" color="primary" variant="tonal" size="small">
            {{ item.badge }}
          </v-chip>
        </template>
      </v-treeview-item>
    </template>
  </v-treeview>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import type { components } from '../schema.ts'
import { mapStores } from 'pinia'
import useAppStore from '../store.ts'
import CreateEquipment from './CreateEquipment.vue'

export default {
  name: 'EquipmentTree',
  components: {
    CreateEquipment,
  },
  props: {
    total: Number as PropType<number>,
  },
  data: () => ({
    list: [] as components['schemas']['EquipmentList'],
  }),
  computed: {
    ...mapStores(useAppStore),
    items() {
      return [
        {
          value: 'equipment_list',
          title: 'Оборудование',
          subtitle: '',
          badge: this.total,
          children: this.list.map((item) => ({
            value: `/equipment/${item.id}`,
            title: item.name,
            subtitle: `ID:${item.id},Units:${item.units}`,
            badge: item.total_qualification_rules,
          })),
        },
      ]
    },
  },
  methods: {
    async refreshList(): Promise<void> {
      const res = await this.appStore.client.GET('/queries/equipment_list')
      this.list = res.data ? res.data : []
    },
  },
  mounted() {
    this.refreshList()
  },
}
</script>
