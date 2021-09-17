zipkin-go 介绍
===

## zipkin-go 中的数据结构

###  Endpoint

代码: `zipkin-go/model/endpoint.go`

网络上下文，包含处理 http/rpc 请求的服务的 ip，端口，名称

zipkin 提供了 api 查询 endpoint，`http://localhost:9411/api/v2/services`，`service_name` 也可以作为参数筛选 Span 和 Trace

### Tracer

代码: `zipkin-go/tracer.go`

服务初始化的时候创建，所有请求应共用一个 tracer。
它包含 sampler (指标发送频率控制器), generator (id 生成器) , reporter (指标发送器)

### Reporter

代码: `zipkin-go/reporter/reporter.go`

将 span 数据发送到 zipkin 的客户端，zipkin-go 提供了5种 reporter, 

  - `http-reporter`: 通过 http 请求异步地发送 span 数据到 collector 中
  - `kafka-reporter`: 将 span 数据发送到 kafka 队列中，再由 collector 收集
  - `amqp-reporter`: 将 span 数据发送到 amqp 队列中，再由 collector 收集
  - `recorder-reporter`: 将 span 数据存储在内存中，提供了获取 span 数据的接口
  - `log-reporter`: 将 span 数据写入到日志中

### Span & Trace

- Span 代表某个节点上某一次操作的记录
- Trace 通常代表一次 http/rpc 请求。 Trace 一个隐含的概念，在 zipkin-go 中没有代码和这一概念对应，但 zipkin 提供了查询 trace 的 api。`"http://localhost:9411/api/v2/traces?limit=10"`

Trace 是由许多 Trace ID 相同的 Span 嵌套组成的一颗树。Span 之间通过 parent_id 属性建立起连接关系。

Root Span 没有 Parent ID，通常情况下它的处理时间是最长的。如果请求中有异步请求，那么异步请求的 Span 的处理时间会超过 Root Span。

#### Span 的数据结构

+ __数据__

- Name: Span 操作的名称
- Kind: 
    - Client: Timestamp 表示 client 发起请求的时间，duration 表示请求结束的时间
    - Server: Timestamp 表示 server 收到请求的时间，duration 表示处理请求的时间
    - Consumer: Timestamp 表示 consumer 收到消息的时间，duration 表示处理消息的时间
    - Producer: Timestamp 表示 producer 发送消息的时间，duration 表示收到消息确认发送成功的时间
- Timestamp: Span 操作的开始时间
- Duration: Span 操作的持续时间
- LocalEndpoint: 本地网络环境 (http middleware 中指 server ip)
- RemoteEndpoint: 远端网络环境 (http middleware 中指 client ip)
- Annotations: Annotation 列表
- Tags: 通过 Tag 记录额外的属性

### Annotation

嵌套在 Span 中的属性，拥有 timestamp 和 value 两个属性，用来描述事件发生的时间。

## zipkin-go 的 http reporter 是如何发送 span 的

http reporter 的数据结构

+ __数据__

- batch: span 数组
- batchInterval: 发送 Span 的时间间隔
- serializer: 发送前压缩 Span 数据的方法

+ __Goroutine__

- loop: 
  1. 将 Send 函数发送的 Span 存储到 batch 数组中，如果 batch 数组达到阈值，则执行一次发送
  2. 每隔 batchInterval 秒，执行一次发送
- sendLoop: 将 batch 数组中的数据取出，通过 http 发送到 zipkin 中

+ __方法__

- Send: 外部函数调用 Send 函数发送 Span
- Close: 外部调用 Close 函数时，http reporter 清理 batch 中的数据，并停止所有 Goroutine

## 参考链接

+ https://pkg.go.dev/github.com/openzipkin/zipkin-go#section-readme
+ https://zipkin.io/zipkin-api/
+ https://zipkin.io/pages/instrumenting.html
+ https://www.cnblogs.com/zhoubaojian/articles/7852358.html
+ https://github.com/bigbully/Dapper-translation/blob/master/dapper%E5%88%86%E5%B8%83%E5%BC%8F%E8%B7%9F%E8%B8%AA%E7%B3%BB%E7%BB%9F%E5%8E%9F%E6%96%87.pdf
+ [zipkin tutorial](https://www.scalyr.com/blog/zipkin-tutorial-distributed-tracing/)
+ [zipkin tutorial翻译](https://zhuanlan.zhihu.com/p/95054724)
+ [openzipkin 文档](https://www.cxyzjd.com/article/lz710117239/89107748)
+ [dapper 论文](https://research.google/pubs/pub36356/)
+ [dapper 论文翻译](https://bigbully.github.io/Dapper-translation/)
+ [dapper 论文翻译](https://github.com/AlphaWang/alpha-dapper-translation-zh)
+ [dapper 论文总结](https://www.cxyzjd.com/article/lz710117239/89107748)
