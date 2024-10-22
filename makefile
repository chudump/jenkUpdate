default: compile run

compile:
	go build

run:
	./bbjenk

clean:
	rm ./bbjenk