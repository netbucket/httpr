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
	"github.com/netbucket/httpr/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Context type holds the desired execution profile for a command
type Context struct {
	Mutex           *sync.Mutex
	HttpService     string
	EnableTLS       bool
	CertFile        string
	KeyFile         string
	UpstreamURL     *url.URL
	Out             io.Writer
	LogJSON         bool
	LogPrettyJSON   bool
	Echo            bool
	HttpCode        int
	Delay           int
	IgnoreTLSErrors bool
	FailureMode     FailureSimulation
}

// FailureSimulation desribes the intended behavior of the transient failure mode in httpr
type FailureSimulation struct {
	Enabled               bool
	FailureCount          int
	SuccessCount          int
	FailureCode           int
	failureIterationCount int
	successIterationCount int
	failureSimulated      bool
}

var singleton *Context

var once sync.Once

func Instance() *Context {
	once.Do(func() {
		singleton = &Context{Mutex: &sync.Mutex{}, Out: os.Stdout, FailureMode: FailureSimulation{Enabled: false}}
	})
	return singleton
}

// Start the HTTP server
func (ctx *Context) StartServer() {

	if ctx.EnableTLS {
		go log.Fatal(tls.StartHTTPSListener(ctx.HttpService, ctx.CertFile, ctx.KeyFile))
	} else {
		go log.Fatal(http.ListenAndServe(ctx.HttpService, nil))
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

// SimulateFailure will run a failure simulation and return an HTTP code representing the outcome
func (ctx *Context) SimulateFailure() int {
	ctx.Mutex.Lock()

	defer ctx.Mutex.Unlock()

	var outcome int = ctx.HttpCode

	if ctx.FailureMode.Enabled {
		if ctx.FailureMode.failureIterationCount < ctx.FailureMode.FailureCount {
			outcome = ctx.FailureMode.FailureCode
			ctx.FailureMode.failureSimulated = true

			ctx.FailureMode.failureIterationCount++

			if ctx.FailureMode.failureIterationCount == ctx.FailureMode.FailureCount {
				// Done with the failure sequence, next call will return success if needed
				// Otherwise, continue with the failure sequence if success count is set to 0
				if ctx.FailureMode.SuccessCount > 0 {
					ctx.FailureMode.successIterationCount = 0
				} else {
					ctx.FailureMode.failureIterationCount = 0
				}
			}

		} else if ctx.FailureMode.successIterationCount < ctx.FailureMode.SuccessCount {
			ctx.FailureMode.successIterationCount++
			ctx.FailureMode.failureSimulated = false

			if ctx.FailureMode.successIterationCount == ctx.FailureMode.SuccessCount {
				// Done with the success sequence, next call will return failure if needed
				// Otherwise, continue with the success sequence if failure count is set to 0
				if ctx.FailureMode.FailureCount > 0 {
					ctx.FailureMode.failureIterationCount = 0
				} else {
					ctx.FailureMode.successIterationCount = 0
				}
			}
		}
	}

	return outcome
}

// SimulateDelay will introduce a timed delay if specified
func (ctx *Context) SimulateDelay() {
	if ctx.Delay > 0 {
		time.Sleep(time.Duration(ctx.Delay) * time.Millisecond)
	}
}

//
func (ctx *Context) FailureSimulated() bool {
	return ctx.FailureSimulationEnabled() && ctx.FailureMode.failureSimulated
}

// Determine if the failure simulation mode is enabled for this invocation
func (ctx *Context) FailureSimulationEnabled() bool {
	return ctx.FailureMode.Enabled
}

// Close the context
func (ctx *Context) Close() {
	// noop for now
}
