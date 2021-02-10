package log

// Config holds application log configurations.
type Config struct {
	Level     string `yaml:"level"`
	Colors    bool   `yaml:"colors"`
	Console   bool   `yaml:"console"`
	File      bool   `yaml:"file"`
	Directory string `yaml:"directory"`
}
