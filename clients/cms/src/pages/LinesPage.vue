<template>
   <v-app>
        <v-app-bar density="compact" class="border-b" elevation="0">
            <v-app-bar-title>Lines</v-app-bar-title>
        </v-app-bar>
        <v-navigation-drawer>
            <v-card-item>
                <v-card-subtitle>Explorer</v-card-subtitle>
                <template v-slot:append>
                   <create-line/>
                </template>
            </v-card-item>
            <v-list nav>
                <v-list-item v-for="item in lines" :to="`/lines/${item.id}`" :title="item.name">
                    <template v-slot:append>
                        <v-chip :text="item.total_equipment" color="primary" rounded="lg"/>
                    </template>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>
        <v-main>
            <router-view/>
        </v-main>
   </v-app>
    <create-equipment-line/>
    <expand-equipment-line/>
</template>

<script lang="ts">
import CreateLine from '@/components/CreateLine.vue';
import ExpandEquipmentLine from '@/components/ExpandEquipmentLine.vue';
import type { components } from '@/schema';
import useAppStore from '@/store';
import { mapStores } from 'pinia';
export default {
    name: "EquipmentLinesPage",
    components: {
        CreateLine,
        ExpandEquipmentLine,
    },
    data: () => ({
        lines: [] as components["schemas"]["EquipmentLinesListItem"][],
    }),
    computed: {...mapStores(useAppStore)},
    methods: {
        async refresh_lines(): Promise<void> {
            const data = await this.appStore.client.GET('/queries/equipment_lines_list')
            this.lines = data.data ? data.data.items : []
        }
    },
    mounted() {this.refresh_lines()}
}
</script>