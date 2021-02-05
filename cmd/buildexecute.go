/*
Copyright © 2021 SUDIPTO BARAL <sudiptobaral.me@gmail.com>

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
	"fmt"
	"os"
	"regexp"

	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// buildexecuteCmd represents the buildexecute command
var buildexecuteCmd = &cobra.Command{
	Use:   "buildexecute",
	Short: "build and execute the project of specified directory",
	Long:  `Builds the porject from the specified directory. First copies the files from the specified prohject, then builds the project accordingly. Currently only the copy function is supported`,
	Run: func(cmd *cobra.Command, args []string) {
		sourchePath := getSourcePath(cmd)
		destinationPath := getDestinationPath()
		copyDir(sourchePath, destinationPath)

	},
}

func init() {
	rootCmd.AddCommand(buildexecuteCmd)

	// Flags
	buildexecuteCmd.PersistentFlags().StringP("copydir", "c", "", "Copies the content of the specified dirtectory")

}

func copyDir(sourcePath, destinationPath string) {

	var destinationFolder string = ""
	re := regexp.MustCompile(`[/\\]`)
	sourceSplit := re.Split(sourcePath, -1)
	destinationSplit := re.Split(destinationPath, -1)

	// Build the destination folder for Copy function
	// Delimeter of path is different in differnet operating system
	// Build the path using forward slash because copy function accepts '/'
	for i := range destinationSplit {
		destinationFolder += destinationSplit[i] + "/"
	}
	destinationFolder += sourceSplit[len(sourceSplit)-1]

	//fmt.Println("Source folder :", sourcePath)
	//fmt.Println("Destination  folder :", destinationFolder)

	err := copy.Copy(sourcePath, destinationFolder)
	if err != nil {
		logrus.Error(err)
	} else {
		fmt.Println("Files copied successfully")
	}
}

func getSourcePath(cmd *cobra.Command) string {
	sourcePath, _ := cmd.Flags().GetString("copydir")

	if sourcePath == "" {
		fmt.Println(`Please specify the source path. Ex: gobuildercli buildexecute --copydir "configs/foldertocopy"`)
		os.Exit(0)
	}
	return sourcePath
}

func getDestinationPath() string {

	destinationPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return destinationPath
}
