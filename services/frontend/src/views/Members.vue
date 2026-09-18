<template>
  <div class="members-view">
    <section class="hero is-light">
      <div class="hero-body">
        <div class="container">
          <div class="level">
            <div class="level-left">
              <div class="level-item">
                <h1 class="title">Members</h1>
              </div>
            </div>
            <div class="level-right">
              <div class="level-item">
                <router-link to="/members/add" class="button is-primary">
                  <span class="icon is-small">
                    <i class="fas fa-plus"></i>
                  </span>
                  <span>Add Member</span>
                </router-link>
              </div>
            </div>
          </div>
          <p class="subtitle">Manage team members and their status</p>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <!-- Loading State -->
        <div v-if="store.loading" class="box has-text-centered">
          <div class="is-loading">
            <p class="has-text-info">
              <span class="icon">
                <i class="fas fa-spinner fa-spin"></i>
              </span>
              <span>Loading members...</span>
            </p>
          </div>
        </div>

        <!-- Error Toast -->
        <ErrorToast
          v-if="store.error"
          :message="store.error"
          @close="store.clearError"
        />

        <!-- Main Content -->
        <div v-if="!store.loading">
          <!-- Tab Navigation -->
          <MemberTabs
            :current-tab="store.currentTab"
            @select-tab="handleTabChange"
          />

          <!-- Members List -->
          <MemberList
            :members="store.filteredMembers"
            @deactivate="handleDeactivate"
            @reactivate="handleReactivate"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMembersStore } from '../stores/members'
import type { CurrentTab } from '../stores/members'
import MemberTabs from '../components/MemberTabs.vue'
import MemberList from '../components/MemberList.vue'
import ErrorToast from '../components/ErrorToast.vue'

const store = useMembersStore()

/**
 * Load members when component mounts
 */
onMounted(async () => {
  try {
    await store.loadMembers()
  } catch (err) {
    // Error is already set in store
  }
})

/**
 * Handle tab change
 */
async function handleTabChange(tab: CurrentTab): Promise<void> {
  store.setCurrentTab(tab)
  try {
    await store.loadMembers()
  } catch (err) {
    // Error is already set in store
  }
}

/**
 * Handle deactivate member
 */
async function handleDeactivate(id: string): Promise<void> {
  try {
    await store.deactivateMember(id)
  } catch (err) {
    // Error is already set in store
  }
}

/**
 * Handle reactivate member
 */
async function handleReactivate(id: string): Promise<void> {
  try {
    await store.reactivateMember(id)
  } catch (err) {
    // Error is already set in store
  }
}
</script>

<style scoped>
/* All styling handled by Bulma classes */
</style>
