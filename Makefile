NAME := overlay
DESCRIPTION := A simple tool to write text on top of an image.
COPYRIGHT := 2024 © Andrea Funtò
LICENSE := MIT
LICENSE_URL := https://opensource.org/license/mit/
VERSION_MAJOR := 0
VERSION_MINOR := 0
VERSION_PATCH := 2
VERSION=$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)
MAINTAINER=dihedron.dev@gmail.com
VENDOR=dihedron.dev@gmail.com
PRODUCER_URL=https://github.com/dihedron/
DOWNLOAD_URL=$(PRODUCER_URL)overlay

METADATA_PACKAGE=$$(grep "module .*" go.mod | sed 's/module //gi')/version

_RULES_MK_MINIMUM_VERSION=202412061025
#_RULES_MK_ENABLE_CGO=0
#_RULES_MK_ENABLE_GOGEN=0
#_RULES_MK_ENABLE_RACE=1
#_RULES_MK_STATIC_LINK=0
#_RULES_MK_ENABLE_NETGO=0
#_RULES_MK_STRIP_SYMBOLS=1
#_RULES_MK_STRIP_DBG_INFO=1
#_RULES_MK_FORCE_DEP_REBUILD=0

include rules.mk

.phony: test
test: compile
	@OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=650,100 --size=72 --font=${HOME}/.fonts/Economica/Economica-Regular.ttf --color=#FFFFFF --input=_test/test.jpg --output=dist/linux/amd64/out.png --text="HALLO, WORLD!"

.phony: test-pipe
test-pipe: compile
	@cat _test/test.jpg | \
	OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=600,100 --size=72 --font=${HOME}/.fonts/Economica/Economica-Regular.ttf --color=#FFFFFF --format=jpg --text="HALLO, WORLD..." | \
	OVERLAY_LOG_LEVEL=d dist/linux/amd64/overlay --point=700,160 --size=48 --font=${HOME}/.fonts/Economica/Economica-Regular.ttf --color=#00FF0033 --output=dist/linux/amd64/out.jpg --text="... from me!"
