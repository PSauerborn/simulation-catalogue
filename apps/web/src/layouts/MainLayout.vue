<template>
  <q-layout view="hHh lpR fFf">
    <!-- Header -->
    <q-header class="app-header">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleDrawer" class="lt-md" />

        <q-toolbar-title class="toolbar-title">
          <router-link to="/catalogue" class="logo-link">
            <q-icon name="science" size="28px" class="logo-icon" />
            <span class="logo-text">Simulation <span class="gradient-text">Catalogue</span></span>
          </router-link>
        </q-toolbar-title>

        <q-space />

        <!-- Status Indicator -->
        <div class="status-indicator" v-if="clientId">
          <div class="status-dot"></div>
          <span class="status-text gt-sm">Connected</span>
        </div>

        <!-- Client Info -->
        <q-btn flat round icon="person" class="gt-sm">
          <q-tooltip>Client: {{ clientIdShort }}</q-tooltip>
        </q-btn>
      </q-toolbar>
    </q-header>

    <!-- Drawer for mobile -->
    <q-drawer v-model="drawerOpen" :width="280" :breakpoint="1024" bordered class="app-drawer">
      <q-list padding>
        <q-item-label header class="drawer-header"> Navigation </q-item-label>

        <q-item clickable v-ripple to="/catalogue" active-class="nav-active">
          <q-item-section avatar>
            <q-icon name="grid_view" />
          </q-item-section>
          <q-item-section>Catalogue</q-item-section>
        </q-item>

        <q-separator spaced dark />

        <q-item-label header class="drawer-header"> Resources </q-item-label>

        <q-item clickable v-ripple href="https://github.com" target="_blank">
          <q-item-section avatar>
            <q-icon name="code" />
          </q-item-section>
          <q-item-section>Source Code</q-item-section>
          <q-item-section side>
            <q-icon name="open_in_new" size="16px" />
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple href="https://fortran-lang.org" target="_blank">
          <q-item-section avatar>
            <q-icon name="school" />
          </q-item-section>
          <q-item-section>Fortran Docs</q-item-section>
          <q-item-section side>
            <q-icon name="open_in_new" size="16px" />
          </q-item-section>
        </q-item>
      </q-list>

      <q-space />

      <!-- Drawer Footer -->
      <div class="drawer-footer">
        <p class="version-text">Simulation Catalogue</p>
        <p class="copyright-text">Powered by Fortran & Quasar</p>
      </div>
    </q-drawer>

    <!-- Main Content -->
    <q-page-container>
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useClientStore } from 'stores/client-store'
import { storeToRefs } from 'pinia'

const clientStore = useClientStore()
const { clientId } = storeToRefs(clientStore)

const drawerOpen = ref(false)

const clientIdShort = computed(() => {
  return clientId.value ? `${clientId.value.substring(0, 8)}...` : 'Unknown'
})

function toggleDrawer() {
  drawerOpen.value = !drawerOpen.value
}
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.app-header {
  background: rgba(0, 0, 0, 0.7) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid $glass-border;
}

.toolbar-title {
  .logo-link {
    display: flex;
    align-items: center;
    gap: 12px;
    text-decoration: none;
    color: #f5f5f7;

    .logo-icon {
      color: $primary;
    }

    .logo-text {
      font-weight: 600;
      font-size: 1.1rem;
      letter-spacing: -0.01em;

      @media (max-width: 600px) {
        display: none;
      }
    }
  }
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: rgba($positive, 0.1);
  border-radius: 100px;
  margin-right: 12px;

  .status-dot {
    width: 8px;
    height: 8px;
    background: $positive;
    border-radius: 50%;
    animation: pulse-dot 2s ease-in-out infinite;
  }

  .status-text {
    font-size: 0.8rem;
    color: $positive;
    font-weight: 500;
  }
}

@keyframes pulse-dot {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.app-drawer {
  background: rgba(28, 28, 30, 0.95) !important;
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
  border-right: 1px solid $glass-border !important;

  .drawer-header {
    font-size: 0.75rem;
    font-weight: 600;
    color: #636366;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .nav-active {
    background: rgba($primary, 0.15);
    color: $primary;

    .q-icon {
      color: $primary;
    }
  }

  .drawer-footer {
    padding: 24px 16px;
    border-top: 1px solid $glass-border;

    .version-text {
      font-size: 0.85rem;
      color: #a1a1a6;
      margin: 0 0 4px 0;
    }

    .copyright-text {
      font-size: 0.75rem;
      color: #48484a;
      margin: 0;
    }
  }
}

// Page transitions
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.3s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
