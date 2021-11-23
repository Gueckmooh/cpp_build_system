package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"tools/pkg/config"
	"tools/pkg/modules"
	"tools/pkg/options"
	"tools/pkg/sources"
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

	fmt.Println("Reading config file...")
	conf, err := config.GetConfig(cwd)
	if err != nil {
		return fmt.Errorf("tryMain: %s", err.Error())
	}

	c := utils.MakeContext()
	newConfig := conf.Expand(c)

	var name string
	var thirdParty bool
	var modType string
	var baseDir string
	var exportDir string
	var dependencies string

	flag.StringVar(&name, "n", "", "Specify the name of the module, mandatory option")
	flag.BoolVar(&thirdParty, "T", false, "Make the module a third party module")
	flag.StringVar(&modType, "t", "shared_library", "Specify the type of the module, can be 'executable', 'shared_library' or 'headers_only'")
	flag.StringVar(&baseDir, "d", "", "Specify the base directory of the module (default: same as module name)")
	flag.StringVar(&exportDir, "e", "", "Specify the base export directory of the module (default: same as module name)")
	flag.StringVar(&dependencies, "D", "", "Specify the dependencies of the module, comma separated list")

	flag.Parse()

	if name == "" {
		flag.Usage()
		return fmt.Errorf("A name must be provided using -n")
	}
	switch modType {
	case "executable":
	case "shared_library":
	case "headers_only":
		break
	default:
		return fmt.Errorf("Unknown module type: %s", modType)
	}
	if baseDir == "" {
		baseDir = name
	}
	if exportDir == "" {
		exportDir = name
	}

	module := modules.Module{
		File:         path.Join(newConfig.ModulesDir, name+".xml"),
		Name:         name,
		ThirdParty:   thirdParty,
		Type:         modType,
		BaseDir:      baseDir,
		ExportDir:    exportDir,
		Dependencies: strings.Split(dependencies, ","),
	}

	if utils.FileExists(module.File) {
		return fmt.Errorf("Could not create module %s, file '%s' already exists",
			module.Name, module.File)
	}

	err = sources.GenLibModuleRepository(newConfig, &module)
	if err != nil {
		return err
	}

	if !utils.DirExists(path.Dir(module.File)) {
		fmt.Printf("Creating directory '%s'...\n", path.Dir(module.File))
		err = utils.Mkdir(path.Dir(module.File))
		if err != nil {
			return err
		}
	}

	out, err := xml.MarshalIndent(module, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("Writing file '%s'...\n", module.File)
	err = ioutil.WriteFile(module.File, []byte(xml.Header+string(out)), 0600)

	if err != nil {
		return err
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
