build:
	@go build -o bin/api .

dev:
	@JWT_SECRET="DBIUAB!!312312ADHI"
	@air \
	 --build.cmd "clear && make -s build" \
	 --build.bin "./bin/api" \
	 --build.exclude_dir "bin,tmp"

test:
	@go test ./tests