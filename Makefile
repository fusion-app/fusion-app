DOCKER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app:201911241623
CONSUMER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/consumer:201911212110
WATCHER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app-watcher:201911222242

image:
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile .

consumer:
	docker build -t $(CONSUMER_IMAGE) -f build/Dockerfile.subscriber .

watcher:
	docker build -t $(WATCHER_IMAGE) -f build/Dockerfile.watcher .

push:
	docker push $(DOCKER_IMAGE)
	docker push $(CONSUMER_IMAGE)
	docker push $(WATCHER_IMAGE)
