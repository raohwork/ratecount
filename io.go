// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ratecount

import "io"

// NewWriter creates a RatedWriter
func NewWriter(w io.Writer, c *Counter) *RatedWriter {
	return &RatedWriter{
		w: w,
		c: c,
	}
}

// RatedWriter computes transfer rate when you write to it
type RatedWriter struct {
	w io.Writer
	c *Counter
}

// Write implements io.Writer.Write
func (w *RatedWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	if n > 0 {
		w.c.Incr(int64(n))
	}

	return
}

// Rate wraps Counter.Rate()
func (w *RatedWriter) Rate() int64 {
	return w.c.Rate()
}

// In wraps Counter.In()
func (w *RatedWriter) In(i int64) int64 {
	return w.c.In(i)
}

// NewReader creates a RatedReader
func NewReader(r io.Reader, c *Counter) *RatedReader {
	return &RatedReader{
		r: r,
		c: c,
	}
}

// RatedReader computes transfer rate when you read from it
type RatedReader struct {
	r io.Reader
	c *Counter
}

// Read implements io.Reader.Read()
func (r *RatedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		r.c.Incr(int64(n))
	}

	return
}

// Rate wraps Counter.Rate
func (r *RatedReader) Rate() int64 {
	return r.c.Rate()
}

// In wraps Counter.In
func (r *RatedReader) In(i int64) int64 {
	return r.c.In(i)
}
