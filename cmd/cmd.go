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

package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/palantir/go-nobadfuncs/nobadfuncs"
	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "nobadfuncs [flags] [packages]",
		Short: "verifies that blacklisted functions are not called",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrapf(err, "failed to determine working directory")
			}

			if printAllFlagVal {
				// if print-all flag is specified, perform print all action
				return nobadfuncs.PrintAllFuncRefs(args, wd, cmd.OutOrStdout())
			}
			return printBadFuncRefsJSONConfig(args, configJSONFlagVal, wd, cmd.OutOrStdout())
		},
	}

	printAllFlagVal   bool
	configJSONFlagVal string
)

func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func init() {
	rootCmd.Flags().BoolVar(&printAllFlagVal, "print-all", false, "print all function references in the provided package (useful for determining format of forbidden references)")
	rootCmd.Flags().StringVar(&configJSONFlagVal, "config-json", "", "the JSON configuration for the check")
}

func printBadFuncRefsJSONConfig(pkgs []string, jsonConfig string, dir string, w io.Writer) error {
	var sigs map[string]string
	if jsonConfig != "" {
		if err := json.Unmarshal([]byte(jsonConfig), &sigs); err != nil {
			return errors.Wrapf(err, "failed to unmarshal configuration as JSON: %q", jsonConfig)
		}
	}
	return nobadfuncs.PrintBadFuncRefs(pkgs, sigs, dir, w)
}
