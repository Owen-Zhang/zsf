package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	s           Severity
	backend     ILog
	mu          sync.Mutex
	freeList    *buffer
	freeListMu  sync.Mutex
	logToStderr bool
}

//NewLogger new logger
func NewLogger(level interface{}, log ILog) *Logger {
	l := new(Logger)
	l.SetSeverity(level)
	l.backend = log
	return l
}

func (l *Logger) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

func (l *Logger) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// Let big buffers die a natural death.
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

func (l *Logger) formatHeader(s Severity, file string, line int) *buffer {
	now := time.Now()
	if line < 0 {
		line = 0 // not a real line number, but acceptable to someDigits
	}
	buf := l.getBuffer()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	buf.nDigits(4, 0, year, '0')
	buf.tmp[4] = '-'
	buf.twoDigits(5, int(month))
	buf.tmp[7] = '-'
	buf.twoDigits(8, day)
	buf.tmp[10] = ' '
	buf.twoDigits(11, hour)
	buf.tmp[13] = ':'
	buf.twoDigits(14, minute)
	buf.tmp[16] = ':'
	buf.twoDigits(17, second)
	buf.tmp[19] = '.'
	buf.nDigits(6, 20, now.Nanosecond()/1000, '0')
	buf.tmp[26] = ' '
	buf.Write(buf.tmp[:27])
	buf.WriteString(severityName[s])
	buf.WriteByte(' ')
	buf.WriteString(file)
	buf.tmp[0] = ':'
	n := buf.someDigits(1, line)
	buf.tmp[n+1] = ' '
	buf.Write(buf.tmp[:n+2])
	return buf
}

func (l *Logger) header(s Severity, depth int) *buffer {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		dirs := strings.Split(file, "/")
		if len(dirs) >= 2 {
			file = dirs[len(dirs)-2] + "/" + dirs[len(dirs)-1]
		} else {
			file = dirs[len(dirs)-1]
		}
	}
	return l.formatHeader(s, file, line)
}

func (l *Logger) print(s Severity, args ...interface{}) {
	l.printDepth(s, 1, args...)
}

func (l *Logger) printf(s Severity, format string, args ...interface{}) {
	l.printfDepth(s, 1, format, args...)
}

func (l *Logger) printDepth(s Severity, depth int, args ...interface{}) {
	if l.s < s {
		return
	}
	buf := l.header(s, depth)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf)
}

func (l *Logger) printfDepth(s Severity, depth int, format string, args ...interface{}) {
	if l.s < s {
		return
	}
	buf := l.header(s, depth)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf)
}

func (l *Logger) printfSimple(format string, args ...interface{}) {
	buf := l.getBuffer()
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(INFO, buf)
}

//output 将生成的日志信息转交给要保存或者显示的对象
func (l *Logger) output(s Severity, buf *buffer) {
	if l.s < s {
		return
	}
	if l.logToStderr {
		os.Stderr.Write(buf.Bytes())
	} else {
		l.backend.Log(s, buf.Bytes())
	}
	if s == FATAL {
		trace := stacks(true)
		os.Stderr.Write(trace)
		os.Exit(255)
	}
	l.putBuffer(buf)
}

//SetSeverity 设置输出级别(如:INFO)
func (l *Logger) SetSeverity(level interface{}) {
	if s, ok := level.(Severity); ok {
		l.s = s
	} else {
		if s, ok := level.(string); ok {
			for i, name := range severityName {
				if name == s {
					l.s = Severity(i)
				}
			}
		}
	}
}

//Close 关闭日志
func (l *Logger) Close() {
	if l.backend != nil {
		l.backend.close()
	}
}

func (l *Logger) LogToStderr() {
	l.logToStderr = true
}

func (l *Logger) Debug(args ...interface{}) {
	l.print(DEBUG, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.printf(DEBUG, format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.print(INFO, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.printf(INFO, format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.print(WARNING, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.printf(WARNING, format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.print(ERROR, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.printf(ERROR, format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.print(FATAL, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.printf(FATAL, format, args...)
}

func (l *Logger) SetLogging(level interface{}, log ILog) {
	l.SetSeverity(level)
	l.backend = log
}

func (l *Logger) LogDepth(s Severity, depth int, format string, args ...interface{}) {
	l.printfDepth(s, depth+1, format, args...)
}

func (l *Logger) PrintfSimple(format string, args ...interface{}) {
	l.printfSimple(format, args...)
}

func stacks(all bool) []byte {
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

//---------------------------------------------------------
//---------------------------------------------------------
type buffer struct {
	bytes.Buffer
	tmp  [64]byte
	next *buffer
}

const digits = "0123456789"

//twoDigits formats a zero-prefixed two-digit integer at buf.tmp[i]
func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

// nDigits formats an n-digit integer at buf.tmp[i],
// padding with pad on the left.
// It assumes d >= 0.
func (buf *buffer) nDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.tmp[i+j] = pad
	}
}

// someDigits formats a zero-prefixed variable-width integer at buf.tmp[i].
func (buf *buffer) someDigits(i, d int) int {
	// Print into the top, then copy down. We know there's space for at least
	// a 10-digit number.
	j := len(buf.tmp)
	for {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}
