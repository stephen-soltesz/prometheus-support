// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/common/model"

	"github.com/prometheus/prometheus/promql"
	//	"github.com/prometheus/prometheus/rules"
)

func TestAlertingRule(t *testing.T) {
	suite, err := promql.NewTest(t, `
		load 5m
			http_requests{job="app-server", instance="0", group="canary"}	75 85  95 105 105  95  85
			http_requests{job="app-server", instance="1", group="canary"}   80 90 100 110 120 130 140
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer suite.Close()

	fmt.Printf("Running suite\n")
	if err := suite.Run(); err != nil {
		t.Fatal(err)
	}

	engine := suite.QueryEngine()
	query, err := engine.NewInstantQuery(`http_requests{group="canary", job="app-server"} < 100`, model.Time(0).Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	result := query.Exec(suite.Context())
	if result.Err != nil {
		t.Fatal(result.Err)
	}

	fmt.Printf("Type: %q\n", result.Value.Type())
	actual := strings.Split(result.Value.String(), "\n")
	fmt.Printf("Value: %q\n", actual)
}

func TestEvaluations(t *testing.T) {
	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range files {
		content, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Fatal(err)
		}
		test, err := promql.NewTest(t, string(content))
		if err != nil {
			t.Errorf("error creating test for %s: %s", fn, err)
		}
		err = test.Run()
		if err != nil {
			t.Errorf("error running test %s: %s", fn, err)
		}
		test.Close()
	}
}
