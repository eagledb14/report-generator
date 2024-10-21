all: run

build:
ifeq ($(os), win)
	GOOS=windows GOARCH=amd64 go build -o shodan-form.exe .
else ifeq ($(os), linux)
	GOOS=linux GOARCH=amd64 go build -o shodan-form .
endif

build-all:
	GOOS=linux GOARCH=amd64 go build -o shodan-form .
	GOOS=windows GOARCH=amd64 go build -o shodan-form.exe .

run: 
	go run .

clean:
	rm shodan-form*

