.SECONDEXPANSION:

-include $(shell find $(DEPS_DIR) -name "*.d" -print 2>/dev/null)

############################## CONFIG ##############################
### COMPUTE OPTIONS
CC:=$(TARGET_ARCH)gcc
CXX:=$(TARGET_ARCH)g++

CXXSTD=-std=c++20

INCLUDE_FLAGS+=-I$(INCLUDE_DIR) -I$(MODULE_SRC_PATH)
WARNING_FLAGS?=-Wall -Wextra
CXXFLAGS+=$(INCLUDE_FLAGS) $(WARNING_FLAGS)
CXXFLAGS+=$(CXXSTD)

ifeq ($(DEBUG),1)
CXXFLAGS+=-g -O0
else
CXXFLAGS+=-O3
endif

ifeq ($(MODULE_TARGET_KIND),shared_library)
CXXFLAGS+=-fPIC									#@todo change that
endif

LDFLAGS+=$(addprefix -L,$(LIB_DIR)) $(addprefix -l,$(MODULE_LIB_DEPENDENCIES))

CXX_MSG=$(ECHO) "\tCXX\t$(shell realpath --relative-to="$(PWD)" $<)"
CXXLD_MSG=$(ECHO) "\tCXXLD\t$(shell realpath --relative-to="$(PWD)" $@)"
GEN_MSG=$(ECHO) "\tGEN\t$(shell realpath --relative-to="$(PWD)" $@)"

### COMPUTE FILES
SOURCE_FILES?=$(shell find $(MODULE_SOURCE_PATH) -name "*.cpp" -print)
SOURCE_FILES:=$(subst $(MODULE_SOURCE_PATH)/,,$(SOURCE_FILES))
MODULE_DEPS_PATH:=$(DEPS_DIR)/$(MODULE_BASE_DIR)
MODULE_OBJS_PATH:=$(OBJS_DIR)/$(MODULE_BASE_DIR)
DEP_FILES:=$(addprefix $(MODULE_DEPS_PATH)/, $(SOURCE_FILES:.cpp=.d))
OBJ_FILES:=$(addprefix $(MODULE_OBJS_PATH)/, $(SOURCE_FILES:.cpp=.o))

ifneq ($(MODULE_TARGET_KIND),executable)
HEADER_FILES:=$(shell find $(MODULE_HEADERS_PATH) -type f -print)
HEADER_FILES:=$(subst $(MODULE_HEADERS_PATH)/,,$(HEADER_FILES))
EXPORTED_HEADER_FILES:=$(addprefix $(HEADERS_EXPORT_PATH)/, $(HEADER_FILES))
endif


############################## BUILD TARGETS ##############################
### EXPORT HEADER FILES
.PHONY: export_headers
export_headers: $(EXPORTED_HEADER_FILES)

$(HEADERS_EXPORT_PATH)/%: $(MODULE_HEADERS_PATH)/% $$(@D)/.f
	$(GEN_MSG)
	$(QAT)$(SCRIPTS_DIR)/export_header $< $@

### DEPS FILES
.PHONY: dep_files
dep_files: $(DEP_FILES)

$(MODULE_DEPS_PATH)/%.d: $(MODULE_SOURCE_PATH)/%.cpp $$(@D)/.f
	$(GEN_MSG)
	$(QAT)$(CXX) $(CXXFLAGS) -MM -MT '$(subst $(MODULE_DEPS_PATH),$(MODULE_OBJS_PATH),$(@:.d=.o))' $< -o $@
	$(QAT)awk -i inplace -f $(SCRIPTS_DIR)/sanitize_deps.awk $@

### OBJECT FILES
.PHONY: object_files
object_files: $(OBJ_FILES)

$(MODULE_OBJS_PATH)/%.o: $(MODULE_SOURCE_PATH)/%.cpp $$(@D)/.f
	$(CXX_MSG)
	$(QAT)$(CXX) $(CXXFLAGS) $< -c -o $@

### TARGET FILE
$(BIN_DIR)/%: $(OBJ_FILES) $$(@D)/.f
	$(CXXLD_MSG)
	$(QAT)$(CXX) -o $@ $(OBJ_FILES) $(LDFLAGS)

$(LIB_DIR)/%.so: $(OBJ_FILES) $$(@D)/.f
	$(CXXLD_MSG)
	$(QAT)$(CXX) -shared -o $@ $(OBJ_FILES) $(LDFLAGS)

##### BUILD TARGETS
.PHONY: prebuild
prebuild: dep_files export_headers

.PHONY: build_prehook
build_prehook:

.PHONY: build
build: build_prehook real_build

.PHONY: real_build
ifeq ($(MODULE_TARGET_KIND),executable)
real_build: $(BIN_DIR)/$(MODULE_TARGET)
endif
ifeq ($(MODULE_TARGET_KIND),shared_library)
real_build: $(LIB_DIR)/$(MODULE_TARGET).so
endif