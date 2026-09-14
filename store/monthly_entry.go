package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"leadpulse/engine/domain"
)

// SQLiteMonthlyEntryRepository implements the MonthlyEntryRepository interface
// using SQLite as the backing store.
type SQLiteMonthlyEntryRepository struct {
	db *sql.DB
}

// NewSQLiteMonthlyEntryRepository creates a new SQLite-backed monthly entry repository.
func NewSQLiteMonthlyEntryRepository(db *sql.DB) *SQLiteMonthlyEntryRepository {
	return &SQLiteMonthlyEntryRepository{db: db}
}

// Save persists a monthly entry aggregate to the database.
func valueOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// nullableInt converts a *int to sql.NullInt64 for DB writes.
func nullableInt(v *int) interface{} {
	if v == nil {
		return nil
	}
	return int64(*v)
}

func (r *SQLiteMonthlyEntryRepository) Save(entry *domain.MonthlyEntry) error {
	if entry == nil {
		return fmt.Errorf("cannot save nil monthly entry")
	}

	id := entry.ID()
	signals := entry.Signals()

	// Serialize computed scores to JSON if they exist
	var computedScoresJSON sql.NullString
	if entry.HasComputedScores() {
		computed := entry.ComputedScores()
		jsonBytes, err := json.Marshal(computed)
		if err != nil {
			return fmt.Errorf("failed to marshal computed scores: %w", err)
		}
		computedScoresJSON.String = string(jsonBytes)
		computedScoresJSON.Valid = true
	}

	// Serialize computed_at timestamp
	var computedAtStr sql.NullString
	if entry.ComputedAt() != nil {
		computedAtStr.String = entry.ComputedAt().UTC().Format(time.RFC3339)
		computedAtStr.Valid = true
	}

	// Check if entry exists
	existing, err := r.FindByID(id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing entry: %w", err)
	}

	if existing == nil {
		// Insert new entry
		_, err := r.db.Exec(`
            INSERT INTO monthly_entries (
                member_id, month,
                morale, billability, csat, net_margin,
                positive_feedback, critical_feedback,
                overtime_hours, delivery_reliability,
                mentoring_hours, evidence_notes_count,
                computed_scores, created_at, computed_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `,
			int64(id.MemberID),
			id.Month,
			nullableInt(signals.Morale),
			nullableInt(signals.Billability),
			nullableInt(signals.CSAT),
			nullableInt(signals.NetMargin),
			nullableInt(signals.PositiveFeedback),
			nullableInt(signals.CriticalFeedback),
			nullableInt(signals.OvertimeHours),
			nullableInt(signals.DeliveryReliability),
			nullableInt(signals.MentoringHours),
			nullableInt(signals.EvidenceNotesCount),
			computedScoresJSON,
			entry.CreatedAt().UTC().Format(time.RFC3339),
			computedAtStr,
		)
		if err != nil {
			return fmt.Errorf("failed to insert monthly entry: %w", err)
		}
	} else {
		// Update existing entry
		_, err := r.db.Exec(`
            UPDATE monthly_entries SET
                morale = ?, billability = ?, csat = ?, net_margin = ?,
                positive_feedback = ?, critical_feedback = ?,
                overtime_hours = ?, delivery_reliability = ?,
                mentoring_hours = ?, evidence_notes_count = ?,
                computed_scores = ?, computed_at = ?
            WHERE member_id = ? AND month = ?
        `,
			nullableInt(signals.Morale),
			nullableInt(signals.Billability),
			nullableInt(signals.CSAT),
			nullableInt(signals.NetMargin),
			nullableInt(signals.PositiveFeedback),
			nullableInt(signals.CriticalFeedback),
			nullableInt(signals.OvertimeHours),
			nullableInt(signals.DeliveryReliability),
			nullableInt(signals.MentoringHours),
			nullableInt(signals.EvidenceNotesCount),
			computedScoresJSON,
			computedAtStr,
			int64(id.MemberID),
			id.Month,
		)
		if err != nil {
			return fmt.Errorf("failed to update monthly entry: %w", err)
		}
	}

	return nil
}

