PLATFORMS = linux/ppc64le,linux/arm64,linux/amd64,linux/s390x
ENGINE ?= podman
ORG ?= cloud-bulldozer
REGISTRY ?= quay.io
REG = $(REGISTRY)/$(ORG)
REPOS = etcd-perf nginx

all: build push

build:
	for repo in $(REPOS); do \
		$(ENGINE) build --jobs=4 --platform=$(PLATFORMS) --manifest=$(REG)/$$repo:latest $$repo; \
	done

push:
	for repo in $(REPOS); do \
		$(ENGINE) manifest push $(REG)/$$repo:latest $(REG)/$$repo:latest; \
	done
