package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tools/pkg/config"
	"tools/pkg/git"
)

type Module struct {
	File         string
	Name         string `xml:"name,attr"`
	ThirdParty   bool   `xml:"third_party,attr"`
	Type         string `xml:"type"`
	BaseDir      string `xml:"baseDir"`
	ExportDir    string `xml:"exportDir"`
	Dependencies struct {
		Dependency []string `xml:"dependency"`
	} `xml:"dependencies"`
	Sources struct {
		Git *git.GitRepository `xml:"git"`
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

func (mb *ModuleBundle) GetModuleByName(name string) *Module {
	for _, m := range mb.Modules {
		if m.Name == name {
			return m
		}
	}
	return nil
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

func Get3PModuleFiles(conf *config.Config) (*ModuleFileBundle, error) {
	var filePaths []string
	err := filepath.Walk(conf.ModulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk in %s: %s", conf.ModulesDir, err.Error())
		}

		if info.Mode().IsRegular() {
			if strings.HasSuffix(info.Name(), ".xml") && strings.HasPrefix(path, conf.TPModulesDir) {
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

func CloneModuleRepository(m *Module) error {
	return m.Sources.Git.Clone()
}
