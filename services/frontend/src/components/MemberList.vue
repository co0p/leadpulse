<template>
  <div class="member-list">
    <div v-if="members.length === 0" class="box has-text-centered">
      <p class="has-text-grey">No members found</p>
    </div>

    <div v-else class="table-container">
      <table class="table is-fullwidth is-striped is-hoverable">
        <thead>
          <tr>
            <th>Name</th>
            <th>Seniority</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="member in members" :key="member.id">
            <td>
              <router-link :to="`/members/${member.id}/edit`">
                {{ member.firstName }} {{ member.lastName }}
              </router-link>
            </td>
            <td>{{ member.seniority }}</td>
            <td>
              <span :class="['tag', statusClass(member.status)]">
                {{ formatStatus(member.status) }}
              </span>
            </td>
            <td>
              <div class="buttons are-small">
                <router-link :to="`/members/${member.id}/edit`" class="button is-info">
                  <span class="icon is-small">
                    <i class="fas fa-edit"></i>
                  </span>
                  <span>Edit</span>
                </router-link>
                <button
                  v-if="member.status === 'active'"
                  type="button"
                  class="button is-warning"
                  @click="requestDeactivate(member)"
                >
                  <span class="icon is-small">
                    <i class="fas fa-ban"></i>
                  </span>
                  <span>Deactivate</span>
                </button>
                <button
                  v-else
                  type="button"
                  class="button is-success"
                  @click="requestReactivate(member)"
                >
                  <span class="icon is-small">
                    <i class="fas fa-redo"></i>
                  </span>
                  <span>Reactivate</span>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <ConfirmDialog
      :is-open="showConfirmDialog"
      :title="confirmDialogTitle"
      :message="confirmDialogMessage"
      :confirm-button-label="confirmButtonLabel"
      :is-loading="isConfirmingAction"
      @confirm="handleConfirm"
      @cancel="cancelAction"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMembersStore } from '../stores/members'
import type { MemberDTO } from '../stores/members'
import ConfirmDialog from './ConfirmDialog.vue'

export interface MemberListProps {
  members: MemberDTO[]
}

const props = defineProps<MemberListProps>()

const emit = defineEmits<{
  deactivate: [id: string]
  reactivate: [id: string]
}>()

const router = useRouter()

const showConfirmDialog = ref(false)
const selectedMember = ref<MemberDTO | null>(null)
const actionType = ref<'deactivate' | 'reactivate' | null>(null)
const isConfirmingAction = ref(false)

const confirmDialogTitle = computed(() => {
  return actionType.value === 'deactivate' ? 'Deactivate Member' : 'Reactivate Member'
})

const confirmDialogMessage = computed(() => {
  if (!selectedMember.value) return ''
  const name = `${selectedMember.value.firstName} ${selectedMember.value.lastName}`
  return actionType.value === 'deactivate'
    ? `Are you sure you want to deactivate ${name}?`
    : `Are you sure you want to reactivate ${name}?`
})

const confirmButtonLabel = computed(() => {
  return actionType.value === 'deactivate' ? 'Deactivate' : 'Reactivate'
})

function formatStatus(status: 'active' | 'deactivated'): string {
  return status === 'active' ? 'Active' : 'Deactivated'
}

function statusClass(status: 'active' | 'deactivated'): string {
  return status === 'active' ? 'is-success' : 'is-warning'
}

function requestDeactivate(member: MemberDTO): void {
  selectedMember.value = member
  actionType.value = 'deactivate'
  showConfirmDialog.value = true
}

function requestReactivate(member: MemberDTO): void {
  selectedMember.value = member
  actionType.value = 'reactivate'
  showConfirmDialog.value = true
}

async function handleConfirm(): Promise<void> {
  if (!selectedMember.value || !actionType.value) return

  isConfirmingAction.value = true
  try {
    if (actionType.value === 'deactivate') {
      emit('deactivate', selectedMember.value.id)
    } else {
      emit('reactivate', selectedMember.value.id)
    }
    showConfirmDialog.value = false
  } finally {
    isConfirmingAction.value = false
  }
}

function cancelAction(): void {
  showConfirmDialog.value = false
  selectedMember.value = null
  actionType.value = null
}
</script>

<style scoped>
.table-container {
  overflow-x: auto;
}

a {
  color: var(--bulma-link);
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}
</style>
