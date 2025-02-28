/*
Copyright 2022 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tektoncd/pipeline/pkg/apis/config"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/test/diff"
	"knative.dev/pkg/apis"
)

func TestResultsValidate(t *testing.T) {
	tests := []struct {
		name      string
		Result    v1.TaskResult
		apiFields string
	}{{
		name: "valid result type empty",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Description: "my great result",
		},
		apiFields: "stable",
	}, {
		name: "valid result type string",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeString,
			Description: "my great result",
		},

		apiFields: "stable",
	}, {
		name: "valid result type array",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeArray,
			Description: "my great result",
		},

		apiFields: "alpha",
	}, {
		name: "valid result type object",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeObject,
			Description: "my great result",
			Properties:  map[string]v1.PropertySpec{"hello": {Type: v1.ParamTypeString}},
		},
		apiFields: "alpha",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.apiFields == "alpha" {
				ctx = config.EnableAlphaAPIFields(ctx)
			}
			if err := tt.Result.Validate(ctx); err != nil {
				t.Errorf("TaskSpec.Validate() = %v", err)
			}
		})
	}
}

func TestResultsValidateError(t *testing.T) {
	tests := []struct {
		name          string
		Result        v1.TaskResult
		apiFields     string
		expectedError apis.FieldError
	}{{
		name: "invalid result type in stable",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        "wrong",
			Description: "my great result",
		},
		apiFields: "stable",
		expectedError: apis.FieldError{
			Message: `invalid value: wrong`,
			Paths:   []string{"type"},
			Details: "type must be string",
		},
	}, {
		name: "invalid result type in alpha",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        "wrong",
			Description: "my great result",
		},
		apiFields: "alpha",
		expectedError: apis.FieldError{
			Message: `invalid value: wrong`,
			Paths:   []string{"type"},
			Details: "type must be string",
		},
	}, {
		name: "invalid array result type in stable",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeArray,
			Description: "my great result",
		},
		apiFields: "stable",
		expectedError: apis.FieldError{
			Message: "results type requires \"enable-api-fields\" feature gate to be \"alpha\" or \"beta\" but it is \"stable\"",
		},
	}, {
		name: "invalid object result type in stable",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeObject,
			Description: "my great result",
			Properties:  map[string]v1.PropertySpec{"hello": {Type: v1.ParamTypeString}},
		},
		apiFields: "stable",
		expectedError: apis.FieldError{
			Message: "results type requires \"enable-api-fields\" feature gate to be \"alpha\" or \"beta\" but it is \"stable\"",
		},
	}, {
		name: "invalid object properties type",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeObject,
			Description: "my great result",
			Properties:  map[string]v1.PropertySpec{"hello": {Type: "wrong type"}},
		},
		apiFields: "alpha",
		expectedError: apis.FieldError{
			Message: "The value type specified for these keys [hello] is invalid, the type must be string",
			Paths:   []string{"MY-RESULT.properties"},
		},
	}, {
		name: "invalid object properties empty",
		Result: v1.TaskResult{
			Name:        "MY-RESULT",
			Type:        v1.ResultsTypeObject,
			Description: "my great result",
		},
		apiFields: "alpha",
		expectedError: apis.FieldError{
			Message: "missing field(s)",
			Paths:   []string{"MY-RESULT.properties"},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.apiFields == "alpha" {
				ctx = config.EnableAlphaAPIFields(ctx)
			}
			err := tt.Result.Validate(ctx)
			if err == nil {
				t.Fatalf("Expected an error, got nothing for %v", tt.Result)
			}
			if d := cmp.Diff(tt.expectedError.Error(), err.Error(), cmpopts.IgnoreUnexported(apis.FieldError{})); d != "" {
				t.Errorf("TaskSpec.Validate() errors diff %s", diff.PrintWantGot(d))
			}
		})
	}
}
