package main

import (
	"bytes"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"io/ioutil"

	"strings"

	"bufio"

	"fmt"

	"go.uber.org/zap"
)

var dirNames []string

func main() {
	target := flag.String("t", "/tmp/", "Parse Target")
	domain := flag.String("d", "localhost", "Search Domain")
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

	result.Headers = append(result.Headers, "ProjectName")

	for _, d := range dirNames {
		result.Headers = append(result.Headers, d)
	}

	// 次に各プロジェクトのglide.yamlから依存プロジェクトを探索
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}

		gyamlName := filepath.Join(*target, dir.Name(), "glide.yaml")
		//fmt.Println(gyamlName)
		fp, err := os.Open(gyamlName)
		defer fp.Close()
		if err != nil {
			//logger.Warn("Failed to open glide.yaml",
			//	zap.String("filepath", *target),
			//	zap.String("glide.yaml.name", gyamlName))
			continue
		}

		oneRec := make([]string, len(result.Headers))
		oneRec[0] = dir.Name()

		//fmt.Println(gyamlName)
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.HasPrefix(txt, "- package: ") && strings.Contains(txt, *domain) {
				//fmt.Println(txt)
				txtParts := strings.Split(txt, "/")
				pname := txtParts[len(txtParts)-1]
				//fmt.Println(pname)
				for i, pn := range dirNames {
					if pn == pname {
						oneRec[i+1] = "o"
					} else {
						oneRec[i+1] = "x"
					}
				}
			}
		}

		result.Bodies = append(result.Bodies, oneRec)
	}

	tmpl := template.Must(template.ParseFiles("tmpl.md"))
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, result)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}

type Result struct {
	//Branch   string
	Datetime string
	Headers  []string
	Bodies   [][]string
}

var result = &Result{
	Datetime: time.Now().Format("2006-01-02 15:04"),
	Headers:  []string{},
	Bodies:   [][]string{},
}
