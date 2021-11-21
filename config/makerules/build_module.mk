.SECONDEXPANSION:

include $(MAKERULES_DIR)/common.mk

SCRIPTS_DIR:=$(ROOT)/scripts

include $(MAKE_INCLUDE_DIR)/config.mk
include $(MAKEFILES_DIR)/$(MODULE_NAME).mk

.PHONY: dependencies
dependencies: $(MODULE_DEPENDENCIES)

MODULE_PATH?=$(SOURCE_DIR)/$(MODULE_BASE_DIR)
MODULE_SOURCE_PATH?=$(MODULE_PATH)/src
MODULE_HEADERS_PATH?=$(MODULE_PATH)/include
HEADERS_EXPORT_PATH?=$(INCLUDE_DIR)/$(MODULE_HEADERS_EXPORT_DIR)

.DEFAULT_GOAL := build

ifeq ($(OS),windows)
include $(MAKERULES_DIR)/windows.mk
endif

ifeq ($(MODULE_TYPE),cpp)
include $(MAKERULES_DIR)/cpp.mk
endif
ifeq ($(MODULE_TYPE),headers)
include $(MAKERULES_DIR)/headers.mk
endif
