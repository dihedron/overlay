_APPLICATION_NAME := overlay
_APPLICATION_DESCRIPTION := A simple tool to write text on top of an image.
_APPLICATION_COPYRIGHT := 2025 © Andrea Funtò
_APPLICATION_LICENSE := MIT
_APPLICATION_LICENSE_URL := https://opensource.org/license/mit/
_APPLICATION_VERSION_MAJOR := 1
_APPLICATION_VERSION_MINOR := 0
_APPLICATION_VERSION_PATCH := 0
_APPLICATION_VERSION=$(_APPLICATION_VERSION_MAJOR).$(_APPLICATION_VERSION_MINOR).$(_APPLICATION_VERSION_PATCH)
_APPLICATION_MAINTAINER=dihedron.dev@gmail.com
_APPLICATION_VENDOR=dihedron.dev@gmail.com
_APPLICATION_PRODUCER_URL=https://github.com/dihedron/
_APPLICATION_DOWNLOAD_URL=$(_APPLICATION_PRODUCER_URL)$(_APPLICATION_NAME)
_APPLICATION_METADATA_PACKAGE=$$(grep "module .*" go.mod | sed 's/module //gi')/metadata
#_APPLICATION_DOTENV_VAR_NAME=


_GOLANG_MK_FLAG_ENABLE_CGO=0
_GOLANG_MK_FLAG_ENABLE_GOGEN=0
_GOLANG_MK_FLAG_ENABLE_RACE=0
#_GOLANG_MK_FLAG_STATIC_LINK=1
#_GOLANG_MK_FLAG_ENABLE_NETGO=1
#_GOLANG_MK_FLAG_STRIP_SYMBOLS=1
#_GOLANG_MK_FLAG_STRIP_DBG_INFO=1
#_GOLANG_MK_FLAG_FORCE_DEP_REBUILD=1
#_GOLANG_MK_FLAG_OMIT_VCS_INFO=1

include golang.mk
include nfpm.mk
include help.mk
include piped.mk

# Add custom targets below...

#
# compile is the default target; it builds the 
# application for the default platform (linux/amd64)
#
.DEFAULT_GOAL := compile

.PHONY: compile 
compile: linux/amd64 ## build for the default linux/amd64 platform

.PHONY: clean 
clean: golang-clean ## remove all build artifacts

.PHONY: deb
deb: nfpm-deb ## build a DEB package

.PHONY: rpm
rpm: nfpm-rpm ## build a RPM package

.PHONY: apk
apk: nfpm-apk ## build a APK package

.phony: test
test: compile
	@OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=650,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --color=#FFFFFF --input=_test/test.jpg --output=dist/linux/amd64/out.png --text="HALLO, WORLD!"

.phony: test-pipe
test-pipe: compile
	@cat _test/test.jpg | \
	OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=600,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --color=#FFFFFF --format=jpg --text="HALLO, WORLD..." | \
	OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=700,160 --size=48 --font=_test/Economica/Economica-Regular.ttf --color=#00FF0033 --output=dist/linux/amd64/out.jpg --text="... from me!"

