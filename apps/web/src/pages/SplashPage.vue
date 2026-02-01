<template>
  <div class="splash-screen">
    <!-- Animated Background -->
    <div class="splash-bg">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <!-- Content -->
    <div class="splash-content">
      <!-- Logo/Icon -->

      <!-- Title -->
      <h1 class="splash-title"><span class="gradient-text">Simulation</span> Catalogue</h1>

      <!-- Subtitle -->
      <p class="splash-subtitle">Fortran Physics Simulations</p>

      <!-- Loading -->
      <div class="splash-loader" v-if="isLoading">
        <q-spinner size="40px" color="primary" />
        <p class="loader-text">{{ loadingMessage }}</p>
      </div>

      <!-- Error State -->
      <div class="splash-error" v-else-if="error">
        <q-icon name="error_outline" size="32px" color="negative" />
        <p class="error-text">{{ error }}</p>
        <q-btn outline color="primary" label="Try Again" @click="initializeApp" />
      </div>
    </div>

    <!-- Footer -->
    <div class="splash-footer">
      <p>Powered by Fortran, Go, Kubernetes & Vue</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from 'stores/client-store'
import { useSimulationStore } from 'stores/simulation-store'
import { initializeClient, fetchClient, fetchSimulations } from 'src/api'

const router = useRouter()
const clientStore = useClientStore()
const simulationStore = useSimulationStore()

const isLoading = ref(true)
const error = ref(null)
const loadingMessage = ref('Initializing...')

onMounted(() => {
  initializeApp()
})

async function initializeApp() {
  isLoading.value = true
  error.value = null

  try {
    // Step 1: Initialize client
    loadingMessage.value = 'Connecting to server...'
    await initializeClientSession()

    // Step 2: Fetch simulations
    loadingMessage.value = 'Loading simulations...'
    await loadSimulations()

    // Step 3: Brief pause for visual effect
    loadingMessage.value = 'Ready!'
    await new Promise((resolve) => setTimeout(resolve, 800))

    // Navigate to catalogue
    router.replace('/catalogue')
  } catch (err) {
    console.error('Initialization error:', err)
    error.value = err.message || 'Failed to initialize application'
    isLoading.value = false
  }
}

async function initializeClientSession() {
  clientStore.setIsLoading(true)
  clientStore.clearError()

  try {
    const storedClientId = clientStore.getStoredClientId()

    if (storedClientId) {
      // Try to fetch existing client
      try {
        const response = await fetchClient()
        clientStore.setClient(response.data)
        clientStore.setIsInitialized(true)
        return response.data
      } catch {
        // Client not found or invalid, create new one
        console.warn('Stored client not found, creating new one...')
        clientStore.clearStoredClientId()
      }
    }

    // Create new client
    const response = await initializeClient()
    clientStore.setClient(response.data)
    clientStore.storeClientId(response.data?.id)
    clientStore.setIsInitialized(true)
    return response.data
  } catch (err) {
    clientStore.setError(
      err.response?.data?.message || err.message || 'Failed to initialize client',
    )
    throw err
  } finally {
    clientStore.setIsLoading(false)
  }
}

async function loadSimulations() {
  simulationStore.setIsLoadingSimulations(true)
  simulationStore.clearError()

  try {
    const response = await fetchSimulations()
    const data = response.data.data || []
    simulationStore.setSimulations(data)
    return data
  } catch (err) {
    simulationStore.setError(err.response?.data?.message || 'Failed to fetch simulations')
    throw err
  } finally {
    simulationStore.setIsLoadingSimulations(false)
  }
}
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.splash-screen {
  position: fixed;
  inset: 0;
  background: $dark-page;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.splash-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;

  .orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.4;
    animation: float 8s ease-in-out infinite;
  }

  .orb-1 {
    width: 400px;
    height: 400px;
    background: linear-gradient(135deg, $primary, $secondary);
    top: -100px;
    right: -100px;
    animation-delay: 0s;
  }

  .orb-2 {
    width: 300px;
    height: 300px;
    background: linear-gradient(135deg, $secondary, $accent);
    bottom: -50px;
    left: -50px;
    animation-delay: -2s;
  }

  .orb-3 {
    width: 250px;
    height: 250px;
    background: linear-gradient(135deg, $accent, $primary);
    top: 40%;
    left: 50%;
    transform: translateX(-50%);
    animation-delay: -4s;
  }
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-30px) scale(1.05);
  }
}

.splash-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 24px;
}

.splash-logo {
  position: relative;
  width: 160px;
  height: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 40px;

  .logo-icon {
    color: $primary;
    z-index: 2;
    animation: pulse-icon 2s ease-in-out infinite;
  }

  .logo-ring {
    position: absolute;
    inset: 0;
    border: 2px solid rgba($primary, 0.3);
    border-radius: 50%;
    animation: ring-expand 3s ease-out infinite;

    &.ring-2 {
      animation-delay: 1.5s;
    }
  }
}

@keyframes pulse-icon {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes ring-expand {
  0% {
    transform: scale(0.8);
    opacity: 1;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

.splash-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 12px 0;
  letter-spacing: -0.03em;
  animation: fadeSlideUp 1s ease-out 0.2s both;

  @media (max-width: 768px) {
    font-size: 2rem;
  }
}

.splash-subtitle {
  font-size: 1.2rem;
  color: #a1a1a6;
  margin: 0;
  animation: fadeSlideUp 1s ease-out 0.4s both;

  @media (max-width: 768px) {
    font-size: 1rem;
  }
}

.splash-loader {
  margin-top: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  animation: fadeSlideUp 1s ease-out 0.6s both;

  .loader-text {
    color: #636366;
    font-size: 0.9rem;
    margin: 0;
  }
}

.splash-error {
  margin-top: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  animation: fadeSlideUp 0.5s ease-out both;

  .error-text {
    color: $negative;
    font-size: 0.95rem;
    margin: 0;
    max-width: 400px;
  }
}

.splash-footer {
  position: absolute;
  bottom: 32px;
  left: 0;
  right: 0;
  text-align: center;
  z-index: 1;

  p {
    color: #48484a;
    font-size: 0.8rem;
    margin: 0;
  }
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
