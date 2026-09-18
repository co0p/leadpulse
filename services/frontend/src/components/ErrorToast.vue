<template>
  <transition name="fade">
    <div v-if="isVisible" class="notification is-danger is-fixed-bottom is-fixed-right" role="alert" style="margin: 1rem; z-index: 1000;">
      <button
        type="button"
        class="delete"
        aria-label="Close error message"
        @click="close"
      ></button>
      <div class="content">
        <p>
          <span class="icon is-small">
            <i class="fas fa-exclamation-circle"></i>
          </span>
          <span>{{ message }}</span>
        </p>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

export interface ErrorToastProps {
  message: string
  duration?: number // milliseconds, 0 = no auto-dismiss
  autoShow?: boolean // whether to auto-show on mount (default true)
}

const props = withDefaults(defineProps<ErrorToastProps>(), {
  duration: 5000,
  autoShow: true,
})

const emit = defineEmits<{
  close: []
}>()

const isVisible = ref(props.autoShow)
let timeoutId: ReturnType<typeof setTimeout> | null = null

function close(): void {
  isVisible.value = false
  if (timeoutId !== null) {
    clearTimeout(timeoutId)
  }
  emit('close')
}

onMounted(() => {
  if (props.autoShow && props.duration > 0) {
    timeoutId = setTimeout(() => {
      close()
    }, props.duration)
  }
})
</script>

<style scoped>
.is-fixed-bottom {
  position: fixed;
  bottom: 0;
}

.is-fixed-right {
  right: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
