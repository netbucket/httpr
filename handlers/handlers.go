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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

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

func RawRequestLoggingHandler(out io.Writer, echo bool, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := httputil.DumpRequest(r, true)
		if err == nil {
			body = append(body, []byte("\n")...)

			out.Write(body)

			if echo {
				w.Write(body)
			}
		}
	})
}

func JSONRequestLoggingHandler(out io.Writer, prettyPrint bool, echo bool, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h.ServeHTTP(w, r)
		}

		body, err := encodeAsJSON(r, prettyPrint)

		if err != nil {
			log.Fatal(err)
		} else {
			body = append(body, []byte("\n")...)
			out.Write(body)

			if echo {
				w.Header().Set("Content-Type", "application/json")
				w.Write(body)
			}
		}
	})
}
