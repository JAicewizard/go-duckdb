DUCKDB_REPO=https://github.com/duckdb/duckdb.git
DUCKDB_BRANCH=v1.0.0

.PHONY: install
install:
	go install .

.PHONY: examples
examples:
	go run examples/simple/simple.go

.PHONY: test
test:
	go test -v -race -count=1 .

.PHONY: deps.header
deps.header:
	git clone -b ${DUCKDB_BRANCH} --depth 1 ${DUCKDB_REPO}
	cp duckdb/src/include/duckdb.h duckdb.h

.PHONY: duckdb
duckdb:
	rm -rf duckdb
	git clone -b ${DUCKDB_BRANCH} --depth 1 ${DUCKDB_REPO}

DUCKDB_COMMON_BUILD_FLAGS := BUILD_SHELL=0 BUILD_UNITTESTS=0 DUCKDB_PLATFORM=any

.PHONY: deps.darwin.amd64
deps.darwin.amd64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "darwin" ]; then echo "Error: must run build on darwin"; false; fi
	mkdir -p deps/darwin_amd64

	cd duckdb && \
	CFLAGS="-target x86_64-apple-macos11 -O3" CXXFLAGS="-target x86_64-apple-macos11 -O3" ${DUCKDB_COMMON_BUILD_FLAGS} make bundle-library -j 2
	cp duckdb/build/release/libduckdb_bundle.a deps/darwin_amd64/libduckdb.a

.PHONY: deps.darwin.arm64
deps.darwin.arm64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "darwin" ]; then echo "Error: must run build on darwin"; false; fi
	mkdir -p deps/darwin_arm64

	cd duckdb && \
	CFLAGS="-target arm64-apple-macos11 -O3" CXXFLAGS="-target arm64-apple-macos11 -O3" ${DUCKDB_COMMON_BUILD_FLAGS}  make bundle-library -j 2
	cp duckdb/build/release/libduckdb_bundle.a deps/darwin_arm64/libduckdb.a

.PHONY: deps.linux.amd64
deps.linux.amd64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "linux" ]; then echo "Error: must run build on linux"; false; fi
	mkdir -p deps/linux_amd64

	cd duckdb && \
	CFLAGS="-O3" CXXFLAGS="-O3" ${DUCKDB_COMMON_BUILD_FLAGS} make bundle-library -j 2
	cp duckdb/build/release/libduckdb_bundle.a deps/linux_amd64/libduckdb.a

.PHONY: deps.linux.arm64
deps.linux.arm64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "linux" ]; then echo "Error: must run build on linux"; false; fi
	mkdir -p deps/linux_arm64

	cd duckdb && \
	CC="aarch64-linux-gnu-gcc" CXX="aarch64-linux-gnu-g++" CFLAGS="-O3" CXXFLAGS="-O3" ${DUCKDB_COMMON_BUILD_FLAGS} make bundle-library -j 2
	cp duckdb/build/release/libduckdb_bundle.a deps/linux_arm64/libduckdb.a

.PHONY: deps.freebsd.amd64
deps.windows.amd64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "mingw64_nt-10.0-20348" ]; then echo "Error: must run build on windows"; false; fi
	mkdir -p deps/windows_amd64
	
	# this is just code copied from duckdb and fixed for windows. Would like to not change this, and use `make bundle-library` once its fixed.
	cd duckdb && \
	${DUCKDB_COMMON_BUILD_FLAGS} gmake release -j 2
	cd duckdb/build/release && \
		mkdir -p bundle && \
		cp src/Release/duckdb_static.lib bundle/. && \
		cp third_party/*/Release/duckdb_*.lib bundle/. && \
		cp extension/*/Release/*_extension.lib bundle/.
	cd duckdb/build/release/bundle && \
		find . -name '*.lib' -exec ${AR} -x {} \;
	cd duckdb/build/release/bundle && \
		${AR} cr ../libduckdb_bundle.a *.obj

	mkdir tmp
	mv duckdb/build/release/libduckdb_bundle.a tmp/libduckdb_bundle.a
	cd tmp && ${AR} -x libduckdb_bundle.a
	rm tmp/libduckdb_bundle.a	
	ls tmp
	
	cat cgo_static.go
	sed -i '11s/libduckdb_*.a//' cgo_static.go
	cat cgo_static.go

	num=0; for file in tmp/*.obj; do sed -i '11s/LDFLAGS: /LDFLAGS: -lduckdb_$$num/ ' cgo_static.go; echo $$file+hey; ${AR} cr tmp/libduckdb_$$num.a $$file; num=$$((num+1)); done
	cat cgo_static.go

	ls tmp
	cp tmp/libduckdb_*.a deps/windows_amd64/



	

.PHONY: deps.freebsd.amd64
deps.freebsd.amd64: duckdb
	if [ "$(shell uname -s | tr '[:upper:]' '[:lower:]')" != "freebsd" ]; then echo "Error: must run build on freebsd"; false; fi
	mkdir -p deps/freebsd_amd64

	cd duckdb && \
	CFLAGS="-O3" CXXFLAGS="-O3" ${DUCKDB_COMMON_BUILD_FLAGS} gmake bundle-library -j 2
	cp duckdb/build/release/libduckdb_bundle.a deps/freebsd_amd64/libduckdb.a
