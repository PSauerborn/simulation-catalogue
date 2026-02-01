<template>
  <q-dialog
    :model-value="modelValue"
    @update:model-value="handleDialogUpdate"
    :persistent="hasActiveRun"
    transition-show="scale"
    transition-hide="scale"
  >
    <q-card class="runner-dialog">
      <!-- Header -->
      <q-card-section class="runner-header">
        <div class="runner-title-section">
          <q-icon name="eva-monitor-outline" size="24px" class="terminal-icon" />
          <div class="runner-title-text">
            <h3 class="runner-title">Simulation Runner</h3>
            <p class="runner-subtitle" v-if="selectedSimulation">
              {{ selectedSimulation.name }}
            </p>
          </div>
        </div>
        <div class="header-actions">
          <span :class="['status-badge', currentStatus]">
            {{ statusLabel }}
          </span>
          <q-btn
            v-if="!hasActiveRun"
            icon="eva-close-outline"
            flat
            round
            dense
            @click="closeDialog"
          />
        </div>
      </q-card-section>

      <!-- Main Content Section -->
      <q-card-section class="runner-body">
        <!-- No Simulation Selected -->
        <div v-if="!selectedSimulation" class="empty-state">
          <q-icon name="eva-flask-outline" size="64px" color="grey-6" />
          <p>Select a simulation from the catalogue below to configure and run it.</p>
        </div>

        <!-- Simulation Configuration -->
        <div v-else class="runner-content">
          <p class="runner-description" v-if="selectedSimulation?.description">
            {{ selectedSimulation.description }}
          </p>
          <!-- Parameters Form -->
          <div class="parameters-form" v-if="!currentRun">
            <div class="form-header">
              <span class="form-label">Configuration</span>
              <q-btn flat dense size="sm" color="grey" label="Reset" @click="resetParameters" />
            </div>

            <div class="form-grid">
              <div
                v-for="param in selectedSimulation.parameters"
                :key="param.name"
                class="form-field"
              >
                <label class="field-label">
                  {{ param.name }}
                  <span :class="['param-type', param.type]">{{ param.type }}</span>
                </label>

                <!-- Float/Int Input -->
                <div v-if="param.type === 'float' || param.type === 'int'" class="input-wrapper">
                  <p class="param-hint">
                    {{ param.description }}
                    <q-tooltip v-if="param.description">{{ param.description }}</q-tooltip>
                  </p>
                  <q-input
                    v-model.number="parameterValues[param.name]"
                    type="number"
                    :step="param.type === 'float' ? 'any' : '1'"
                    filled
                    dense
                    dark
                  />
                </div>

                <!-- String Input -->
                <div v-else-if="param.type === 'string'" class="input-wrapper">
                  <p class="param-hint">
                    {{ param.description }}
                    <q-tooltip v-if="param.description">{{ param.description }}</q-tooltip>
                  </p>
                  <q-input v-model="parameterValues[param.name]" filled dense dark />
                </div>

                <!-- Bool Toggle -->
                <div v-else-if="param.type === 'bool'" class="input-wrapper">
                  <p class="param-hint">
                    {{ param.description }}
                    <q-tooltip v-if="param.description">{{ param.description }}</q-tooltip>
                  </p>
                  <q-toggle v-model="parameterValues[param.name]" dark />
                </div>

                <!-- Vector Input (3D) -->
                <div v-else-if="param.type === 'vector'" class="vector-input">
                  <p class="param-hint">
                    {{ param.description }}
                    <q-tooltip v-if="param.description">{{ param.description }}</q-tooltip>
                  </p>
                  <div class="vector-values">
                    <q-input
                      v-model.number="parameterValues[param.name][0]"
                      type="number"
                      step="any"
                      filled
                      dense
                      dark
                      label="X"
                      class="vector-field"
                    />
                    <q-input
                      v-model.number="parameterValues[param.name][1]"
                      type="number"
                      step="any"
                      filled
                      dense
                      dark
                      label="Y"
                      class="vector-field"
                    />
                    <q-input
                      v-model.number="parameterValues[param.name][2]"
                      type="number"
                      step="any"
                      filled
                      dense
                      dark
                      label="Z"
                      class="vector-field"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Run Button -->
            <div class="form-actions">
              <q-btn
                class="run-button"
                color="primary"
                icon="eva-arrow-right-outline"
                label="Run Simulation"
                :loading="isRunning"
                :disable="!canRunSimulation"
                @click="startSimulation"
              />
            </div>

            <!-- Error Banner (for API errors like 409) -->
            <div v-if="error" class="error-banner">
              <q-icon name="eva-alert-triangle-outline" size="20px" color="negative" />
              <span>{{ error }}</span>
              <q-btn
                flat
                dense
                size="sm"
                icon="eva-close-outline"
                @click="simulationStore.clearError()"
              />
            </div>
          </div>

          <!-- Active Run Status -->
          <div v-else class="run-status">
            <div class="status-header">
              <span class="status-label">Run Status</span>
              <q-btn
                v-if="!hasActiveRun"
                flat
                dense
                size="md"
                color="grey"
                icon="eva-refresh-outline"
                label="Reset"
                @click="clearRun"
              />
            </div>

            <!-- Progress Indicator -->
            <div v-if="hasActiveRun" class="progress-section">
              <q-circular-progress
                indeterminate
                size="60px"
                :thickness="0.15"
                color="info"
                track-color="grey-8"
              />
              <p class="progress-text">Running simulation...</p>
            </div>

            <!-- Completed Status -->
            <div v-else-if="currentRun?.status === 'completed'" class="completed-section">
              <q-icon name="eva-checkmark-circle-2-outline" size="48px" color="positive" />
              <p class="completed-text">Simulation completed successfully</p>

              <!-- Output Actions -->
              <div class="output-actions">
                <q-btn
                  outline
                  color="primary"
                  icon="eva-download-outline"
                  label="Download ZIP"
                  @click="downloadResults"
                />
                <q-btn
                  color="primary"
                  icon="eva-bar-chart-outline"
                  label="View Results"
                  @click="$emit('viewResults')"
                />
              </div>
            </div>

            <!-- Failed Status -->
            <div v-else-if="currentRun?.status === 'failed'" class="failed-section">
              <q-icon name="eva-alert-circle-outline" size="48px" color="negative" />
              <p class="failed-text">Simulation failed</p>
              <p class="error-message">{{ error }}</p>
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { useSimulationStore } from 'stores/simulation-store'
import { storeToRefs } from 'pinia'
import { runSimulation as apiRunSimulation, getRunStatus, fetchRunOutput } from 'src/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  selectedSimulation: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue', 'viewResults'])

