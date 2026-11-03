<template>
  <v-container v-if="details">
    <v-row>
      <v-col cols="12">
        <v-card-item>
          <v-card-title class="text-display-large font-weight-black">
            {{ details.name }}
          </v-card-title>
          <v-card-subtitle class="text-body-large">#{{ id }}</v-card-subtitle>
        </v-card-item>
      </v-col>
    </v-row>
    <v-row class="align-center">
      <v-col cols="12">
        <v-label>Оборудование</v-label>
      </v-col>
      <v-col v-for="item in details.equipment_list" cols="12" md="4" xl="3">
        <v-card :to="`/equipment/${item.id}`"
        :title="item.name" 
        :subtitle="item.id" 
        prepend-icon="mdi-link-variant" 
        rounded="lg"/>
      </v-col>
      <expand-line 
      :equipment-line-id="id" 
      :equipment-line-name="details.name"
      @expanded="refresh(id)"/>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { mapStores } from 'pinia'
import useAppStore from '../store.ts'
import type { components } from '../schema.ts'
import ExpandLine from './ExpandLine.vue';

export default {
  name: 'AppLinePage',
  components: {
    ExpandLine,
  },
  data: () => ({
    details: null as components['schemas']['EquipmentLineDetails'] | null,
  }),
  computed: {
    ...mapStores(useAppStore),
    id(): components['schemas']['EquipmentLineId'] {
      return this.$route.params['id'] as components['schemas']['EquipmentLineId']
    },
  },
  methods: {
    async refresh(id: components['schemas']['EquipmentLineId']): Promise<void> {
      const res = await this.appStore.client.GET('/queries/equipment_line_details', {
        params: { query: { id: id } },
      })
      this.details = res.data ? res.data : null
    },
  },
  mounted() { this.refresh(this.id) },
  beforeRouteUpdate(to, from, next) {
    this.refresh(to.params["id"] as components['schemas']['EquipmentLineId'])
    next()
  }
}
</script>
