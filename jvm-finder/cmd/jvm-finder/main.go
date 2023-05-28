package main

import (
	"fmt"
	"jvm-finder/internal/config"
	"jvm-finder/internal/console"
	"jvm-finder/internal/discovery"
	"jvm-finder/internal/jvm"
	"jvm-finder/internal/log"
	"jvm-finder/internal/rules"
	"jvm-finder/internal/selection"
	"os"
	"path/filepath"
)

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		log.Die(err)
	}
	config, err := config.LoadConfig("/etc/jvm-finder/config.json", args.ConfigKey)
	if err != nil {
		log.Die(err)
	}
	javaExecutables, err := discovery.FindAllJavaExecutables(&config.JvmsLookupPaths)
	if err != nil {
		log.Die(err)
	}
	jvmInfos, err := jvm.LoadJvmsInfos(config.JvmsMetadataCachePath, &javaExecutables)
	if err != nil {
		log.Die(err)
	}
	rules := rules.SelectionRules(config, args.MinJavaVersion, args.MaxJavaVersion, args.Vendors, args.Programs)
	if jvms := selection.Select(rules, &jvmInfos); len(jvms) > 0 {
		jvm := jvms[0]
		selection.LogJvmList("[SELECTED]", jvms[0:1])
		if err := processOutput(args, jvm); err != nil {
			log.Die(err)
		}
	} else {
		log.Die(fmt.Errorf("unable to find a JVM matching requirements %s", rules))
	}
}

func processOutput(args *Args, jvm jvm.Jvm) error {
	if args.OutputMode == outputModeJavaHome {
		console.Writer.Printf("%s\n", jvm.JavaHome)
		return nil
	}
	if args.OutputMode == outputModeBinary {
		for _, program := range args.Programs {
			console.Writer.Printf("%s\n", filepath.Join(jvm.JavaHome, "bin", program))
		}
		return nil
	}
	return fmt.Errorf("unsupported output-mode \"%s\"", args.OutputMode)
}