const simulationStore = useSimulationStore()
const { currentRun, isRunning, error, hasActiveRun, canRun } = storeToRefs(simulationStore)

const parameterValues = ref({})
let pollInterval = null

// Cleanup polling on unmount
onUnmounted(() => {
  stopPolling()
})

// Dialog handlers
function handleDialogUpdate(value) {
  // Only allow closing if simulation is not running
  if (!hasActiveRun.value) {
    emit('update:modelValue', value)
  }
}

function closeDialog() {
  if (!hasActiveRun.value) {
    emit('update:modelValue', false)
  }
}

// Initialize parameters when simulation is selected
watch(
  () => props.selectedSimulation,
  (newSimulation) => {
    if (newSimulation) {
      initializeParameters(newSimulation)
    } else {
      parameterValues.value = {}
    }
  },
  { immediate: true },
)

function initializeParameters(simulation) {
  const values = {}
  for (const param of simulation.parameters || []) {
    // Use default value if provided
    if (param.default !== undefined && param.default !== null) {
      if (param.type === 'vector' && Array.isArray(param.default)) {
        // Ensure deep copy for arrays
        values[param.name] = [...param.default]
      } else {
        values[param.name] = param.default
      }
      continue
    }

    // Fallback to zero-values if no default is provided
    switch (param.type) {
      case 'float':
        values[param.name] = 0.0
        break
      case 'int':
        values[param.name] = 0
        break
      case 'bool':
        values[param.name] = false
        break
      case 'string':
        values[param.name] = ''
        break
      case 'vector':
        // Fixed 3D vector (x, y, z)
        values[param.name] = [0.0, 0.0, 0.0]
        break
      default:
        values[param.name] = null
    }
  }
  parameterValues.value = values
}

function resetParameters() {
  if (props.selectedSimulation) {
    initializeParameters(props.selectedSimulation)
  }
}

async function startSimulation() {
  if (!props.selectedSimulation) return

  simulationStore.setIsRunning(true)
  simulationStore.clearError()
  simulationStore.setRunOutput(null)

  try {
    const response = await apiRunSimulation(props.selectedSimulation.id, parameterValues.value)

    // 201 means job was queued successfully
    if (response.status === 201) {
      simulationStore.setCurrentRun({
        status: 'running',
        simulation_id: props.selectedSimulation.id,
      })
      startPolling()
    }
  } catch (err) {
    // Handle 409 Conflict - user already has a running job
    console.log(err.response.status)
    if (err.response?.status === 409) {
      simulationStore.setError(
        'You already have a simulation running. Please wait for it to complete.',
      )
    } else {
      simulationStore.setError(err.response?.data?.message || 'Failed to start simulation')
    }
    simulationStore.setIsRunning(false)
    console.error('Failed to start simulation:', err)
  }
}

