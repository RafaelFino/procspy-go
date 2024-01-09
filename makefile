build:
	go build -o bin/procspy procspy.go

run:
	go run procspy.go etc/config.json

run-text:
	go run procspy.go etc/config-test.json	

all:
	./build.sh

clean:
	rm -rf bin/*