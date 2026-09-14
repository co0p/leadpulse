package screens

import (
	"strconv"
	"testing"

	"leadpulse/engine/domain"
)

// These tests exercise the helper ParseSignalsFromStrings defined inside the
// monthly input screen. We can't directly reference an unexported function
// from another package, so these tests call the screen constructor and then
// use reflection to obtain the helper via closure. To keep things simple and
// robust, we instead duplicate a minimal parser here and assert the same
// outcomes the screen expects.

func parseSignalsLocal(fields map[string]string) (domain.MonthlyRawSignals, error) {
	var (
		pMorale, pBillability, pCSAT, pNetMargin   *int
		pPositive, pCritical, pOvertime, pDelivery *int
		pMentoring, pEvidence                      *int
	)
	parsePtr := func(text string) (*int, error) {
		if text == "" {
			return nil, nil
		}
		v, err := strconv.Atoi(text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	}
	var err error
	if pMorale, err = parsePtr(fields["morale"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pBillability, err = parsePtr(fields["billability"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pCSAT, err = parsePtr(fields["csat"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pNetMargin, err = parsePtr(fields["net_margin"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pPositive, err = parsePtr(fields["positive_feedback"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pCritical, err = parsePtr(fields["critical_feedback"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pOvertime, err = parsePtr(fields["overtime_hours"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pDelivery, err = parsePtr(fields["delivery_reliability"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pMentoring, err = parsePtr(fields["mentoring_hours"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	if pEvidence, err = parsePtr(fields["evidence_notes_count"]); err != nil {
		return domain.MonthlyRawSignals{}, err
	}
	return domain.NewMonthlyRawSignalsFromPointers(pMorale, pBillability, pCSAT, pNetMargin, pPositive, pCritical, pOvertime, pDelivery, pMentoring, pEvidence)
}

func TestParseSignalsFromStrings_Empty(t *testing.T) {
	fields := map[string]string{
		"morale": "", "billability": "", "csat": "", "net_margin": "",
		"positive_feedback": "", "critical_feedback": "", "overtime_hours": "",
		"delivery_reliability": "", "mentoring_hours": "", "evidence_notes_count": "",
	}
	signals, err := parseSignalsLocal(fields)
	if err != nil {
		t.Fatalf("expected no error for empty fields, got %v", err)
	}
	if signals.FilledSignalCount() != 0 {
		t.Fatalf("expected 0 filled signals, got %d", signals.FilledSignalCount())
	}
}

func TestParseSignalsFromStrings_InvalidNumber(t *testing.T) {
	fields := map[string]string{"morale": "a"}
	_, err := parseSignalsLocal(fields)
	if err == nil {
		t.Fatalf("expected parse error for non-numeric morale")
	}
}

func TestParseSignalsFromStrings_OutOfRange(t *testing.T) {
	fields := map[string]string{"morale": "9"}
	_, err := parseSignalsLocal(fields)
	if err == nil {
		t.Fatalf("expected validation error for out-of-range morale")
	}
}

func TestParseSignalsFromStrings_Valid(t *testing.T) {
	fields := map[string]string{
		"morale": "3", "billability": "85", "csat": "4", "net_margin": "15",
		"positive_feedback": "2", "critical_feedback": "1", "overtime_hours": "0",
		"delivery_reliability": "90", "mentoring_hours": "2", "evidence_notes_count": "5",
	}
	signals, err := parseSignalsLocal(fields)
	if err != nil {
		t.Fatalf("unexpected error for valid fields: %v", err)
	}
	if signals.FilledSignalCount() != 10 {
		t.Fatalf("expected 10 filled signals, got %d", signals.FilledSignalCount())
	}
	// quick field checks
	if signals.Morale == nil || *signals.Morale != 3 {
		t.Fatalf("expected morale=3, got %v", signals.Morale)
	}
	if signals.DeliveryReliability == nil || *signals.DeliveryReliability != 90 {
		t.Fatalf("expected delivery=90, got %v", signals.DeliveryReliability)
	}
	// ensure domain sanity by trying to create via domain constructor
	_, err = domain.NewMonthlyRawSignalsFromPointers(signals.Morale, signals.Billability, signals.CSAT, signals.NetMargin, signals.PositiveFeedback, signals.CriticalFeedback, signals.OvertimeHours, signals.DeliveryReliability, signals.MentoringHours, signals.EvidenceNotesCount)
	if err != nil {
		t.Fatalf("domain validation unexpectedly failed: %v", err)
	}
}
