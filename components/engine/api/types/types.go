package types

// Version contains response of Engine API:
// GET "/version"
type Version struct {
	Version       string
	APIVersion    string `json:"ApiVersion"`
	MinAPIVersion string `json:"MinAPIVersion,omitempty"`
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
	Experimental  bool   `json:",omitempty"`
	BuildTime     string `json:",omitempty"`
}

// ErrorResponse Represents an error.
type ErrorResponse struct {

	// The error message.
	// Required: true
	Message string `json:"message"`
}

// Ping contains response of Engine API:
// GET "/_ping"
type Ping struct {
	APIVersion string
	OSType     string
}
