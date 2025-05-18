
build:
	go build -ldflags="-s -w" -buildvcs=false -trimpath

build_compressed:
	go build -ldflags="-s -w" -buildvcs=false -trimpath -a
	upx sslchecker

clean:
	rm -fv sslchecker
