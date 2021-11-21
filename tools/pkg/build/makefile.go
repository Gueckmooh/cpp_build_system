package build

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"regexp"
	"strings"
	"tools/pkg/config"
	"tools/pkg/modules"
	"tools/pkg/options"
	"tools/pkg/utils"
)

const makefileHeader = `# This file has been automatically generated.
# edit this at your own risk.

`
const configMakefileName = "config.mk"

func GetMakefileName(conf *config.Config, module *modules.Module) string {
	filename := strings.TrimSuffix(module.Name, ".xml") + ".mk"
	return path.Join(conf.MakefilesDir, filename)
}

func getLibName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}

func getTarget(module *modules.Module) string {
	var target string
	switch module.Type {
	case "shared_library":
		target += fmt.Sprintf("MODULE_TARGET=lib%s\n", getLibName(module.Name))
	}
	target += fmt.Sprintf("MODULE_TARGET_KIND=%s\n", module.Type)
	return target
}

func getBaseDir(module *modules.Module) string {
	return fmt.Sprintf("MODULE_BASE_DIR=%s\n", module.BaseDir)
}

func getHeadersExportDir(module *modules.Module) string {
	return fmt.Sprintf("MODULE_HEADERS_EXPORT_DIR=%s\n", module.ExportDir)
}

func GetModuleMakefileContent(conf *config.Config, module *modules.Module,
	modBundle *modules.ModuleBundle) string {
	/*
	  MODULE_TARGET=libhello
	  MODULE_TARGET_KIND=shared_library
	  MODULE_TYPE=cpp
	  MODULE_BASE_DIR=hello
	  MODULE_HEADERS_EXPORT_DIR=hello
	*/
	var content string = makefileHeader

	content += `INCLUDE_DEPENDENCY?=0

ifeq ($(INCLUDE_DEPENDENCY),0)
`

	// Compute target
	content += getTarget(module)

	// Compute type
	content += "MODULE_TYPE=cpp\n"

	// Compute dirs
	content += getBaseDir(module)
	content += getHeadersExportDir(module)

	var modDeps string
	var modLibDeps string
	{
		var modDepsL []string
		var modLibDepsL []string
		for _, dep := range module.Dependencies.Dependency {
			m := modBundle.GetModuleByName(dep)
			if m != nil {
				modDepsL = append(modDepsL, m.Name)
				if m.Type != "headers_only" {
					modLibDepsL = append(modLibDepsL, getLibName(m.Name))
				}
			}
		}
		modDeps = strings.Join(modDepsL, " ")
		modLibDeps = strings.Join(modLibDepsL, " ")
	}

	content += fmt.Sprintf("MODULE_DEPENDENCIES=%s\n", modDeps)
	content += fmt.Sprintf("MODULE_LIB_DEPENDENCIES=%s\n", modLibDeps)

	content += "else\n"
	content += fmt.Sprintf(".PHONY: %s_upstream\n", module.Name)
	content += fmt.Sprintf("%s_upstream:\n\t$(MAKE) -C $(SOURCE_DIR)/%s all BUILD_UPSTREAM=1\n",
		module.Name, module.BaseDir)
	content += "endif"

	return content
}

func GenerateModuleMakefile(conf *config.Config, module *modules.Module,
	modBundle *modules.ModuleBundle, filename string) error {
	content := GetModuleMakefileContent(conf, module, modBundle)

	err := utils.Mkdir(path.Dir(filename))
	if err != nil {
		return fmt.Errorf("build.GenerateModuleMakefile: %s", err.Error())
	}

	err = ioutil.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("build.GenerateModuleMakefile: %s", err.Error())
	}

	return nil
}

func GenrateModuleMakefileBundle(conf *config.Config, modBundle *modules.ModuleBundle, thirdParty bool) error {
	for _, module := range modBundle.Modules {
		if module.ThirdParty != thirdParty {
			continue
		}
		filename := GetMakefileName(conf, module)
		err := GenerateModuleMakefile(conf, module, modBundle, filename)
		if err != nil {
			return fmt.Errorf("build.GenrateModuleMakefileBundle: %s", err.Error())
		}
	}
	return nil
}

