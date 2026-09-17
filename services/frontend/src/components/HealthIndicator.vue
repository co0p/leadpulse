<template>
  <div class="health-indicator" :class="`is-${currentStatus}`" :title="statusTitle">
    <span class="health-icon">
      <i
        :class="[
          'fas',
          currentStatus === 'healthy' ? 'fa-check-circle' : 
          currentStatus === 'unhealthy' ? 'fa-times-circle' :
          'fa-spinner fa-spin'
        ]"
        aria-hidden="true"
      ></i>
    </span>
    <span class="sr-only">{{ statusText }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useHealthStore } from '../stores/health'

const healthStore = useHealthStore()

const currentStatus = computed(() => healthStore.status)

const statusText = computed(() => {
  switch (currentStatus.value) {
    case 'healthy':
      return 'Backend is healthy'
    case 'unhealthy':
      return `Backend is unreachable: ${healthStore.checkError || 'Unknown error'}`
    case 'checking':
      return 'Checking backend health...'
    default:
      return 'Unknown status'
  }
})

const statusTitle = computed(() => statusText.value)
</script>

<style scoped>
.health-indicator {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.health-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.health-indicator.is-healthy .health-icon {
  color: #48c774;
}

.health-indicator.is-unhealthy .health-icon {
  color: #f14668;
}

.health-indicator.is-checking .health-icon {
  color: #ffdd57;
  animation: pulse 1s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* Visually hidden text for screen readers */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
