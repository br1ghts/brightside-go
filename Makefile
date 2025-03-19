build:
	go build -o brightside

install: build
	sudo mv brightside /usr/local/bin/

run: install
	brightside

clean:
	rm -f brightside