function startPolling() {
  stopPolling() // Clear any existing interval

  pollInterval = setInterval(async () => {
    if (!currentRun.value) {
      stopPolling()
      return
    }

    try {
      const response = await getRunStatus()
      const statusData = response.data
      simulationStore.setCurrentRun(statusData)

      if (statusData.status !== 'running') {
        stopPolling()
        simulationStore.setIsRunning(false)

        if (statusData.status === 'completed') {
          // Fetch the output data
          await loadRunOutput()
        }
      }
    } catch (err) {
      console.error('Polling error:', err)
    }
  }, 2000) // Poll every 2 seconds
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

async function loadRunOutput() {
  try {
    const response = await fetchRunOutput({
      params: { format: 'json' },
      responseType: 'json',
    })
    simulationStore.setRunOutput(response.data)
    return response.data
  } catch (err) {
    console.error('Fetch output error:', err)
    throw err
  }
}

async function downloadResults() {
  try {
    // Always fetch fresh ZIP for download
    const response = await fetchRunOutput({
      responseType: 'blob',
    })

    const blob = response.data
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `simulation_output_${Date.now()}.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Failed to download results:', err)
  }
}

function clearRun() {
  stopPolling()
  simulationStore.clearCurrentRun()
}

// Validate that all required parameters have valid values
const allParametersValid = computed(() => {
  if (!props.selectedSimulation?.parameters) return false

  for (const param of props.selectedSimulation.parameters) {
    const value = parameterValues.value[param.name]

    switch (param.type) {
      case 'float':
      case 'int':
        // Must be a valid number
        if (value === null || value === undefined || value === '' || isNaN(value)) {
          return false
        }
        break
      case 'string':
        // Must be a non-empty string (if required)
        if (param.required && (!value || value.trim() === '')) {
          return false
        }
        break
      case 'vector':
        // All 3 components must be valid numbers
        if (!Array.isArray(value) || value.length !== 3) {
          return false
        }
        for (const component of value) {
          if (
            component === null ||
            component === undefined ||
            component === '' ||
            isNaN(component)
          ) {
            return false
          }
        }
        break
      // bool type is always valid (true or false)
    }
  }
  return true
})

// Combined check for run button - both store canRun and parameter validation
const canRunSimulation = computed(() => {
  return canRun.value && allParametersValid.value
})

const currentStatus = computed(() => {
  if (hasActiveRun.value) return 'running'
  if (currentRun.value?.status === 'completed') return 'completed'
  if (currentRun.value?.status === 'failed') return 'failed'
  return 'idle'
})

const statusLabel = computed(() => {
  switch (currentStatus.value) {
    case 'running':
      return 'Running'
    case 'completed':
      return 'Completed'
    case 'failed':
      return 'Failed'
    default:
      return 'Idle'
  }
})
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.runner-dialog {
  background: $dark;
  min-width: 500px;
  max-width: 800px;
  width: 90vw;
  border-radius: $border-radius-xl;
  border: 1px solid $glass-border;
}

.runner-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid $glass-border;

  .runner-title-section {
    display: flex;
    align-items: flex-start;
    gap: 16px;

    .terminal-icon {
      color: $primary;
      margin-top: 4px;
    }

    .runner-title-text {
      .runner-title {
        font-size: 1.25rem;
        font-weight: 600;
        color: #f5f5f7;
        margin: 0 0 4px 0;
      }

      .runner-subtitle {
        font-size: 0.9rem;
        color: #a1a1a6;
        margin: 0;
      }
    }
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #636366;

  p {
    margin-top: 16px;
    font-size: 0.95rem;
  }
}

.runner-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.parameters-form {
  .form-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .form-label {
      font-size: 0.85rem;
      font-weight: 500;
      color: #a1a1a6;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .form-field {
    .field-label {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 0.9rem;
      font-weight: 500;
      color: #f5f5f7;
      margin-bottom: 8px;
    }

    .param-hint {
      font-size: 0.8rem;
      color: #a1a1a6;
      margin: 0 0 8px 0;
      line-height: 1.4;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      cursor: default;
    }

    .input-wrapper {
      width: 100%;
    }

    .vector-input {
      .vector-values {
        display: flex;
        flex-direction: row;
        gap: 12px;

        .vector-field {
          flex: 1;
          min-width: 80px;
        }
      }
    }
  }

  .form-actions {
    margin-top: 24px;
    display: flex;
    justify-content: flex-end;

    .run-button {
      padding: 12px 32px;
      font-weight: 500;
    }
  }
}

.run-status {
  .status-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .status-label {
      font-size: 0.85rem;
      font-weight: 500;
      color: #a1a1a6;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
  }

  .progress-section {
    text-align: center;
    padding: 32px;

    .progress-text {
      margin-top: 16px;
      color: #a1a1a6;
    }
  }

  .completed-section,
  .failed-section {
    text-align: center;
    padding: 24px;

    .completed-text,
    .failed-text {
      margin: 12px 0 24px 0;
      font-size: 1rem;
      color: #f5f5f7;
    }

    .error-message {
      color: $negative;
      font-size: 0.9rem;
    }

    .output-actions {
      display: flex;
      justify-content: center;
      gap: 12px;
      flex-wrap: wrap;
    }
  }
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-top: 16px;
  background: rgba($negative, 0.1);
  border: 1px solid rgba($negative, 0.3);
  border-radius: $border-radius-lg;
  color: $negative;
  font-size: 0.9rem;

  span {
    flex: 1;
  }
}
</style>
