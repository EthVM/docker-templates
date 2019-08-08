// Stolen from dockerize (with some minor changes): https://github.com/jwilder/dockerize/blob/master/template.go

package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	tpl "text/template"

	"github.com/jwilder/gojq"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func contains(item map[string]string, key string) bool {
	if _, ok := item[key]; ok {
		return true
	}
	return false
}

func defaultValue(args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("default called with no values")
	}

	if len(args) > 0 {
		if args[0] != nil {
			return args[0].(string), nil
		}
	}

	if len(args) > 1 {
		if args[1] == nil {
			return "", fmt.Errorf("default called with nil default value")
		}

		if _, ok := args[1].(string); !ok {
			return "", fmt.Errorf("default is not a string value. hint: surround it w/ double quotes")
		}

		return args[1].(string), nil
	}

	return "", fmt.Errorf("default called with no default value")
}

func parseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		l.Fatalf("unable to parse url %s: %s", rawurl, err)
	}
	return u
}

func add(arg1, arg2 int) int {
	return arg1 + arg2
}

func isTrue(s string) bool {
	b, err := strconv.ParseBool(strings.ToLower(s))
	if err == nil {
		return b
	}
	return false
}

func jsonQuery(jsonObj string, query string) (interface{}, error) {
	parser, err := gojq.NewStringQuery(jsonObj)
	if err != nil {
		return "", err
	}
	res, err := parser.Query(query)
	if err != nil {
		return "", err
	}
	return res, nil
}

func loop(args ...int) (<-chan int, error) {
	var start, stop, step int
	switch len(args) {
	case 1:
		start, stop, step = 0, args[0], 1
	case 2:
		start, stop, step = args[0], args[1], 1
	case 3:
		start, stop, step = args[0], args[1], args[2]
	default:
		return nil, fmt.Errorf("wrong number of arguments, expected 1-3"+
			", but got %d", len(args))
	}

	c := make(chan int)
	go func() {
		for i := start; i < stop; i += step {
			c <- i
		}
		close(c)
	}()
	return c, nil
}

func renderFile(t template) bool {
	tmpl := tpl.New(filepath.Base(t.Dest)).Funcs(tpl.FuncMap{
		"contains":  contains,
		"exists":    exists,
		"split":     strings.Split,
		"replace":   strings.Replace,
		"default":   defaultValue,
		"parseUrl":  parseURL,
		"atoi":      strconv.Atoi,
		"add":       add,
		"isTrue":    isTrue,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"jsonQuery": jsonQuery,
		"loop":      loop,
	})

	tmpl = tmpl.Delims(delims[0], delims[1])

	isSrcAbs := filepath.IsAbs(t.Src)
	if !isSrcAbs {
		var err error
		t.Src, err = filepath.Abs(t.Src)
		if err != nil {
			l.Fatalf("invalid src path passed: %s", err)
		}
	}

	isDestAbs := filepath.IsAbs(t.Dest)
	if t.Dest != "" && !isDestAbs {
		var err error
		t.Dest, err = filepath.Abs(t.Dest)
		if err != nil {
			l.Fatalf("invalid dest path passed: %s", err)
		}
	}

	tmpl, err := tmpl.ParseFiles(t.Src)
	if err != nil {
		l.Fatalf("unable to parse template: %s", err)
	}

	if _, err := os.Stat(t.Dest); err == nil && noOverwriteFlag {
		l.Println("File already exists! Ignoring overwrite")
		return false
	}

	dest := os.Stdout
	if !forceStdOutFlag && t.Dest != "" {
		ensureDir(t.Dest)
		dest, err = os.Create(t.Dest)
		if err != nil {
			l.Fatalf("unable to %s", err)
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(t.Src), &t.LocalVars)
	if err != nil {
		l.Fatalf("template error: %s\n", err)
	}

	if fi, err := os.Stat(t.Dest); err == nil {
		if err := dest.Chmod(fi.Mode()); err != nil {
			l.Fatalf("unable to chmod temp file: %s\n", err)
		}
		if err := dest.Chown(int(fi.Sys().(*syscall.Stat_t).Uid), int(fi.Sys().(*syscall.Stat_t).Gid)); err != nil {
			l.Fatalf("unable to chown temp file: %s\n", err)
		}
	}

	return true
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}
