build:
	./build.sh

run:
	go run cmd/procspy.go etc/config.json

run-test:
	go run cmd/procspy.go etc/config-test.json	

all:
	./build.sh all

clean:
	rm -rf bin/*