build:
	./build.sh

run:
	go run procspy.go etc/config.json

run-text:
	go run procspy.go etc/config-test.json	

all:
	./build.sh all

clean:
	rm -rf bin/*