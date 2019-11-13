DOCKER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app:latest
CONSUMER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/consumer:latest

image:
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile .

consumer:
	docker build -t $(CONSUMER_IMAGE) -f build/Dockerfile.subscriber .

push:
	docker push $(DOCKER_IMAGE)
	docker push $(CONSUMER_IMAGE)
