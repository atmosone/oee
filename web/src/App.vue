<template>
  <v-app>
    <bar-widget @menu="menu = !menu"/>
    <navigation-widget v-model="menu"/>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import NavigationWidget from '@/widgets/NavigationWidget.vue';
import BarWidget from './widgets/BarWidget.vue';

export default {
  name: "App",
  components: {
    BarWidget,
    NavigationWidget,
  },
  data: () => ({
    menu: true,
  })
}
</script>
