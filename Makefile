default: get fmt vet test install_front_dep clean_front install_front install

test:
	go test -v -timeout 60s -race ./...

vet:
	go vet ./...

fmt:
	@if [ -n "$$(go fmt ./...)" ]; then echo 'Please run go fmt on your code.' && exit 1; fi

get:
	go get -t

install_front_dep:
	@if type "npm" &> /dev/null ; then mkdir -p public/node_modules;npm install --prefix ./public/node_modules; else echo 'Please install NPM: https://www.npmjs.com/'; fi

clean_front:
	rm -rf public/node_modules && rm -rf public/styles && rm -rf public/app && rm -f public/index.html && rm -f public/systemjs.config.js

install_front:
	npm run tsc && cp index.html public && cp systemjs.config.js public && cp -r styles public && cp -r app public/app && cp -r node_modules public/node_modules

install:
	go install

