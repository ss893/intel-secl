VERSION := v4.0.0

all:
	cd ../db && rm -rf out && mkdir -p out
	cd ../db/src/flavor-template && go build -ldflags "-X main.BuildVersion=$(VERSION)" -o $(VERSION)-flavor-convert && mv $(VERSION)-flavor-convert ../../out
	$ echo "upgrades/hvs/db/out/$(VERSION)-flavor-convert successfully created"