package cmd

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"
)

func BenchmarkCopy(b *testing.B) {
	forceNonZeroTestVal = 5 * time.Millisecond
	defer func() {
		forceNonZeroTestVal = 0
	}()
	bs := bytes.Repeat([]byte{'a'}, 2<<12+1)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	lr := &smallReader{R: rd, N: 57}
	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := CopyTime(buf, lr, time.Now().Add(-50*time.Millisecond))
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
		rd.Reset(bs)
	}
}

type smallReader struct {
	R *bytes.Reader
	N int
}

func (s *smallReader) Read(p []byte) (int, error) {
	if len(p) >= s.N {
		return s.R.Read(p[:s.N])
	}
	return s.R.Read(p)
}

func BenchmarkWriter(b *testing.B) {
	forceNonZeroTestVal = 5 * time.Millisecond
	defer func() {
		forceNonZeroTestVal = 0
	}()
	bs := bytes.Repeat([]byte{'a'}, 2<<12)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	lr := &smallReader{R: rd, N: 57}
	buf := new(bytes.Buffer)
	w := NewWriter(buf, time.Now().Add(-50*time.Millisecond))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := io.Copy(w, lr)
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
		rd.Reset(bs)
	}
}

func BenchmarkWriterBig(b *testing.B) {
	forceNonZeroTestVal = 5 * time.Millisecond
	defer func() {
		forceNonZeroTestVal = 0
	}()
	bs := bytes.Repeat([]byte{'a'}, 2<<12)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	buf := new(bytes.Buffer)
	w := NewWriter(buf, time.Now().Add(-50*time.Millisecond))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := io.Copy(w, rd)
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
		rd.Reset(bs)
	}
}
