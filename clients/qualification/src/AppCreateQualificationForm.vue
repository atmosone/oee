<template>
    <v-dialog width="350">
        <template v-slot:activator="{props}">
            <v-card v-bind="props" color="error" variant="tonal">
                <v-card-item class="py-0">
                    <template v-slot:prepend>
                        <v-timeline density="compact" class="mr-4">
                            <v-timeline-item size="15" 
                            dot-color="error" 
                            class="text-label-large">
                                {{ start }}
                            </v-timeline-item>
                            <v-timeline-item size="20" 
                            dot-color="error" 
                            class="text-label-large">
                                {{ finish }}
                            </v-timeline-item>
                        </v-timeline>
                        <v-divider vertical/>
                    </template>
                    <v-card-title class="font-weight-black text-headline-large">
                        {{ duration }}
                    </v-card-title>
                </v-card-item>
            </v-card>
        </template>
        <v-card>
            <v-card-item class="pa-6">
                <v-card-title class="font-weight-bold text-headline-large">{{ duration }}</v-card-title>
                <v-card-subtitle>Квалификация простоя</v-card-subtitle>
            </v-card-item>
            <v-timeline density="compact" class="my-4" direction="horizontal">
                <v-timeline-item size="30" 
                dot-color="primary" class="text-primary">
                    {{ start }}
                </v-timeline-item>
                <v-timeline-item size="30" 
                dot-color="primary" class="text-primary">
                    {{ finish }}
                </v-timeline-item>
            </v-timeline>
            <v-card-text class="border-t pa-4">
                <v-select v-model="specification"
                :items="specifications"
                label="Укажите спецификацию"
                item-title="description" 
                item-value="id"
                color="primary"
                clearable/>
                <v-textarea v-model="comment"
                label="Добавьте комментарий"
                color="primary"
                rows="2">
                </v-textarea>
            </v-card-text>
            <v-card-actions class="pa-0">
                <v-btn :disabled="disabled"
                v-on:click="submit"
                text="Квалифицировать"
                variant="tonal"
                color="primary"
                rounded="0"
                height="60"
                block/>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import type { PropType } from 'vue';
import type { components } from '@/schema';

export default {
    name: "AppCreateQualificationForm",
    emits: { submit: (event: components['schemas']['CreateQualificationCommand']) => event },
    props: {
        downtime: {
            type: Object as PropType<components['schemas']['UnqualifiedDowntimesListItem']>,
            required: true
        },
        specifications: {
            type: Array as PropType<components['schemas']['EquipmentSpecificationsList']>,
            required: true,
            default: () => [] as components['schemas']['EquipmentSpecificationsList']
        },
    },
    data: () => ({
        specification: null as string | null,
        comment: null as string | null,
    }),
    computed: {
        duration(): string {
            return this.downtime.duration.
            replace('h', 'ч ').
            replace('m', 'м ').
            replace('s', 'c')
        },
        start(): string {
            return new Date(this.downtime.started_at).
                toLocaleTimeString('ru-RU')
        },
        finish(): string {
            return new Date(this.downtime.finished_at).
                toLocaleTimeString('ru-RU')
        },
        disabled(): boolean { return this.specification === null }
    },
    methods: {
        submit() {
            if (this.specification === null) return
            this.$emit('submit', {
                downtime_id: this.downtime.id,
                specification_id: this.specification,
                qualification_comment: this.comment ? this.comment : ''
            })
            this.specification = null
            this.comment = null
        }
    }
}
</script>