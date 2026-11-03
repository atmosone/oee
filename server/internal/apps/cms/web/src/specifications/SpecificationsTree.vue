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
      <v-progress-linear v-show="loading" color="primary" indeterminate />
      <v-treeview-item v-bind="props">
        <template #prepend>
          <v-icon>mdi-table</v-icon>
        </template>
        <template #append>
          <create-specification @created="refreshList"/>
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
import type { components } from '@/schema'
import { mapStores } from 'pinia'
import useAppStore from '@/store'
import CreateSpecification from './CreateSpecification.vue';

export default {
  name: 'AppEquipmentTree',
  components: {
    CreateSpecification,
  },
  props: {
    total: Number as PropType<number>,
  },
  data: () => ({
    loading: false,
    list: [] as components['schemas']['SpecificationsList'],
  }),
  computed: {
    ...mapStores(useAppStore),
    items() {
      return [
        {
          value: 'specifications_list',
          title: 'Спецификации',
          subtitle: 'Table',
          badge: this.total,
          children: this.list.map((item) => ({
            value: `/specifications/${item.id}`,
            subtitle: `ID:${item.id},Code:${item.code}`,
            title: item.description,
            badge: item.total_qualification_rules,
          })),
        },
      ]
    },
  },
  methods: {
    async refreshList(): Promise<void> {
      this.loading = true
      const res = await this.appStore.client.GET('/queries/specifications_list')
      this.list = res.data ? res.data : []
      this.loading = false
    },
  },
  mounted() {
    this.refreshList()
  },
}
</script>
