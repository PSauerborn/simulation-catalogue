<template>
  <q-page class="catalogue-page">
    <!-- Header Section -->
    <header class="page-header">
      <div class="header-content">
        <div class="header-text">
          <h1 class="page-title"><span class="gradient-text">Physics</span> Simulations</h1>
          <p class="page-subtitle">Explore and run Fortran-powered physics simulations</p>
        </div>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ simulations.length }}</span>
            <span class="stat-label">Simulations</span>
          </div>
        </div>
      </div>
    </header>

    <!-- Simulation Runner Dialog -->
    <SimulationRunner
      v-model="showRunner"
      :selected-simulation="selectedSimulation"
      @view-results="showResults = true"
    />

    <!-- Search & Filter Bar -->
    <div class="filter-bar">
      <q-input
        v-model="searchQuery"
        placeholder="Search simulations..."
        filled
        dense
        dark
        class="search-input"
      >
        <template #prepend>
          <q-icon name="eva-search-outline" color="grey" />
        </template>
        <template #append v-if="searchQuery">
          <q-icon name="eva-close-outline" @click="searchQuery = ''" class="cursor-pointer" />
        </template>
      </q-input>

      <q-btn-toggle
        v-model="viewMode"
        toggle-color="primary"
        :options="[
          { value: 'grid', slot: 'grid' },
          { value: 'list', slot: 'list' },
        ]"
        flat
        dense
        class="view-toggle"
      >
        <template #grid>
          <q-icon name="eva-grid-outline" />
        </template>
        <template #list>
          <q-icon name="eva-list-outline" />
        </template>
      </q-btn-toggle>
    </div>

    <!-- Loading State -->
    <div v-if="isLoadingSimulations" class="loading-state">
      <q-spinner-orbit size="64px" color="primary" />
      <p>Loading simulations...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredSimulations.length === 0" class="empty-state">
      <q-icon name="eva-flask-outline" size="80px" color="grey-6" />
      <h3>No Simulations Found</h3>
      <p v-if="searchQuery">No simulations match "{{ searchQuery }}". Try a different search.</p>
      <p v-else>No simulations are available at this time.</p>
    </div>

    <!-- Simulations Grid -->

    <div v-else :class="['simulations-container', viewMode]">
      <TransitionGroup name="card-fade">
        <SimulationCard
          v-for="simulation in filteredSimulations"
          :key="simulation.id"
          :simulation="simulation"
          :can-run="canRun"
          @select="selectSimulation"
          @run="runSelectedSimulation"
        />
      </TransitionGroup>
    </div>

    <!-- Results Viewer Dialog -->
    <ResultsViewer
      v-model="showResults"
      :output-data="runOutput"
      :simulation-name="selectedSimulation?.name"
    />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSimulationStore } from 'stores/simulation-store'
import { storeToRefs } from 'pinia'
import SimulationCard from 'components/SimulationCard.vue'
import SimulationRunner from 'components/SimulationRunner.vue'
import ResultsViewer from 'components/ResultsViewer.vue'
import { fetchSimulations } from 'src/api'

const simulationStore = useSimulationStore()

const { simulations, isLoadingSimulations, canRun, runOutput } = storeToRefs(simulationStore)

const searchQuery = ref('')
const viewMode = ref('grid')
const selectedSimulation = ref(null)
const showRunner = ref(false)
const showResults = ref(false)

const filteredSimulations = computed(() => {
  if (!searchQuery.value) return simulations.value

  const query = searchQuery.value.toLowerCase()
  return simulations.value.filter((sim) => {
    return (
      sim.name?.toLowerCase().includes(query) ||
      sim.description?.toLowerCase().includes(query) ||
      sim.parameters?.some((p) => p.name.toLowerCase().includes(query))
    )
  })
})

function selectSimulation(simulation) {
  selectedSimulation.value = simulation
  showRunner.value = true
}

async function runSelectedSimulation(simulation) {
  selectSimulation(simulation)
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

onMounted(async () => {
  // Refresh simulations if needed
  if (simulations.value.length === 0) {
    await loadSimulations()
  }
})
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.catalogue-page {
  min-height: 100vh;
  padding-bottom: 48px;
}

.page-header {
  padding: 48px 24px 24px;
  background: linear-gradient(180deg, rgba(48, 209, 88, 0.08) 0%, transparent 100%);

  .header-content {
    max-width: 1400px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    flex-wrap: wrap;
    gap: 24px;
  }

  .header-text {
    .page-title {
      font-size: 2.5rem;
      font-weight: 700;
      margin: 0 0 8px 0;
      letter-spacing: -0.02em;

      @media (max-width: 768px) {
        font-size: 1.75rem;
      }
    }

    .page-subtitle {
      font-size: 1.1rem;
      color: #a1a1a6;
      margin: 0;

      @media (max-width: 768px) {
        font-size: 0.95rem;
      }
    }
  }

  .header-stats {
    display: flex;
    gap: 32px;

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;

      .stat-value {
        font-size: 2rem;
        font-weight: 600;
        color: $primary;
      }

      .stat-label {
        font-size: 0.8rem;
        color: #636366;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }
    }
  }
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 24px;
  margin: 24px auto;
  max-width: 1400px;

  .search-input {
    flex: 1;
    max-width: 400px;
  }

  .view-toggle {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 4px;
    flex-shrink: 0;
  }

  @media (max-width: 600px) {
    padding: 0 12px;
    gap: 12px;

    .search-input {
      max-width: none;
    }
  }
}

.loading-state {
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  text-align: center;
  padding: 48px 24px;

  h3 {
    color: #f5f5f7;
    margin: 24px 0 8px 0;
  }

  p {
    color: #a1a1a6;
    max-width: 400px;
  }
}

.simulations-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;

  @media (max-width: 600px) {
    padding: 0 12px;
  }

  &.grid {
    display: grid;
    // Wider cards for desktop
    grid-template-columns: repeat(auto-fill, minmax(450px, 1fr));
    gap: 24px;

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
      gap: 16px;
    }
  }

  &.list {
    display: flex;
    flex-direction: column;
    gap: 16px;

    :deep(.simulation-card) {
      max-width: none;
    }
  }
}

// Card transition animations
.card-fade-enter-active {
  transition: all 0.4s ease-out;
}

.card-fade-leave-active {
  transition: all 0.3s ease-in;
}

.card-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.card-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.card-fade-move {
  transition: transform 0.4s ease;
}
</style>
