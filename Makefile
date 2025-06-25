PWD=$(shell pwd)
ALSA_PREFIX=$(PWD)/alsa-local

build: alsa-lib-1.2.8
	cd alsa-lib-1.2.8 && \
	./configure --enable-static --disable-shared --prefix=$(ALSA_PREFIX)
	cd alsa-lib-1.2.8 && \
	make -j32 && make install
	CGO_CFLAGS="-I$(ALSA_PREFIX)/include" \
	CGO_LDFLAGS="-L$(ALSA_PREFIX)/lib -static" \
	go build -a -ldflags '-extldflags "-static"' -v

alsa-lib-1.2.8:
	wget ftp://ftp.alsa-project.org/pub/lib/alsa-lib-1.2.8.tar.bz2
	tar -xjf alsa-lib-1.2.8.tar.bz2

clean:
	rm -rf alsa-lib-1.2.8 alsa-local
	rm -f alsa-lib-1.2.8.tar.bz2

.PHONY: build clean
serve:
	cd docs && npm install
	cd docs && npx vitepress dev