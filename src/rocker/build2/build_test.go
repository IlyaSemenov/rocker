/*-
 * Copyright 2015 Grammarly, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package build2

import (
	"rocker/template"
	"runtime"
	"strings"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBuild(t *testing.T) {
	b, _ := makeBuild(t, "FROM ubuntu", BuildConfig{})
	assert.IsType(t, &Rockerfile{}, b.rockerfile)
}

// internal helpers

func makeBuild(t *testing.T, rockerfileContent string, cfg BuildConfig) (*Build, *MockClient) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	r, err := NewRockerfile(fn.Name(), strings.NewReader(rockerfileContent), template.Vars{}, template.Funs{})
	if err != nil {
		t.Fatal(err)
	}

	c := &MockClient{}

	b, err := New(c, r, cfg)
	if err != nil {
		t.Fatal(err)
	}

	return b, c
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) InspectImage(name string) (*docker.Image, error) {
	args := m.Called(name)
	return args.Get(0).(*docker.Image), args.Error(1)
}

func (m *MockClient) PullImage(name string) error {
	args := m.Called(name)
	return args.Error(0)
}