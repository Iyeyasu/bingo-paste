package log

// Level is the logging level used for logs.
type Level int

const (
	// PanicLevel logs only panics.
	PanicLevel Level = iota

	// FatalLevel logs fatals and more serious incidents.
	FatalLevel

	// ErrorLevel logs errors and more serious incidents.
	ErrorLevel

	// WarnLevel logs warnings and more serious incidents.
	WarnLevel

	// InfoLevel logs infos and more serious incidents.
	InfoLevel

	// DebugLevel logs debugs and more serious incidents.
	DebugLevel

	// TraceLevel logs traces and more serious incidents.
	TraceLevel
)
