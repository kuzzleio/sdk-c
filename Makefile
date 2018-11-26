VERSION = 1.0.0

ROOT_DIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

SDK_FOLDER_NAME=kuzzle-c-sdk

ifeq ($(OS),Windows_NT)
	STATICLIB = .lib
	DYNLIB = .dll
	GOROOT ?= C:/Go
	GOCC ?= $(GOROOT)bin\go
	SEP = \\
	RM = del /Q /F /S
	RRM = rmdir /S /Q
	MV = rename
	CMDSEP = &
	ROOT_DIR_CLEAN = $(subst /,\,$(ROOT_DIR))
	LIB_PREFIX =
else
	STATICLIB = .a
	DYNLIB = .so
	GOROOT ?= /usr/local/go
	GOCC ?= $(GOROOT)/bin/go
	SEP = /
	RM = rm -f
	RRM = rm -f -r
	MV = mv -f
	ROOT_DIR_CLEAN = $(ROOT_DIR)
	LIB_PREFIX = lib
	ARCH=$(shell uname -p)
endif

SDKGOPATH = go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)kuzzleio$(PATHSEP)sdk-go
CGOPATH = cgo$(PATHSEP)kuzzle
PATHSEP = $(strip $(SEP))
BUILD_DIR = $(ROOT_DIR)build
GOFLAGS ?= -buildmode=c-archive
GOFLAGSSHARED = -buildmode=c-shared
GOSRC = .$(PATHSEP)cgo$(PATHSEP)kuzzle$(PATHSEP)
GOTARGET = $(BUILD_DIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(STATICLIB)
GOTARGETSO = $(BUILD_DIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(DYNLIB)

export GOPATH = $(ROOT_DIR)go

CORE_SRC = $(wildcard $(GOSRC)*.go)
all:  $(BUILD_DIR)/sdk

$(BUILD_DIR)/pre_core:
	cd $(SDKGOPATH) && go get .$(PATHSEP)...
	@touch $@

$(BUILD_DIR)/core: $(CORE_SRC) Makefile
ifneq ($(OS),Windows_NT)
ifeq ($(wildcard $(GOCC)),)
	$(error "Unable to find go compiler in $(GOCC)")
endif
	cd $(BUILD_DIR) && rm -f $(GOTARGET).* $(GOTARGETSO).*
endif
ifeq ($(GOOS), android)
	$(GOCC) build -o $(GOTARGET) $(GOFLAGSSHARED) $(GOSRC)
else
	$(GOCC) build -o $(GOTARGET) $(GOFLAGS) $(GOSRC)
endif
	$(GOCC) build -o $(GOTARGETSO) $(GOFLAGSSHARED) $(GOSRC)

	cd $(BUILD_DIR) && mv $(GOTARGET) $(GOTARGET).$(VERSION) && mv $(GOTARGETSO) $(GOTARGETSO).$(VERSION)

ifeq ($(OS),Windows_NT)
	$(MV) $(subst /,\,$(BUILD_DIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk.h) kuzzle.h
else
	$(MV) $(BUILD_DIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk.h $(BUILD_DIR)$(PATHSEP)kuzzle.h
endif
	 @touch $@

$(BUILD_DIR):
ifeq ($(OS),Windows_NT)
	@if not exist $(subst /,\,$(BUILD_DIR)) mkdir $(subst /,\,$(BUILD_DIR))
else
	mkdir -p $@
endif

$(BUILD_DIR)/libs: $(BUILD_DIR) $(BUILD_DIR)/pre_core $(BUILD_DIR)/core
	cd $(BUILD_DIR) && ln -srf $(LIB_PREFIX)kuzzlesdk$(STATICLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(STATICLIB)
	cd $(BUILD_DIR) && ln -srf $(LIB_PREFIX)kuzzlesdk$(DYNLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(DYNLIB)
	@touch $@

$(BUILD_DIR)/sdk: $(BUILD_DIR)/libs
	mkdir -p $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)/lib
	mkdir -p $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)/include/internal
	cp $(BUILD_DIR)$(PATHSEP)kuzzle.h  $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)/include/internal
	cp -fr $(ROOT_DIR)$(PATHSEP)include $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)
	cp -fr $(ROOT_DIR)$(PATHSEP)include $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)
	cp -a $(BUILD_DIR)$(PATHSEP)*.a* $(BUILD_DIR)$(PATHSEP)*.so*  $(BUILD_DIR)$(PATHSEP)$(SDK_FOLDER_NAME)/lib
	@touch $@

package: $(BUILD_DIR)/sdk
	mkdir -p deploy && cd $(BUILD_DIR) && tar cfz ..$(PATHSEP)deploy$(PATHSEP)kuzzlesdk-c-$(VERSION)-$(ARCH).tar.gz $(SDK_FOLDER_NAME)

clean:
ifeq ($(OS),Windows_NT)
	if exist build $(RRM) build
	$(RRM) $(ROOT_DIR)$(PATHSEP)deploy
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)pkg
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)bin
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)gorilla
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)satori
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)stretchr
else
	$(RRM) build
	$(RRM) deploy
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)pkg
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)bin
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)gorilla
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)satori
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)stretchr
endif
.PHONY: all clean 


.DEFAULT_GOAL := all
