OUT=out.put
CLUSTERS=2

all:
	go run main.go $(CLUSTERS) | tee $(OUT)

debug:
	./debug.sh

test:
	./test.sh

.PHONY: clean
clean:
	rm *.put
