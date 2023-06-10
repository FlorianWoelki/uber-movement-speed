IMAGE_NAME ?= uber-movement-speed
CONTAINER_NAME ?= uber-movement-speed

docker-build:
	docker build --rm -t ${IMAGE_NAME} .

docker-run:
	docker run --name ${CONTAINER_NAME} ${IMAGE_NAME}

install:
	go mod download
	cd ./services/kinesis_data_forwarder && pnpm install

build:
	cd ./services/dynamo_getter && ./build.sh
	cd ./services/preprocessing && ./build.sh
	cd ./services/kinesis_data_forwarder && pnpm run build
	go build -o main ./

run:
	./main
	