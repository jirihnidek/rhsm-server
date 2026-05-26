.ONESHELL:
.SHELLFLAGS := -e -c

VERSION := $(shell rpmspec rhsm-server.spec --query --queryformat '%{version}')
LDFLAGS := -ldflags "-X github.com/jirihnidek/rhsm-server/pkg/version.Version=$(VERSION)"
GO_BUILD := go build $(LDFLAGS)

# The 'build' target is not used during packaging; it is present for upstream development purposes.
.PHONY: build
build:
	$(GO_BUILD) -o rhsm-server ./cmd/rhsm-server

.PHONY: archive
archive:
	git archive --prefix rhsm-server-$(VERSION)/ --format tar.gz HEAD > rhsm-server-$(VERSION).tar.gz
	go_vendor_archive create --output rhsm-server-$(VERSION)-vendor.tar.bz2 .

.PHONY: srpm
srpm: archive
	rpmbuild --define "_sourcedir $$(pwd)" -bs rhsm-server.spec

.PHONY: rpm
rpm: archive
	rpmbuild --define "_sourcedir $$(pwd)" -bb rhsm-server.spec

# The 'clean' target removes build artifacts.
.PHONY: clean
clean:
	rm -f rhsm-server
	rm -f rhsm-server-*.tar*