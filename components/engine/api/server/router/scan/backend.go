package scan

import "github.com/maliceio/engine/api/types/scan"

// Backend for Plugin
type Backend interface {
	Scan(path string, config *scan.Config) (scan.Result, error)
}
