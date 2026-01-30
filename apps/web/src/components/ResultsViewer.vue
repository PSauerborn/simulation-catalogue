<template>
  <q-dialog
    v-model="dialogVisible"
    maximized
    transition-show="slide-up"
    transition-hide="slide-down"
  >
    <q-card class="results-dialog">
      <!-- Header -->
      <q-card-section class="dialog-header">
        <div class="header-content">
          <div class="header-title">
            <q-icon name="analytics" size="28px" color="primary" />
            <div>
              <h2>Simulation Results</h2>
              <p v-if="simulationName">{{ simulationName }}</p>
            </div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </div>

        <q-tabs
          v-if="Object.keys(allDatasets).length > 0"
          v-model="selectedFile"
          dense
          class="text-grey q-mt-sm"
          active-color="primary"
          indicator-color="primary"
          align="left"
          narrow-indicator
          no-caps
        >
          <q-tab v-for="file in Object.keys(allDatasets)" :key="file" :name="file" :label="file" />
        </q-tabs>
      </q-card-section>

      <q-separator dark />

      <!-- Content -->
      <q-card-section class="dialog-content">
        <!-- Loading State -->
        <div v-if="isLoading" class="loading-state">
          <q-spinner-orbit size="64px" color="primary" />
          <p>Loading results...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="error-state">
          <q-icon name="error_outline" size="64px" color="negative" />
          <p>{{ error }}</p>
          <q-btn outline color="primary" label="Retry" @click="loadData" />
        </div>

        <!-- Charts -->
        <div v-else class="charts-container">
          <!-- Trajectory Chart -->
          <div v-if="hasTrajectoryData" class="chart-section">
            <h3 class="chart-title">
              <q-icon name="timeline" size="20px" />
              3D Trajectory
            </h3>
            <div class="chart-wrapper">
              <div ref="trajectoryChartRef" class="plotly-chart"></div>
            </div>
          </div>

          <!-- Position vs Time Charts -->
          <div class="chart-grid">
            <div v-if="chartData.x?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="show_chart" size="18px" />
                X Position vs Time
              </h3>
              <div class="chart-wrapper">
                <div ref="xChartRef" class="plotly-chart"></div>
              </div>
            </div>

            <div v-if="chartData.y?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="show_chart" size="18px" />
                Y Position vs Time
              </h3>
              <div class="chart-wrapper">
                <div ref="yChartRef" class="plotly-chart"></div>
              </div>
            </div>

            <div v-if="chartData.z?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="show_chart" size="18px" />
                Z Position vs Time
              </h3>
              <div class="chart-wrapper">
                <div ref="zChartRef" class="plotly-chart"></div>
              </div>
            </div>
          </div>

          <!-- Data Table -->
          <div v-if="tableData.length" class="table-section">
            <h3 class="chart-title">
              <q-icon name="table_chart" size="20px" />
              Raw Data
            </h3>
            <q-table
              :rows="tableData"
              :columns="tableColumns"
              row-key="index"
              dark
              flat
              dense
              :pagination="{ rowsPerPage: 20 }"
              class="data-table"
            />
          </div>
        </div>
      </q-card-section>

      <!-- Footer Actions -->
      <q-card-actions class="dialog-footer">
        <q-btn outline color="grey" icon="download" label="Download CSV" @click="downloadCSV" />
        <q-btn
          outline
          color="grey"
          icon="archive"
          label="Download ZIP"
          @click="$emit('downloadZip')"
        />
        <q-space />
        <q-btn color="primary" label="Close" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import Plotly from 'plotly.js-dist-min'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  outputData: {
    type: [ArrayBuffer, Object, Array, null],
    default: null,
  },
  simulationName: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'downloadZip'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const isLoading = ref(false)
const error = ref(null)
const allDatasets = ref({})
const selectedFile = ref(null)

const chartData = computed(() => {
  if (!selectedFile.value || !allDatasets.value[selectedFile.value]) {
    return { t: [], x: [], y: [], z: [] }
  }
  return allDatasets.value[selectedFile.value]
})

const tableData = computed(() => createTableData(chartData.value))

// Chart refs
const trajectoryChartRef = ref(null)
const xChartRef = ref(null)
const yChartRef = ref(null)
const zChartRef = ref(null)

// Chart instances

const hasTrajectoryData = computed(() => {
  return chartData.value.x?.length > 0 && chartData.value.y?.length > 0
})

