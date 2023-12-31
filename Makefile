EXE = deb/yj
EXE_DIR := $(shell dirname $(EXE))
VERSION := $(shell git describe --tags --dirty | sed -e 's/v\(.*\)/\1/' -e 's/-g.*//')

.PHONY: all clean deb install

all::
	./build.sh $(VERSION)

deb:
	echo "yj ($(VERSION)) stable; urgency=medium\n\n  * See https://github.com/sclevine/yj/releases\n\n -- Stephen Levine <stephen.levine@gmail.com>  $$(date -R)" >debian/changelog
	dpkg-buildpackage -us -uc -b

clean:
	-rm -rf $(EXE_DIR)

$(EXE):
	mkdir -p $(EXE_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(EXE) .

install: $(EXE)
	install  -m 0755 -d $(DESTDIR)/usr/bin
	install -D -m 0555 $(EXE) $(DESTDIR)/usr/bin/$(shell filename $(EXE))

# vi: ts=8:sw=8:noai:noexpandtab:filetype=make
