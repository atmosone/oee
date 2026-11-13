<template>
    <v-container>
        <v-row>
            <v-col cols="12">
                <line-topology-widget :id="id"/>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts">
import type { components } from '@/api';
import LineTopologyWidget from '@/widgets/LineTopologyWidget.vue';

export default {
    name: "LineScreen",
    components: {
        LineTopologyWidget,
    },
    computed: {
        id(): components['schemas']['LineID'] {
            return this.$route.params['id'] as components['schemas']['LineID']
        },
    }
}
</script>