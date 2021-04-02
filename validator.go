package log

import "fmt"

// validateCfg checks whether given configurations are valid.
func validateCfg(cfg Config) error {

	_, ok := granularity[cfg.Level]
	if !ok {
		return fmt.Errorf("Unknown log level '%s'", cfg.Level)
	}

	return nil
}
