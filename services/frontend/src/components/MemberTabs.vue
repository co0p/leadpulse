<template>
  <div class="tabs">
    <ul>
      <li :class="{ 'is-active': currentTab === 'active' }">
        <a @click="selectTab('active')">
          <span class="icon is-small">
            <i class="fas fa-check-circle"></i>
          </span>
          <span>Active</span>
        </a>
      </li>
      <li :class="{ 'is-active': currentTab === 'deactivated' }">
        <a @click="selectTab('deactivated')">
          <span class="icon is-small">
            <i class="fas fa-ban"></i>
          </span>
          <span>Deactivated</span>
        </a>
      </li>
      <li :class="{ 'is-active': currentTab === 'all' }">
        <a @click="selectTab('all')">
          <span class="icon is-small">
            <i class="fas fa-list"></i>
          </span>
          <span>All</span>
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
export type TabType = 'active' | 'deactivated' | 'all'

export interface MemberTabsProps {
  currentTab: TabType
}

const props = defineProps<MemberTabsProps>()

const emit = defineEmits<{
  selectTab: [tab: TabType]
}>()

function selectTab(tab: TabType): void {
  emit('selectTab', tab)
}
</script>

<style scoped>
/* Bulma handles all tab styling via .tabs class */
</style>
