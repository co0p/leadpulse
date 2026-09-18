<template>
  <div class="member-form">
    <form @submit.prevent="handleSubmit">
      <!-- First Name Field -->
      <div class="field">
        <label class="label">First Name</label>
        <div class="control">
          <input
            v-model="form.firstName"
            type="text"
            class="input"
            :class="{ 'is-danger': fieldErrors.firstName }"
            placeholder="Enter first name"
            @blur="validateField('firstName')"
          />
        </div>
        <p v-if="fieldErrors.firstName" class="help is-danger">
          {{ fieldErrors.firstName }}
        </p>
      </div>

      <!-- Last Name Field -->
      <div class="field">
        <label class="label">Last Name</label>
        <div class="control">
          <input
            v-model="form.lastName"
            type="text"
            class="input"
            :class="{ 'is-danger': fieldErrors.lastName }"
            placeholder="Enter last name"
            @blur="validateField('lastName')"
          />
        </div>
        <p v-if="fieldErrors.lastName" class="help is-danger">
          {{ fieldErrors.lastName }}
        </p>
      </div>

      <!-- Seniority Field -->
      <div class="field">
        <label class="label">Seniority</label>
        <div class="control">
          <div class="select" :class="{ 'is-danger': fieldErrors.seniority }">
            <select
              v-model="form.seniority"
              @blur="validateField('seniority')"
            >
              <option value="">Select a seniority level</option>
              <option value="Junior">Junior</option>
              <option value="Mid">Mid</option>
              <option value="Senior">Senior</option>
              <option value="Lead">Lead</option>
            </select>
          </div>
        </div>
        <p v-if="fieldErrors.seniority" class="help is-danger">
          {{ fieldErrors.seniority }}
        </p>
      </div>

      <!-- Form Actions -->
      <div class="field is-grouped mt-5">
        <div class="control">
          <button
            type="submit"
            class="button is-primary"
            :disabled="!isFormValid || isSubmitting"
          >
            <span v-if="isSubmitting" class="icon is-small">
              <i class="fas fa-spinner fa-spin"></i>
            </span>
            <span>{{ submitButtonLabel }}</span>
          </button>
        </div>
        <div class="control">
          <button
            type="button"
            class="button is-light"
            @click="$emit('cancel')"
            :disabled="isSubmitting"
          >
            Cancel
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

export interface MemberFormData {
  firstName: string
  lastName: string
  seniority: string
}

export interface MemberFormProps {
  initialData?: MemberFormData
  isSubmitting?: boolean
  submitButtonLabel?: string
}

const props = withDefaults(defineProps<MemberFormProps>(), {
  isSubmitting: false,
  submitButtonLabel: 'Save',
})

const emit = defineEmits<{
  submit: [data: MemberFormData]
  cancel: []
}>()

const form = ref<MemberFormData>({
  firstName: '',
  lastName: '',
  seniority: '',
})

const fieldErrors = ref<Record<string, string>>({
  firstName: '',
  lastName: '',
  seniority: '',
})

const VALID_SENIORITIES = ['Junior', 'Mid', 'Senior', 'Lead']

/**
 * Validate a single field
 */
function validateField(fieldName: keyof MemberFormData): void {
  if (fieldName === 'firstName') {
    if (!form.value.firstName.trim()) {
      fieldErrors.value.firstName = 'First name is required'
    } else {
      fieldErrors.value.firstName = ''
    }
  } else if (fieldName === 'lastName') {
    if (!form.value.lastName.trim()) {
      fieldErrors.value.lastName = 'Last name is required'
    } else {
      fieldErrors.value.lastName = ''
    }
  } else if (fieldName === 'seniority') {
    if (!form.value.seniority) {
      fieldErrors.value.seniority = 'Seniority level is required'
    } else if (!VALID_SENIORITIES.includes(form.value.seniority)) {
      fieldErrors.value.seniority = 'Invalid seniority level'
    } else {
      fieldErrors.value.seniority = ''
    }
  }
}

/**
 * Validate all fields
 */
function validateForm(): boolean {
  validateField('firstName')
  validateField('lastName')
  validateField('seniority')

  return (
    !fieldErrors.value.firstName &&
    !fieldErrors.value.lastName &&
    !fieldErrors.value.seniority
  )
}

/**
 * Computed: is form valid?
 */
const isFormValid = computed(() => {
  return (
    form.value.firstName.trim() !== '' &&
    form.value.lastName.trim() !== '' &&
    VALID_SENIORITIES.includes(form.value.seniority) &&
    !fieldErrors.value.firstName &&
    !fieldErrors.value.lastName &&
    !fieldErrors.value.seniority
  )
})

/**
 * Handle form submission
 */
async function handleSubmit(): Promise<void> {
  if (!validateForm()) {
    return
  }

  await emit('submit', {
    firstName: form.value.firstName.trim(),
    lastName: form.value.lastName.trim(),
    seniority: form.value.seniority,
  })
}

/**
 * Pre-fill form when initialData prop changes
 */
watch(
  () => props.initialData,
  (newData) => {
    if (newData) {
      form.value = {
        firstName: newData.firstName || '',
        lastName: newData.lastName || '',
        seniority: newData.seniority || '',
      }
      // Clear errors after pre-fill
      fieldErrors.value = {
        firstName: '',
        lastName: '',
        seniority: '',
      }
    }
  },
  { immediate: true }
)

onMounted(() => {
  if (props.initialData) {
    form.value = {
      firstName: props.initialData.firstName || '',
      lastName: props.initialData.lastName || '',
      seniority: props.initialData.seniority || '',
    }
  }
})
</script>

<style scoped>
/* No custom styles - Bulma classes handle all layout and styling */
</style>
