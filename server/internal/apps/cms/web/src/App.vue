<template>
  <v-app>
    <v-app-bar class="border-b" elevation="0">
      <v-app-bar-nav-icon />
      <v-app-bar-title>App</v-app-bar-title>
    </v-app-bar>
    <v-navigation-drawer order="-1" width="320">
      <v-card-item class="pa-4">
        <v-card-title>Карта</v-card-title>
      </v-card-item>
      <v-divider />
      <lines-tree />
      <v-divider />
      <equipment-tree />
      <v-divider />
      <specifications-tree />
      <v-divider />
    </v-navigation-drawer>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import LinesTree from '@/lines/LinesTree.vue'
import EquipmentTree from '@/equipment/EquipmentTree.vue'
import SpecificationsTree from '@/specifications/SpecificationsTree.vue'

export default {
  name: 'CMS',
  components: {
    LinesTree,
    EquipmentTree,
    SpecificationsTree,
  },
  data: () => ({}),
}
</script>
