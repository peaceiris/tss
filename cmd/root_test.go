package cmd_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/peaceiris/tss/cmd"
	"github.com/stretchr/testify/assert"
)

type sleepReader struct {
	count    int
	max      int
	sleepFor time.Duration
}

func (s *sleepReader) Read(p []byte) (int, error) {
	s.count++
	if s.count > s.max {
		return 0, io.EOF
	}
	if s.count == 1 {
		copy(p[:6], "hello\n")
		return 6, nil
	}
	time.Sleep(s.sleepFor)
	if s.count == 2 {
		copy(p[:3], "hel")
		return 3, nil
	}
	if s.count == 3 {
		copy(p[:3], "lo\n")
		return 3, nil
	}
	if s.count == 4 {
		copy(p[:15], "hello\nhello\nhel")
		return 15, nil
	}
	if s.count == 5 {
		copy(p[:3], "lo\n")
		return 3, nil
	}
	copy(p[:6], "hello\n")
	return 6, nil
}

func TestWriter(t *testing.T) {
	t.Parallel()
	max := 6
	s := &sleepReader{max: max, sleepFor: 2 * time.Millisecond}
	buf := new(bytes.Buffer)
	w := cmd.NewWriter(buf, time.Time{})
	n, err := io.Copy(w, s)
	assert.Nil(t, err)
	assert.Equal(t, len("hello\n")*max, int(n), "n is 36")
}

func TestCopy(t *testing.T) {
	t.Parallel()

	max := 6
	s := &sleepReader{max: max, sleepFor: 2 * time.Millisecond}
	buf := new(bytes.Buffer)
	n, err := cmd.Copy(buf, s)
	assert.Nil(t, err)
	want := len("hello\n") * 6
	assert.Equal(t, want, int(n), "expected line length n")

	parts := strings.Split(buf.String(), "\n")
	assert.Equal(t, 7, len(parts), "incorrect number of parts")

	line1 := parts[0]
	assert.Equal(t, 23, len(line1), "line1 length is 23")

	lineParts := strings.Fields(line1)
	assert.Equal(t, 3, len(lineParts), "first line has 3 parts")

	part, err := time.ParseDuration(lineParts[0])
	assert.Nil(t, err)
	assert.Truef(t, part <= 100*time.Millisecond, "part <= 100*time.Millisecond: %s <= %s", part, 100*time.Millisecond)

	lineParts = strings.Fields(parts[1])
	assert.Equal(t, 3, len(lineParts), "wrong number of line parts in line 2")
}

func TestTimeScaler(t *testing.T) {
	scalerTests := []struct {
		in  time.Duration
		out string
	}{
		{100 * time.Microsecond, "0.1ms"},
		{500 * time.Microsecond, "0.5ms"},
		{99 * time.Microsecond, "0.1ms"},
		{49 * time.Microsecond, "49.0µs"},
		{time.Millisecond, "1.0ms"},
		{56*time.Millisecond + 290*time.Microsecond, "56.3ms"},
		{56*time.Millisecond + 251*time.Microsecond, "56.3ms"},
		{56*time.Millisecond + 100*time.Microsecond, "56.1ms"},
		{3*time.Minute + 4*time.Second + 100*time.Millisecond, "3m04s"},
		{3*time.Minute + 14*time.Second + 100*time.Millisecond, "3m14s"},
		{14*time.Second + 100*time.Millisecond, "14.10s"},
		{2*time.Minute - 1*time.Microsecond, "2m00s"},
		{0, "0.0ms"},
	}

	for _, tt := range scalerTests {
		v := cmd.TimeScaler(tt.in)
		assert.Equal(t, tt.out, v, "timeScaler")
	}
}
