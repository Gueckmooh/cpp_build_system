package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tools/pkg/config"
)

type Action struct {
	Type string `xml:"type,attr"`
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
	Text string `xml:",chardata"`
}

type Module struct {
	File         string
	Name         string `xml:"name,attr"`
	P3           bool   `xml:"p3,attr"`
	Type         string `xml:"type"`
	BaseDir      string `xml:"baseDir"`
	Dependancies struct {
		Dependancy []string `xml:"dependancy"`
	} `xml:"dependancies"`
	Sources struct {
		Git     string `xml:"git"`
		Actions struct {
			Action []Action `xml:"action"`
		} `xml:"actions"`
	} `xml:"sources"`
}

type ModuleFileBundle struct {
	basePath string
	files    []string
}

type ModuleBundle struct {
	basePath string
	Modules  []*Module
}

func GetModuleFiles(conf *config.Config) (*ModuleFileBundle, error) {
	var filePaths []string
	err := filepath.Walk(conf.ModulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk in %s: %s", conf.ModulesDir, err.Error())
		}

		if info.Mode().IsRegular() {
			if strings.HasSuffix(info.Name(), ".xml") && !strings.HasPrefix(path, conf.TPModulesDir) {
				filePaths = append(filePaths, path)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("modules.GetModuleFiles: %s", err.Error())
	}

	replacer := strings.NewReplacer(conf.ModulesDir+"/", "")
	for i, fn := range filePaths {
		fn = replacer.Replace(fn)
		filePaths[i] = fn
	}

	modFiles := new(ModuleFileBundle)
	modFiles.basePath = conf.ModulesDir
	modFiles.files = filePaths

	return modFiles, nil
}
