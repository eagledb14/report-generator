all: run

build:
	GOOS=linux GOARCH=amd64 go build -o form-scanner .
	GOOS=windows GOARCH=amd64 go build -o form-scanner.exe .
	-rm ./resources/event_cache.db
	tar -czf form-scanner.tar.gz ./resources ./form-scanner*
	rm form-scanner
	rm form-scanner.exe

auto:
	@go run . -auto

run: 
	@go run .

clean:
	rm form-scanner*

