<template>
  <div v-if="isOpen" class="modal is-active">
    <div class="modal-background" @click="cancel"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{{ title }}</p>
        <button
          type="button"
          class="delete"
          aria-label="Close dialog"
          @click="cancel"
        ></button>
      </header>
      <section class="modal-card-body">
        <p>{{ message }}</p>
      </section>
      <footer class="modal-card-foot">
        <button
          type="button"
          class="button is-light"
          @click="cancel"
          :disabled="isLoading"
        >
          {{ cancelButtonLabel }}
        </button>
        <button
          type="button"
          class="button is-danger"
          @click="confirm"
          :disabled="isLoading"
        >
          <span v-if="isLoading" class="icon is-small">
            <i class="fas fa-spinner fa-spin"></i>
          </span>
          <span>{{ confirmButtonLabel }}</span>
        </button>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface ConfirmDialogProps {
  isOpen: boolean
  title: string
  message: string
  confirmButtonLabel?: string
  cancelButtonLabel?: string
  isLoading?: boolean
}

const props = withDefaults(defineProps<ConfirmDialogProps>(), {
  confirmButtonLabel: 'Confirm',
  cancelButtonLabel: 'Cancel',
  isLoading: false,
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

function confirm(): void {
  emit('confirm')
}

function cancel(): void {
  emit('cancel')
}
</script>

<style scoped>
/* Bulma handles all modal styling */
</style>
