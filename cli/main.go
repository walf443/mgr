package main

import (
	"flag"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/walf443/mgr/sqlparser/mysql"
	"github.com/walf443/mgr/diff"
	"io/ioutil"
	"os"
)

func main() {
	var beforeFile = flag.String("before", "", "before schema filename")
	var afterFile  = flag.String("after", "",  "after schema filename")
	flag.Parse()
	if os.Getenv("DEBUG") == "" {
		pp.SetDefaultOutput(ioutil.Discard)
	}

	if *beforeFile == "" || *afterFile == "" {
		fmt.Fprintf(os.Stderr, "-before or -after are missing\n")
		os.Exit(1)
	}

	beforeSchema := ""
	var err error
	if *beforeFile == "stdin" {
		beforeSchema,  err = loadStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load file: %q\n", err)
			os.Exit(1)
		}
	} else {
		beforeSchema, err = loadFile(*beforeFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load file: %q\n", err)
			os.Exit(1)
		}
	}
	afterSchema, err := loadFile(*afterFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load file: %q\n", err)
		os.Exit(1)
	}

	beforeStmts, err := parseSchema(beforeSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse file: %s\n%s\n", *beforeFile, err)
		os.Exit(1)
	}
	afterStmts, err := parseSchema(afterSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse file: %s\n%s\n", *afterFile, err)
		os.Exit(1)
	}
	result := diff.Extract(beforeStmts, afterStmts)
	for _, stmt := range(result.Changes()) {
		fmt.Println(stmt)
	}
	pp.Print(result)
}

func loadFile(fname string) (string, error) {
	result, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func loadStdin() (string, error) {
	result, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func parseSchema(schema string) ([]mysql.Statement, error) {
	s := new(mysql.Scanner)
	s.Init(schema)
	return mysql.Parse(s)
}
