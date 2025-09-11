//go:build !windows

package windows

import (
	"github.com/shirou/gopsutil/v4/process"
)

func initializeProcessInfo(p *process.Process, info *model.Process) error {
	return nil
}
