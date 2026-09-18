<template>
  <div class="app-shell">
    <!-- Top bar -->
    <header class="app-topbar" role="banner">
      <div class="app-topbar-left">
        <button
          class="button is-white sidebar-toggle"
          type="button"
          aria-label="Toggle sidebar"
          :aria-expanded="sidebarOpen"
          @click="toggleSidebar"
        >
          <span class="icon">
            <i class="fas fa-bars"></i>
          </span>
        </button>
        <span class="app-page-title">{{ pageTitle }}</span>
      </div>

      <div class="app-topbar-center">
        <div class="field is-fullwidth">
          <p class="control has-icons-left">
            <input
              class="input"
              type="text"
              placeholder="Search..."
              aria-label="Search members and content"
            />
            <span class="icon is-left">
              <i class="fas fa-search"></i>
            </span>
          </p>
        </div>
      </div>

      <div class="app-topbar-right">
        <div class="buttons">
          <button
            class="button is-ghost"
            type="button"
            aria-label="Notifications"
            title="Notifications"
          >
            <span class="icon">
              <i class="fas fa-bell"></i>
            </span>
          </button>
          <button
            class="button is-ghost"
            type="button"
            aria-label="Help"
            title="Help"
          >
            <span class="icon">
              <i class="fas fa-question-circle"></i>
            </span>
          </button>
        </div>
      </div>
    </header>

    <div class="app-container">
      <!-- Sidebar -->
      <aside
        class="app-sidebar"
        :class="{ 'app-sidebar--open': sidebarOpen }"
        role="navigation"
        aria-label="Main navigation"
      >
        <div class="app-sidebar-header">
          <p class="is-size-5 has-text-weight-semibold">LeadPulse</p>
        </div>

        <nav class="app-sidebar-nav">
          <aside class="menu">
            <ul class="menu-list">
              <li>
                <RouterLink to="/" @click="closeSidebar">
                  <span class="icon-text">
                    <span class="icon"><i class="fas fa-home"></i></span>
                    <span>Home</span>
                  </span>
                </RouterLink>
              </li>
              <li>
                <RouterLink to="/members" @click="closeSidebar">
                  <span class="icon-text">
                    <span class="icon"><i class="fas fa-users"></i></span>
                    <span>Members</span>
                  </span>
                </RouterLink>
              </li>
              <li>
                <RouterLink to="/alerts" @click="closeSidebar">
                  <span class="icon-text">
                    <span class="icon"><i class="fas fa-exclamation-circle"></i></span>
                    <span>Alerts</span>
                  </span>
                </RouterLink>
              </li>
              <li>
                <RouterLink to="/reports" @click="closeSidebar">
                  <span class="icon-text">
                    <span class="icon"><i class="fas fa-chart-bar"></i></span>
                    <span>Reports</span>
                  </span>
                </RouterLink>
              </li>
            </ul>
          </aside>
        </nav>

        <div class="app-sidebar-footer">
          <button
            v-if="isOnMembersPage"
            class="button is-primary is-fullwidth"
            type="button"
            @click="navigateToAddMember"
          >
            <span class="icon-text">
              <span class="icon"><i class="fas fa-plus"></i></span>
              <span>Add Member</span>
            </span>
          </button>
        </div>
      </aside>

      <!-- Main content -->
      <main class="app-main" role="main">
        <RouterView />
      </main>
    </div>

    <!-- Footer -->
    <footer class="app-footer" role="contentinfo">
      <nav class="level is-mobile">
        <div class="level-left">
          <div class="level-item">
            <p class="is-size-7 has-text-grey">&copy; 2026 LeadPulse. All rights reserved.</p>
          </div>
        </div>
        <div class="level-right">
          <div class="level-item">
            <HealthIndicator />
          </div>
        </div>
      </nav>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import HealthIndicator from './HealthIndicator.vue'
import { useHealthStore } from '../stores/health'

const route = useRoute()
const router = useRouter()
const healthStore = useHealthStore()

const sidebarOpen = ref(window.innerWidth >= 1024)

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/': 'Home',
    '/members': 'Members',
    '/members/add': 'Add Member',
    '/alerts': 'Alerts',
    '/reports': 'Reports'
  }
  return titles[route.path] || 'Team Impact Scorecard'
})

