VERSION := $(shell cat ./VERSION)

release:
	git tag -a $(VERSION) -m "release" -s || true
	git push origin main --tags
.PHONY: release
