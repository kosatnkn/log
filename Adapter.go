package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/kosatnkn/log/internal"
)

// Adapter is used to provide structured log messages.
type Adapter struct {
	cfg Config
	lf  *os.File
}

// NewAdapter creates a new Log adapter instance.
func NewAdapter(cfg Config) (AdapterInterface, error) {

	err := validateCfg(cfg)
	if err != nil {
		return nil, err
	}

	a := &Adapter{
		cfg: cfg,
	}

	err = a.initLogFile()
	if err != nil {
		return nil, err
	}

	return a, nil
}

// AddTraceID attaches a trace id to context that can be later read by the logger.
func (a *Adapter) AddTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, internal.ID, id)
}

// AppendTracePoint appends the given trace point to a trace path in context that can be later read by the logger.
func (a *Adapter) AppendTracePoint(ctx context.Context, point string) context.Context {

	path := ctx.Value(internal.TraceKey)

	if path == nil {
		return context.WithValue(ctx, internal.TraceKey, point)
	}

	return context.WithValue(ctx, internal.TraceKey, path.(string)+">"+point)
}

// Error logs a message as of error type.
func (a *Adapter) Error(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, levelError, message, options...)
}

// Debug logs a message as of debug type.
func (a *Adapter) Debug(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, levelDebug, message, options...)
}

// Info logs a message as of information type.
func (a *Adapter) Info(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, levelInfo, message, options...)
}

// Warn logs a message as of warning type.
func (a *Adapter) Warn(ctx context.Context, message string, options ...interface{}) {
	a.log(ctx, levelWarn, message, options...)
}

// Destruct will close the logger gracefully releasing all resources.
func (a *Adapter) Destruct() {

	if a.cfg.File {
		a.lf.Close()
	}
}

// initLogFile initialize the log file.
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

// log logs a message using the following format.
// <date> <time_in_24h_foramt_plus_milliseconds> [<log_level>] [<uuid>] [<trace_points>] [<message>] : [<additional_information>]
// ex:
//		2019/01/14 12:13:29.435517 [ERROR] [b2e1bfc7-11ed-40e5-ab08-abeadef079e6] [usecases.TestUsecase.TestFunc] [error message] : [key1: value1, ...]
func (a *Adapter) log(ctx context.Context, logLevel string, message string, options ...interface{}) {

	// check whether the message should be logged
	if !a.isLoggable(logLevel) {
		return
	}

	m := a.formatMessage(ctx, logLevel, message, options...)

	a.toConsole(m)
	a.toFile(m)
}

// formatMessage formats the log message.
func (a *Adapter) formatMessage(ctx context.Context, logLevel string, message string, options ...interface{}) string {

	var now = time.Now().Format("2006/01/02 15:04:05.000000")
	var level = a.setTag(logLevel)
	var uuid = "NONE"
	var trace = "NONE"

	id, ok := ctx.Value(internal.ID).(string)
	if ok {
		uuid = id
	}

	points, ok := ctx.Value(internal.TraceKey).(string)
	if ok {
		trace = points
	}

	if len(options) == 0 {
		return fmt.Sprintf("%s %s [%s] [%s] [%s]", now, level, uuid, trace, message)
	}

	return fmt.Sprintf("%s %s [%s] [%s] [%s] : %v", now, level, uuid, trace, message, options)
}

// Check whether the message should be logged depending on the granularity of the log level.
func (a *Adapter) isLoggable(level string) bool {
	return granularity[level] >= granularity[a.cfg.Level]
}

// Generate tags based on color configuration settings.
func (a *Adapter) setTag(level string) interface{} {

	if a.cfg.Colors {

		switch level {
		case levelError:
			return color.New(color.FgRed).Sprint("[ERROR]")
		case levelDebug:
			return color.Debug.Sprint("[DEBUG]")
		case levelInfo:
			return color.Info.Sprint("[INFO]")
		case levelWarn:
			return color.New(color.FgYellow).Sprint("[WARN]")
		default:
			return "[" + level + "]"
		}
	}

	return "[" + level + "]"
}

// toConsole logs a message to the console.
func (a *Adapter) toConsole(message string) {

	if a.cfg.Console {
		fmt.Println(message)
	}
}

// toFile logs a message to a file.
func (a *Adapter) toFile(message string) {

	if !a.cfg.File {
		return
	}

	_, err := a.lf.WriteString(message + "\n")
	if err != nil {
		fmt.Println(err)
	}
}
