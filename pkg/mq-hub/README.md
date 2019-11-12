
## kafka-rest
### 获取topics
```
curl 192.168.1.10:32015/topics
```

### Produce 一条数据
```
curl -X POST -H "Content-Type: application/vnd.kafka.json.v2+json"   --data '{"records":[{"value":{"name": "testUser"}}]}' "http://192.168.1.10:32015/topics/test"
```

这样会往kafka中传入一条数据  {"name": "testUser"}

如果是传入多条数据，可以这样
```
curl -X POST -H "Content-Type: application/vnd.kafka.json.v2+json"   --data '{"records":[{"value":{"name": "testUser"}},{"value":{"name": "testUser2"}}]}' "http://192.168.1.10:32015/topics/test"
```
如果成功，便会返回下面的日志
```
{"offsets":[{"partition":0,"offset":10,"error_code":null,"error":null},{"partition":0,"offset":11,"error_code":null,"error":null}],"key_schema_id":null,"value_schema_id":null}
```
这样便依次传入了 {"name": "testUser"}和{"name": "testUser2"}

## kafka-consumer 封装
如何用：
```
k := &KafkaSubscriber{}

// 初始化EventSource
s := &core.EventSource{}
s.Component = core.ExternalWorkflowEngine
s.Host = "127.0.0.1"

// consumer config配置，broker_list 为kafka访问接口， group则是consuemr group id
k.broker_list = []string{"114.212.87.225:32590"}
k.group = "testworkflowEngine127.0.0.1"

// 通过Subscribe最后得到一个 chan 接口，可以通过valueChan来获取数据
valueChan, _ := k.SubscribeTo(s)
```

具体使用可以参考 subsriber_test.go