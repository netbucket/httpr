// Copyright © 2017 Igor Bondarenko <igor@context7.com>
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
  "fmt"
  "net/http"
  "net/http/httputil"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log the incoming HTTP requests",
	Long: `Log the incoming HTTP requests to the standard output or a file, in raw or structured JSON format.
See detailed help for options to modify the HTTP response behavior.`,
	Run: executeLog,
}

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolP("raw", "r", false, "Log HTTP requests in raw text format")
  logCmd.Flags().BoolP("json", "j", false, "Log HTTP requests in JSON format")
}

func executeLog(cmd *cobra.Command, args []string) {
  http.HandleFunc("/", logRequests)
  startServer()
}

func logRequests(w http.ResponseWriter, r *http.Request) {
  requestData, err := httputil.DumpRequest(r, true)
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(string(requestData))
}
