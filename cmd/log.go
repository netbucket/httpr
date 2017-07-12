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
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log the incoming HTTP requests",
	Long: `Start the HTTP server and log the incoming HTTP requests to the standard output.
See options to modify the HTTP response behavior.`,
	Run: executeLog,
}

func init() {
	RootCmd.AddCommand(logCmd)

	ctx := context.Instance()

	logCmd.Flags().BoolVarP(&ctx.LogJSON, "json", "j", false, "Log HTTP requests in JSON format")
	logCmd.Flags().BoolVarP(&ctx.LogPrettyJSON, "json-pp", "p", false, "Log HTTP requests in pretty-printed (indented) JSON format")
	logCmd.Flags().BoolVarP(&ctx.Echo, "echo", "e", false, "Send the logged contents back to the HTTP client")
	logCmd.Flags().IntVarP(&ctx.HttpCode, "response-code", "r", 200, "Send the specified HTTP status code back to the client")
	logCmd.Flags().IntVarP(&ctx.Delay, "delay", "d", 0, "Delay, in milliseconds, when replying to incoming HTTP requests")
	logCmd.Flags().BoolVarP(&ctx.FailureMode.Enabled, "simulate-failure", "f", false, "Simulate a transient failure: return an error code before a successful response")
	logCmd.Flags().IntVarP(&ctx.FailureMode.FailureCount, "simulate-failure-count", "", 1, "For --simulate-failure, determines how many errors are returned before a successful response")
	logCmd.Flags().IntVarP(&ctx.FailureMode.SuccessCount, "simulate-success-count", "", 1, "For --simulate-failure, determines how many successful responses are returned before returning a error code")
	logCmd.Flags().IntVarP(&ctx.FailureMode.FailureCode, "simulate-failure-code", "", 500, "For --simulate-failure, determines the HTTP status code for an error response")
}

func executeLog(cmd *cobra.Command, args []string) {
	ctx := context.Instance()

	h := setupLogHandlerChain(ctx)

	http.Handle("/", h)

	// Start the HTTP server and handle the command
	ctx.StartServer()

	ctx.Close()
}

func setupLogHandlerChain(ctx *context.Context) http.Handler {
	var h http.Handler
	{
		h = handlers.DelayHandler(ctx, nil)

		if ctx.LogJSON || ctx.LogPrettyJSON {
			h = handlers.JSONRequestLoggingHandler(ctx, h)
		} else {
			h = handlers.RawRequestLoggingHandler(ctx, h)
		}

		if ctx.FailureMode.Enabled {
			h = handlers.FailureSimulationHandler(ctx, h)
		} else {
			h = handlers.ResponseCodeHandler(ctx, h)
		}

		h = handlers.ContentTypeHandler(ctx, h)
	}

	return h
}
