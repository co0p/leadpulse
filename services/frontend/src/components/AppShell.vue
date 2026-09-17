<template>
  <div class="app-shell">
    <!-- Sidebar -->
    <aside class="sidebar" data-testid="sidebar">
      <nav class="navbar is-primary">
        <div class="navbar-brand">
          <div class="navbar-item">
            <h1 class="title is-4" style="color: white">LeadPulse</h1>
          </div>
        </div>
      </nav>
      <div class="sidebar-menu">
        <p class="menu-label">Main</p>
        <ul class="menu-list">
          <li><a href="/">Dashboard</a></li>
          <li><a href="/members">Members</a></li>
          <li><a href="/analytics">Analytics</a></li>
        </ul>
      </div>
    </aside>

    <!-- Top Bar -->
    <div class="top-bar" data-testid="topbar">
      <nav class="navbar is-light">
        <div class="navbar-end">
          <div class="navbar-item">
            <span>Welcome to LeadPulse</span>
          </div>
        </div>
      </nav>
    </div>

    <!-- Main Content -->
    <main class="main-content">
      <slot></slot>
    </main>

    <!-- Footer -->
    <footer class="footer" data-testid="footer">
      <div class="content has-text-centered">
        <p>
          <span data-testid="health-indicator" class="health-indicator" :class="{ 'is-healthy': isBackendHealthy, 'is-unhealthy': !isBackendHealthy }">
            <span v-if="isBackendHealthy" class="tag is-success">✓ Backend Healthy</span>
            <span v-else class="tag is-danger">✗ Backend Unreachable</span>
          </span>
        </p>
        <p><small>&copy; 2026 LeadPulse. All rights reserved.</small></p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isBackendHealthy = ref(false)

const checkBackendHealth = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/health`, {
      method: 'GET',
    })
    isBackendHealthy.value = response.ok
  } catch (error) {
    console.error('Health check error:', error)
    isBackendHealthy.value = false
  }
}

onMounted(() => {
  checkBackendHealth()
  // Check health every 5 seconds
  setInterval(checkBackendHealth, 5000)
})
</script>

<style scoped>
.app-shell {
  display: grid;
  grid-template-columns: 250px 1fr;
  grid-template-rows: auto 1fr auto;
  height: 100vh;
}

.sidebar {
  grid-column: 1 / 2;
  grid-row: 1 / 4;
  background-color: #f5f5f5;
  border-right: 1px solid #ddd;
  overflow-y: auto;
}

.sidebar-menu {
  padding: 1.5rem;
}

.menu-label {
  font-weight: 600;
  color: #7a7a7a;
  padding: 0.5rem 0;
}

.menu-list li a {
  display: block;
  padding: 0.75rem;
  color: #363636;
  text-decoration: none;
  border-radius: 4px;
}

.menu-list li a:hover {
  background-color: #e8e8e8;
  color: #00d1b2;
}

.top-bar {
  grid-column: 2 / 3;
  grid-row: 1 / 2;
  border-bottom: 1px solid #ddd;
  background-color: white;
}

.main-content {
  grid-column: 2 / 3;
  grid-row: 2 / 3;
  padding: 2rem;
  overflow-y: auto;
  background-color: #fafafa;
}

.footer {
  grid-column: 1 / 3;
  grid-row: 3 / 4;
  border-top: 1px solid #ddd;
  background-color: #f5f5f5;
  padding: 1rem;
}

.health-indicator {
  font-weight: 600;
}

.health-indicator.is-healthy {
  color: #48c774;
}

.health-indicator.is-unhealthy {
  color: #f14668;
}
</style>
