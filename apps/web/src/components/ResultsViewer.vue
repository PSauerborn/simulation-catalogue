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
            <q-icon name="eva-activity-outline" size="28px" color="primary" />
            <div>
              <h2>Simulation Results</h2>
              <p v-if="simulationName">{{ simulationName }}</p>
            </div>
          </div>
          <q-btn icon="eva-close-outline" flat round dense v-close-popup />
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
          <q-icon name="eva-alert-circle-outline" size="64px" color="negative" />
          <p>{{ error }}</p>
          <q-btn outline color="primary" label="Retry" @click="loadData" />
        </div>

        <!-- Charts -->
        <div v-else class="charts-container">
          <!-- Trajectory Chart -->
          <div v-if="hasTrajectoryData" class="chart-section">
            <h3 class="chart-title">
              <q-icon name="eva-trending-up-outline" size="20px" />
              3D Trajectory
            </h3>
            <div class="chart-wrapper">
              <div ref="trajectoryChartRef" class="plotly-chart"></div>
            </div>
          </div>

          <!-- Position vs Time Charts -->
          <div class="chart-grid">
            <div v-if="chartData.x1?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="eva-trending-up-outline" size="18px" />
                X1 Evolution
              </h3>
              <div class="chart-wrapper">
                <div ref="x1ChartRef" class="plotly-chart"></div>
              </div>
            </div>

            <div v-if="chartData.x2?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="eva-trending-up-outline" size="18px" />
                X2 Evolution
              </h3>
              <div class="chart-wrapper">
                <div ref="x2ChartRef" class="plotly-chart"></div>
              </div>
            </div>

            <div v-if="chartData.x3?.length" class="chart-section">
              <h3 class="chart-title">
                <q-icon name="eva-trending-up-outline" size="18px" />
                X3 Evolution
              </h3>
              <div class="chart-wrapper">
                <div ref="x3ChartRef" class="plotly-chart"></div>
              </div>
            </div>
          </div>

          <!-- Data Table -->
          <div v-if="tableData.length" class="table-section">
            <h3 class="chart-title">
              <q-icon name="eva-grid-outline" size="20px" />
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
        <q-btn
          outline
          color="grey"
          icon="eva-download-outline"
          label="Download CSV"
          @click="downloadCSV"
        />
        <q-btn
          outline
          color="grey"
          icon="eva-archive-outline"
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
    return { x1: [], x2: [], x3: [] }
  }
  return allDatasets.value[selectedFile.value]
})

const tableData = computed(() => createTableData(chartData.value))

// Chart refs
const trajectoryChartRef = ref(null)
const x1ChartRef = ref(null)
const x2ChartRef = ref(null)
const x3ChartRef = ref(null)

// Chart instances

const hasTrajectoryData = computed(() => {
  return chartData.value.x1?.length > 0 && chartData.value.x2?.length > 0
})

const tableColumns = computed(() => {
  const columns = []
  columns.push({ name: 'index', label: 'Step', field: 'index', sortable: true, align: 'left' })

  if (chartData.value.x1?.length) {
    columns.push({ name: 'x1', label: 'X1 (m)', field: 'x1', sortable: true, align: 'left' })
  }
  if (chartData.value.x2?.length) {
    columns.push({ name: 'x2', label: 'X2 (m)', field: 'x2', sortable: true, align: 'left' })
  }
  if (chartData.value.x3?.length) {
    columns.push({ name: 'x3', label: 'X3 (m)', field: 'x3', sortable: true, align: 'left' })
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
  // Expected format: { output: { "filename.csv": [{x1, x2, x3}, ...] } }
  const output = data.output || data
  const datasets = {}

  for (const filename of Object.keys(output)) {
    const points = output[filename]

    if (!Array.isArray(points)) {
      console.warn('Expected array of points for file:', filename)
      continue
    }

    // Transform to chart format { x1: [], x2: [], x3: [] }
    const result = { x1: [], x2: [], x3: [] }

    for (const point of points) {
      // Handle x1 (compat with x)
      if (point.x1 !== undefined) result.x1.push(Number(point.x1))
      else if (point.x !== undefined) result.x1.push(Number(point.x))
      else if (point.pos_x !== undefined) result.x1.push(Number(point.pos_x))

      // Handle x2 (compat with y)
      if (point.x2 !== undefined) result.x2.push(Number(point.x2))
      else if (point.y !== undefined) result.x2.push(Number(point.y))
      else if (point.pos_y !== undefined) result.x2.push(Number(point.pos_y))

      // Handle x3 (compat with z)
      if (point.x3 !== undefined) result.x3.push(Number(point.x3))
      else if (point.z !== undefined) result.x3.push(Number(point.z))
      else if (point.pos_z !== undefined) result.x3.push(Number(point.pos_z))
    }

    if (result.x1.length > 0) {
      datasets[filename] = result
    }
  }

  return datasets
}

function createTableData(data) {
  const length = data.x1?.length || 0
  const rows = []

  for (let i = 0; i < Math.min(length, 1000); i++) {
    const x1 = Number(data.x1?.[i])
    const x2 = Number(data.x2?.[i])
    const x3 = Number(data.x3?.[i])

    rows.push({
      index: i,
      x1: !isNaN(x1) ? x1.toFixed(6) : '',
      x2: !isNaN(x2) ? x2.toFixed(6) : '',
      x3: !isNaN(x3) ? x3.toFixed(6) : '',
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
  if (trajectoryChartRef.value && chartData.value.x1?.length) {
    Plotly.newPlot(
      trajectoryChartRef.value,
      [
        {
          x: chartData.value.x1,
          y: chartData.value.x2,
          z: chartData.value.x3,
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
            title: 'X1',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          yaxis: {
            title: 'X2',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          zaxis: {
            title: 'X3',
            gridcolor: 'rgba(255, 255, 255, 0.1)',
            backgroundcolor: 'rgba(0,0,0,0)',
          },
          bgcolor: 'rgba(0,0,0,0)',
        },
      },
      commonConfig,
    )
  }

  // X1 vs Time (Step)
  if (x1ChartRef.value && chartData.value.x1?.length) {
    Plotly.newPlot(
      x1ChartRef.value,
      [
        {
          y: chartData.value.x1,
          mode: 'lines',
          line: { color: '#ff453a', width: 2 },
        },
      ],
      { ...commonLayout },
      commonConfig,
    )
  }

  // X2 vs Time (Step)
  if (x2ChartRef.value && chartData.value.x2?.length) {
    Plotly.newPlot(
      x2ChartRef.value,
      [
        {
          y: chartData.value.x2,
          mode: 'lines',
          line: { color: '#30d158', width: 2 },
        },
      ],
      { ...commonLayout },
      commonConfig,
    )
  }

  // X3 vs Time (Step)
  if (x3ChartRef.value && chartData.value.x3?.length) {
    Plotly.newPlot(
      x3ChartRef.value,
      [
        {
          y: chartData.value.x3,
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
  if (x1ChartRef.value)
    try {
      Plotly.purge(x1ChartRef.value)
    } catch (e) {
      console.warn('Failed to purge x1 chart', e)
    }
  if (x2ChartRef.value)
    try {
      Plotly.purge(x2ChartRef.value)
    } catch (e) {
      console.warn('Failed to purge x2 chart', e)
    }
  if (x3ChartRef.value)
    try {
      Plotly.purge(x3ChartRef.value)
    } catch (e) {
      console.warn('Failed to purge x3 chart', e)
    }
}

function downloadCSV() {
  const headers = ['x1', 'x2', 'x3'].filter((h) => chartData.value[h]?.length)
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
