package contex

// Context key type to be used with contexts.
type ctxKey string

// UUIDKey is the universally unique identifier key to be used with context.
const UUIDKey ctxKey = "uuid"

// TraceKey is the key to add an additional trace values to the context.
const TraceKey ctxKey = "trace"
