/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"time"

	errs "github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/netpoll"
)

type Conn struct {
	readTimeout time.Duration
	zr          network.Reader
	zw          network.ReadWriter
}

type SlowReadConn struct {
	*Conn
}

func SlowReadDialer(addr string) (network.Conn, error) {
	return NewSlowReadConn(""), nil
}

func (m *Conn) ReadBinary(n int) (p []byte, err error) {
	return m.zr.(netpoll.Reader).ReadBinary(n)
}

func (m *Conn) Read(b []byte) (n int, err error) {
	return netpoll.NewIOReader(m.zr.(netpoll.Reader)).Read(b)
}

func (m *Conn) Write(b []byte) (n int, err error) {
	return netpoll.NewIOWriter(m.zw.(netpoll.ReadWriter)).Write(b)
}

func (m *Conn) Release() error {
	return nil
}

func (m *Conn) Peek(i int) ([]byte, error) {
	b, err := m.zr.Peek(i)
	if err != nil || len(b) != i {
		if m.readTimeout <= 0 {
			// simulate timeout forever
			select {}
		}
		time.Sleep(m.readTimeout)
		return nil, errs.ErrTimeout
	}
	return b, err
}

func (m *Conn) Skip(n int) error {
	return m.zr.Skip(n)
}

func (m *Conn) ReadByte() (byte, error) {
	return m.zr.ReadByte()
}

func (m *Conn) Len() int {
	return m.zr.Len()
}

func (m *Conn) Malloc(n int) (buf []byte, err error) {
	return m.zw.Malloc(n)
}

func (m *Conn) WriteBinary(b []byte) (n int, err error) {
	return m.zw.WriteBinary(b)
}

func (m *Conn) Flush() error {
	return m.zw.Flush()
}

func (m *Conn) WriterRecorder() network.Reader {
	return m.zw
}

func (m *SlowReadConn) Peek(i int) ([]byte, error) {
	b, err := m.zr.Peek(i)
	time.Sleep(100 * time.Millisecond)
	if err != nil || len(b) != i {
		time.Sleep(m.readTimeout)
		return nil, errs.ErrTimeout
	}
	return b, err
}

func NewConn(source string) *Conn {
	zr := netpoll.NewReader(strings.NewReader(source))
	zw := netpoll.NewReadWriter(&bytes.Buffer{})

	return &Conn{
		zr: zr,
		zw: zw,
	}
}

func NewSlowReadConn(source string) *SlowReadConn {
	return &SlowReadConn{NewConn(source)}
}

func (m *Conn) Close() error {
	return nil
}

func (m *Conn) LocalAddr() net.Addr {
	return nil
}

func (m *Conn) RemoteAddr() net.Addr {
	return nil
}

func (m *Conn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (m *Conn) SetReadDeadline(t time.Time) error {
	m.readTimeout = -time.Since(t)
	return nil
}

func (m *Conn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}

func (m *Conn) Reader() network.Reader {
	return m.zr
}

func (m *Conn) Writer() network.Writer {
	return m.zw
}

func (m *Conn) IsActive() bool {
	panic("implement me")
}

func (m *Conn) SetIdleTimeout(timeout time.Duration) error {
	return nil
}

func (m *Conn) SetReadTimeout(t time.Duration) error {
	m.readTimeout = t
	return nil
}

func (m *Conn) SetOnRequest(on netpoll.OnRequest) error {
	panic("implement me")
}

func (m *Conn) AddCloseCallback(callback netpoll.CloseCallback) error {
	panic("implement me")
}

type StreamConn struct {
	Data []byte
}

func NewStreamConn() *StreamConn {
	return &StreamConn{
		Data: make([]byte, 1<<15, 1<<16),
	}
}

func (m *StreamConn) Peek(n int) ([]byte, error) {
	if len(m.Data) >= n {
		return m.Data[:n], nil
	}
	if n == 1 {
		m.Data = m.Data[:cap(m.Data)]
		return m.Data[:1], nil
	}
	return nil, errors.New("not enough data")
}

func (m *StreamConn) Skip(n int) error {
	if len(m.Data) >= n {
		m.Data = m.Data[n:]
		return nil
	}
	return errors.New("not enough data")
}

func (m *StreamConn) Release() error {
	panic("implement me")
}

func (m *StreamConn) Len() int {
	return len(m.Data)
}

func (m *StreamConn) ReadByte() (byte, error) {
	panic("implement me")
}

func (m *StreamConn) ReadBinary(n int) (p []byte, err error) {
	panic("implement me")
}

func DialerFun(addr string) (network.Conn, error) {
	return NewConn(""), nil
}
