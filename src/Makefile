OUT=out.put

all:
	go run main.go | tee $(OUT)

debug:
	./debug.sh

test:
	./test.sh

.PHONY: clean
clean:
	rm *.put
