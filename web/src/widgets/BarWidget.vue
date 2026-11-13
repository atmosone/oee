<template>
    <v-app-bar>
        <v-app-bar-nav-icon icon="mdi-menu" @click="$emit('menu')"/>
        <v-app-bar-title class="font-weight-black">
            Workflow<span class="text-grey">Tracker</span>
        </v-app-bar-title>
    </v-app-bar>
</template>

<script lang="ts">
export default {
    name: "BarWidget",
    emits: {  menu: (event: void) => event },
}
</script>