// Copyright 2018 Drone.IO Inc
//
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
//
// Taken from github.com/drone/drone-template-lib
package plugin

import (
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/drone-plugins/drone-plugin-lib/drone"
	"github.com/flowchartsman/handlebars/v3"
)

func TestToDuration(t *testing.T) {
	from := time.Date(2017, time.November, 15, 23, 0, 0, 0, time.UTC).Unix()
	vals := map[int64]string{
		time.Date(2018, time.November, 15, 23, 0, 0, 0, time.UTC).Unix():   "8760h0m0s",
		time.Date(2017, time.November, 16, 23, 0, 0, 0, time.UTC).Unix():   "24h0m0s",
		time.Date(2017, time.November, 15, 23, 30, 0, 0, time.UTC).Unix():  "30m0s",
		time.Date(2017, time.November, 15, 23, 10, 15, 0, time.UTC).Unix(): "10m15s",
		time.Date(2017, time.October, 15, 23, 0, 0, 0, time.UTC).Unix():    "-744h0m0s",
	}
	for input, want := range vals {
		if got := toDuration(from, input); got != want {
			t.Errorf("Want transform %d-%d to %s, got %s", from, input, want, got)
		}
	}
}

func TestTruncate(t *testing.T) {
	vals := map[string]string{
		"foobarz": "fooba",
		"foöäüüu": "foöäü",
		"üpsßßßk": "üpsßß",
		"1234567": "12345",
		"!'§$%&/": "!'§$%",
	}
	for input, want := range vals {
		if got := truncate(input, 5); got != want {
			t.Errorf("Want transform %s to %s, got %s", input, want, got)
		}
	}
}

func TestNegativeTruncate(t *testing.T) {
	vals := map[string]string{
		"foobarz": "rz",
		"foöäüüu": "üu",
		"üpsßßßk": "ßk",
		"1234567": "67",
		"!'§$%&/": "&/",
	}
	for input, want := range vals {
		if got := truncate(input, -5); got != want {
			t.Errorf("Want transform %s to %s, got %s", input, want, got)
		}
	}
}

func TestUppercaseFirst(t *testing.T) {
	vals := map[string]string{
		"hello":  "Hello",
		"ßqwert": "ßqwert",
		"üps":    "Üps",
		"12345":  "12345",
		"Foobar": "Foobar",
	}
	for input, want := range vals {
		if got := uppercaseFirst(input); got != want {
			t.Errorf("Want transform %s to %s, got %s", input, want, got)
		}
	}
}

func TestRegexReplace(t *testing.T) {
	expected := "hello-my-String-123"
	actual := regexReplace("(.*?)\\/(.*)", "hello/my-String-123", "$1-$2")
	if actual != "hello-my-String-123" {
		t.Errorf("error, expected %s, got %s", expected, actual)
	}
}

func TestRender(t *testing.T) {
	s, err := handlebars.Render(string(testTpl), testPipeline())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf(">>>\n%s\n<<<", strings.TrimSpace(s))
}

func testPipeline() drone.Pipeline {
	return drone.Pipeline{
		Repo: drone.Repo{
			Owner: "octocat",
			Name:  "hello-world",
		},
		Build: drone.Build{
			Tag:      "1.0.0",
			Event:    "push",
			Number:   1,
			Branch:   "master",
			DeployTo: "",
			Status:   "success",
			Link:     "http://github.com/octocat/hello-world",
			Started:  time.Unix(1546340400, 0), // 2019-01-01 12:00:00
			Created:  time.Unix(1546340400, 0), // 2019-01-01 12:00:00
		},
		Commit: drone.Commit{
			SHA:     "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Ref:     "",
			Branch:  "master",
			Link:    "https://github.com/octocat/hello-world/compare/0000000...7fd1a60b01f91b31",
			Message: drone.ParseMessage("Initial commit\n\nMessage body line1\nmessage body line 2"),
			Author: drone.Author{
				Username: "octocat",
				Name:     "The Octocat",
				Email:    "octocat@github.com",
				Avatar:   "https://avatars0.githubusercontent.com/u/583231?s=460&v=4",
			},
		},
	}
}

//go:embed test.tpl
var testTpl []byte
