package sources

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"tools/pkg/config"
	"tools/pkg/modules"
	"tools/pkg/utils"
)

const dummyHeaderContent = `#pragma once

namespace %s {

void hello();

}`

const dummySourceContent = `#include <%s>

#include <iostream>

namespace %s {

void hello() {
std::cout << "Hello" << std::endl;
}

}`

const dummyMakefileContent = `MODULE_NAME=%s

include config.mk

include $(MAKE_INCLUDE_DIR)/config.mk
include $(MAKE_INCLUDE_DIR)/build_module.mk
`

func GenLibModuleRepository(conf *config.Config, module *modules.Module) error {
	moduleDir := path.Join(conf.SrcDir, module.BaseDir)
	if utils.DirExists(moduleDir) {
		return fmt.Errorf("module dir '%s' already exists", moduleDir)
	}

	fmt.Printf("Creating directory '%s'...\n", moduleDir)
	err := utils.Mkdir(moduleDir)
	if err != nil {
		return err
	}

	headerDir := path.Join(moduleDir, "include")
	sourceDir := path.Join(moduleDir, "src")
	headerFile := path.Join(headerDir, "hello.hpp")
	sourceFile := path.Join(sourceDir, "hello.cpp")
	makefileFile := path.Join(moduleDir, "Makefile")

	headerContent := fmt.Sprintf(dummyHeaderContent, strings.ReplaceAll(module.Name, "/", "_"))
	sourceContent := fmt.Sprintf(dummySourceContent,
		path.Join(module.Name, "hello.hpp"), strings.ReplaceAll(module.Name, "/", "_"))
	makefileContent := fmt.Sprintf(dummyMakefileContent, module.Name)

	fmt.Printf("Creating directory '%s'...\n", headerDir)
	err = utils.Mkdir(headerDir)
	if err != nil {
		return err
	}
	fmt.Printf("Creating directory '%s'...\n", sourceDir)
	err = utils.Mkdir(sourceDir)
	if err != nil {
		return err
	}

	fmt.Printf("Writing file '%s'...\n", headerFile)
	err = ioutil.WriteFile(headerFile, []byte(headerContent), 0600)
	if err != nil {
		return err
	}
	fmt.Printf("Writing file '%s'...\n", sourceFile)
	err = ioutil.WriteFile(sourceFile, []byte(sourceContent), 0600)
	if err != nil {
		return err
	}
	fmt.Printf("Writing file '%s'...\n", makefileFile)
	err = ioutil.WriteFile(makefileFile, []byte(makefileContent), 0600)
	if err != nil {
		return err
	}

	return nil
}
