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
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var excludeTestFlag bool

// buildexecuteCmd represents the buildexecute command
var buildexecuteCmd = &cobra.Command{
	Use:   "buildexecute",
	Short: "build and execute the project of specified directory",
	Long:  `Builds the porject from the specified directory. First copies the files from the specified prohject, then builds the project accordingly. Currently only the copy function is supported`,
	Run: func(cmd *cobra.Command, args []string) {
		sourchePath := getSourcePath(cmd)
		destinationPath := getDestinationPath(cmd, sourchePath)
		copyDir(sourchePath, destinationPath)
		compileAsBin(cmd, sourchePath)

	},
}

func init() {
	rootCmd.AddCommand(buildexecuteCmd)

	// Flags
	buildexecuteCmd.PersistentFlags().StringP("copydir", "c", "", "Copies the content of the specified dirtectory")
	buildexecuteCmd.PersistentFlags().StringP("builddir", "b", "", "Specify the destination directory")
	buildexecuteCmd.PersistentFlags().StringP("exe", "e", "", "compile the code of the directory as a binary")

	buildexecuteCmd.PersistentFlags().BoolVarP(&excludeTestFlag, "exclude-tests", "x", false, "excludes the golang test files")
}

func compileAsBin(cmd *cobra.Command, sourcePath string) {
	sourcePath = formatPath(sourcePath)
	binaryName, _ := cmd.Flags().GetString("exe")
	if binaryName == "" {
		return
	}

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	_, err := exec.Command("go", "build", "-o", binaryName, sourcePath).CombinedOutput()
	if err != nil {
		logrus.Error("No compileable go code found! ", err)
		return
	}
	fmt.Println(binaryName, "created in ", sourcePath)

}

func copyDir(sourcePath, destinationPath string) {

	// Format the paths into a common delimiter '/'
	sourcePath = formatPath(sourcePath)
	destinationPath = formatPath(destinationPath)

	sourceSplit := strings.Split(sourcePath, "/")
	// if a folder is specified,
	// append the name of the folder to the destination path
	// to  create a folder of same name
	// so that otiai10/copy can copy everyting into that folder
	// in the destination
	// if content of the file specified we do not append
	if sourceSplit[len(sourceSplit)-1] != "" {

		destinationPath += "/" + sourceSplit[len(sourceSplit)-1]
	}

	//fmt.Println("Source folder :", sourcePath)
	//fmt.Println("Destination  folder :", destinationPath)

	//Do not perform copy if the source and the destination is the same
	if isDestinationUnderSource(sourcePath, destinationPath) {
		fmt.Println("Destination folder is under the source directory. Did not copy.")
		return
	}
	opt := copy.Options{
		OnDirExists: func(src, dest string) copy.DirExistsAction {
			// Replace if directory exist
			return 1
		},
	}
	if excludeTestFlag == true {
		opt = copy.Options{
			// Skip the test files
			Skip: func(src string) (bool, error) {
				return strings.HasSuffix(src, "_test.go"), nil
			},
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				// Replace if directory exist
				return 1
			},
		}
	}
	err := copy.Copy(sourcePath, destinationPath, opt)
	if err != nil {
		logrus.Error(err)
	} else {
		fmt.Println("Copied to ", destinationPath)
	}
}

func getSourcePath(cmd *cobra.Command) string {
	sourcePath, _ := cmd.Flags().GetString("copydir")

	if sourcePath == "" {
		fmt.Println("Please specify a valid source path. Ex: gobuildercli buildexecute --copydir configs/foldertocopy")
		os.Exit(0)
	}

	return sourcePath
}

func getDestinationPath(cmd *cobra.Command, sourcePath string) string {

	// if destination path is sepecified return the specified path
	destinationPath, _ := cmd.Flags().GetString("builddir")

	// if destination path is not specified
	// return the current directory as destination
	if destinationPath == "" {

		return getCurrentDir()
	}

	return destinationPath
}

// Checks if the destination directory is under the source directory or not

func isDestinationUnderSource(sourcePath, destinationPath string) bool {

	sourceSplit := strings.Split(sourcePath, "/")
	destinationSplit := strings.Split(destinationPath, "/")

	// if the length of destination split is lower than source split
	// then surely destination folder is not under the source folder
	if len(destinationSplit) < len(sourceSplit) {
		return false
	}

	flag := true
	limit := len(sourceSplit)

	for i := 0; i < limit; i++ {
		if sourceSplit[i] != destinationSplit[i] {
			flag = false
			break
		}
	}
	return flag

}

// Returns the path in the linux format
func getCurrentDir() string {

	currentWorkingDirtectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return currentWorkingDirtectory

}

func formatPath(path string) string {
	var formattedPath string = ""
	// Delimiter used in path string
	re := regexp.MustCompile(`[/\\]`)
	token := re.Split(path, -1)

	// Build the path in the linux format
	for i := range token {
		formattedPath += token[i]
		if i < len(token)-1 {
			formattedPath += "/"
		}
	}

	return formattedPath
}
