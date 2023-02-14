package resp

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/http1/ext"
)

type chunkedBodyWriter struct {
	sync.Once
	finalizeErr error
	wroteHeader bool
	r           *protocol.Response
	w           network.Writer
}

// Write will encode chunked p before writing
// It will only return the length of p and a nil error if the writing is successful or 0, error otherwise.
//
// NOTE: Write will use the user buffer to flush.
// Before flush successfully, the buffer b should be valid.
func (c *chunkedBodyWriter) Write(p []byte) (n int, err error) {
	if !c.wroteHeader {
		c.r.Header.SetContentLength(-1)
		if err = WriteHeader(&c.r.Header, c.w); err != nil {
			return
		}
		c.wroteHeader = true
	}
	if err = ext.WriteChunk(c.w, p, false); err != nil {
		return
	}
	return len(p), nil
}

func (c *chunkedBodyWriter) Flush() error {
	return c.w.Flush()
}

// Finalize will write the ending chunk as well as trailer and flush the writer.
// Warning: do not call this method by yourself, unless you know what you are doing.
func (c *chunkedBodyWriter) Finalize() error {
	c.Do(func() {
		c.finalizeErr = ext.WriteChunk(c.w, nil, true)
		if c.finalizeErr != nil {
			return
		}
		c.finalizeErr = ext.WriteTrailer(c.r.Header.Trailer(), c.w)
	})
	return c.finalizeErr
}

func NewChunkedBodyWriter(r *protocol.Response, w network.Writer) network.ExtWriter {
	return &chunkedBodyWriter{
		r: r,
		w: w,
	}
}