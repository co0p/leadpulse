package scoring

// ComputeDimensionGrowth computes the Personal Growth dimension score (DG).
// PRD 5.4: DG = Σ(weight_i * C_i) / 100
// Where C_i = contribution scores (0-100) and weights sum to 100.
//
// Weights per PRD 5.4:
// - MoraleN 20
// - CriticalN 15
// - PositiveN 10
// - MentoringN 15
// - DeliveryN 10
// - BillabilityN 10
// - OvertimeN 15
// - EvidenceN 5
func ComputeDimensionGrowth(morale, critical, positive, mentoring, delivery, billability, overtime, evidence float64) float64 {
	weighted := 20.0*morale + 15.0*critical + 10.0*positive + 15.0*mentoring +
		10.0*delivery + 10.0*billability + 15.0*overtime + 5.0*evidence
	return weighted / 100.0
}

// ComputeDimensionProject computes the Project Impact dimension score (DP).
// PRD 5.4: DP = Σ(weight_i * C_i) / 100
//
// Weights per PRD 5.4:
// - DeliveryN 25
// - CSATN 20
// - MarginN 20
// - BillabilityN 15
// - CriticalN 10
// - PositiveN 5
// - MoraleN 5
func ComputeDimensionProject(delivery, csat, margin, billability, critical, positive, morale float64) float64 {
	weighted := 25.0*delivery + 20.0*csat + 20.0*margin + 15.0*billability +
		10.0*critical + 5.0*positive + 5.0*morale
	return weighted / 100.0
}

// ComputeDimensionTeam computes the Team Impact dimension score (DT).
// PRD 5.4: DT = Σ(weight_i * C_i) / 100
//
// Weights per PRD 5.4:
// - MentoringN 25
// - PositiveN 20
// - CriticalN 20
// - MoraleN 15
// - DeliveryN 10
// - OvertimeN 10
func ComputeDimensionTeam(mentoring, positive, critical, morale, delivery, overtime float64) float64 {
	weighted := 25.0*mentoring + 20.0*positive + 20.0*critical + 15.0*morale +
		10.0*delivery + 10.0*overtime
	return weighted / 100.0
}

// ComputeDimensionOrg computes the Organization Impact dimension score (DO).
// PRD 5.4: DO = Σ(weight_i * C_i) / 100
//
// Weights per PRD 5.4:
// - MarginN 25
// - CSATN 20
// - BillabilityN 15
// - DeliveryN 15
// - MentoringN 10
// - PositiveN 10
// - EvidenceN 5
func ComputeDimensionOrg(margin, csat, billability, delivery, mentoring, positive, evidence float64) float64 {
	weighted := 25.0*margin + 20.0*csat + 15.0*billability + 15.0*delivery +
		10.0*mentoring + 10.0*positive + 5.0*evidence
	return weighted / 100.0
}
