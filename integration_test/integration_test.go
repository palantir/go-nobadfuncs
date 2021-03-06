// Copyright 2016 Palantir Technologies, Inc.
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

package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/pkg/products/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoBadFuncs(t *testing.T) {
	// explicitly unset GOFLAGS environment variable during tests -- these tests perform module/package resolution and
	// assume that vendor mode is not enabled.
	prevValue := os.Getenv("GOFLAGS")
	defer func() {
		_ = os.Setenv("GOFLAGS", prevValue)
	}()
	err := os.Setenv("GOFLAGS", "")
	require.NoError(t, err)

	cli, err := products.Bin("go-nobadfuncs")
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name          string
		filesToCreate []gofiles.GoFileSpec
		args          []string
		expectErr     bool
		wantStdout    func(currTestCaseDir string) string
	}{
		{
			name: "Empty configuration has blank output",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `
package foo

import (
	"net/http"
)

func MyFunction() {
	http.DefaultClient.Do(nil)
}
`,
				},
			},
			args: []string{
				"./foo",
			},
			expectErr: false,
			wantStdout: func(currTestCaseDir string) string {
				return ""
			},
		},
		{
			name: "Basic case",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `
package foo

import (
	"net/http"
)

func MyFunction() {
	http.DefaultClient.Do(nil)
}
`,
				},
			},
			args: []string{
				"--config-json",
				`{"func (*net/http.Client).Do(*net/http.Request) (*net/http.Response, error)": ""}`,
				"./foo",
			},
			expectErr: true,
			wantStdout: func(currTestCaseDir string) string {
				return fmt.Sprintf("%s/foo/foo.go:9:21: references to \"func (*net/http.Client).Do(*net/http.Request) (*net/http.Response, error)\" are not allowed. Remove this reference or whitelist it by adding a comment of the form '// OK: [reason]' to the line before it.\n", currTestCaseDir)
			},
		},
		{
			name: "All flag",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `
package foo

import (
	"net/http"
)

func MyFunction() {
	http.DefaultClient.Do(nil)
}
`,
				},
			},
			args: []string{
				"--print-all",
				"./foo",
			},
			expectErr: false,
			wantStdout: func(currTestCaseDir string) string {
				return fmt.Sprintf("%s/foo/foo.go:9:21: func (*net/http.Client).Do(*net/http.Request) (*net/http.Response, error)\n", currTestCaseDir)
			},
		},
	} {
		t.Run(currCase.name, func(t *testing.T) {
			currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
			require.NoError(t, err)

			_, err = gofiles.Write(currCaseTmpDir, append(currCase.filesToCreate, gofiles.GoFileSpec{
				RelPath: "go.mod",
				Src:     "module github.com/palantir/go-nobadfuncs-test",
			}))
			require.NoError(t, err, "Case %d", i)

			var output []byte
			func() {
				err := os.Chdir(currCaseTmpDir)
				defer func() {
					err := os.Chdir(wd)
					require.NoError(t, err)
				}()
				require.NoError(t, err)

				cmd := exec.Command(cli, currCase.args...)
				output, err = cmd.CombinedOutput()

				if currCase.expectErr {
					require.Error(t, err, fmt.Sprintf("Case %d: %s\nOutput: %s", i, currCase.name, string(output)))
				} else {
					require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
				}
			}()

			// make expected dir location canonical
			currCaseTmpDir, err = filepath.EvalSymlinks(currCaseTmpDir)
			require.NoError(t, err)

			assert.Equal(t, currCase.wantStdout(currCaseTmpDir), string(output), "Case %d: %s\nOutput:\n%s", i, currCase.name, string(output))
		})
	}
}
