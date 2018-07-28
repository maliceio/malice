/*
Package database provides libraries to write plugin results to a database.
*/
package database

// PluginResults a malice plugin results object
type PluginResults struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	Category string                 `json:"category,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// Database is a Malice Database interface
type Database interface {
	Init() error
	TestConnection() error
	StoreFileInfo(sample map[string]interface{}) error
	StoreHash(hash string) error
	StorePluginResults(results PluginResults) error
}
