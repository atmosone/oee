<template>
Equipment Screen {{ $route.params["id"] }}
    <v-text-field v-model="workshift_id" clearable>
        <template #append>
            <v-btn @click="refresh"
            icon="mdi-refresh"/>
        </template>
    </v-text-field>
    <timeline-widget :series="timeline" 
    @select="console.log($event)"
    @hover="console.log($event)"
    />
</template>

<script lang="ts">
import type { components } from '@/api';
import useAppStore from '@/store';
import TimelineWidget from '@/widgets/TimelineWidget.vue';
import { mapStores } from 'pinia';

export default {
    name: "EquipmentScreen",
    components: {
        TimelineWidget,
    },
    data: () => ({
        workshift_id: null as components['schemas']['WorkshiftID'] | null,
        timeline: [] as components['schemas']['ApexchartsTimeline'][],
    }),
    computed: {
        ...mapStores(useAppStore),
        equipment_id(): components['schemas']['EquipmentID'] {
            return this.$route.params['id'] as components['schemas']['EquipmentID']
        },
    },
    methods: {
        async refresh(): Promise<void> {
            if (!this.workshift_id) return
            const res = await this.appStore.api.GET('/api/intervals/queries/apexcharts_timeline', {
                params: {query: {workshift_id: this.workshift_id, equipment_id: this.equipment_id}}
            })
            if (!res.response.ok) return
            this.timeline = res.data ? res.data : []
        }
    }
}
</script>