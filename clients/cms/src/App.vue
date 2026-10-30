<template>
  <v-app>
    <v-app-bar>
      <v-app-bar-nav-icon @click="switch_nav"/>
      <v-app-bar-title>CMS</v-app-bar-title>
    </v-app-bar>
    <v-navigation-drawer v-model="nav">
      <v-list nav>
        <v-list-item to="/equipment">Equipment</v-list-item>
        <v-list-item to="/lines">Lines</v-list-item>
        <v-list-item to="/specifications">Specifications</v-list-item>
        <v-list-item to="/rules">Rules</v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script lang="ts">
export default{
  name: "CMS",
  data: () => ({nav: true}),
  methods: {
    switch_nav(): void {this.nav = !this.nav}
  },

}
</script>
