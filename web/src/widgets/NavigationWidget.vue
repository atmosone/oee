<template>
    <v-navigation-drawer order="-1">
        <v-list nav>
            <v-list-item to="/" prepend-icon="mdi-home" color="primary">
                Главная
            </v-list-item>
            <v-list-item to="/equipment" prepend-icon="mdi-hexagon" color="primary">
                Оборудование
            </v-list-item>
            <v-list-item to="/lines" prepend-icon="mdi-hexagon-multiple" color="primary">
                Линии
            </v-list-item>
            <v-list-item to="/specifications" prepend-icon="mdi-file-check" color="primary">
                Спецификации
            </v-list-item>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">
export default {
    name: "NavigationWidget"
}
</script>