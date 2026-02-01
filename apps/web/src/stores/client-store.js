import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const CLIENT_ID_KEY = 'clientId'

export const useClientStore = defineStore('client', () => {
  // State - only data storage
  const client = ref(null)
  const isInitialized = ref(false)
  const isLoading = ref(false)
  const error = ref(null)

  // SSR-safe localStorage check
  const isClient = typeof window !== 'undefined'

  // Getters
  const clientId = computed(() => {
    if (client.value?.id) return client.value.id
    return isClient ? localStorage.getItem(CLIENT_ID_KEY) : null
  })
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
    if (id && isClient) {
      localStorage.setItem(CLIENT_ID_KEY, id)
    }
  }

  function getStoredClientId() {
    return isClient ? localStorage.getItem(CLIENT_ID_KEY) : null
  }

  function clearStoredClientId() {
    if (isClient) {
      localStorage.removeItem(CLIENT_ID_KEY)
    }
  }

  function clearClient() {
    client.value = null
    isInitialized.value = false
    if (isClient) {
      localStorage.removeItem(CLIENT_ID_KEY)
    }
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
