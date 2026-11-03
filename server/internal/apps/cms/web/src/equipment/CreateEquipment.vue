<template>
  <v-dialog width="350">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        color="primary"
        variant="text"
        icon="mdi-plus"
        rounded="lg"
        size="small"
      />
    </template>
    <v-card rounded="lg">
      <v-card-item class="border-b">
        <v-card-title>Новое оборудование</v-card-title>
        <v-card-subtitle>Регистрация нового оборудования</v-card-subtitle>
      </v-card-item>
      <v-card-text class="pa-4">
        <v-text-field v-model="id" label="Идентификатор" color="primary" />
        <v-text-field v-model="name" label="Наименование" color="primary" />
        <v-text-field v-model="units" label="Единицы измерения продукции" color="primary" />
      </v-card-text>
      <v-card-actions class="pa-0">
        <v-btn
          @click="submit"
          :disabled="disabled"
          variant="flat"
          text="Создать"
          color="primary"
          height="52"
          rounded="0"
          block
        />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { mapStores } from 'pinia'
import type { components } from '@/schema'
import useAppStore from '@/store'

export default {
  name: 'AppCreateEquipment',
  emits: {
    created: (event: void) => event,
  },
  data: () => ({
    id: null as components['schemas']['EquipmentId'] | null,
    name: null as components['schemas']['EquipmentName'] | null,
    units: null as components['schemas']['EquipmentUnits'] | null,
  }),
  computed: {
    ...mapStores(useAppStore),
    disabled(): boolean {
      return !this.id || !this.name || !this.units
    },
  },
  methods: {
    async submit() {
      if (!this.id || !this.name || !this.units) return
      const res = await this.appStore.client.POST('/commands/create_equipment', {
        body: {
          equipment_id: this.id,
          equipment_name: this.name,
          equipment_units: this.units,
        },
      })
      if (!res.response.ok) return
      this.$emit('created')
    },
  },
}
</script>
