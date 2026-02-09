<template>
  <q-layout view="hHh lpR fFf">
    <!-- Header -->
    <q-header class="app-header q-pa-sm">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="eva-menu-outline"
          aria-label="Menu"
          @click="toggleDrawer"
          class="lt-md"
        />

        <q-toolbar-title class="toolbar-title">
          <router-link to="/simulations" class="logo-link">
            <span class="logo-text">Simulation <span class="gradient-text">Simulations</span></span>
          </router-link>
        </q-toolbar-title>

        <!-- Centered Text Links -->
        <div class="gt-sm header-center">
          <q-btn flat no-caps label="Simulations" to="/simulations" class="header-text-btn" />
          <q-btn flat no-caps label="Blog" to="/blog" class="header-text-btn" />
        </div>

        <q-space />

        <!-- Right Icons -->
        <div class="gt-sm row items-center no-wrap q-gutter-sm">
          <q-btn class="q-mr-lg" flat round icon="eva-github-outline">
            <q-tooltip>Source Code</q-tooltip>
            <q-menu class="github-menu">
              <q-list dense style="min-width: 200px">
                <q-item
                  clickable
                  v-close-popup
                  href="https://github.com/psauerborn/simulation-catalogue"
                  target="_blank"
                >
                  <q-item-section avatar>
                    <q-icon name="eva-code-outline" size="20px" />
                  </q-item-section>
                  <q-item-section>App Source Code</q-item-section>
                </q-item>
                <q-item
                  clickable
                  v-close-popup
                  href="https://github.com/psauerborn/simulations"
                  target="_blank"
                >
                  <q-item-section avatar>
                    <q-icon name="eva-cube-outline" size="20px" />
                  </q-item-section>
                  <q-item-section>Simulation Source Code</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>

        <!-- Status Indicator -->
        <div class="status-indicator" v-if="clientId">
          <div class="status-dot"></div>
          <span class="status-text gt-sm">Connected</span>
        </div>
      </q-toolbar>
    </q-header>

    <!-- Drawer for mobile -->
    <q-drawer v-model="drawerOpen" :width="280" :breakpoint="1024" bordered class="app-drawer">
      <q-list padding>
        <q-item-label header class="drawer-header"> Navigation </q-item-label>

        <q-item clickable v-ripple to="/simulations" active-class="nav-active">
          <q-item-section avatar>
            <q-icon name="eva-grid-outline" />
          </q-item-section>
          <q-item-section>Catalogue</q-item-section>
        </q-item>

        <q-item clickable v-ripple to="/blog" active-class="nav-active">
          <q-item-section avatar>
            <q-icon name="eva-file-text-outline" />
          </q-item-section>
          <q-item-section>Blog</q-item-section>
        </q-item>

        <q-separator spaced dark />

        <q-item-label header class="drawer-header"> Resources </q-item-label>

        <q-item
          clickable
          v-ripple
          href="https://github.com/psauerborn/simulation-catalogue"
          target="_blank"
        >
          <q-item-section avatar>
            <q-icon name="eva-github-outline" />
          </q-item-section>
          <q-item-section>App Source Code</q-item-section>
          <q-item-section side>
            <q-icon name="eva-external-link-outline" size="16px" />
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple href="https://github.com/psauerborn/simulations" target="_blank">
          <q-item-section avatar>
            <q-icon name="eva-cube-outline" />
          </q-item-section>
          <q-item-section>Simulation Source Code</q-item-section>
          <q-item-section side>
            <q-icon name="eva-external-link-outline" size="16px" />
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple href="https://fortran-lang.org" target="_blank">
          <q-item-section avatar>
            <q-icon name="eva-book-open-outline" />
          </q-item-section>
          <q-item-section>Fortran Docs</q-item-section>
          <q-item-section side>
            <q-icon name="eva-external-link-outline" size="16px" />
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
import { ref } from 'vue'
import { useClientStore } from 'stores/client-store'
import { storeToRefs } from 'pinia'

const clientStore = useClientStore()
const { clientId } = storeToRefs(clientStore)

const drawerOpen = ref(false)

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

.header-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 4px;
}

.header-text-btn {
  font-size: 0.9rem;
  font-weight: 500;
  color: #a1a1a6;
  letter-spacing: 0.01em;

  &.q-btn--active,
  &:hover {
    color: #f5f5f7;
  }
}

.github-menu {
  background: rgba(28, 28, 30, 0.95) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid $glass-border;
  border-radius: $border-radius-md;
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
