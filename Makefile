NAME=tbdmud

build: scripting
	mkdir -p bin && go build -o bin/$(NAME) internal/cmd/Main.go

build-debug: scripting
	mkdir -p bin && go build -o bin/$(NAME) -gcflags "all=-N -l" internal/cmd/Main.go

scripting:
	tsc --project runtime/tsconfig.json 
	cp runtime/dist/engine.js ./internal/scripting
	cp runtime/dist/engine.d.ts ./sample/src

module:
	tsc --project ./sample/src

all: build module

test:
	go test -count=1 ./...

clean:
	rm bin/* || true
	find . -name "*.js" -type f -delete