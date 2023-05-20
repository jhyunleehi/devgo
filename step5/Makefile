# SO Version Information
SO_VERSION_INFO = "-X 'main.version=$$(git describe --tags --abbrev=0)' -X 'main.date=$$(date +%F)' -X 'main.commit=$$(git rev-parse --short HEAD)'"

# Image URL to use all building/pushing image targets
SO_VER = "$$(git describe --tags --abbrev=0)"
IMG_VER = "$$(echo $(SO_VER)| sed -e s/v//g)"
IMG_NAME = docker.io/jhyunleehi/gtrend:${IMG_VER}

# Build th SO main binary
main:	
	go build -mod vendor --ldflags ${SO_VERSION_INFO} main.go

# Build the docker image
docker-build:
	docker build . -t ${IMG_NAME}

# Push the docker image
docker-push:
	docker push ${IMG_NAME}

docker-clean:
	docker rmi -f "$$(docker images | grep "none" | awk '{print $3}')"
	docker rmi -f "$$(docker images | grep "sos" | awk '{print $3}')"



# Delete files
clean:
	rm -f main
	rm -f main.exe

# Default
all: main