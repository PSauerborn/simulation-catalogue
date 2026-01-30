import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const CLIENT_ID_KEY = 'clientId'

export const useClientStore = defineStore('client', () => {
  // State - only data storage
  const client = ref(null)
  const isInitialized = ref(false)
  const isLoading = ref(false)
  const error = ref(null)

  // Getters
  const clientId = computed(() => client.value?.id || localStorage.getItem(CLIENT_ID_KEY))
  const isReady = computed(() => isInitialized.value && !isLoading.value)

  // Setters
  function setClient(data) {
    client.value = data
  }

  function setIsInitialized(value) {
    isInitialized.value = value
  }

  function setIsLoading(value) {
    isLoading.value = value
  }

  function setError(message) {
    error.value = message
  }

  function clearError() {
    error.value = null
  }

  function storeClientId(id) {
    if (id) {
      localStorage.setItem(CLIENT_ID_KEY, id)
    }
  }

  function getStoredClientId() {
    return localStorage.getItem(CLIENT_ID_KEY)
  }

  function clearStoredClientId() {
    localStorage.removeItem(CLIENT_ID_KEY)
  }

  function clearClient() {
    client.value = null
    isInitialized.value = false
    localStorage.removeItem(CLIENT_ID_KEY)
  }

  return {
    // State
    client,
    isInitialized,
    isLoading,
    error,
    // Getters
    clientId,
    isReady,
    // Setters
    setClient,
    setIsInitialized,
    setIsLoading,
    setError,
    clearError,
    storeClientId,
    getStoredClientId,
    clearStoredClientId,
    clearClient,
  }
})
