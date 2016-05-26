OBSERVER_GO_EXECUTABLE ?= go
VERSION := $(shell git describe --tags)
DIST_DIRS := find * -type d -exec

build:
	cd src; CGO_ENABLED=0 go install -ldflags "-X main.version=${VERSION}" -v ./observer.miclle.com/...

# 安装开发环境
install:
	go get github.com/codegangsta/gin


clean:
	cd src; go clean -i ./observer.miclle.com/...; cd ..; rm -rf pkg bin
	find ./src -name '*.git*' | xargs rm -rf
	find ./src -name '.travis.yml' | xargs rm -rf
	find ./src -name 'gin-bin' | xargs rm -rf

# 检验代码风格
style:
	go vet ./src/observer.miclle.com/...

# 测试
test:
	go get github.com/axw/gocov/gocov
	go get github.com/AlekSi/gocov-xml
	go get gopkg.in/matm/v1/gocov-html
	rm -rf coverage.html coverage.xml
	sh coverage.sh --html

# 启动后端服务
dev: clean-backend
	cd ./src/observer.miclle.com/app; gin -p '9000' -a '9090'
