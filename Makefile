.PHONY: build
build:
	mage build

.PHONY: buildrace
buildrace:
	mage buildrace

.PHONY: test
test:
	mage test

.PHONY: testrace
testrace:
	mage testrace

.PHONY: install
install:
	mage install

.PHONY: uninstall
uninstall:
	mage uninstall

.PHONY: check
check:
	mage check

.PHONY: fmt
fmt:
	mage fmt
