package main

import (
	"bytes"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"strings"

	"bufio"

	"regexp"

	"fmt"

	"strconv"

	"go.uber.org/zap"
)

var dirNames []string
var domain = flag.String("d", "github.com", "Search Domain")

var summaries []*Summary

type Summary struct {
	baseProject string
	useProjects []string
}

type Header struct {
	Display string
	URL     string
}

type Result struct {
	//Branch   string
	Datetime string
	Headers  []Header
	Bodies   [][]string
}

var result = &Result{
	Datetime: time.Now().Format("2006-01-02 15:04"),
	Headers:  []Header{},
	Bodies:   [][]string{},
}

func main() {
	target := flag.String("t", "./samples", "Parse Target")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	err = filepath.Walk(*target, Apply)
	if err != nil {
		logger.Error("Failed to walk", zap.String("error", err.Error()))
		os.Exit(-1)
	}

	// summary = baseProject + useProjects
	for _, s := range summaries {
		seps := strings.Split(s.baseProject, "/")
		result.Headers = append(
			result.Headers,
			Header{Display: seps[len(seps)-1], URL: s.baseProject},
		)
	}

	for idx, s := range summaries {
		seps := strings.Split(s.baseProject, "/")
		body := []string{strconv.Itoa(idx + 1), seps[len(seps)-1]}
		for _, h := range result.Headers {
			var isHit bool = false
			for _, u := range s.useProjects {
				if h.URL == u {
					isHit = true
				}
			}
			if isHit {
				body = append(body, "o")
			} else {
				body = append(body, "-")
			}
		}
		result.Bodies = append(result.Bodies, body)
	}

	result.Headers = append(
		[]Header{
			Header{Display: "No"},
			Header{Display: "Projects"},
		},
		result.Headers...)

	tmpl := template.Must(template.ParseFiles("tmpl.md"))
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, result)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}

func Apply(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !filter(path, info) {
		return nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if fp != nil {
			fp.Close()
		}
	}()

	s := &Summary{}
	useProjects := []string{}

	// scan glide.yaml
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "package: ") {
			s.baseProject = strings.Replace(line, "package: ", "", -1)
		}

		if strings.HasPrefix(line, "- package: ") {
			if strings.Contains(line, *domain) {
				useProjects = append(useProjects, strings.Replace(line, "- package: ", "", -1))
			}
		}
	}

	s.useProjects = useProjects
	summaries = append(summaries, s)

	return nil
}

func filter(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	outDirExp, err := regexp.Compile("vendor")
	if err != nil {
		return false
	}
	if outDirExp.MatchString(absPath) {
		return false
	}

	inFileExp, err := regexp.Compile("glide.yaml")
	if err != nil {
		return false
	}
	if !inFileExp.MatchString(path) {
		return false
	}

	return true
}
