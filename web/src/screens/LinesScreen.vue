<template>
    <v-navigation-drawer>
        <v-card-item>
            <template #append>
                <v-btn icon="mdi-plus"/>
            </template>
            <v-card-title>Линии</v-card-title>
        </v-card-item>
        <v-list nav>
            <v-list-item v-for="item in list" 
            prepend-icon="mdi-table-row"
            append-icon="mdi-chevron-right"
            :title="item.name" 
            :to="`/lines/${item.id}`"/>
        </v-list>
    </v-navigation-drawer>
    <router-view/>
</template>

<script lang="ts">
import type { components } from '@/api';
import useAppStore from '@/store';
import { mapStores } from 'pinia';

export default {
    name: "LinesPage",
    data: () => ({
        list: [] as components['schemas']['LinesList'],
    }),
    computed: {
        ...mapStores(useAppStore),
    },
    methods: {
        async refresh(): Promise<void> {
            const res = await this.appStore.api.GET('/lines/queries/list')
            this.list = res.data ? res.data : []
        },
    },
    mounted() {
        this.refresh()
    }
}
</script>