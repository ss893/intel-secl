TEMPLATE_REL_VERSION := v4.0.0

flavor-template:
	rm -rf out && mkdir -p out
	cd src/flavor-template && go build -ldflags "-X main.BuildVersion=$(TEMPLATE_REL_VERSION)" -o $(TEMPLATE_REL_VERSION)_flavor_convert && mv $(TEMPLATE_REL_VERSION)_flavor_convert ../../out
	$ echo "upgrades/hvs/db/out/$(TEMPLATE_REL_VERSION)_flavor_convert successfully created"

all: flavor-template
