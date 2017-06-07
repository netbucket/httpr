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
  "log"
  "net/http"
	"os"
  "os/signal"
  "syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
  httpService string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "httpr",
	Short: "A simple HTTP server for examining and testing HTTP requests",
	Long: `HTTP Rake (or httpr) is a compact and flexible HTTP server.
This application is a useful tool for examining and testing HTTP requests
without the need to configure and deploy a full-fledged Web server or a proxy.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

  RootCmd.PersistentFlags().StringVarP(&httpService, "http", "s", ":3115", "HTTP service address")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

// Start the HTTP server
func startServer() {
  go log.Fatal(http.ListenAndServe(httpService, nil))

  ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
