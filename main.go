package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/kballard/go-shellquote"
)

func installedInPath(name string) bool {
	cmd := exec.Command("which", name)
	outBytes, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(outBytes)) != ""
}

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	packages := os.Getenv("packages")
	additionalArgs := os.Getenv("additional_arguments")

	log.Infof("Configs:")
	log.Printf("- packages: %s", packages)
	log.Printf("- additional_arguments: %s", additionalArgs)

	if packages == "" {
		failf("Required input not defined: packages")
	}

	var args []string
	if len(additionalArgs) > 0 {
		var err error
		args, err = shellquote.Split(additionalArgs)
		if err != nil {
			failf("Invalid additional_arguments: %s", err)
		}
	}



	if !installedInPath("errcheck") {
		cmd := command.New("go", "get", "-u", "github.com/kisielk/errcheck")
		cmd.SetDir("/") // workaround for https://github.com/golang/go/issues/30515 to install the package globally

		log.Infof("\nInstalling errcheck")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			failf("Failed to install errcheck: %s", out)
		}
	}

	log.Infof("\nRunning errcheck...")

	for _, p := range strings.Split(packages, "\n") {
		args = append(args, p)
		cmd := command.NewWithStandardOuts("errcheck", args...)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("errcheck failed: %s", err)
		}
	}
}
