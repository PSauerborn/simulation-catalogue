import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useSimulationStore = defineStore('simulation', () => {
  // State - only data storage
  const simulations = ref([])
  const currentRun = ref(null)
  const isLoadingSimulations = ref(false)
  const isRunning = ref(false)
  const runOutput = ref(null)
  const error = ref(null)

  // Getters
  const hasActiveRun = computed(
    () => currentRun.value !== null && currentRun.value.status === 'running',
  )
  const canRun = computed(() => !isRunning.value && !hasActiveRun.value)

  const sortedSimulations = computed(() => {
    return [...simulations.value].sort((a, b) => {
      return new Date(b.created_at) - new Date(a.created_at)
    })
  })

  // Setters for state updates
  function setSimulations(data) {
    simulations.value = data
  }

  function setCurrentRun(data) {
    currentRun.value = data
  }

  function setIsLoadingSimulations(value) {
    isLoadingSimulations.value = value
  }

  function setIsRunning(value) {
    isRunning.value = value
  }

  function setRunOutput(data) {
    runOutput.value = data
  }

  function setError(message) {
    error.value = message
  }

  function clearError() {
    error.value = null
  }

  function clearCurrentRun() {
    currentRun.value = null
    runOutput.value = null
    isRunning.value = false
  }

  /**
   * Parse CSV content from the output
   */
  function parseCSVFromOutput(csvContent) {
    const lines = csvContent.trim().split('\n')
    if (lines.length < 2) return { headers: [], data: [] }

    const headers = lines[0].split(',').map((h) => h.trim())
    const data = lines.slice(1).map((line) => {
      const values = line.split(',').map((v) => parseFloat(v.trim()))
      return headers.reduce((obj, header, index) => {
        obj[header] = values[index]
        return obj
      }, {})
    })

    return { headers, data }
  }

  return {
    // State
    simulations,
    currentRun,
    isLoadingSimulations,
    isRunning,
    runOutput,
    error,
    // Getters
    hasActiveRun,
    canRun,
    sortedSimulations,
    // Setters
    setSimulations,
    setCurrentRun,
    setIsLoadingSimulations,
    setIsRunning,
    setRunOutput,
    setError,
    clearError,
    clearCurrentRun,
    // Utilities
    parseCSVFromOutput,
  }
})
