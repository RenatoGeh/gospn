OUT=out.put
ARGS=-p=0.7 -rseed=123456 -clusters=3 -iterations=1 -mode=class -dataset=digits

completion:
	go run main.go -mode=cmpl | tee $(OUT)

completion_single:
	go run main.go -mode=cmpl -concurrents=1 | tee $(OUT)

classification:
	go run main.go $(ARGS) | tee $(OUT)

data:
	go run main.go -mode=cmpl | tee $(OUT)

debug:
	./debug.sh 0

build:
	go build -o gospn

.PHONY: clean
clean:
	rm *.put *.pbm *.ppm *.pgm test_after_*.spn test_before_*.spn
