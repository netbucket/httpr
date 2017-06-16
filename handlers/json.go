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
	"encoding/json"
	"net/http"
)

type requestModel struct {
	RemoteAddr string      `json:"remoteAddr,omitempty"`
	Host       string      `json:"host,omitempty"`
	Method     string      `json:"method,omitempty"`
	URL        string      `json:"url,omitempty"`
	Proto      string      `json:"proto,omitempty"`
	Header     http.Header `json:"header,omitempty"`
	//Body io.ReadCloser
	ContentLength    int64    `json:"content_length,omitempty"`
	TransferEncoding []string `json:"transfer_encoding,omitempty"`
	Body             string   `json:"body,omitempty"`
}

func EncodeAsJSON(r *http.Request, prettyPrint bool) ([]byte, error) {
	model := requestModel{
		RemoteAddr: r.RemoteAddr, Host: r.Host, Method: r.Method,
		URL: r.RequestURI, Proto: r.Proto, Header: r.Header,
		ContentLength: r.ContentLength, TransferEncoding: r.TransferEncoding,
	}

	if r.Body != nil {
		model.Body = string(copyRequestBody(r))
	}

	if prettyPrint {
		return json.MarshalIndent(model, "", "    ")
	}

	return json.Marshal(model)
}
