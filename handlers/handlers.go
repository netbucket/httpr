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

package handlers

import (
	"bytes"
	"github.com/netbucket/httpr/context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// RawRequestLoggingHandler returns a handler function that logs the incoming
// HTTP request in plain text format
func RawRequestLoggingHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := httputil.DumpRequest(r, true)
		if err == nil {
			body = append(body, []byte("\n")...)

			ctx.Out.Write(body)

			if ctx.Echo {
				w.Write(body)
			}
		}
	})
}

// JSONLoggingHandler returns a handler function that logs the incoming
// HTTP request in a compact or formatted JSON format
func JSONRequestLoggingHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := encodeAsJSON(r, ctx.LogPrettyJSON)

		if err != nil {
			log.Fatal(err)
		} else {
			body = append(body, []byte("\n")...)
			ctx.Out.Write(body)

			if ctx.Echo {
				w.Header().Set("Content-Type", "application/json")
				w.Write(body)
			}
		}
	})
}

// ResponseCodeHandler returns a handler function that returns the specified
// HTTP status code
func ResponseCodeHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		w.WriteHeader(ctx.HttpCode)
	})
}

// DelayHandler returns a handler function that introduces a delay in
// responding to the HTTP request
func DelayHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		time.Sleep(time.Duration(ctx.Delay) * time.Millisecond)
	})
}

// FailureSimulationHandler returns a handler function that simulates a transient
// HTTP failure scenario, e.g. a series of failure response codes followed
// by a series of successful HTTP status codes
func FailureSimulationHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		w.WriteHeader(ctx.SimulateFailure())
	})
}

// copyRequestBody makes a non-destructive copy of the HTTP request body contents
// to make the contents available for repeated use by multiple HTTP handlers
func copyRequestBody(r *http.Request) []byte {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return []byte("")
	}

	// Reset the body to allow other HTTP handlers to read the contents as well
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	return data
}
