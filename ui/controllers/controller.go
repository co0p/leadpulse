package controllers

// Controller is the base interface for all UI controllers.
// Controllers own business logic and state; Fyne screens are thin renderers.
type Controller interface {
	// Load initializes the controller with data from services.
	Load() error
	// Validate checks if the current state is valid.
	Validate() error
}
