// Code generated by nixmod2go. DO NOT EDIT.
package e2clickermodule

// BackendConfig is the struct type for `config`.
type BackendConfig struct {
	// Enable: whether to enable e2clicker backend.
	Enable  bool   `json:"enable"`
	Package string `json:"package"`
	// API: configuration for the API server.
	API API `json:"api"`
	// Debug: enable debug logging and other debug features.
	Debug bool `json:"debug"`
	// LogFormat: format of the log output.
	LogFormat LogFormat `json:"logFormat"`
	// PostgreSQL: configuration for the PostgreSQL database.
	PostgreSQL PostgreSQL `json:"postgresql"`
}

// PostgreSQL is the struct type for `config.postgresql`.
type PostgreSQL struct {
	// DatabaseURI: URI of the database to use.
	DatabaseURI string `json:"databaseURI"`
}

// API is the struct type for `config.api`.
type API struct {
	// ListenAddress address the API server should listen on.
	ListenAddress string `json:"listenAddress"`
}

// LogFormat is the enum type for `config.logFormat`.
type LogFormat string

const (
	LogFormatColor LogFormat = "color"
	LogFormatJSON  LogFormat = "json"
	LogFormatText  LogFormat = "text"
)
