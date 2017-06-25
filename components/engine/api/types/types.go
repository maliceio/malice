package types

// Scan returns malice scan
type Scan struct {
	ID string `json:"id,omitempty"`
}

// Ping contains response of Engine API:
// GET "/_ping"
type Ping struct {
	APIVersion   string
	OSType       string
	Experimental bool
}
