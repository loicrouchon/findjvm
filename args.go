package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	logLevel       string
	configKey      string
	minJavaVersion uint
	maxJavaVersion uint
	vendors        list
	programs       list
}

type list []string

func (i *list) String() string {
	return "[" + strings.Join(*i, ", ") + "]"
}

func (i *list) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func ParseArgs(commandArgs []string) (*Args, error) {
	args := Args{}
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	output := bytes.NewBufferString("")
	cmd.SetOutput(output)
	cmd.StringVar(&args.logLevel, "log-level", "error",
		"Sets the log level to one of: debug, info, warn, error")
	cmd.StringVar(&args.configKey, "config-key", "",
		"If specified, will look for an optional config.<KEY>.json to load before loading the default configuration")
	cmd.UintVar(&args.minJavaVersion, "min-java-version", allVersions,
		"The minimum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.UintVar(&args.maxJavaVersion, "max-java-version", allVersions,
		"The maximum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.Var(&args.vendors, "vendors",
		"The vendors to filter on. If empty, no vendor filtering will be done")
	cmd.Var(&args.programs, "programs",
		"The programs the JVM should provide in its \"${java.home}/bin\" directory. If empty, defaults to java")
	if err := cmd.Parse(commandArgs); err != nil {
		return nil, fmt.Errorf("%s\n%s", err, output)
	}
	if unresolvedArgs := cmd.Args(); len(unresolvedArgs) > 0 {
		cmd.Usage()
		return nil, fmt.Errorf("unresolved arguments: %v\n%s", unresolvedArgs, output)
	}
	if err := setLogLevel(args.logLevel); err != nil {
		return nil, err
	}
	if len(args.programs) == 0 {
		args.programs = append(args.programs, "java")
	}
	return &args, nil
}
