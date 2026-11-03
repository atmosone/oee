<template>
  <v-container v-if="details">
    <v-row>
      <v-col cols="12">
        <v-card-item>
          <v-card-title class="text-display-large font-weight-black">
            {{ details.name }}
          </v-card-title>
          <v-card-subtitle class="text-body-large">
            {{ id }}
          </v-card-subtitle>
          <v-chip
            v-if="details.line"
            :to="`/lines/${details.line.id}`"
            class="mr-2"
            prepend-icon="mdi-link-variant"
            color="primary"
          >
            {{ details.line.name }}
          </v-chip>
          <v-chip color="primary">
            {{ details.units }}
          </v-chip>
        </v-card-item>
      </v-col>
    </v-row>
    <v-row class="align-center">
      <v-col cols="12">
        <v-label>Спецификации</v-label>
      </v-col>
      <v-col v-for="item in details.specifications_list" 
      cols="12" 
      md="4" 
      xl="3">
        <v-card :to="`/specifications/${item.id}`" 
        :title="item.code" 
        :subtitle="item.id" 
        :text="item.description"
        prepend-icon="mdi-link-variant"
        rounded="lg"/>
      </v-col>
      <attach-specification 
      :equipment-id="id" 
      :equipment-name="details.name" 
      @attached="refresh(id)"/>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { mapStores } from 'pinia'
import useAppStore from '../store.ts'
import type { components } from '../schema.ts'
import AttachSpecification from './AttachSpecification.vue'

export default {
  name: 'AppEquipmentPage',
  components: {
    AttachSpecification,
  },
  data: () => ({
    details: null as components['schemas']['EquipmentDetails'] | null,
  }),
  computed: {
    ...mapStores(useAppStore),
    id(): components['schemas']['EquipmentId'] {
      return this.$route.params['id'] as components['schemas']['EquipmentId']
    },
  },
  methods: {
    async refresh(id: components['schemas']['EquipmentId']): Promise<void> {
      const res = await this.appStore.client.GET('/queries/equipment_details', {
        params: { query: { id: id } },
      })
      this.details = res.data ? res.data : null
    },
  },
  mounted() {this.refresh(this.id)},
  beforeRouteUpdate(to, from, next) {
    this.refresh(to.params["id"] as components['schemas']['EquipmentId'])
    next()
  }
}
</script>
