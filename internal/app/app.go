package app

import (
	"context"
	"time"

	"leadpulse/internal/employees"
	"leadpulse/internal/overview"
	"leadpulse/internal/performances"
)

type App struct {
	ctx             context.Context
	employeeRepo    employees.Repository
	performanceRepo performances.Repository
}

// New creates the Wails application binding.
func New(employeeRepo employees.Repository, performanceRepo performances.Repository) *App {
	return &App{
		employeeRepo:    employeeRepo,
		performanceRepo: performanceRepo,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListEmployees() ([]employees.Employee, error) {
	return a.employeeRepo.List()
}

func (a *App) AddEmployee(input employees.Input) (employees.Employee, error) {
	return a.employeeRepo.Add(input)
}

func (a *App) RemoveEmployee(id int64) error {
	return a.employeeRepo.Remove(id)
}

// TeamPulseOverview returns the current team pulse overview.
func (a *App) TeamPulseOverview() (overview.Overview, error) {
	now := time.Now()
	useCase := overview.NewUseCase(a.employeeRepo, a.performanceRepo, now.Year(), int(now.Month()))
	return useCase.Get()
}
