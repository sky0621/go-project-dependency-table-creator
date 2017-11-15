package main

import (
	"flag"
	"os"
	"path/filepath"

	"fmt"

	"io/ioutil"

	"strings"

	"bufio"

	"go.uber.org/zap"
)

var dirNames []string

func main() {
	target := flag.String("t", "/tmp/", "Parse Target")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	dirs, err := ioutil.ReadDir(*target)
	if err != nil {
		logger.Error("Failed to read directory",
			zap.String("filepath", *target))
	}

	// ひとまずプロジェクト名を収集
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		dirNames = append(dirNames, dir.Name())
	}
	//fmt.Println(dirNames)

	// 次に各プロジェクトのglide.yamlから依存プロジェクトを探索
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}

		gyamlName := filepath.Join(*target, dir.Name(), "glide.yaml")
		fmt.Println(gyamlName)
		fp, err := os.Open(gyamlName)
		defer fp.Close()
		if err != nil {
			logger.Warn("Failed to open glide.yaml",
				zap.String("filepath", *target),
				zap.String("glide.yaml.name", gyamlName))
			continue
		}

		fmt.Println(gyamlName)
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			// FIXME
		}
	}
	//rootDir := fmt.Sprintf("%v*", *target)
	rootDir := *target
	err = filepath.Walk(rootDir, Apply)
	if err != nil {
		logger.Error("Failed to walk in filepath",
			zap.String("filepath", *target))
	}

	//tmpl := template.Must(template.ParseFiles("tmpl.csv"))
	//buf := &bytes.Buffer{}
	//err = tmpl.Execute(buf, result)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(buf.String())
}

func Apply(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		//fmt.Println(info.Name())
		return filepath.SkipDir
	}

	return nil
}