const isOnMembersPage = computed(() => {
  return route.path.startsWith('/members')
})

function toggleSidebar() {
  if (window.innerWidth < 1024) {
    sidebarOpen.value = !sidebarOpen.value
  }
}

function closeSidebar() {
  if (window.innerWidth < 1024) {
    sidebarOpen.value = false
  }
}

function navigateToAddMember() {
  router.push('/members/add')
  closeSidebar()
}

onMounted(() => {
  window.addEventListener('resize', () => {
    sidebarOpen.value = window.innerWidth >= 1024
  })

  window.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeSidebar()
  })

  healthStore.checkHealth()
})
</script>

<style scoped>
/* ── Shell frame ─────────────────────────────────────────────────── */
.app-shell {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background-color: var(--bulma-scheme-main-ter, #fafafa);
}

/* ── Top bar ─────────────────────────────────────────────────────── */
.app-topbar {
  height: 56px;
  background-color: var(--bulma-scheme-main, #fff);
  border-bottom: 1px solid var(--bulma-border, #dbdbdb);
  display: flex;
  align-items: center;
  padding: 0 1.25rem;
  gap: 1rem;
  flex-shrink: 0;
  z-index: 10;
}

.app-topbar-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 200px;
}

.app-page-title {
  font-size: 0.95rem;
  color: var(--bulma-text-weak, #7a7a7a);
}

.app-topbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.app-topbar-center .field {
  width: 100%;
  max-width: 400px;
}

.app-topbar-right {
  display: flex;
  align-items: center;
  min-width: 120px;
  justify-content: flex-end;
}

.sidebar-toggle {
  min-width: 44px;
}

/* ── App body (sidebar + main) ───────────────────────────────────── */
.app-container {
  display: flex;
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* ── Sidebar ─────────────────────────────────────────────────────── */
.app-sidebar {
  width: 260px;
  background-color: var(--bulma-scheme-main-bis, #f5f5f5);
  border-right: 1px solid var(--bulma-border, #dbdbdb);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-shrink: 0;
  z-index: 100;
}

.app-sidebar-header {
  padding: 1.25rem 1rem;
  border-bottom: 1px solid var(--bulma-border, #dbdbdb);
  flex-shrink: 0;
}

.app-sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 0.75rem 0.5rem;
}

.app-sidebar-footer {
  padding: 0.75rem 0.5rem;
  border-top: 1px solid var(--bulma-border, #dbdbdb);
  flex-shrink: 0;
}

/* ── Sidebar overlay (mobile only) ──────────────────────────────── */
.app-sidebar-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 98;
  pointer-events: auto;
  top: 56px; /* below topbar */
}

/* ── Main content ────────────────────────────────────────────────── */
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  background-color: var(--bulma-scheme-main-ter, #fafafa);
}

/* ── Footer ──────────────────────────────────────────────────────── */
.app-footer {
  background-color: var(--bulma-scheme-main, #fff);
  border-top: 1px solid var(--bulma-border, #dbdbdb);
  padding: 0 1.25rem;
  flex-shrink: 0;
  min-height: 52px;
  display: flex;
  align-items: center;
}

.app-footer .level {
  width: 100%;
  margin-bottom: 0;
}

/* ── Responsive: sidebar slides in on mobile ─────────────────────── */
@media screen and (max-width: 1023px) {
  .app-sidebar {
    position: fixed;
    left: 0;
    top: 56px;
    bottom: 52px;
    width: 260px;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.12);
    z-index: 1000;
  }

  .app-sidebar.app-sidebar--open {
    transform: translateX(0);
  }

  /* Add scrim/backdrop when sidebar is open */
  .app-sidebar.app-sidebar--open::before {
    content: '';
    position: fixed;
    inset: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 999;
  }
}

@media screen and (max-width: 767px) {
  .app-topbar {
    padding: 0 0.75rem;
    gap: 0.5rem;
  }

  .app-topbar-center .field {
    max-width: 180px;
  }

  .app-footer {
    min-height: auto;
    padding: 0.75rem;
  }

  .app-footer .level {
    flex-direction: column;
    gap: 0.25rem;
  }
}
</style>
