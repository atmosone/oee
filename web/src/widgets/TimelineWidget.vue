<template>
   <div>
     <apexchart
     type="rangeBar"
    :options="options"
     :series="series"
        width="100%"
        height="180"
        @dataPointSelection="onSelect"
        @dataPointMouseEnter="onHover"
        />
   </div>
</template>

<script lang="ts">
import type { components } from '@/api';
import type { ApexChart, ApexOptions } from 'apexcharts';
import type { PropType } from 'vue';
import VueApexCharts from 'vue3-apexcharts';
import ruLocale from 'apexcharts/dist/locales/ru.json'

interface ApexBarConfig {
    seriesIndex: number,
    dataPointIndex: number,
}
export default {
    name: "TimelineWidget",
    components: {
        apexchart: VueApexCharts,
    },
    emits: {
        select: (event: components['schemas']['IntervalID']) => event,
        hover: (event: components['schemas']['IntervalID']) => event,
    },
    props: {
        series: {
            type: Array as PropType<Array<components['schemas']['ApexchartsTimeline']>>,
            default: () => new Array<components['schemas']['ApexchartsTimeline']>(),
        },
    },
    data: () => ({
        options: {
            chart: {id: 'vuechart-example', type: 'rangeBar', locales: [ruLocale], defaultLocale: 'ru', parentHeightOffset: 0, sparkline: {enabled: false}},
            tooltip: {x: {show: true, format: 'HH:mm:ss'}},
            plotOptions: {bar: {horizontal: true, barHeight: '100%', rangeBarGroupRows: true, borderRadius: 2 }},
            xaxis: {type: 'datetime', labels: {datetimeUTC: false, datetimeFormatter: {hour: 'HH:mm',minute: 'HH:mm:ss'}}},
            yaxis: {show: false, labels: {show: false}},
            fill: {type: 'solid', opacity: 0.6},
            grid: {xaxis: {lines: {show: true}}, yaxis: {lines: {show: false}}},
            colors: ['#901215', '#04af54']
        } as ApexOptions
    }),
    methods: {
        onSelect(event:unknown, chartContext:unknown, config: ApexBarConfig) {
            const item = this.series[config.seriesIndex]?.data[config.dataPointIndex]
            if (item) this.$emit('select', item.id)
        },
        onHover(event:unknown, chartContext:unknown, config: ApexBarConfig) {
            const item = this.series[config.seriesIndex]?.data[config.dataPointIndex]
            if (item) this.$emit('hover', item.id)
        }
    }
}
</script>