// Copyright © 2017 Igor Bondarenko <ibondare@protonmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package context

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Context
type Context struct {
	Mutex         *sync.Mutex
	HttpService   string
	Out           io.Writer
	LogJSON       bool
	LogPrettyJSON bool
	Echo          bool
	HttpCode      int
	Delay         int
	FailureMode   failureSimulation
}

// FailureSimulation desribes the intended behavior of the transient failure mode in httpr
type failureSimulation struct {
	Enabled               bool
	FailureCount          int
	SuccessCount          int
	FailureCode           int
	SuccessCode           int
	FailureIterationCount int
	SuccessIterationCount int
}

var singleton *Context

var once sync.Once

func Instance() *Context {
	once.Do(func() {
		singleton = &Context{Mutex: &sync.Mutex{}, FailureMode: failureSimulation{Enabled: false}}
	})
	return singleton
}

// Start the HTTP server
func (ctx *Context) StartServer() {
	go log.Fatal(http.ListenAndServe(ctx.HttpService, nil))

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

// Execute a failure simulation and return an HTTP code representing the outcome
func (ctx *Context) SimulateFailure() int {
	var outcome int

	ctx.Mutex.Lock()

	defer ctx.Mutex.Unlock()

	if ctx.FailureMode.FailureIterationCount < ctx.FailureMode.FailureCount {
		outcome = ctx.FailureMode.FailureCode

		ctx.FailureMode.FailureIterationCount++

		if ctx.FailureMode.FailureIterationCount == ctx.FailureMode.FailureCount {
			// Done with the failure sequence, next call will return success
			ctx.FailureMode.SuccessIterationCount = 0
		}

	} else if ctx.FailureMode.SuccessIterationCount < ctx.FailureMode.SuccessCount {
		outcome = ctx.FailureMode.SuccessCode

		ctx.FailureMode.SuccessIterationCount++

		if ctx.FailureMode.SuccessIterationCount == ctx.FailureMode.SuccessCount {
			// Done with the failure sequence, next call will return success
			ctx.FailureMode.FailureIterationCount = 0
		}

	}

	return outcome
}
