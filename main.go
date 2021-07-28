package main

import (
	"fmt"
	"github.com/jinzhu/copier"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var dir string

type Env struct {
	Key   string
	Value string
}

func init() {
	flag.StringVar(&dir, "dir", "./", "directory to operate within (relative or absolute)")
}

func main() {
	flag.Parse()

	vars := parseEnvs()
	if len(vars) == 0 {
		fmt.Println("No env vars prefixed with 'REACT_APP_', doing nothing")
		os.Exit(0)
	}

	files := findFiles(dir)

	var changeList []string
	for _, file := range files {
		updatedFile := changeFile(file, vars)
		if updatedFile {
			changeList = append(changeList, file)
		}
	}

	for _, changed := range changeList {
		fmt.Printf("changed file: %s\n", changed)
	}

}

// split all env vars into key/values and search for keys prefixed with REACT_APP_
func parseEnvs() []Env {
	var result []Env

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		val := pair[1]
		obj := Env{
			Key:   key,
			Value: val,
		}
		if strings.HasPrefix(key, "REACT_APP_") {
			result = append(result, obj)
		}
	}
	count := len(result)
	fmt.Printf("doing work for env vars (%d):\n", count)
	for _, x := range result {
		fmt.Printf("- %s\n", x.Key)
	}
	return result
}

func findFiles(searchPath string) []string {
	var result []string

	var searchFunc = func(pathX string, infoX os.FileInfo, errX error) error {
		if errX != nil {
			//log.Warnw("findFiles error",
			//	"path", pathX,
			//	"err", errX,
			//)
			return errX
		}
		if isFile(pathX) {
			result = append(result, pathX)
		}
		return nil
	}

	realPath := getFileAbsolutePath(searchPath)
	err := filepath.Walk(realPath, searchFunc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}

func getFileAbsolutePath(fileName string) (result string) {

	if strings.HasPrefix(fileName, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		fileName = filepath.Join(dir, fileName[2:])
	}

	result, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	return result
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(6)
		return false
	}
	return fileInfo.Mode().IsRegular()
}

func readFile(filePath string) ([]byte, string) {
	r, err := ioutil.ReadFile(filePath)
	stringContents := string(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(7)
	}
	return r, stringContents
}

func changeFile(fileName string, envVars []Env) (result bool) {
	effected := false
	for _, envVar := range envVars {
		change := changeStringInFile(fileName, envVar)
		if change {
			effected = true
		}
	}
	return effected
}

func changeStringInFile(fileName string, envVar Env) bool {
	r, stringContents := readFile(fileName)

	var originalContents string
	changed := false

	err := copier.Copy(&originalContents, &stringContents) // makes a copy of the original string contents
	if err != nil {
		fmt.Println(err)
	}

	key := envVar.Key

	// calculate the trailing part of the key name.. EG REACT_APP_FOO > FOO
	suffix := strings.ReplaceAll(key, "REACT_APP_", "")

	// now we know the string we want to search for and replace..
	oldString := fmt.Sprintf("DEFAULT_VALUE_%s", suffix)

	newString := envVar.Value

	stringContents = strings.ReplaceAll(string(r), oldString, newString)
	if originalContents != stringContents {
		changed = true
	}
	if changed {
		err = ioutil.WriteFile(fileName, []byte(stringContents), 0)
		if err != nil {
			fmt.Println(err)
			os.Exit(10)
		}
	}
	return changed
}
