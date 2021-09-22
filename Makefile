binary:
	# go build -o ./dist/depth ./cmd/depth
	go build -o ./dist/multi ./cmd/multi
	# go build -o ./dist/single ./cmd/single
	go build -o ./dist/skiplist ./cmd/skiplist
	go build -o ./dist/success-rate ./cmd/success-rate