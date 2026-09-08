<script lang="ts" setup>
import { onMounted, reactive } from 'vue'
import { employeeApi, type Employee } from './employeeApi'
import { addEmployeeToState, emptyEmployeeListState, employeeInput, removeEmployeeFromState } from './employeeListState'
import { addView, detailsView, emptyView, type ViewState } from './viewState'
import { emptyOverviewState, formatScore, hasEvidenceGap } from './overviewState'
import { overviewMarkers } from './overviewMarkers'

const state = reactive(emptyEmployeeListState())
const overview = reactive(emptyOverviewState())
const form = reactive({ firstName: '', secondName: '', seniority: 'Junior', startDate: '' })
const view = reactive<{ current: ViewState }>({ current: emptyView() })

const selectedEmployee = () => {
  const currentView = view.current
  if (currentView.kind !== 'details') return undefined
  return state.employees.find((employee) => employee.id === currentView.employeeId)
}

const loadEmployees = async () => {
  state.loading = true
  state.error = ''
  try {
    state.employees = (await employeeApi.list()) ?? []
  } catch (error) {
    state.error = error instanceof Error ? error.message : 'Unable to load employees.'
  } finally {
    state.loading = false
  }
}

const loadOverview = async () => {
  overview.loading = true
  overview.error = ''
  try {
    overview.overview = await employeeApi.teamPulseOverview()
  } catch (error) {
    overview.error = error instanceof Error ? error.message : 'Unable to load overview.'
  } finally {
    overview.loading = false
  }
}

const progressClass = (value: number): string => {
  if (value >= 4) return 'is-success'
  if (value >= 3) return 'is-warning'
  return 'is-danger'
}

const progressPercent = (value: number): number => Math.min(Math.max((value / 5) * 100, 0), 100)

const anomalyClass = (type: string): string => {
  if (type === 'BurnoutRisk') return 'is-danger'
  if (type === 'MoraleDrop') return 'is-warning'
  return 'is-info'
}

const addEmployee = async () => {
  state.error = ''
  try {
    const employee = await employeeApi.add(employeeInput(form.firstName, form.secondName, form.seniority, form.startDate))
    Object.assign(state, addEmployeeToState(state, employee))
    Object.assign(form, { firstName: '', secondName: '', seniority: 'Junior', startDate: '' })
    view.current = emptyView()
  } catch (error) {
    state.error = error instanceof Error ? error.message : 'Unable to add employee.'
  }
}

const removeEmployee = async (employee: Employee) => {
  state.error = ''
  try {
    await employeeApi.remove(employee.id)
    Object.assign(state, removeEmployeeFromState(state, employee.id))
    if (view.current.kind === 'details' && view.current.employeeId === employee.id) view.current = emptyView()
  } catch (error) {
    state.error = error instanceof Error ? error.message : 'Unable to remove employee.'
  }
}

onMounted(() => {
  loadEmployees()
  loadOverview()
})
</script>

