.PHONY: test release clean

packagebeat:
	go build

test:
	go test ./...

clean:
	go clean
	rm -f packagebeat packagebeat-*.tar.gz

VERSION=$(shell grep 'var Version' main.go | sed 's/.*"\([^"]*\)"/\1/')
release: clean packagebeat test
	tar -c --transform 's,^,packagebeat-$(VERSION)/,' \
      -zvf packagebeat-$(VERSION)-x86_64.tar.gz \
	  packagebeat packagebeat.template.json packagebeat.yml
