NAME=tbdmud

build:
	mkdir -p bin && go build -o bin/$(NAME) internal/cmd/Main.go

test:
	go test -count=1 ./...

clean:
	rm bin/*