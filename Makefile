.SUFFIXES:

TEST_USER := $(shell echo ${TEST_USER})

INSTALL_TARGETS := \
					/etc/pam.d/pam_https \
					/lib/x86_64-linux-gnu/security/pam_https.so \
					/lib/x86_64-linux-gnu/libnss_https.so.2

default: integrate

clean:
	@sudo rm -rf $(INSTALL_TARGETS)

FORCE:

modules/pam_https.so: $(wildcard pam-https/*.go)
	@mkdir -p $(@D)
	@go build -buildmode=c-shared -o $@ $^

modules/libnss_%.so.2: $(wildcard nss-https/*.go)
	@CGO_CFLAGS="-g -O2 -D __LIB_NSS_NAME=$*" go build --buildmode=c-shared -o $@ $^

/lib/x86_64-linux-gnu/%: modules/%
	@sudo mv $< $@

/lib/x86_64-linux-gnu/security/%.so: modules/%.so
	@sudo mv $< $@

/etc/pam.d/%: pam.d/%
	@sudo cp $< $@

install: $(INSTALL_TARGETS)

integrate: install
	@sudo getent passwd $(TEST_USER)
	@sudo pamtester -v -I rhost=localhost pam_https $(TEST_USER) authenticate
