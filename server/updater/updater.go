package updater

import (
	"fmt"
	"runtime"
)

var (
	BuildVersion = ""
	BuildDate    = ""
	BuildOS      = runtime.GOOS
	BuildARCH    = runtime.GOARCH
	BuildGo      = runtime.Version()
	BuildMode    = ""
)

func PrintBuildInfo() {
	fmt.Println("Crosshairs - backend for a Counter-Strike crosshair picker, generator and saver.")
	fmt.Println()

	switch BuildVersion {
	case "":
		fmt.Printf("Build version: \t\tunset\n")
	default:
		fmt.Printf("Build version: \t\t%s\n", BuildVersion)
	}

	switch BuildDate {
	case "":
		fmt.Printf("Build date: \t\tunset\n")
	default:
		fmt.Printf("Build date: \t\t%s\n", BuildDate)
	}

	fmt.Printf("Build os: \t\t%s\n", BuildOS)
	fmt.Printf("Build arch: \t\t%s\n", BuildARCH)
	fmt.Printf("Build Go version: \t%s\n", BuildGo)

	switch BuildMode {
	case "":
		fmt.Printf("Build mode: \t\tunset\n")
	default:
		fmt.Printf("Build mode: \t\t%s\n", BuildMode)
	}
}
