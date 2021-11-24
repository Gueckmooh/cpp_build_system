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

ifneq ($(MODULE_TARGET_KIND),headers_only)
.DEFAULT_GOAL := build
else
.DEFAULT_GOAL := prebuild
endif

ifneq ($(MODULE_TARGET_KIND),headers_only)
all: prebuild build check
else
all: prebuild check
endif

ifeq ($(OS),windows)
include $(MAKERULES_DIR)/windows.mk
endif

ifeq ($(MODULE_TYPE),cpp)
include $(MAKERULES_DIR)/cpp.mk
endif
ifeq ($(MODULE_TYPE),headers)
include $(MAKERULES_DIR)/headers.mk
endif

include $(MAKERULES_DIR)/build_upstream.mk

include $(MAKERULES_DIR)/unittest.mk
