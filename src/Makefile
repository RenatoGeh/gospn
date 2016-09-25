OUT=out.put
ARGS=0.7 123456 -1 1

all:
	go run main.go $(ARGS) | tee $(OUT)

debug:
	./debug.sh

test:
	./test.sh

.PHONY: clean
clean:
	rm *.put *.pbm
