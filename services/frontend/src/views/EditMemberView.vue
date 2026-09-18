<template>
  <div>
    <section class="hero is-light">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">Edit Member</h1>
          <p class="subtitle">Update member information</p>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container" style="max-width: 800px;">
        <!-- Loading State -->
        <div v-if="isLoadingMember" class="box has-text-centered">
          <p class="has-text-info">
            <span class="icon">
              <i class="fas fa-spinner fa-spin"></i>
            </span>
            <span>Loading member...</span>
          </p>
        </div>

        <!-- Error Toast -->
        <ErrorToast
          v-if="store.error"
          :message="store.error"
          @close="store.clearError"
        />

        <!-- Not Found -->
        <div v-if="!isLoadingMember && !member" class="box has-text-centered has-text-danger">
          <p>Member not found</p>
          <router-link to="/members" class="button is-light">
            Back to Members
          </router-link>
        </div>

        <!-- Form Container -->
        <div v-if="member" class="box">
          <MemberForm
            :initial-data="formData"
            :is-submitting="store.loading"
            submit-button-label="Update Member"
            @submit="handleSubmit"
            @cancel="navigateBack"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMembersStore } from '../stores/members'
import MemberForm from '../components/MemberForm.vue'
import ErrorToast from '../components/ErrorToast.vue'
import type { MemberFormData } from '../components/MemberForm.vue'

const router = useRouter()
const route = useRoute()
const store = useMembersStore()

const isLoadingMember = ref(false)

const memberId = computed(() => route.params.id as string)

const member = computed(() => {
  return store.findMemberById(memberId.value)
})

const formData = computed(() => {
  if (!member.value) return undefined
  return {
    firstName: member.value.firstName,
    lastName: member.value.lastName,
    seniority: member.value.seniority,
  }
})

/**
 * Load members list if not already loaded
 */
onMounted(async () => {
  // If members list is empty, load it
  if (store.members.length === 0) {
    isLoadingMember.value = true
    try {
      await store.loadMembers()
    } catch (err) {
      // Error is already set in store
    } finally {
      isLoadingMember.value = false
    }
  }
})

/**
 * Handle form submission
 */
async function handleSubmit(data: MemberFormData): Promise<void> {
  try {
    await store.editMember(memberId.value, {
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