const tableColumns = computed(() => {
  const columns = []
  if (chartData.value.t?.length) {
    columns.push({ name: 't', label: 'Time (s)', field: 't', sortable: true, align: 'left' })
  }
  if (chartData.value.x?.length) {
    columns.push({ name: 'x', label: 'X (m)', field: 'x', sortable: true, align: 'left' })
  }
  if (chartData.value.y?.length) {
    columns.push({ name: 'y', label: 'Y (m)', field: 'y', sortable: true, align: 'left' })
  }
  if (chartData.value.z?.length) {
    columns.push({ name: 'z', label: 'Z (m)', field: 'z', sortable: true, align: 'left' })
  }
  return columns
})

watch(
  () => props.modelValue,
  async (visible) => {
    if (visible && props.outputData) {
      await loadData()
    }
  },
)

watch(selectedFile, async () => {
  await nextTick()
  createCharts()
})

async function loadData() {
  if (!props.outputData) {
    error.value = 'No output data available'
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const datasets = processJSONData(props.outputData)
    const files = Object.keys(datasets)

    if (files.length === 0) {
      throw new Error('No valid data found in output')
    }

    allDatasets.value = datasets
    // Reset selection to first file or keep if exists? prefer reset on new load
    selectedFile.value = files[0]

    // Show content before creating charts
    isLoading.value = false

    // Create charts after DOM update
    await nextTick()
    createCharts()
  } catch (err) {
    console.error('Error loading results:', err)
    error.value = 'Failed to parse simulation output'
    isLoading.value = false
  }
}

function processJSONData(data) {
  // Expected format: { output: { "filename.csv": [{time, x, y, z}, ...] } }
  const output = data.output || data
  const datasets = {}

  for (const filename of Object.keys(output)) {
    const points = output[filename]

    if (!Array.isArray(points)) {
      console.warn('Expected array of points for file:', filename)
      continue
    }

    // Transform to chart format { t: [], x: [], y: [], z: [] }
    const result = { t: [], x: [], y: [], z: [] }

    for (const point of points) {
      // Handle time/t
      if (point.t !== undefined) result.t.push(Number(point.t))
      else if (point.time !== undefined) result.t.push(Number(point.time))

      // Handle x
      if (point.x !== undefined) result.x.push(Number(point.x))
      else if (point.pos_x !== undefined) result.x.push(Number(point.pos_x))

      // Handle y
      if (point.y !== undefined) result.y.push(Number(point.y))
      else if (point.pos_y !== undefined) result.y.push(Number(point.pos_y))

      // Handle z
      if (point.z !== undefined) result.z.push(Number(point.z))
      else if (point.pos_z !== undefined) result.z.push(Number(point.pos_z))
    }

    if (result.t.length > 0) {
      datasets[filename] = result
    }
  }

  return datasets
}

function createTableData(data) {
  const length = Math.max(data.t?.length || 0, data.x?.length || 0)
  const rows = []

  for (let i = 0; i < Math.min(length, 1000); i++) {
    const t = Number(data.t?.[i])
    const x = Number(data.x?.[i])
    const y = Number(data.y?.[i])
    const z = Number(data.z?.[i])

    rows.push({
      index: i,
      t: !isNaN(t) ? t.toFixed(6) : '',
      x: !isNaN(x) ? x.toFixed(6) : '',
      y: !isNaN(y) ? y.toFixed(6) : '',
      z: !isNaN(z) ? z.toFixed(6) : '',
    })
  }

  return rows
}

