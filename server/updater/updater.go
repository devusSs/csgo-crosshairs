package updater

import (
	"fmt"
	"runtime"
)

var (
	buildVersion = ""
	buildDate    = ""
	buildOS      = runtime.GOOS
	buildARCH    = runtime.GOARCH
	buildGo      = runtime.Version()
	buildMode    = ""
)

func PrintBuildInfo() {
	fmt.Println("Crosshairs - backend for a CSGO crosshair picker and generator.")
	fmt.Println()

	switch buildVersion {
	case "":
		fmt.Printf("Build version: \t\tunset\n")
	default:
		fmt.Printf("Build version: \t\t%s\n", buildVersion)
	}

	switch buildDate {
	case "":
		fmt.Printf("Build date: \t\tunset\n")
	default:
		fmt.Printf("Build date: \t\t%s\n", buildDate)
	}

	fmt.Printf("Build os: \t\t%s\n", buildOS)
	fmt.Printf("Build arch: \t\t%s\n", buildARCH)
	fmt.Printf("Build Go version: \t%s\n", buildGo)

	switch buildMode {
	case "":
		fmt.Printf("Build mode: \t\tunset\n")
	default:
		fmt.Printf("Build mode: \t\t%s\n", buildMode)
	}
}
