<template>
  <div class="app-shell">
    <!-- Top bar -->
    <header class="topbar" role="banner">
      <div class="topbar-left">
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
        <div class="breadcrumb-area">
          <span class="breadcrumb-text">{{ pageTitle }}</span>
        </div>
      </div>

      <div class="topbar-center">
        <div class="field has-addons is-fullwidth">
          <p class="control is-expanded has-icons-left">
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

      <div class="topbar-right">
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
        class="sidebar"
        :class="{ 'is-active': sidebarOpen }"
        role="navigation"
      >
        <div class="sidebar-header">
          <div class="logo">
            <h2>LeadPulse</h2>
          </div>
        </div>

        <nav class="sidebar-nav">
          <ul class="menu-list">
            <li>
              <RouterLink to="/" @click="closeSidebar">
                <span class="icon"><i class="fas fa-home"></i></span>
                <span>Home</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/members" @click="closeSidebar">
                <span class="icon"><i class="fas fa-users"></i></span>
                <span>Members</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/alerts" @click="closeSidebar">
                <span class="icon"><i class="fas fa-exclamation-circle"></i></span>
                <span>Alerts</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/reports" @click="closeSidebar">
                <span class="icon"><i class="fas fa-chart-bar"></i></span>
                <span>Reports</span>
              </RouterLink>
            </li>
          </ul>
        </nav>

        <div class="sidebar-footer">
          <button class="button is-fullwidth is-primary" type="button">
            <span class="icon-text">
              <span class="icon"><i class="fas fa-plus"></i></span>
              <span>Add Item</span>
            </span>
          </button>
        </div>
      </aside>

      <!-- Sidebar overlay for mobile -->
      <div
        class="sidebar-overlay"
        :class="{ 'is-active': sidebarOpen }"
        role="presentation"
        @click="closeSidebar"
      ></div>

      <!-- Main content -->
      <main class="main-content" role="main">
        <RouterView />
      </main>
    </div>

    <!-- Footer -->
    <footer class="app-footer" role="contentinfo">
      <div class="footer-content">
        <div class="footer-left">
          <p>&copy; 2026 LeadPulse. All rights reserved.</p>
        </div>
        <div class="footer-right">
          <HealthIndicator />
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import HealthIndicator from './HealthIndicator.vue'
import { useHealthStore } from '../stores/health'

const route = useRoute()
const healthStore = useHealthStore()

const sidebarOpen = ref(window.innerWidth >= 1024)

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/': 'Home',
    '/members': 'Members',
    '/alerts': 'Alerts',
    '/reports': 'Reports'
  }
  return titles[route.path] || 'Team Impact Scorecard'
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

// Handle window resize
onMounted(() => {
  window.addEventListener('resize', () => {
    if (window.innerWidth >= 1024) {
      sidebarOpen.value = true
    } else {
      sidebarOpen.value = false
    }
  })

  // Handle escape key to close sidebar on mobile
  window.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      closeSidebar()
    }
  })

  // Trigger health check on app mount
  healthStore.checkHealth()
})
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100vh;
  overflow: hidden;
}

.app-shell {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #fafafa;
}

.topbar {
  height: 56px;
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  padding: 0 1.25rem;
  gap: 1rem;
  flex-shrink: 0;
  z-index: 10;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 200px;
}

.topbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.topbar-center .field {
  width: 100%;
  max-width: 400px;
}

.topbar-right {
  display: flex;
  align-items: center;
  min-width: 120px;
}

.sidebar-toggle {
  padding: 0.5rem;
  min-width: 44px;
}

.sidebar-toggle:focus {
  outline: 2px solid #3273dc;
  outline-offset: 2px;
}

.button:focus {
  outline: 2px solid #3273dc;
  outline-offset: 2px;
}

.input:focus {
  outline: 2px solid #3273dc;
  outline-offset: 2px;
  border-color: #3273dc;
}

.breadcrumb-area {
  display: flex;
  align-items: center;
  font-size: 0.95rem;
  color: #7a7a7a;
}

.app-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 260px;
  background-color: #f5f5f5;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-shrink: 0;
  z-index: 5;
}

.sidebar-header {
  padding: 1.5rem 1.25rem;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.sidebar-header .logo h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

.menu-list {
  list-style: none;
  padding: 0 0.5rem;
}

.menu-list li {
  margin: 0;
}

.menu-list a {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  color: #4a4a4a;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.menu-list a:hover {
  background-color: #ebebeb;
}

.menu-list a:focus {
  outline: 2px solid #3273dc;
  outline-offset: 2px;
}

.sidebar-footer {
  padding: 1rem 0.5rem;
  border-top: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.sidebar-footer .button {
  font-size: 0.9rem;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  background-color: #fafafa;
}

.sidebar-overlay {
  display: none;
  position: absolute;
  top: 56px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 4;
}

.sidebar-overlay.is-active {
  display: block;
}

.app-footer {
  height: 56px;
  background-color: #fff;
  border-top: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  padding: 0 1.25rem;
  flex-shrink: 0;
}

.footer-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-left p {
  font-size: 0.85rem;
  color: #7a7a7a;
  margin: 0;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

/* Responsive */
@media screen and (max-width: 1023px) {
  .sidebar {
    position: absolute;
    left: 0;
    top: 56px;
    height: calc(100vh - 56px);
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  }

  .sidebar.is-active {
    transform: translateX(0);
  }
}

@media screen and (max-width: 767px) {
  .topbar {
    padding: 0 0.75rem;
    gap: 0.5rem;
  }

  .topbar-center .field {
    max-width: 200px;
  }

  .topbar-center .input {
    font-size: 14px;
  }

  .topbar-right .buttons {
    display: flex;
    gap: 0.5rem;
  }

  .topbar-right .button {
    padding: 0.5rem;
    min-width: auto;
  }

  .app-footer {
    flex-direction: column;
    height: auto;
    padding: 0.75rem 1.25rem;
    gap: 0.5rem;
  }

  .footer-content {
    flex-direction: column;
    gap: 0.5rem;
  }

  .footer-left p {
    font-size: 0.75rem;
  }
}
</style>
