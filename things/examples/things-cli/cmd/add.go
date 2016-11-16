// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"fmt"

	"github.com/spf13/cobra"
	"log"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new thing",
	Long: `Add a new thing`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := createConn()
		if err != nil { er(err) }
		log.Println(conn)

		if len(args) > 2 || len(args) < 2{
			er("invalid args count")
		}

		// if file
		// load content from file

		// validate json content
		// call add


		fmt.Println("add called")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&thingId, "id", "t", "", "")
	addCmd.Flags().StringVarP(&content, "content", "c", "", "")
	addCmd.Flags().StringVarP(&file, "file", "f", "", "")
}
