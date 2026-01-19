package applycfg

import (
	"io"
	"fmt"
	"os"
	"path/filepath"
)

func ApplyCfg(routerID, gnsProjectName string) {
	// routerID: example: R4 => routerID = "4"
	home, err := os.UserHomeDir()
	filename := "i" + routerID + "_startup-config.cfg"

	if err != nil {
		panic(err)
	}

	root := filepath.Join(
		home,
		"GNS3",
		"projects",
		gnsProjectName,
		"project-files",
		"dynamips",
	)
	fmt.Println(root)

	pattern := filepath.Join(root, "*", "configs", filename)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	if len(matches) == 0 {
		panic("R4.cfg not found")
	}

	wd, _ := os.Getwd()
	cfgFilename := "R" + routerID + "_configs_i" + routerID + "_startup-config.cfg"

	cfgFilepath := filepath.Join(
		wd,
		cfgFilename,
	)
	
	for _, pathToFile := range matches {
		rewriteFile(pathToFile, cfgFilepath)
	}
}

func rewriteFile(pathToFile string, cfgFilepath string) {
	wd, _ := os.Getwd()
	fmt.Println("WD =", wd)
	fmt.Println("Trying to open:", cfgFilepath)

	cfgf, err := os.Open(cfgFilepath)
	if err != nil {
		panic("rewriteFile(): error when opening cfg file")
	}
	defer cfgf.Close()

	newfile, err := os.Create(pathToFile)

	if err != nil {
		panic("rewriteFile(): error when overwriting file")
	}
	defer newfile.Close()

	_, err = io.Copy(newfile, cfgf)
	if err != nil {
		panic("rewriteFile(): error when copying cfg file to new file")
	}
}
