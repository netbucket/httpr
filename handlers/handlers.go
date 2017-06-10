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

func ResponeCodeHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		w.WriteHeader(ctx.HttpCode)
	})
}

func DelayHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		time.Sleep(time.Duration(ctx.Delay) * time.Millisecond)
	})
}

func FailureSimulationHandler(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		w.WriteHeader(ctx.SimulateFailure())
	})
}

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
