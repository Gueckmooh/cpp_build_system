MODULE_NAME=hello

include config.mk

include $(MAKE_INCLUDE_DIR)/base.mk

ifeq ($(TOTO),1)
SOURCE_FILES=$(MODULE_SOURCE_DIR)/hello.cpp
endif

include $(MAKE_INCLUDE_DIR)/build_module.mk

$(info SOURCE_FILES = $(SOURCE_FILES))
$(info DEP_FILES = $(DEP_FILES))
