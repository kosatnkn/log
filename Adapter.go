package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	logContext "github.com/kosatnkn/log/context"
	"github.com/logrusorgru/aurora"
)

// Adapter is used to provide structured log messages.
type Adapter struct {
	cfg Config
	lf  *os.File
}

// NewAdapter creates a new Log adapter instance.
func NewAdapter(cfg Config) (AdapterInterface, error) {

	a := &Adapter{
		cfg: cfg,
	}

	err := a.initLogFile()
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Error logs a message as of error type.
func (a *Adapter) Error(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, "ERROR", message, options)
}

// Debug logs a message as of debug type.
func (a *Adapter) Debug(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, "DEBUG", message)
}

// Info logs a message as of information type.
func (a *Adapter) Info(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, "INFO", message, options)
}

// Warn logs a message as of warning type.
func (a *Adapter) Warn(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, "WARN", message, options)
}

// Destruct will close the logger gracefully releasing all resources.
func (a *Adapter) Destruct() {

	if a.cfg.File {
		a.lf.Close()
	}
}

// Initialize the log file.
func (a *Adapter) initLogFile() error {

	if !a.cfg.File {
		return nil
	}

	ld := a.cfg.Directory

	f, err := os.OpenFile(filepath.Join(ld, "out.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	a.lf = f

	return nil
}

// Logs a message using the following format.
// <date> <time_in_24h_foramt_plus_milliseconds> [<message_type>] [<uuid>] [<prefix>] [<message>] [<additional_information>]
// ex:
//		2019/01/14 12:13:29.435517 [ERROR] [b2e1bfc7-11ed-40e5-ab08-abeadef079e6] [usecases.TestUsecase.TestFunc] [error message] [key1: value1, ...]
func (a *Adapter) log(ctx context.Context, logLevel string, message string, options ...interface{}) {

	// check whether the message should be logged
	if !a.isLoggable(logLevel) {
		return
	}

	m := a.formatMessage(ctx, logLevel, message, options)

	a.logToConsole(m)
	a.logToFile(m)
}

// formatMessage formats the log message.
func (a *Adapter) formatMessage(ctx context.Context, logLevel string, message string, options ...interface{}) string {

	now := time.Now().Format("2006/01/02 15:04:05.000000")
	uuid := ctx.Value(logContext.UUIDKey)
	prefix := ctx.Value(logContext.TraceKey)
	level := a.setTag(logLevel)

	return fmt.Sprintf("%s %s [%s] [%v] [%v] [%v]", now, level, uuid, prefix, message, options)
}

// Check whether the message should be logged depending on the log level setting.
func (a *Adapter) isLoggable(logLevel string) bool {

	l := map[string]int{
		"ERROR": 1,
		"DEBUG": 2,
		"WARN":  3,
		"INFO":  4,
	}

	return l[logLevel] >= l[a.cfg.Level]
}

// Generate tags based on color configuration settings.
func (a *Adapter) setTag(logLevel string) interface{} {

	if a.cfg.Colors {
		var logLevelVal aurora.Value

		switch logLevel {
		case "ERROR":
			logLevelVal = aurora.Red("[ERROR]")
			break
		case "DEBUG":
			logLevelVal = aurora.Green("[DEBUG]")
			break
		case "INFO":
			logLevelVal = aurora.Cyan("[INFO]")
			break
		case "WARN":
			logLevelVal = aurora.Brown("[WARN]")
			break
		default:
			logLevelVal = aurora.Magenta("[UNKNOWN]")
			break
		}

		return logLevelVal
	}

	return "[" + logLevel + "]"
}

// Logs a message to the console.
func (a *Adapter) logToConsole(message string) {

	if a.cfg.Console {
		fmt.Println(message)
	}
}

// Logs a message to a file.
func (a *Adapter) logToFile(message string) {

	if !a.cfg.File {
		return
	}

	_, err := a.lf.WriteString(message + "\n")
	if err != nil {
		fmt.Println(err)
	}
}