TEST=.

all: test bench

test:
	go test $(TEST) -count=1 -race

bench:
	go test -benchmem -bench $(TEST)

profile:
	go test -benchmem -cpuprofile=/tmp/profile.out -bench=$(TEST)
	go tool pprof -http=:8080 /tmp/profile.out
