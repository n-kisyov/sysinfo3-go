package main

import (
	"context"
	"fmt"
	"sysinfo3-go/pkg/collector"

	"github.com/shirou/gopsutil/v3/process"
)

// App struct
type App struct {
	ctx           context.Context
	staticSnap    *collector.SystemSnapshot
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.staticSnap = collector.CollectAll(nil)
}

// GetSystemSnapshot returns the current system snapshot.
func (a *App) GetSystemSnapshot() *collector.SystemSnapshot {
	return collector.CollectAll(a.staticSnap)
}

// KillProcess attempts to terminate a process by its PID.
func (a *App) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}
	return p.Kill()
}
