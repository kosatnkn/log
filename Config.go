package log

// Config holds application log configurations.
type Config struct {
	// Level defines the log level that determines whether to log a message or not.
	//
	// There are four log levels to choose from with associated granularity.
	//  ERROR: log messages with log level 'ERROR'
	//  WARN: log messages with log level 'ERROR' and 'WARN'
	//  DEBUG: log messages with log level 'ERROR', 'WARN' and 'DEBUG'
	//  INFO: log messages with log level 'ERROR', 'WARN', 'DEBUG' and 'INFO'
	Level string `yaml:"level"`

	// Colors determine whether logs are printed in colours or not.
	Colors bool `yaml:"colors"`

	// Console determine whether logs are printed to the console or not.
	Console bool `yaml:"console"`

	// File determine whether logs are printed to a file or not.
	File bool `yaml:"file"`

	// Directory specifies the directory to which the logs are written
	// if the `File` option is true.
	Directory string `yaml:"directory"`
}
