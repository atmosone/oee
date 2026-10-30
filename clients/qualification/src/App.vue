<template>
  <v-app>
    <v-main>
      <v-container>
          <v-row>
            <v-col cols="12">
              <v-card-item v-if="equipment">
                <v-card-title class="text-display-medium font-weight-black">
                  {{ equipment.name }}
                </v-card-title>
                <v-card-subtitle>
                  {{ equipment.line }}
                </v-card-subtitle>
                <template v-slot:append>
                  <v-chip variant="flat"
                    class="mr-2"
                    rounded="lg"
                    :color="current_downtime ? 'error' : 'success'"
                    :text="current_downtime ? 
                      `Простой: ${new Date(current_downtime.started_at).toLocaleTimeString('ru-RU')}` :
                      'В работе'"
                  />
                  <v-chip v-if="current_workshift"
                    rounded="lg"
                    variant="flat"
                    color="primary"
                    :text="current_workshift.name"
                  />
                </template>
              </v-card-item>
            </v-col>
          </v-row>
          <v-row class="pa-6 bg-surface rounded-xl">
            <v-col cols="12">
              <v-label>Неквалифицированные простои</v-label>
            </v-col>
            <v-col v-for="downtime in unqualified_downtimes"
            cols="12" md="6" xl="4">
              <app-create-qualification-form 
                :downtime="downtime" 
                :specifications="specifications"
                v-on:submit="createQualification($event)"
              />
            </v-col>
          </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>
<script lang="ts">
import createClient from 'openapi-fetch';
import type { components, paths } from './schema';
import AppCreateQualificationForm from './AppCreateQualificationForm.vue';
export default{
  name: "Qualification",
  components: {
    AppCreateQualificationForm,
  },
  data: () => ({
    es: null as EventSource | null,
    client: createClient<paths>({baseUrl: '/api'}),
    equipment: null as components['schemas']['Equipment'] | null,
    current_downtime: null as components['schemas']['CurrentEquipmentDowntime'] | null,
    current_workshift: null as components['schemas']['CurrentEquipmentWorkshift'] | null,
    unqualified_downtimes: [] as components['schemas']['UnqualifiedDowntimesList'],
    specifications: [] as components['schemas']['EquipmentSpecificationsList'],
  }),
  computed: {
    equipment_id(): string { return window.location.pathname.replace('/', '') }
  },
  methods: {
    async refreshEquipment(): Promise<void> {
      const data = await this.client.GET(
        '/queries/equipment', {
          params: {query: {equipment_id: this.equipment_id}}
        })
      this.equipment = data.data ? data.data : null
    },
    async refreshCurrentDowntime(): Promise<void> {
      const data = await this.client.GET(
        '/queries/current_equipment_downtime', {
        params: {query: {equipment_id: this.equipment_id}}
      })
      this.current_downtime = data.data ? data.data : null
    },
    async refreshCurrentWorkshift(): Promise<void> {
      const data = await this.client.GET(
        '/queries/current_equipment_workshift')
      this.current_workshift = data.data ? data.data : null
    },
    async refreshUnqualifiedDowntimes(): Promise<void> {
      if (this.current_workshift === null) return
      const data = await this.client.GET(
        '/queries/unqualified_downtimes_list', {
        params: {query: {
          equipment_id: this.equipment_id, 
          workshift_id: this.current_workshift.id
        }},
      })
      this.unqualified_downtimes = data.data ? data.data : []
    },
    async refreshSpecifications(): Promise<void> {
      const data = await this.client.GET(
        '/queries/equipment_specifications_list', {
        params: {query: {equipment_id: this.equipment_id}}
      })
      this.specifications = data.data ? data.data : []
    },
    async createQualification(cmd: components['schemas']['CreateQualificationCommand']): Promise<void> {
      const res = await this.client.POST(
        '/commands/create_qualification', {
        body: cmd
      })
      await this.refreshUnqualifiedDowntimes()
    },
    async onWorkshiftStarted(_: MessageEvent): Promise<void> { 
      await this.refreshCurrentWorkshift() 
      await this.refreshUnqualifiedDowntimes()
    },
    async onDowntimeStarted(e: MessageEvent): Promise<void> {
      await this.refreshCurrentDowntime()
    },
    async onDowntimeFinished(e: MessageEvent): Promise<void> {
      await this.refreshCurrentDowntime()
      await this.refreshUnqualifiedDowntimes()
    },
  },
  async mounted() {
    this.es = new EventSource(`/api/events?equipment_id=${this.equipment_id}`)
    this.es.addEventListener("workshift_started", this.onWorkshiftStarted)
    this.es.addEventListener("downtime_started", this.onDowntimeStarted)
    this.es.addEventListener("downtime_finished", this.onDowntimeFinished)
    await this.refreshEquipment()
    await this.refreshCurrentWorkshift()
    await this.refreshCurrentDowntime()
    await this.refreshUnqualifiedDowntimes()
    await this.refreshSpecifications()
  },
  beforeUnmount() {
    if (!this.es) return
    this.es.removeEventListener("workshift_started", this.onWorkshiftStarted)
    this.es.removeEventListener("downtime_started", this.onDowntimeStarted)
    this.es.removeEventListener("downtime_finished", this.onDowntimeFinished)
    this.es.close()
    this.es = null
  }
}
</script>
