<template>
    <v-dialog width="350">
        <template v-slot:activator="{props}">
            <v-btn v-bind="props"
            icon="mdi-plus"
            color="primary"
            size="small"
            variant="tonal"
            rounded="lg"/>
        </template>
        <v-card>
            <v-card-item class="border-b">
                <v-card-title>Create line</v-card-title>
                <v-card-subtitle>Provide line info</v-card-subtitle>
            </v-card-item>
            <v-card-text class="pa-4">
                <v-text-field 
                    label="Line ID"
                    color="primary"
                />
                <v-text-field
                    label="Line name"
                    color="primary"
                />
            </v-card-text>
            <v-card-actions class="pa-0">
                <v-btn 
                    text="Create"
                    color="primary"
                    variant="tonal"
                    rounded="0"
                    height="52"
                    block
                />
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
export default{
    name: "CreateEquipmentLine",
}
</script>