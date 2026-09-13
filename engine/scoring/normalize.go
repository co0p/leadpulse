package scoring

// NormalizeMorale converts morale self-score (0-5) to a 0-100 scale.
// PRD 5.2: MoraleN = (Morale / 5) * 100
func NormalizeMorale(morale int) float64 {
	return float64(morale) / 5.0 * 100.0
}

// NormalizeBillability converts billability percentage to a 0-100 score.
// PRD 5.2: BillabilityN = clamp(100 - (abs(BillabilityPct - 75) / 15) * 100, 0, 100)
// Optimal zone is 75%; deviation in either direction reduces the score.
func NormalizeBillability(billabilityPercent int) float64 {
	deviation := float64(billabilityPercent - 75)
	if deviation < 0 {
		deviation = -deviation
	}
	return Clamp(100.0-(deviation/15.0)*100.0, 0.0, 100.0)
}

// NormalizeCSAT converts CSAT (1-5) to a 0-100 scale.
// PRD 5.2: CSATN = ((CSAT - 1) / 4) * 100
func NormalizeCSAT(csat int) float64 {
	return float64(csat-1) / 4.0 * 100.0
}

// NormalizeMargin converts net margin % to a 0-100 scale.
// PRD 5.2: MarginN = norm01(NetMarginPct, -20, 60) * 100
// Expected range: -20 to +60
func NormalizeMargin(netMarginPercent int) float64 {
	return Norm01(float64(netMarginPercent), -20.0, 60.0) * 100.0
}

// NormalizePositive converts positive feedback count to a 0-100 scale.
// PRD 5.2: PositiveN = min(PositiveFeedbackCount, 8) / 8 * 100
// Capped at 8 items; more than 8 counts as 100.
func NormalizePositive(positiveFeedbackCount int) float64 {
	count := float64(positiveFeedbackCount)
	if count > 8.0 {
		count = 8.0
	}
	return count / 8.0 * 100.0
}

// NormalizeCritical converts critical feedback count to a 0-100 score.
// PRD 5.2: CriticalN = 100 - (min(CriticalFeedbackCount, 6) / 6 * 100)
// Higher critical feedback reduces the score. Capped at 6 items.
func NormalizeCritical(criticalFeedbackCount int) float64 {
	count := float64(criticalFeedbackCount)
	if count > 6.0 {
		count = 6.0
	}
	return 100.0 - (count / 6.0 * 100.0)
}

// NormalizeOvertime converts overtime hours to a 0-100 score.
// PRD 5.2: OvertimeN = 100 - (min(OvertimeHours, 30) / 30 * 100)
// More overtime reduces the score. Capped at 30 hours.
func NormalizeOvertime(overtimeHours int) float64 {
	hours := float64(overtimeHours)
	if hours > 30.0 {
		hours = 30.0
	}
	return 100.0 - (hours / 30.0 * 100.0)
}

// NormalizeDelivery converts delivery reliability % to a 0-100 scale.
// PRD 5.2: DeliveryN = clamp(DeliveryReliabilityPct, 0, 100)
// Input is already 0-100, just clamped to ensure bounds.
func NormalizeDelivery(deliveryReliabilityPercent int) float64 {
	return Clamp(float64(deliveryReliabilityPercent), 0.0, 100.0)
}

// NormalizeMentoring converts mentoring/enablement hours to a 0-100 scale.
// PRD 5.2: MentoringN = min(MentoringHours, 12) / 12 * 100
// Capped at 12 hours; more than 12 counts as 100.
func NormalizeMentoring(mentoringHours int) float64 {
	hours := float64(mentoringHours)
	if hours > 12.0 {
		hours = 12.0
	}
	return hours / 12.0 * 100.0
}

// NormalizeEvidence converts evidence notes count to a 0-100 scale.
// PRD 5.2: EvidenceN = min(EvidenceNotesCount, 6) / 6 * 100
// Capped at 6 notes; more than 6 counts as 100.
func NormalizeEvidence(evidenceNotesCount int) float64 {
	count := float64(evidenceNotesCount)
	if count > 6.0 {
		count = 6.0
	}
	return count / 6.0 * 100.0
}
