<template>
    <v-dialog width="350" @update:model-value="onOpen">
        <template #activator="{props}">
            <v-btn v-bind="props"
            icon="mdi-plus"
            variant="tonal"
            color="primary"
            />
        </template>
        <v-card>
            <v-card-item>
                <v-card-title>
                    {{ equipmentLineName }}
                </v-card-title>
                <v-card-subtitle>
                    Расширить линию
                </v-card-subtitle>
            </v-card-item>
            <v-card-text class="pa-0">
                <v-select v-model="equipmentId"
                :items="equipmentList"
                item-title="name"
                item-value="id"
                hide-details>
                    <template #item="{props, item}">
                        <v-list-item v-bind="props" 
                        :subtitle="item.id" 
                        prepend-icon="mdi-link-variant"/>
                    </template>
                </v-select>
            </v-card-text>
            <v-card-actions class="pa-0">
                <v-btn
                  @click="submit"
                  :disabled="disabled"
                  variant="flat"
                  text="Расширить"
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
import useAppStore from '../store';
import { type PropType } from 'vue';
import type { components } from '../schema';

export default {
    name: "AppExpandLine",
    emits: {
        expanded: (event: void) => event,
    },
    props: {
        equipmentLineName: {
            type: String as PropType<components['schemas']['EquipmentLineName']>,
            required: true
        },
        equipmentLineId: {
            type: String as PropType<components['schemas']['EquipmentLineId']>,
            required: true
        }
    },
    data: () => ({
        equipmentId: null as components['schemas']['EquipmentId'] | null,
        equipmentList: [] as components['schemas']['ExpandableEquipmentList'],
    }),
    computed: {
        ...mapStores(useAppStore),
        disabled(): boolean { return !this.equipmentId }
    },
    methods: {
        onOpen(isOpened: boolean) {
            if (!isOpened) return
            this.refreshEquipmentList()
        },
        async refreshEquipmentList(): Promise<void> {
            const res = await this.appStore.client.GET('/queries/expandable_equipment', {
                params: {query: {to: this.equipmentLineId}}
            })
            this.equipmentList = res.data ? res.data : []
        },
        async submit(): Promise<void> {
            if (!this.equipmentId) return
            const res = await this.appStore.client.POST('/commands/expand_equipment_line', {
                body: {equipment_line_id: this.equipmentLineId, equipment_id: this.equipmentId}
            })
            if (!res.response.ok) return
            this.$emit('expanded')
        }
    },
}
</script>