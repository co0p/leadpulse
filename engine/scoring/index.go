package scoring

// ComputeTII computes the Total Impact Index (TII).
// PRD 5.5: TII = 0.20*DG + 0.35*DP + 0.25*DT + 0.20*DO
// Range: 0–100. Overall monthly score for a team member.
//
// Dimension weights:
// - DG (Personal Growth) 20%
// - DP (Project Impact) 35%
// - DT (Team Impact) 25%
// - DO (Organization Impact) 20%
func ComputeTII(dg, dp, dt, do float64) float64 {
	return 0.20*dg + 0.35*dp + 0.25*dt + 0.20*do
}

// ComputeCompleteness computes the completeness percentage.
// PRD 4.3: CompletenessPct = (FilledRequiredFields / 11) * 100
// Range: 0–100. Measures how fully the Team Lead populated the monthly entry.
// The 11 required fields are the 10 raw signals plus confirmation that
// at least one impact rating has been set for those signals.
func ComputeCompleteness(filledRequiredFields int) float64 {
	return float64(filledRequiredFields) / 11.0 * 100.0
}
