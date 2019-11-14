DOCKER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app:201911142009
CONSUMER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/consumer:201911142009

image:
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile .

consumer:
	docker build -t $(CONSUMER_IMAGE) -f build/Dockerfile.subscriber .

push:
	docker push $(DOCKER_IMAGE)
	docker push $(CONSUMER_IMAGE)
