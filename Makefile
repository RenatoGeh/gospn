OUT=out.put
ARGS=-p=0.7 -rseed=123456 -clusters=-1 -iterations=1

completion:
	go run main.go -p=0 | tee $(OUT)

completion_single:
	go run main.go -p=0 -concurrents=1 | tee $(OUT)

all:
	go run main.go $(ARGS) | tee $(OUT)

data:
	go run main.go -p=-1 | tee $(OUT)

debug:
	./debug.sh 0

build:
	go build -o gospn

.PHONY: clean
clean:
	rm *.put *.pbm *.ppm *.pgm