func GenerateModuleConfigMakefileContent(conf *config.Config, module *modules.Module) string {
	var content string = makefileHeader

	makeIncludeDir := path.Join(conf.SandboxRoot, conf.MakerulesDir)

	content += fmt.Sprintf("MAKE_INCLUDE_DIR=%s", makeIncludeDir)

	return content
}

func GenerateModuleConfigMakefile(conf *config.Config, module *modules.Module) error {
	content := GenerateModuleConfigMakefileContent(conf, module)

	configFileName := path.Join(conf.SandboxRoot, conf.SrcDir, module.BaseDir, configMakefileName)

	fmt.Printf("configFileName = %s\n", configFileName)

	if utils.FileExists(configFileName) && !options.GetOptionBool("overwrite-config-file") {
		return fmt.Errorf("file %s already exists, use --allow-overwrite to overwrite", configFileName)
	}

	if !utils.DirExists(path.Dir(configFileName)) {
		return fmt.Errorf("dir %s does not exist", path.Dir(configFileName))
	}

	err := ioutil.WriteFile(configFileName, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("build.GenerateModuleConfigMakefile: %s", err.Error())
	}

	return nil
}

func GenerateModuleBundleConfigMakefile(conf *config.Config, modBundle *modules.ModuleBundle,
	thirdParty bool) error {
	for _, module := range modBundle.Modules {
		if module.ThirdParty != thirdParty {
			continue
		}
		err := GenerateModuleConfigMakefile(conf, module)
		if err != nil {
			return fmt.Errorf("build.GenerateModuleBundleConfigMakefile: %s", err.Error())
		}
	}
	return nil
}

var pouetRe = regexp.MustCompile(`\$\(([^)]*)\)`)

func createMakefileValue(value string) string {
	if pouetRe.MatchString(value) {
		for _, m := range pouetRe.FindAllStringSubmatch(value, -1) {
			value = strings.ReplaceAll(value, m[1], strings.ToUpper(m[1]))
		}
		return value
	} else {
		return path.Join("$(ROOT)", value)
	}
}

func genConfigMakefileAttributes(conf *config.Config) string {
	var content string
	vConf := reflect.ValueOf(conf)

	for i := 0; i < vConf.Elem().NumField(); i++ {
		if vConf.Elem().Type().Field(i).Tag.Get("type") != "path" ||
			vConf.Elem().Type().Field(i).Tag.Get("dump_to_mk") != "true" {
			continue
		}
		key := vConf.Elem().Type().Field(i).Tag.Get("json")
		value := vConf.Elem().Field(i).String()
		key = strings.ToUpper(key)
		value = createMakefileValue(value)
		if key != "" {
			content += fmt.Sprintf("%s=%s\n", key, value)
		}
	}

	return content
}

func GenerateConfigMakefileContent(conf *config.Config) string {
	var content string = makefileHeader

	content += fmt.Sprintf("ROOT=%s\n\n", conf.SandboxRoot)
	content += genConfigMakefileAttributes(conf)

	content += "\n"

	content += fmt.Sprintf("LIB_DIR:=$(BUILD_DIR)/lib\n")
	content += fmt.Sprintf("BIN_DIR:=$(BUILD_DIR)/bin\n")

	content += "\n"

	content += fmt.Sprintf("OBJS_DIR:=$(BUILD_DIR)/objs\n")
	content += fmt.Sprintf("DEPS_DIR:=$(BUILD_DIR)/deps\n")

	fmt.Println(content)

	return content
}

func GenerateConfigMakefile(conf *config.Config, confToDump *config.Config) error {
	configFileName := path.Join(conf.SandboxRoot, conf.MakerulesDir, configMakefileName)

	if utils.FileExists(configFileName) && !options.GetOptionBool("overwrite-config-file") {
		return fmt.Errorf("file %s already exists, use --allow-overwrite to overwrite", configFileName)
	}

	if !utils.DirExists(path.Dir(configFileName)) {
		return fmt.Errorf("dir %s does not exist", path.Dir(configFileName))
	}

	content := GenerateConfigMakefileContent(confToDump)

	err := ioutil.WriteFile(configFileName, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("build.GenerateModuleConfigMakefile: %s", err.Error())
	}
	return nil
}
