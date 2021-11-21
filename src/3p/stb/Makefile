MODULE_NAME=3p/stb

include config.mk

HEADER_FILES=$(shell find stb -maxdepth 1 -type f -name "*.h" | sed 's_stb/__g')
MODULE_HEADERS_PATH=$(MODULE_PATH)/stb

include $(MAKE_INCLUDE_DIR)/config.mk
include $(MAKE_INCLUDE_DIR)/build_module.mk
