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
GOFLAGS ?= -buildmode=c-archive
GOFLAGSSHARED = -buildmode=c-shared
GOSRC = .$(PATHSEP)cgo$(PATHSEP)kuzzle$(PATHSEP)
GOTARGET = $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(STATICLIB)
GOTARGETSO = $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(DYNLIB)

all: c

pre_core:
	cd $(SDKGOPATH) && go get .$(PATHSEP)...

core:
ifneq ($(OS),Windows_NT)
ifeq ($(wildcard $(GOCC)),)
	$(error "Unable to find go compiler")
endif
endif
ifeq ($(GOOS), android)
	$(GOCC) build -o $(GOTARGET) $(GOFLAGSSHARED) $(GOSRC)
else
	$(GOCC) build -o $(GOTARGET) $(GOFLAGS) $(GOSRC)
endif
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
c: makedir pre_core core
	 cd $(ROOTOUTDIR) && mv $(GOTARGET) $(GOTARGET).$(VERSION) && mv $(GOTARGETSO) $(GOTARGETSO).$(VERSION)
	 cd $(ROOTOUTDIR) && ln -sr $(LIB_PREFIX)kuzzlesdk$(STATICLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(STATICLIB)
	 cd $(ROOTOUTDIR) && ln -sr $(LIB_PREFIX)kuzzlesdk$(DYNLIB).$(VERSION) $(LIB_PREFIX)kuzzlesdk$(DYNLIB)

package: $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(STATICLIB).$(VERSION) $(ROOTOUTDIR)$(PATHSEP)$(LIB_PREFIX)kuzzlesdk$(DYNLIB).$(VERSION)
	mkdir $(ROOTOUTDIR)$(PATHSEP)lib
	mkdir $(ROOTOUTDIR)$(PATHSEP)include
	cp -fr $(ROOT_DIR)$(PATHSEP)include$(PATHSEP)*.h $(ROOTOUTDIR)$(PATHSEP)include
	cp $(ROOTOUTDIR)$(PATHSEP)*.so  $(ROOTOUTDIR)$(PATHSEP)lib
	cp $(ROOTOUTDIR)$(PATHSEP)*.a  $(ROOTOUTDIR)$(PATHSEP)lib
	mkdir deploy && cd $(ROOTOUTDIR) && tar cfz ..$(PATHSEP)deploy$(PATHSEP)kuzzlesdk-$(VERSION)-$(ARCH).tar.gz lib include

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
	$(RRM) $(ROOT_DIR)$(PATHSEP)deploy
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)pkg
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)bin
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)gorilla
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)satori
	$(RRM) $(ROOT_DIR)$(PATHSEP)go$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)stretchr
endif
.PHONY: all c core clean


.DEFAULT_GOAL := all
