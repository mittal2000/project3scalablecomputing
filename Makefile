all:
	cd cmd; go build -o ../build/cmd

test:
	./build/cmd --config "app.json"

run: all test