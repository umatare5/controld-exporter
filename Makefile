.PHONY: image force-image build

bin := controld-exporter
src := $(wildcard *.go)

# Default target
${bin}: Makefile ${src}
	go build -v -o "${bin}"

# Docker targets
image:
	docker build -t ${USER}/controld-exporter .

force-image:
	docker build --no-cache -t ${USER}/controld-exporter .
