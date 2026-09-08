import { describe, expect, it } from 'vitest'
import { emptyPageId } from './emptyPage'

describe('empty application page', () => {
  it('renders the initial empty-page landmark', () => {
    expect(emptyPageId).toBe('empty-page')
  })
})