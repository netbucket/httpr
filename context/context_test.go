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
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	ctx1 := Instance()
	ctx2 := Instance()

	if ctx1 != ctx2 {
		t.Error("Context instance is not a singleton")
	}
}

func TestDisabledSimulateFailure(t *testing.T) {
	expectedHttpCode := 200

	ctx := &Context{
		Mutex: &sync.Mutex{},
		FailureMode: FailureSimulation{
			Enabled: false, FailureCount: 2, SuccessCount: 2, FailureCode: 500},
		HttpCode: expectedHttpCode}

	actualHttpCode := ctx.SimulateFailure()

	if expectedHttpCode != actualHttpCode {
		t.Errorf("Expected HTTP status code %d, got %d", expectedHttpCode, actualHttpCode)
	}

}

func TestSimulateFailure(t *testing.T) {

	tests := []*Context{
		&Context{
			Mutex: &sync.Mutex{},
			FailureMode: FailureSimulation{
				Enabled: true, FailureCount: 5, SuccessCount: 10, FailureCode: 500},
		},
		&Context{
			Mutex: &sync.Mutex{},
			FailureMode: FailureSimulation{
				Enabled: true, FailureCount: 1, SuccessCount: 1, FailureCode: 502},
		},
		&Context{
			Mutex: &sync.Mutex{},
			FailureMode: FailureSimulation{
				Enabled: true, FailureCount: 5, SuccessCount: 0, FailureCode: 500},
		},
		&Context{
			Mutex: &sync.Mutex{},
			FailureMode: FailureSimulation{
				Enabled: true, FailureCount: 0, SuccessCount: 5, FailureCode: 500},
		},
		&Context{
			Mutex: &sync.Mutex{},
			FailureMode: FailureSimulation{
				Enabled: true, FailureCount: 0, SuccessCount: 0, FailureCode: 500},
		},
	}

	for _, ctx := range tests {
		for fCount := 0; fCount < ctx.FailureMode.FailureCount; fCount++ {
			if httpCode := ctx.SimulateFailure(); httpCode != ctx.FailureMode.FailureCode {
				t.Errorf("Expected HTTP status code %d, got %d")
			}
		}
		for sCount := 0; sCount < ctx.FailureMode.SuccessCount; sCount++ {
			if httpCode := ctx.SimulateFailure(); httpCode != ctx.HttpCode {
				t.Errorf("Expected HTTP status code %d, got %d")
			}
		}
	}
}