function createCharts() {
  destroyCharts()

  const commonLayout = {
    paper_bgcolor: 'rgba(0,0,0,0)',
    plot_bgcolor: 'rgba(0,0,0,0)',
    font: { color: '#f5f5f7' },
    margin: { t: 30, r: 20, l: 40, b: 40 },
    xaxis: {
      gridcolor: 'rgba(255, 255, 255, 0.1)',
      zerolinecolor: 'rgba(255, 255, 255, 0.2)',
    },
    yaxis: {
      gridcolor: 'rgba(255, 255, 255, 0.1)',
      zerolinecolor: 'rgba(255, 255, 255, 0.2)',
    },
    showlegend: false,
  }

  const commonConfig = { responsive: true, displayModeBar: false }

  // Trajectory 3D
  if (trajectoryChartRef.value && chartData.value.x?.length) {
    Plotly.newPlot(
      trajectoryChartRef.value,
      [
        {
          x: chartData.value.x,
          y: chartData.value.y,
          z: chartData.value.z,
          type: 'scatter3d',
          mode: 'lines',
          line: { color: '#0071e3', width: 4 },
          opacity: 0.8,
        },
      ],
      {
        ...commonLayout,
        title: false,
        margin: { t: 0, r: 0, l: 0, b: 0 },
        scene: {
          xaxis: {
            title: 'X',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          yaxis: {
            title: 'Y',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          zaxis: {
            title: 'Z',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          bgcolor: 'rgba(0,0,0,0)',
        },
      },
      commonConfig,
    )
  }

  // X vs Time
  if (xChartRef.value && chartData.value.x?.length) {
    Plotly.newPlot(
      xChartRef.value,
      [
        {
          x: chartData.value.t,
          y: chartData.value.x,
          mode: 'lines',
          line: { color: '#ff453a', width: 2 },
        },
      ],
      { ...commonLayout },
      commonConfig,
    )
  }

  // Y vs Time
  if (yChartRef.value && chartData.value.y?.length) {
    Plotly.newPlot(
      yChartRef.value,
      [
        {
          x: chartData.value.t,
          y: chartData.value.y,
          mode: 'lines',
          line: { color: '#30d158', width: 2 },
        },
      ],
      { ...commonLayout },
      commonConfig,
    )
  }

  // Z vs Time
  if (zChartRef.value && chartData.value.z?.length) {
    Plotly.newPlot(
      zChartRef.value,
      [
        {
          x: chartData.value.t,
          y: chartData.value.z,
          mode: 'lines',
          line: { color: '#5e5ce6', width: 2 },
        },
      ],
      { ...commonLayout },
      commonConfig,
    )
  }
}

function destroyCharts() {
  if (trajectoryChartRef.value)
    try {
      Plotly.purge(trajectoryChartRef.value)
    } catch (e) {
      console.warn('Failed to purge trajectory chart', e)
    }
  if (xChartRef.value)
    try {
      Plotly.purge(xChartRef.value)
    } catch (e) {
      console.warn('Failed to purge x chart', e)
    }
  if (yChartRef.value)
    try {
      Plotly.purge(yChartRef.value)
    } catch (e) {
      console.warn('Failed to purge y chart', e)
    }
  if (zChartRef.value)
    try {
      Plotly.purge(zChartRef.value)
    } catch (e) {
      console.warn('Failed to purge z chart', e)
    }
}

function downloadCSV() {
  const headers = ['t', 'x', 'y', 'z'].filter((h) => chartData.value[h]?.length)
  const rows = []

  const length = Math.max(...headers.map((h) => chartData.value[h]?.length || 0))
  for (let i = 0; i < length; i++) {
    const row = headers.map((h) => chartData.value[h]?.[i] ?? '')
    rows.push(row.join(','))
  }

  const csvContent = [headers.join(','), ...rows].join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `simulation_results_${Date.now()}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}

onMounted(() => {
  if (props.modelValue && props.outputData) {
    loadData()
  }
})

onBeforeUnmount(() => {
  destroyCharts()
})
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.results-dialog {
  background: $dark-page;
  color: #f5f5f7;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  padding: 20px 24px;
  flex-shrink: 0;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-title {
      display: flex;
      align-items: center;
      gap: 16px;

      h2 {
        font-size: 1.5rem;
        font-weight: 600;
        margin: 0;
      }

      p {
        font-size: 0.9rem;
        color: #a1a1a6;
        margin: 4px 0 0 0;
      }
    }
  }
}

.dialog-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  gap: 16px;

  p {
    color: #a1a1a6;
  }
}

.charts-container {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.chart-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid $glass-border;
  border-radius: $border-radius-lg;
  padding: 24px;

  .chart-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 1rem;
    font-weight: 500;
    color: #f5f5f7;
    margin: 0 0 16px 0;
  }

  .chart-wrapper {
    height: 300px;
    position: relative;
    overflow: hidden;

    .plotly-chart {
      width: 100%;
      height: 100%;
    }
  }
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.table-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid $glass-border;
  border-radius: $border-radius-lg;
  padding: 24px;

  .chart-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 1rem;
    font-weight: 500;
    color: #f5f5f7;
    margin: 0 0 16px 0;
  }
}

.data-table {
  max-height: 400px;
}

.dialog-footer {
  padding: 16px 24px;
  border-top: 1px solid $glass-border;
  flex-shrink: 0;
}
</style>
