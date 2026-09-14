package screens

import "testing"

func TestScreenRegistry_ProvidesMonthlyInputScreen(t *testing.T) {
	registry := NewScreenRegistry(nil, nil)
	if registry.MonthlyInputScreen() == nil {
		t.Fatal("expected monthly input screen")
	}
}
