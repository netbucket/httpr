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
	"github.com/netbucket/httpr/context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestDelayHandler(t *testing.T) {
	const expectedDelay = 500

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	ctx := &context.Context{
		Mutex:       &sync.Mutex{},
		FailureMode: context.FailureSimulation{Enabled: false},
		Delay:       expectedDelay}

	h := DelayHandler(ctx, nil)

	start := time.Now()

	h.ServeHTTP(rec, req)

	durationMillis := time.Since(start) * time.Millisecond

	if durationMillis < expectedDelay {
		t.Errorf("Expected %d milliseconds delay, got %d", expectedDelay, durationMillis)
	}
}

func TestResponeCodeHandler(t *testing.T) {
	const expectedResponse = 419

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	ctx := &context.Context{
		Mutex:       &sync.Mutex{},
		FailureMode: context.FailureSimulation{Enabled: false},
		HttpCode:    expectedResponse}

	ResponseCodeHandler(ctx, nil).ServeHTTP(rec, req)

	if rec.Code != expectedResponse {
		t.Errorf("Expected HTTP status %d, got %d", expectedResponse, rec.Code)
	}
}

func TestFailureSimulationHandler(t *testing.T) {
	const expectedFailureResponse = 503
	const expectedSuccessResponse = 201

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	ctx := &context.Context{
		Mutex: &sync.Mutex{},
		FailureMode: context.FailureSimulation{
			Enabled:      true,
			FailureCount: 1,
			FailureCode:  expectedFailureResponse,
			SuccessCount: 1,
		},
		HttpCode: expectedSuccessResponse}

	h := FailureSimulationHandler(ctx, nil)

	h.ServeHTTP(rec, req)

	if rec.Code != expectedFailureResponse {
		t.Errorf("Expected HTTP status %d, got %d", expectedFailureResponse, rec.Code)
	}

	rec = httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != expectedSuccessResponse {
		t.Errorf("Expected HTTP status %d, got %d", expectedSuccessResponse, rec.Code)
	}
}

func TestJSONContentType(t *testing.T) {
	expectedContentType := "application/json"

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	ctx := &context.Context{
		Mutex:       &sync.Mutex{},
		FailureMode: context.FailureSimulation{Enabled: false},
		HttpCode:    http.StatusOK,
		LogJSON: true}

	ContentTypeHandler(ctx, nil).ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")

	if expectedContentType != contentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
}
