.SUFFIXES:

TEST_USER := $(shell echo ${TEST_USER})

BUILD_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SCRIPT_DIR = $(BUILD_DIR)/scripts
export PATH := $(shell echo $(SCRIPT_DIR):$${PATH})

PACKAGE_NAME = https-user-management
PACKAGE_FILE = $(PACKAGE_NAME).deb
PACKAGE_CONTROL = $(PACKAGE_NAME)/DEBIAN/control

INSTALL_TARGETS := \
					/etc/pam.d/pam_https \
					/lib/x86_64-linux-gnu/security/pam_https.so \
					/lib/x86_64-linux-gnu/libnss_https.so.2

BUILD_TARGETS := $(PACKAGE_NAME)/lib/x86_64-linux-gnu/security/pam_https.so \
				 $(PACKAGE_NAME)/lib/x86_64-linux-gnu/libnss_https.so.2

include templates.mk

default: package

.PHONY: changelog
changelog:
	@./scripts/git-changelog > $@

$(PACKAGE_CONTROL): templates/$(PACKAGE_CONTROL)
	@: $(call render_template,$<,$@)

$(PACKAGE_FILE): $(BUILD_TARGETS) $(PACKAGE_CONTROL)
	@dpkg-deb -v --build --root-owner-group $(PACKAGE_NAME)
	@echo >&2 "Package Info:"
	@dpkg-deb -v --info $@
	@echo >&2 "Package Contents:"
	@dpkg-deb -v --contents $@

$(PACKAGE_NAME)/lib/x86_64-linux-gnu/security/pam_https.so: $(wildcard pam-https/*.go)
	@mkdir -p $(@D)
	@go build -buildmode=c-shared -o $@ $^

$(PACKAGE_NAME)/lib/x86_64-linux-gnu/libnss_https.so.2: $(wildcard nss-https/*.go)
	@mkdir -p $(@D)
	@CGO_CFLAGS="-g -O2 -D __LIB_NSS_NAME=https" go build --buildmode=c-shared -o $@ $^

/lib/%: $(PACKAGE_NAME)/lib/%
	@sudo cp $< $@

/etc/pam.d/%: $(PACKAGE_NAME)/etc/pam.d/%
	@sudo cp $< $@

build: $(BUILD_TARGETS)

install: $(INSTALL_TARGETS)

package: $(PACKAGE_FILE)

integrate: install
	@sudo getent passwd $(TEST_USER)
	@sudo getent shadow $(TEST_USER)
	@sudo pamtester -v -I rhost=localhost pam_https $(TEST_USER) authenticate

clean:
	@rm -rf $(BUILD_TARGETS) $(PACKAGE_FILE) $(PACKAGE_CONTROL)