// FindByID retrieves a monthly entry by its composite ID (member_id, month).
// Returns nil if not found.
func (r *SQLiteMonthlyEntryRepository) FindByID(id domain.MonthlyEntryID) (*domain.MonthlyEntry, error) {
	var morale, billability, csat, netMargin sql.NullInt64
	var positiveFeedback, criticalFeedback sql.NullInt64
	var overtimeHours, deliveryReliability sql.NullInt64
	var mentoringHours, evidenceNotesCount sql.NullInt64
	var createdAtStr string
	var computedScoresJSON sql.NullString
	var computedAtStr sql.NullString

	err := r.db.QueryRow(`
        SELECT
            morale, billability, csat, net_margin,
            positive_feedback, critical_feedback,
            overtime_hours, delivery_reliability,
            mentoring_hours, evidence_notes_count,
            computed_scores, created_at, computed_at
        FROM monthly_entries
        WHERE member_id = ? AND month = ?
    `, int64(id.MemberID), id.Month).Scan(
		&morale, &billability, &csat, &netMargin,
		&positiveFeedback, &criticalFeedback,
		&overtimeHours, &deliveryReliability,
		&mentoringHours, &evidenceNotesCount,
		&computedScoresJSON, &createdAtStr, &computedAtStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly entry: %w", err)
	}

	// Convert sql.NullInt64 -> *int for domain type
	toPtr := func(n sql.NullInt64) *int {
		if !n.Valid {
			return nil
		}
		v := int(n.Int64)
		return &v
	}
	signals, err := domain.NewMonthlyRawSignalsFromPointers(
		toPtr(morale), toPtr(billability), toPtr(csat), toPtr(netMargin),
		toPtr(positiveFeedback), toPtr(criticalFeedback), toPtr(overtimeHours), toPtr(deliveryReliability),
		toPtr(mentoringHours), toPtr(evidenceNotesCount),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct signals: %w", err)
	}

	// Create the monthly entry aggregate
	entry, err := domain.NewMonthlyEntry(int64(id.MemberID), id.Month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct monthly entry: %w", err)
	}

	// Set creation timestamp
	createdAt, err := parseTimestamp(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
	}
	entry.SetCreatedAt(createdAt)

	// Deserialize and set computed scores if present
	if computedScoresJSON.Valid && computedScoresJSON.String != "" {
		var computed domain.ComputedScores
		if err := json.Unmarshal([]byte(computedScoresJSON.String), &computed); err != nil {
			return nil, fmt.Errorf("failed to unmarshal computed scores: %w", err)
		}
		if err := entry.SetComputedScores(computed); err != nil {
			return nil, fmt.Errorf("failed to set computed scores: %w", err)
		}

		// Set computed_at timestamp
		if computedAtStr.Valid && computedAtStr.String != "" {
			computedAt, err := parseTimestamp(computedAtStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse computed_at timestamp: %w", err)
			}
			entry.SetComputedAt(computedAt)
		}
	}

	return entry, nil
}

// FindByMember retrieves all monthly entries for a given team member.
func (r *SQLiteMonthlyEntryRepository) FindByMember(memberID domain.TeamMemberID) ([]*domain.MonthlyEntry, error) {
	rows, err := r.db.Query(`
		SELECT
			morale, billability, csat, net_margin,
			positive_feedback, critical_feedback,
			overtime_hours, delivery_reliability,
			mentoring_hours, evidence_notes_count,
			computed_scores, created_at, computed_at, month
		FROM monthly_entries
		WHERE member_id = ?
		ORDER BY month DESC
	`, int64(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.MonthlyEntry
	for rows.Next() {
		var morale, billability, csat, netMargin sql.NullInt64
		var positiveFeedback, criticalFeedback sql.NullInt64
		var overtimeHours, deliveryReliability sql.NullInt64
		var mentoringHours, evidenceNotesCount sql.NullInt64
		var month, createdAtStr string
		var computedScoresJSON sql.NullString
		var computedAtStr sql.NullString

		err := rows.Scan(
			&morale, &billability, &csat, &netMargin,
			&positiveFeedback, &criticalFeedback,
			&overtimeHours, &deliveryReliability,
			&mentoringHours, &evidenceNotesCount,
			&computedScoresJSON, &createdAtStr, &computedAtStr, &month,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly entry row: %w", err)
		}

		toPtr := func(n sql.NullInt64) *int {
			if !n.Valid {
				return nil
			}
			v := int(n.Int64)
			return &v
		}
		signals, err := domain.NewMonthlyRawSignalsFromPointers(
			toPtr(morale), toPtr(billability), toPtr(csat), toPtr(netMargin),
			toPtr(positiveFeedback), toPtr(criticalFeedback), toPtr(overtimeHours), toPtr(deliveryReliability),
			toPtr(mentoringHours), toPtr(evidenceNotesCount),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct signals: %w", err)
		}

		// Create the monthly entry aggregate
		entry, err := domain.NewMonthlyEntry(int64(memberID), month, signals)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct monthly entry: %w", err)
		}

		// Set creation timestamp
		createdAt, err := parseTimestamp(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}
		entry.SetCreatedAt(createdAt)

		// Deserialize and set computed scores if present
		if computedScoresJSON.Valid && computedScoresJSON.String != "" {
			var computed domain.ComputedScores
			if err := json.Unmarshal([]byte(computedScoresJSON.String), &computed); err != nil {
				return nil, fmt.Errorf("failed to unmarshal computed scores: %w", err)
			}
			if err := entry.SetComputedScores(computed); err != nil {
				return nil, fmt.Errorf("failed to set computed scores: %w", err)
			}

			// Set computed_at timestamp
			if computedAtStr.Valid && computedAtStr.String != "" {
				computedAt, err := parseTimestamp(computedAtStr.String)
				if err != nil {
					return nil, fmt.Errorf("failed to parse computed_at timestamp: %w", err)
				}
				entry.SetComputedAt(computedAt)
			}
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate monthly entries: %w", err)
	}

	return entries, nil
}

// Delete removes a monthly entry.
func (r *SQLiteMonthlyEntryRepository) Delete(id domain.MonthlyEntryID) error {
	result, err := r.db.Exec(
		`DELETE FROM monthly_entries WHERE member_id = ? AND month = ?`,
		int64(id.MemberID), id.Month,
	)
	if err != nil {
		return fmt.Errorf("failed to delete monthly entry: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("monthly entry not found for deletion")
	}

	return nil
}

// === HELPER FUNCTIONS ===

// parseTimestamp parses an RFC3339 timestamp string to time.Time.
// Returns an error if the format is invalid.
func parseTimestamp(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp format: %s", s)
	}
	return t, nil
}
