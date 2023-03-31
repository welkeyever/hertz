/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/pedia/endless"
)

func serve_http(addr string, ln net.Listener) *server.Hertz {
	h := server.New()
	if ext, err := h.TransporterExt(); err == nil {
		ext.SetListener(ln)
	}
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(200, utils.H{"message": "pong"})
	})

	go h.Spin()

	return h
}

func main() {
	addr := ":3030"
	var s *server.Hertz

	endless.Start(
		func(p *endless.Parent) error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			p.AddListener(ln, addr)
			s = serve_http(addr, ln)
			return err
		}, func(c *endless.Child) error {
			nf, ok := c.NamedFiles[addr]
			fmt.Printf("got %#v\n", nf)
			if !ok {
				return fmt.Errorf("inherit %s not found", addr)
			}

			c.AddListener(nf.Listener, addr)
			s = serve_http(addr, nf.Listener)
			return nil
		},
		func(ctx context.Context) error {
			return s.Shutdown(context.Background())
		},
	)
	fmt.Printf("Quit %d\n", os.Getpid())
}
