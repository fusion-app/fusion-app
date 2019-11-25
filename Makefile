BACKEND_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app-backend:201911251306
OPERATOR_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app-operator:201911251306
CONSUMER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app-consumer:201911251306
WATCHER_IMAGE=registry.cn-hangzhou.aliyuncs.com/njuicscn/fusion-app-watcher:201911251306

backend:
	docker build -t $(BACKEND_IMAGE) -f build/Dockerfile.backend .

operator:
	docker build -t $(OPERATOR_IMAGE) -f build/Dockerfile.operator .

consumer:
	docker build -t $(CONSUMER_IMAGE) -f build/Dockerfile.subscriber .

watcher:
	docker build -t $(WATCHER_IMAGE) -f build/Dockerfile.watcher .

push:
	docker push $(BACKEND_IMAGE)
	docker push $(OPERATOR_IMAGE)
	docker push $(CONSUMER_IMAGE)
	docker push $(WATCHER_IMAGE)
