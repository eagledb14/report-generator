all: run

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -o report-generator .
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"  -o report-generator.exe .
	-rm ./resources/event_cache.db
	zip -r report-generator.zip ./resources ./report-generator*
	rm report-generator
	rm report-generator.exe

auto:
	@go run . -auto

run: 
	@go run .

clean:
	rm report-generator*

