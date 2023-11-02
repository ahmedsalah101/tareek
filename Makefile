run: build
	./bin/tareek	
build:
	@ go build -o bin/tareek main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/tareek-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/tareek-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/tareek-windows-386 main.go
live:
	gin -p 3000 run main.go
