SHELL=/bin/bash

.PHONY: init
init:
ifeq ($(shell uname -s),Darwin)
	@grep -r -l remonade * .goreleaser.yml | xargs sed -i "" "s/go-cli-template/$$(basename `git rev-parse --show-toplevel`)/"
else
	@grep -r -l remonade * .goreleaser.yml | xargs sed -i "s/go-cli-template/$$(basename `git rev-parse --show-toplevel`)/"
endif

.PHONY: mock
mock:
	@cd e2e && go build -o nature_remo_mock . && ./nature_remo_mock &

.PHONY: stopmock
stopmock:
	@ps -ef | grep nature_remo_mock | grep -v grep | awk '{print $$2}' | xargs kill -9

.PHONY: withmock
withmock: mock
	@DEBUG=1 NATURE_REMO_ENDPOINT=http://localhost:9999 go run main.go

.PHONY: run
run:
	@DEBUG=1 go run main.go

.PHONY: clean
clean:
	@rm -rf e2e/nature_remo_mock
