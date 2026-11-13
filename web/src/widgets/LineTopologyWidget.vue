<template>
    <v-timeline direction="horizontal">
        <v-timeline-item v-for="item, i in topology" 
        dot-color="primary" 
        side="start" 
        size="x-small">
            <v-speed-dial location="right center">
                <template #activator="{props}">
                    <v-card v-bind="props" 
                    :title="item.name" 
                    :subtitle="item.category"
                    class="border" 
                    rounded="xl"/>
                </template>
                <v-btn v-if="i > 0"
                @click="move(item.id, 'left')"
                icon="mdi-chevron-left" 
                text="Сдвинуть влево"
                color="primary"/>
                <v-btn v-if="i < topology.length - 1"
                @click="move(item.id, 'right')"
                icon="mdi-chevron-right" 
                text="Сдвинуть вправо"
                color="primary"/>
                <v-btn
                @click="remove(item.id)"
                icon="mdi-delete"
                text="Удалить"
                color="error"/>
            </v-speed-dial>
        </v-timeline-item>
        <v-timeline-item side="start" dot-color="primary" size="x-small">
            <v-btn icon="mdi-plus"/>
        </v-timeline-item>
    </v-timeline>
</template>

<script lang="ts">
import type { components } from '@/api';
import useAppStore from '@/store';
import { mapStores } from 'pinia';
import type { PropType } from 'vue';
export default {
    name: "LineTopologyWidget",
    props: {
        id: {
            type: String as PropType<components['schemas']['LineID']>,
            required: true,
        },
    },
    data: () => ({
        topology: [] as components['schemas']['LineTopology'],
    }),
    computed: {
        ...mapStores(useAppStore),
        equipmentAddedEventType(): components['schemas']['EventType'] {
            return 'equipment_added_to_line'
        },
        equipmentMovedEventType(): components['schemas']['EventType'] {
            return 'equipment_moved_in_line'
        },
        equipmentRemovedEventType(): components['schemas']['EventType'] {
            return 'equipment_removed_from_line'
        }
    },
    methods: {
        async move(
            equipment_id: components['schemas']['EquipmentID'], 
            direction: components['schemas']['MoveEquipmentInLineCommand']['direction'],
        ) :Promise<void> {
            const res = await this.appStore.api.POST('/lines/commands/move_equipment', {
                body: {id: this.id, equipment_id, direction},
            })
        },
        async remove(equipment_id: components['schemas']['EquipmentID']): Promise<void> {
            const res = await this.appStore.api.POST('/lines/commands/remove_equipment', {
                body: {id: this.id, equipment_id},
            })
        },
        async refresh(): Promise<void> {
            const res = await this.appStore.api.GET('/lines/queries/topology', {
                params:{query:{id: this.id}}
            })
            this.topology = res.data ? res.data : []
        },
        onUpdate(event: MessageEvent): void { this.refresh() },
    },
    mounted() {
        this.appStore.sse.addEventListener(
            <components['schemas']['EventType']>"equipment_added_to_line", 
            this.onUpdate,
        )
        this.appStore.sse.addEventListener(
            <components['schemas']['EventType']>"equipment_moved_in_line", 
            this.onUpdate,
        )
        this.appStore.sse.addEventListener(
            <components['schemas']['EventType']>"equipment_removed_from_line",
             this.onUpdate,
            )
        this.refresh()
    },
    beforeUnmount() {
        this.appStore.sse.removeEventListener(
            <components['schemas']['EventType']>"equipment_added_to_line", 
            this.onUpdate,
        )
        this.appStore.sse.removeEventListener(
            <components['schemas']['EventType']>"equipment_moved_in_line", 
            this.onUpdate,
        )
        this.appStore.sse.removeEventListener(
            <components['schemas']['EventType']>"equipment_removed_from_line",
            this.onUpdate
        )
    }
}
</script>
