// Copyright © 2017 Igor Bondarenko <igor@context7.com>
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
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func RawRequestLoggingHandler(out io.Writer, echo bool, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := httputil.DumpRequest(r, true)
		if err == nil {
			out.Write(body)

			if echo {
				w.Write(body)
			}
		}

	})
}

func JSONRequestLoggingHandler(out io.Writer, echo bool, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := json.Marshal(r)

		if err != nil {
			log.Fatal(err)
		} else {
			out.Write(body)

			if echo {
				w.Header().Set("Content-Type", "application/json")
				w.Write(body)
			}
		}
	})
}
