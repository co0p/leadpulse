package store

import (
	"database/sql"
	"fmt"
	"time"
)

// logAuditEvent records a change to a member field in the audit trail.
// This is an internal helper used by member CRUD operations.
func logAuditEvent(db *sql.DB, memberID int64, field, oldValue, newValue string) error {
	_, err := db.Exec(`
		INSERT INTO audit_log (member_id, field, old_value, new_value, changed_at)
		VALUES (?, ?, ?, ?, ?)
	`, memberID, field, oldValue, newValue, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to log audit event for member %d field %s: %w", memberID, field, err)
	}
	return nil
}
