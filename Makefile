VERSION=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${V} -X main.BuildTime=${BUILD_TIME}"
NAME=supplier

gen-code:
	protoc --go_out=. --go-grpc_out=. grpc/proto/*.proto

build-docker-img:
	docker build -t $(NAME):dev .

run: build
	./bin/$(NAME) -service $(SER)

build: clear
	go build ${LDFLAGS} -o ./bin/$(NAME) ./main.go
	./bin/$(NAME) -v


clear:
	rm -rf ./bin/$(SER)

clear-untag-image:
	docker rmi -f $(docker images --filter “dangling=true” -q --no-trunc)


# ## Push image

# gcloud auth configure-docker
# docker tag <image-name>:<tag> asia.gcr.io/muulin-universal/<image-name>:<tag>
# docker push asia.gcr.io/muulin-universal/<image-name>:<tag>

# ## get Kubernetes config
# gcloud container clusters get-credentials muulin-gcp-1 --zone=asia-east1

# ## test 
# kubectl get nodes