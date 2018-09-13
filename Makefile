VERSION = 1.0.0

ROOT_DIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
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
endif

SDKGOPATH = go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)kuzzleio$(PATHSEP)sdk-go
CGOPATH = cgo$(PATHSEP)kuzzle
PATHSEP = $(strip $(SEP))
ROOTOUTDIR = $(ROOT_DIR)build
GOFLAGS = -buildmode=c-archive
GOFLAGSSHARED = -buildmode=c-shared
GOSRC = .$(PATHSEP)cgo$(PATHSEP)kuzzle$(PATHSEP)
GOTARGET = $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(STATICLIB)
GOTARGETSO = $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(DYNLIB)

all: c

core:
ifneq ($(OS),Windows_NT)
ifeq ($(wildcard $(GOCC)),)
	$(error "Unable to find go compiler")
endif
endif
	cd $(SDKGOPATH) && go get .$(PATHSEP)...
	$(GOCC) build -o $(GOTARGET) $(GOFLAGS) $(GOSRC)
	$(GOCC) build -o $(GOTARGETSO) $(GOFLAGSSHARED) $(GOSRC)
ifeq ($(OS),Windows_NT)
	$(MV) $(subst /,\,$(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk.h) kuzzle.h
else
	$(MV) $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk.h $(ROOTOUTDIR)$(PATHSEP)kuzzle.h
endif

makedir:
ifeq ($(OS),Windows_NT)
	@if not exist $(subst /,\,$(ROOTOUTDIR)) mkdir $(subst /,\,$(ROOTOUTDIR))
else
	mkdir -p $(ROOTOUTDIR)
endif

c: export GOPATH = $(ROOT_DIR)go
c: makedir core
	 cd $(ROOTOUTDIR) && mv $(GOTARGET) $(GOTARGET).$(VERSION) && mv $(GOTARGETSO) $(GOTARGETSO).$(VERSION)
	 cd $(ROOTOUTDIR) && ln -sr $(LIB_PREFIX)kuzzlesdk$(STATICLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(STATICLIB)
	 cd $(ROOTOUTDIR) && ln -sr $(LIB_PREFIX)kuzzlesdk$(DYNLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(DYNLIB)

clean:
ifeq ($(OS),Windows_NT)
	if exist build $(RRM) build
else
	$(RRM) build
endif
.PHONY: all c core clean


.DEFAULT_GOAL := all
