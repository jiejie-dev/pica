/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"github.com/spf13/cobra"
)

// benchCmd represents the bench command
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Benchmark apis.",
	Run: func(cmd *cobra.Command, args []string) {
		// Returns a new Benchmark pointer with all the defaults assigned
		// b := New()
		// // time to wait before firing the consequent request
		// // WaitPerReq = time.Millisecond * 1
		// // print available stats while the benchmark is running
		// b.ShowProgress = true
		// // Total number of requests to fire
		// b.TotalRequests = 10
		// // Duration in which all the requests have to be finished firing (in milliseconds).
		// b.BenchDuration = 7300

		// // Updates all the necessary fields according to the configuration provided
		// b.Init()
		// // Run the benchmark
		// b.Run(test)
	},
}

func init() {
	rootCmd.AddCommand(benchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
