<template>
  <div>
    <section class="hero is-light">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">Add Member</h1>
          <p class="subtitle">Create a new team member</p>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container" style="max-width: 800px;">
        <!-- Error Toast -->
        <ErrorToast
          v-if="store.error"
          :message="store.error"
          @close="store.clearError"
        />

        <!-- Form Container -->
        <div class="box">
          <MemberForm
            :is-submitting="store.loading"
            submit-button-label="Create Member"
            @submit="handleSubmit"
            @cancel="navigateBack"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useMembersStore } from '../stores/members'
import MemberForm from '../components/MemberForm.vue'
import ErrorToast from '../components/ErrorToast.vue'
import type { MemberFormData } from '../components/MemberForm.vue'

const router = useRouter()
const store = useMembersStore()

/**
 * Handle form submission
 */
async function handleSubmit(data: MemberFormData): Promise<void> {
  try {
    await store.addMember({
      firstName: data.firstName,
      lastName: data.lastName,
      seniority: data.seniority,
    })
    // Navigate back to members list on success
    router.push('/members')
  } catch (err) {
    // Error is already set in store
  }
}

/**
 * Navigate back to members list
 */
function navigateBack(): void {
  router.push('/members')
}
</script>

<style scoped>
/* All styling handled by Bulma classes */
</style>
