package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
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

	log.Infof("Configs:")
	log.Printf("- packages: %s", packages)

	if packages == "" {
		failf("Required input not defined: packages")
	}

	if !installedInPath("errcheck") {
		cmd := command.New("go", "get", "-u", "github.com/kisielk/errcheck").AppendEnvs("GO111MODULE=off")
		cmd.SetDir("/") // workaround for https://github.com/golang/go/issues/30515 to install the package globally

		log.Infof("\nInstalling errcheck")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			failf("Failed to install errcheck: %s", out)
		}
	}

	log.Infof("\nRunning errcheck...")

	for _, p := range strings.Split(packages, "\n") {
		cmd := command.NewWithStandardOuts("errcheck", "-asserts=true", "-blank=true", "-verbose", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("errcheck failed: %s", err)
		}
	}
}
