td-deps-ubuntu:
	sudo apt update && sudo apt install -y cmake build-essential gperf libssl-dev zlib1g-dev libc++-dev libc++abi-dev libbsd-dev

clean-td:
	rm -rf td

build-td: clean-td
	git clone https://github.com/tdlib/td.git
	cd td && mkdir build
	cd td/build && cmake -DCMAKE_BUILD_TYPE=Release .. && cmake --build . -j 6 && sudo make install

build:
	go build --ldflags "-extldflags '-L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd'" -o ./bot main.go

full-build: td-deps-ubuntu build-td
