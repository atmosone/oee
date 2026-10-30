<template>
    <v-container>
        <v-row>
            <v-col cols="12">
                <h1>{{ line?.name }}</h1>
            </v-col>
        </v-row>
        <v-row>
            <v-col v-for="item in line?.items" cols="4">
                {{ item.name }}
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts">
import type { components } from '@/schema';
import useAppStore from '@/store';
import { mapStores } from 'pinia';

export default{
    name: "EquipmentLinePage",
    data: () => ({
        line: undefined as components['schemas']['EquipmentLine'] | undefined,
    }),
    computed: { ...mapStores(useAppStore) },
    methods: {
        async refreshLine(): Promise<void> {
            const data = await this.appStore.client.GET('/queries/equipment_line', {
                params: {
                    query: {
                        equipment_line_id: this.$route.params.line_id as string,
                    }
                }
            })
            this.line = data.data
        }
    },
    mounted() {
        this.refreshLine()
    }
}
</script>