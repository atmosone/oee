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
        <v-card-title>Новая спецификация</v-card-title>
        <v-card-subtitle>Регистрация новой спецификации</v-card-subtitle>
      </v-card-item>
      <v-card-text class="pa-4">
        <v-text-field v-model="id" label="Идентификатор" color="primary" />
        <v-text-field v-model="code" label="Код" color="primary" />
        <v-textarea v-model="description" label="Описание" color="primary" rows="2"/>
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
import { mapStores } from 'pinia';
import type { components } from '../schema';
import useAppStore from '../store';

export default {
    name: "AppCreateEquipmentLine",
    emits: {
        created: (event: void) => event,
    },
    data: () => ({
        id: null as components['schemas']['SpecificationId'] | null,
        code: null as components['schemas']['SpecificationCode'] | null,
        description: null as components['schemas']['SpecificationDescription'] | null,
    }),
    computed: {
        ...mapStores(useAppStore),
        disabled(): boolean { return !this.id || !this.code || !this.description }
    },
    methods: {
        async submit(): Promise<void> {
            if (!this.id || !this.code || !this.description) return
            const res = await this.appStore.client.POST('/commands/create_specification', {
                body: {specification_id: this.id, specification_code: this.code, specification_description: this.description}
            })
            if (!res.response.ok) return
            this.$emit('created')
        }
    }
}
</script>