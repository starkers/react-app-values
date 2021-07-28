
BINARY := rav

# -s and -w strip debug headers
build:
	GCO_ENABLED=0 GOOS=linux go build \
	-v \
	-a \
	-installsuffix cgo \
	-ldflags="-X main.version=${VERSION} -s -w" \
	-o ${BINARY}-linux .

test:
	REACT_APP_FOO=bar go run . --dir ./test-data
	@echo
	@grep ^REACT_APP_FOO= ./test-data/changed.txt
	@echo "^ should be  =bar"

