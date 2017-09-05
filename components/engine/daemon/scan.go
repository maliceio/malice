package daemon

import (
	"context"

	"github.com/maliceio/engine/api/types/scan"
)

// Scan performs a malice scan on a sample
func (daemon *Daemon) Scan(ctx context.Context, path string, config *scan.Config) (scan.Result, error) {
	result := scan.Result{
		Out: "you ran a scan",
	}
	return result, nil
}
