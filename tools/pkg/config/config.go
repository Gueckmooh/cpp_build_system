package config

import (
	"reflect"
	"tools/pkg/utils"
)

// config.json file:
// {
//     "source_dir":"src",
//     "config_dir":"config",
//     "build_dir":"build",
//     "modules_dir":"$(config_dir)/modules",
//     "3p_modules_dir":"$(modules_dir)/3p",
//     "makefiles_dir":"$(config_dir)/makefiles",
//     "makerules_dir":"$(config_dir)/makerules",
//     "scripts_dir":"scripts",
//     "include_dir":"$(build_dir)/include"
// }

type Config struct {
	SandboxRoot  string
	SrcDir       string `json:"source_dir" type:"path" dump_to_mk:"true"`
	ConfigDir    string `json:"config_dir" type:"path" dump_to_mk:"true"`
	BuildDir     string `json:"build_dir" type:"path" dump_to_mk:"true"`
	ModulesDir   string `json:"modules_dir" type:"path"`
	TPModulesDir string `json:"3p_modules_dir" type:"path"`
	MakefilesDir string `json:"makefiles_dir" type:"path" dump_to_mk:"true"`
	MakerulesDir string `json:"makerules_dir" type:"path" dump_to_mk:"true"`
	ScriptsDir   string `json:"scripts_dir" type:"path" dump_to_mk:"true"`
	IncludeDir   string `json:"include_dir" type:"path" dump_to_mk:"true"`
}

func (c *Config) Expand(context *utils.Context) *Config {
	newConfig := new(Config)
	*newConfig = *c
	vConfig := reflect.ValueOf(newConfig)

	nvalid := 0
	for nvalid < vConfig.Elem().NumField() {
		for i := 0; i < vConfig.Elem().NumField(); i++ {
			value := vConfig.Elem().Field(i).String()
			newValue, ok := context.Expand(value)
			if ok {
				nvalid++
				context.SetValue(vConfig.Elem().Type().Field(i).Tag.Get("json"), newValue)
			}
			vConfig.Elem().Field(i).SetString(newValue)
		}
		if nvalid != vConfig.Elem().NumField() {
			nvalid = 0
		}
	}
	return newConfig
}