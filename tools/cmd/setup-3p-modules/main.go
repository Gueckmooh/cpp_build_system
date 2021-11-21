package main

import (
	"fmt"
	"os"
	"path"
	"tools/pkg/build"
	"tools/pkg/config"
	"tools/pkg/modules"
	"tools/pkg/options"
	"tools/pkg/utils"
)

func processOptions() {
	for _, option := range os.Args {
		switch option {
		case "--allow-overwrite":
			options.SetOption("overwrite-config-file", true)
		}
	}
}

func tryMain() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	conf, err := config.GetConfig(cwd)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	c := utils.MakeContext()
	newConfig := conf.Expand(c)

	modFiles, err := modules.Get3PModuleFiles(newConfig)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	modBundle, err := modules.ReadModuleBundle(modFiles)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}
	for _, m := range modBundle.Modules {
		m.Sources.Git.Target = path.Join(conf.SrcDir, m.BaseDir, m.Sources.Git.Target)
	}

	err = build.GenrateModuleMakefileBundle(newConfig, modBundle)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	err = build.GenerateModuleBundleConfigMakefile(newConfig, modBundle)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	err = build.GenerateConfigMakefile(newConfig, conf)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	for _, m := range modBundle.Modules {
		err = modules.CloneModuleRepository(m)
		if err != nil {
			return fmt.Errorf("tryMain: %s", err.Error())
		}
	}

	return nil
}

func main() {
	processOptions()
	err := tryMain()
	if err != nil {
		fmt.Printf("main: %s\n", err.Error())
		os.Exit(1)
	}
}