<template>
  <main id="employee-list" class="columns is-gapless is-fullheight" aria-label="LeadPulse employee list">
    <aside id="employee-navigation" class="column is-one-quarter has-background-light p-4">
      <div class="is-flex is-justify-content-space-between is-align-items-center mb-4">
        <h1 class="title is-5 mb-0">Employees</h1>
        <button class="button is-primary is-small" type="button" aria-label="Add employee" title="Add employee" @click="view.current = addView()">+</button>
      </div>
      <nav class="menu" aria-label="Employees">
        <ul class="menu-list">
          <li v-for="employee in state.employees" :key="employee.id">
            <button class="button is-white is-fullwidth is-justify-content-flex-start" type="button" @click="view.current = detailsView(employee.id)">
              {{ employee.firstName }} {{ employee.secondName }}
            </button>
          </li>
        </ul>
      </nav>
    </aside>

    <section class="column p-6">
      <div v-if="state.loading">Loading employees...</div>
      <div v-else-if="state.error" role="alert">{{ state.error }}</div>
      <section v-else-if="view.current.kind === 'add'" id="add-employee-screen" aria-label="Add employee">
        <h2 class="title">Add employee</h2>
        <form @submit.prevent="addEmployee">
          <div class="field"><label class="label" for="first-name">First name</label><div class="control"><input id="first-name" class="input" v-model="form.firstName" required /></div></div>
          <div class="field"><label class="label" for="second-name">Second name</label><div class="control"><input id="second-name" class="input" v-model="form.secondName" required /></div></div>
          <div class="field"><label class="label" for="seniority">Seniority</label><div class="select is-fullwidth"><select id="seniority" v-model="form.seniority"><option>Junior</option><option>Midlevel</option><option>Senior</option><option>Principal</option></select></div></div>
          <div class="field"><label class="label" for="start-date">Start date</label><div class="control"><input id="start-date" class="input" v-model="form.startDate" type="date" required /></div></div>
          <button class="button is-primary" type="submit">Add employee</button>
        </form>
      </section>
      <section v-else-if="view.current.kind === 'details' && selectedEmployee()" id="employee-details-screen" aria-label="Employee details">
        <h2 class="title">{{ selectedEmployee()?.firstName }} {{ selectedEmployee()?.secondName }}</h2>
        <p class="subtitle">{{ selectedEmployee()?.seniority }}</p>
        <p>Start date: {{ selectedEmployee()?.startDate }}</p>
        <button class="button is-danger mt-5" type="button" @click="removeEmployee(selectedEmployee()!)">Remove</button>
      </section>
      <section v-else :id="overviewMarkers.regionId" :aria-label="overviewMarkers.regionId">
        <div v-if="overview.loading" class="mb-4">Loading team pulse...</div>
        <div v-else-if="overview.error" class="notification is-danger" role="alert">{{ overview.error }}</div>
        <div v-else-if="overview.overview">
          <h2 class="title">Team Pulse</h2>

          <div v-if="hasEvidenceGap(overview.overview)" class="notification is-warning" :id="overviewMarkers.evidenceGapsId" role="status">
            <p class="has-text-weight-bold">Evidence gaps detected</p>
            <ul>
              <li v-for="gap in overview.overview.evidenceGaps" :key="`${gap.employeeId}-${gap.metric}`">{{ gap.name }} — {{ gap.message }}</li>
            </ul>
          </div>

          <div :id="overviewMarkers.metricBarId" class="columns is-multiline mb-5">
            <div class="column is-3">
              <div class="box has-text-centered">
                <p class="heading">Headcount</p>
                <p class="title">{{ overview.overview.headcount }}</p>
              </div>
            </div>
            <div class="column is-3">
              <div class="box has-text-centered">
                <p class="heading">Rolling Team Morale</p>
                <p class="title">{{ formatScore(overview.overview.rollingTeamMorale) }}</p>
              </div>
            </div>
            <div class="column is-3">
              <div class="box has-text-centered">
                <p class="heading">High Capacity Risk</p>
                <p class="title">{{ overview.overview.highCapacityCount }}</p>
              </div>
            </div>
            <div class="column is-3">
              <div class="box has-text-centered">
                <p class="heading">Team Multiplier</p>
                <p class="title">{{ formatScore(overview.overview.teamMultiplier) }}</p>
              </div>
            </div>
          </div>

          <div :id="overviewMarkers.healthDimensionsId" class="columns mb-5">
            <div class="column is-6">
              <h3 class="subtitle is-5">Rolling 3-Month Team Health</h3>
              <div v-for="(value, key) in overview.overview.healthDimensions" :key="key" class="field">
                <label class="label is-capitalized">{{ key }}</label>
                <progress class="progress is-medium" :class="progressClass(value)" :value="progressPercent(value)" max="100">{{ formatScore(value) }}</progress>
                <p class="help">{{ formatScore(value) }} / 5</p>
              </div>
            </div>
            <div class="column is-6">
              <h3 class="subtitle is-5">Actionable Anomalies</h3>
              <div v-if="overview.overview.anomalies.length === 0" class="notification is-success">No anomalies detected.</div>
              <div v-else :id="overviewMarkers.anomaliesId" class="content">
                <div v-for="anomaly in overview.overview.anomalies" :key="`${anomaly.employeeId}-${anomaly.type}`" class="notification" :class="anomalyClass(anomaly.type)" role="alert">
                  <p class="has-text-weight-bold">{{ anomaly.type }}</p>
                  <p>{{ anomaly.message }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </section>
  </main>
</template>
