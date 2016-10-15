OUT=out.put
ARGS=0.7 123456 -1 1

completion:
	go run main.go 0 | tee $(OUT)

all:
	go run main.go $(ARGS) | tee $(OUT)

data:
	go run main.go -1 | tee $(OUT)

debug:
	./debug.sh 0

test:
	./test.sh

.PHONY: clean
clean:
	rm *.put *.pbm *.ppm *.pgm
