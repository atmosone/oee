<template>
  <v-dialog width="350" @update:model-value="onUpdate($event)">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon="mdi-plus"
        color="primary"
        variant="tonal"
      />
    </template>
    <v-card rounded="lg">
      <v-card-item class="border-b">
        <v-card-title>{{ equipmentName }}</v-card-title>
        <v-card-subtitle>Добавить спецификацию</v-card-subtitle>
      </v-card-item>
      <v-card-text class="pa-0">
        <v-select v-model="specificationId"
        :items="specificationsList"
        rounded="0"
        item-value="id"
        item-title="description"
        label="Спецификация"
        color="primary"
        hide-details>
          <template #item="{item, props}">
            <v-list-item v-bind="props" 
            :subtitle="`ID:${item.id},Code:${item.code}`"
            prepend-icon="mdi-link-variant"/>
          </template>
        </v-select>
      </v-card-text>
      <v-card-actions class="pa-0">
        <v-btn @click="submit"
        :disabled="disabled"
        text="Добавить"
        color="primary"
        variant="flat"
        height="52"
        rounded="0"
        block/>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import type { PropType } from 'vue';
import type { components } from '@/schema';
import { mapStores } from 'pinia';
import useAppStore from '@/store';

export default {
  name: 'AttachSpecification',
  props: {
    equipmentId: {
      type: String as PropType<components['schemas']['EquipmentId']>,
      required: true,
    },
    equipmentName: {
      type: String as PropType<components['schemas']['EquipmentName']>,
      required: true
    }
  },
  emits: {
    attached: (event: void) => event,
  },
  data: () => ({
    specificationsList: [] as components['schemas']['AttachableSpecificationsList'],
    specificationId: null as components['schemas']['SpecificationId'] | null,
  }),
  computed: {
    ...mapStores(useAppStore),
    disabled(): boolean {return !this.specificationId}
  },
  methods: {
    onUpdate(opened: boolean) {
      if (!opened) return
      this.refreshSpecificationsList()
    },
    async refreshSpecificationsList(): Promise<void> {
      const res = await this.appStore.client.GET('/queries/attachable_specifications', {
        params: {query: {to: this.equipmentId}},
      })
      this.specificationsList = res.data ? res.data : []
    },
    async submit(): Promise<void> {
      if (!this.specificationId) return
      const res = await this.appStore.client.POST('/commands/create_qualification_rule', {
        body: {equipment_id: this.equipmentId, specification_id: this.specificationId}
      })
      if (!res.response.ok) return
      this.$emit('attached')
    },
  },
}
</script>
