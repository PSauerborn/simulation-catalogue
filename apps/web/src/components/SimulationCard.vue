<template>
  <div class="simulation-card glass-card" @click="$emit('select', simulation)">
    <div class="card-header">
      <div class="card-title-row">
        <h3 class="card-title">{{ simulation.name }}</h3>
        <q-btn
          v-if="canRun"
          icon="eva-arrow-right-outline"
          round
          size="sm"
          color="primary"
          class="run-btn"
          @click.stop="$emit('run', simulation)"
        >
          <q-tooltip>Run Simulation</q-tooltip>
        </q-btn>
      </div>
      <p class="card-description">{{ simulation.description }}</p>
    </div>

    <!-- Model Section -->
    <div class="model-section" v-if="modelDescription">
      <div class="section-label">
        <q-icon name="eva-hash-outline" size="16px" />
        <span>Mathematical Model</span>
      </div>
      <div class="model-content">
        <div class="math-equation" v-html="modelDescription"></div>
      </div>
    </div>

    <!-- Parameters Section -->
    <div class="parameters-section">
      <div class="section-label">
        <q-icon name="eva-options-2-outline" size="16px" />
        <span>Parameters ({{ simulation.parameters?.length || 0 }})</span>
      </div>
      <div class="parameters-list">
        <div v-for="param in displayedParameters" :key="param.name" class="parameter-item">
          <div class="param-header">
            <span class="param-name">{{ param.name }}</span>
            <span :class="['param-type', param.type]">{{ param.type }}</span>
          </div>
          <p class="param-description">{{ param.description }}</p>
        </div>
        <div v-if="hasMoreParameters" class="more-params">
          <q-icon name="eva-more-horizontal-outline" />
          <span>{{ remainingParameterCount }} more parameters</span>
        </div>
      </div>
    </div>

    <!-- Outputs Section -->
    <div class="outputs-section" v-if="simulation.outputs?.length">
      <div class="section-label">
        <q-icon name="eva-bar-chart-outline" size="16px" />
        <span>Outputs</span>
      </div>
      <div class="outputs-list">
        <q-chip
          v-for="output in simulation.outputs"
          :key="output.name"
          size="sm"
          color="primary"
          text-color="white"
          outline
          class="output-chip"
        >
          {{ output.name }}
        </q-chip>
      </div>
    </div>

    <!-- Footer -->
    <div class="card-footer">
      <span class="card-id">ID: {{ simulation.id?.substring(0, 8) }}...</span>
      <span class="card-date">
        {{ formatDate(simulation.created_at) }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import katex from 'katex'
import 'katex/dist/katex.min.css'

const props = defineProps({
  simulation: {
    type: Object,
    required: true,
  },
  canRun: {
    type: Boolean,
    default: true,
  },
  maxDisplayedParams: {
    type: Number,
    default: 3,
  },
})

defineEmits(['select', 'run'])

// Compute a mathematical model description based on simulation type
// Returns HTML string with rendered KaTeX formulas
const modelDescription = computed(() => {
  const options = { throwOnError: false, displayMode: true }
  return katex.renderToString(props.simulation.model || '', options)
})

const displayedParameters = computed(() => {
  return props.simulation.parameters?.slice(0, props.maxDisplayedParams) || []
})

const hasMoreParameters = computed(() => {
  return (props.simulation.parameters?.length || 0) > props.maxDisplayedParams
})

const remainingParameterCount = computed(() => {
  return (props.simulation.parameters?.length || 0) - props.maxDisplayedParams
})

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.simulation-card {
  padding: 24px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header {
  .card-title-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .card-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #f5f5f7;
    margin: 0;
  }

  .card-description {
    font-size: 0.9rem;
    color: #a1a1a6;
    margin: 0;
    line-height: 1.5;
  }

  .run-btn {
    flex-shrink: 0;
  }
}

.section-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8rem;
  font-weight: 500;
  color: #a1a1a6;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 12px;
}

.model-section {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid $glass-border;
  border-radius: $border-radius-md;
  padding: 16px;

  .model-content {
    font-family: 'SF Mono', 'Monaco', 'Inconsolata', monospace;
    font-size: 0.85rem;
    color: #64d2ff;
    line-height: 1.8;
  }
}

.parameters-section {
  .parameters-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .parameter-item {
    background: rgba(255, 255, 255, 0.02);
    padding: 12px;
    border-radius: $border-radius-sm;
    border: 1px solid rgba(255, 255, 255, 0.05);

    .param-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 4px;
    }

    .param-name {
      font-weight: 500;
      color: #f5f5f7;
      font-family: 'SF Mono', 'Monaco', 'Inconsolata', monospace;
      font-size: 0.9rem;
    }

    .param-description {
      font-size: 0.8rem;
      color: #a1a1a6;
      margin: 0;
      line-height: 1.4;
    }
  }

  .more-params {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #a1a1a6;
    font-size: 0.85rem;
    padding: 8px;
  }
}

.outputs-section {
  .outputs-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .output-chip {
    font-size: 0.75rem;
  }
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  font-size: 0.75rem;
  color: #636366;

  .card-id {
    font-family: 'SF Mono', 'Monaco', 'Inconsolata', monospace;
  }
}
</style>
