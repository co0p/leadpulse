import type { Overview } from './employeeApi'

export interface OverviewState {
  overview: Overview | null
  loading: boolean
  error: string
}

export const emptyOverviewState = (): OverviewState => ({
  overview: null,
  loading: false,
  error: '',
})

export const formatScore = (value: number | undefined): string => {
  if (value === undefined || Number.isNaN(value)) return '—'
  return value.toFixed(1)
}

export const hasEvidenceGap = (overview: Overview | null): boolean => {
  if (!overview) return false
  return overview.evidenceGaps.length > 0
}
