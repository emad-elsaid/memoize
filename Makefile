all: test bench

test:
	go test . -count=1 -race

bench:
	go test -benchmem -bench .
