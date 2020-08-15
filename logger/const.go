package logger

type Severity int

const (
	numSeverity = 5
	bufferSize  = 256 * 1024
)

const (
	FATAL Severity = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

var severityName = []string{
	FATAL:   "FATAL",
	ERROR:   "ERROR",
	WARNING: "WARNING",
	INFO:    "INFO",
	DEBUG:   "DEBUG",
}
