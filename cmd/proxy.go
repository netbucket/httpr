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

package cmd

import (
	"net/http"

	"github.com/netbucket/httpr/context"
	"github.com/netbucket/httpr/handlers"
	"github.com/spf13/cobra"
	"log"
	"net/url"
)

// logCmd represents the log command
var proxyCmd = &cobra.Command{
	Use:   "proxy <upstream-url>",
	Short: "Proxy HTTP requests to an upstream server.",
	Long: `Start the HTTP server that will proxy requests to an upstream server indicated by the upstream-url argument,
and log the incoming HTTP requests to the standard output. See options to modify the HTTP response behavior.`,
	Run: executeProxy,
}

func init() {
	RootCmd.AddCommand(proxyCmd)

	ctx := context.Instance()

	proxyCmd.Flags().BoolVarP(&ctx.LogJSON, "json", "j", false, "Log HTTP requests in JSON format")
	proxyCmd.Flags().BoolVarP(&ctx.LogPrettyJSON, "json-pp", "p", false, "Log HTTP requests in pretty-printed (indented) JSON format")
	proxyCmd.Flags().IntVarP(&ctx.Delay, "delay", "d", 0, "Delay, in milliseconds, when replying to incoming HTTP requests")
	proxyCmd.Flags().BoolVarP(&ctx.FailureMode.Enabled, "simulate-failure", "f", false, "Simulate a transient failure: return an error code before proxying the request upstream")
	proxyCmd.Flags().IntVarP(&ctx.FailureMode.FailureCount, "simulate-failure-count", "", 1, "For --simulate-failure, determines how many errors are returned before proxying the request upstream")
	proxyCmd.Flags().IntVarP(&ctx.FailureMode.FailureCode, "simulate-failure-code", "", 500, "For --simulate-failure, determines the HTTP status code for an error response")
	proxyCmd.Flags().BoolVarP(&ctx.IgnoreTLSErrors, "insecure", "k", false, "Ignore upstream TLS certificate errors")
}

func executeProxy(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		log.Fatal("Upstream URL argument missing")
	}

	ctx := context.Instance()

	u, err := url.Parse(args[0])

	if err != nil {
		log.Fatal(err)
	}

	ctx.UpstreamURL = u

	h := setupProxyHandlerChain(ctx)

	http.Handle("/", h)

	// Start the HTTP server and handle the command
	ctx.StartServer()

	ctx.Close()
}

func setupProxyHandlerChain(ctx *context.Context) http.Handler {
	var h http.Handler
	{
		h = handlers.ProxyHandler(ctx, nil)

		h = handlers.DelayHandler(ctx, h)

		if ctx.LogJSON || ctx.LogPrettyJSON {
			h = handlers.JSONRequestLoggingHandler(ctx, h)
		} else {
			h = handlers.RawRequestLoggingHandler(ctx, h)
		}

		if ctx.FailureMode.Enabled {
			h = handlers.FailureSimulationHandler(ctx, h)
		}
	}

	return h
}
