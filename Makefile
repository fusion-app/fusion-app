DOCKER_IMAGE=registry.njuics.cn/fusion-app/fusion-app:201911131930
CONSUMER_IMAGE=registry.njuics.cn/fusion-app/consumer:201911131930

image:
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile .

consumer:
	docker build -t $(CONSUMER_IMAGE) -f build/Dockerfile.subscriber .

push:
	docker push $(DOCKER_IMAGE)
	docker push $(CONSUMER_IMAGE)
