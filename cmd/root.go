package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// NewCmdRoot creates root command
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tss",
		Short: "Annotate stdin with timestamps per line",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := Copy(os.Stdout, os.Stdin); err != nil {
				cmd.PrintErrf("%v", err)
			}
		},
	}

	cmd.AddCommand(NewCmdVersion())

	return cmd
}

// Execute runs root command
func Execute() {
	cmd := NewCmdRoot()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}

// Writer represents output of tss
type Writer struct {
	w         io.Writer
	start     time.Time
	lastLine  time.Time
	buf       bytes.Buffer
	endOfLine bool
}

var oldestStart time.Time

func init() {
	// This is as early as we can get into measuring when the binary starts.
	oldestStart = time.Now()
}

// NewWriter creates a new Writer. If start is the zero value, the program's
// start time will be used.
func NewWriter(w io.Writer, start time.Time) *Writer {
	if start.IsZero() {
		start = oldestStart
	}
	return &Writer{w: w, start: start, endOfLine: true}
}

var padding = bytes.Repeat([]byte{' '}, 9)

// Write writes the contents of p into the buffer. It returns the number of
// bytes written. If nn < len(p), it also returns an error explaining why the
// write is short.
func (w *Writer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	w.buf.Reset()
	w.buf.Grow(len(p) + bytes.Count(p, []byte{'\n'})*len(padding))
	wrote := 0
	now := time.Now()
	pos := 0
	for {
		// write everything up to the next newline
		if w.endOfLine {
			// print timing info
			sinceStart := now.Sub(w.start).Round(100 * time.Microsecond)
			s := TimeScaler(sinceStart)
			for i := 0; i < 8-len(s); i++ {
				w.buf.WriteByte(' ')
			}
			w.buf.WriteString(s)
			w.buf.WriteByte(' ')
			if w.lastLine.IsZero() {
				for i := 0; i < 8-len(s); i++ {
					w.buf.WriteByte(' ')
				}
				w.buf.WriteString(s)
				w.buf.WriteByte(' ')
				w.lastLine = now
			} else {
				sinceLastLine := now.Sub(w.lastLine).Round(100 * time.Microsecond)
				s := TimeScaler(sinceLastLine)
				for i := 0; i < 8-len(s); i++ {
					w.buf.WriteByte(' ')
				}
				w.buf.WriteString(s)
				w.buf.WriteByte(' ')
				w.lastLine = now
			}
			w.endOfLine = false
		}
		idx := bytes.IndexByte(p[pos:], '\n')
		if idx >= 0 {
			w.buf.Write(p[pos : pos+idx+1])
			wrote += idx + 1
			pos += idx + 1
			w.endOfLine = true
			if pos >= len(p) {
				break
			}
		} else {
			w.buf.Write(p[pos:])
			wrote += len(p) - pos
			break
		}
	}
	_, err := w.w.Write(w.buf.Bytes())
	return wrote, err
}

var forceNonZeroTestVal = time.Duration(0)

// TimeScaler returns a format string for the given Duration where all of the
// decimals will line up in the same column (fourth from the end).
func TimeScaler(d time.Duration) string {
	if d == 0 && forceNonZeroTestVal != 0 {
		d = forceNonZeroTestVal
	}
	switch {
	case d == 0:
		return "0.0ms"
	case d >= time.Minute:
		mins := (d + 30*time.Second) / time.Minute
		d = d - mins*time.Minute
		if d < 0 {
			d = 0 // can happen due to rounding
		}
		s := strconv.FormatFloat(float64(d.Nanoseconds())/1e9, 'f', 0, 64)
		return strconv.Itoa(int(mins)) + "m" + fmt.Sprintf("%02s", s) + "s"
	case d >= time.Second:
		return strconv.FormatFloat(float64(d.Nanoseconds())/1e9, 'f', 2, 64) + "s"
	case d >= 50*time.Microsecond:
		return strconv.FormatFloat(float64(d.Nanoseconds())/1e9*1000, 'f', 1, 64) + "ms"
	case d >= time.Microsecond:
		return strconv.FormatFloat(float64(d.Nanoseconds())/1e9*1000*1000, 'f', 1, 64) + "µs"
	default:
		return strconv.FormatFloat(float64(d.Nanoseconds()), 'f', 1, 64) + "ns"
	}
}

// Copy returns CopyTime
func Copy(w io.Writer, r io.Reader) (written int64, err error) {
	return CopyTime(w, r, time.Now())
}

// CopyTime copies from Reader to Writer until either EOF is reached on src or an error occurs.
func CopyTime(w io.Writer, r io.Reader, start time.Time) (written int64, err error) {
	return io.Copy(NewWriter(w, start), r)
}
