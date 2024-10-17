PLATFORMS = linux/amd64,linux/ppc64le,linux/arm64,linux/s390x
ENGINE ?= podman
ORG ?= cloud-bulldozer
REGISTRY ?= quay.io
REG = $(REGISTRY)/$(ORG)
REPOS = perfapp  etcd-perf nginx frr netpol-scraper nginxecho eipvalidator sampleapp netpolvalidator netpolproxy convergencetracker

all: build push

build:
	@for repo in $(REPOS); do \
	  echo -e "\033[2mBuilding $$repo\033[0m"; \
	  $(ENGINE) build --jobs=4 --platform=$(PLATFORMS) --manifest=$(REG)/$$repo:latest $$repo; \
	done

push:
	for repo in $(REPOS); do \
	  echo -e "\033[2mPushing $$repo\033[0m"; \
	  $(ENGINE) manifest push $(REG)/$$repo:latest $(REG)/$$repo:latest; \
	done
