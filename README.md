# notifyGo
消息推送平台

## 核心流程梳理

### 1）架构层面：

一个消息发送的过程可以抽象为三块，即：**Send(Content) → Target**

- Send：发送动作相关的逻辑。主要负责如何将消息发到kafka，同一种类型可以有多个topic，支持一些路由策略，也支持热加载

- Content：内容渲染相关的逻辑。主要负责根据发送的目标对象、消息模版、自定义变量构建发送内容。

- Target：内容发送目标相关的逻辑。主要负责发送目标的获取，比如接入的用户画像服务、手动上传等，同时会进行一些目标过滤

### 2）业务层面：

简单分了三张表：

- Delivery表：推送记录表，一次发送会创建一个记录；
- Target表：target粒度的发送记录表。针对一次发送，每一个发送的对象创建一个推送记录；
- Template表：模版表，存放的发送模版；


> TODO：
1. 架构合理性、代码质量问题，补充单元测试；
2. 短信、push、邮件等各个发送方式的接入；
3. 定时发送任务；
4. 大批量(百万级)消息发送等的一系列相关支持；
5. 消息记录查询；
6. 监控、报警；
7. 全链路压测

> 附代码结构:

```
├── LICENSE
├── README.md
├── cmd  // 进程服务启动入口
│   ├── message_api  // api服务
│   │   └── main.go
│   └── message_worker // kafka消费者服务
│       └── main.go
├── conf
│   ├── app.toml // 应用配置
│   └── kafka_topic.toml // kafka topic配置,支持热加载进行topic扩展
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── handler
│   │   ├── middleware
│   │   └── router
│   ├── config // 配置解析
│   │   └── app.go
│   ├── model // 模型层
│   │   └── model.go
│   ├── service // 核心服务层
│   │   ├── content_service.go // 获取发送内容服务
│   │   ├── core.go
│   │   ├── send_service.go // 执行发送服务
│   │   └── target_service.go // 获取发送目标
│   ├── target.go // 发送目标的定义
│   ├── type.go
│   └── worker // 具体消费kafka消息的worker
│       ├── sender
│       │   ├── base.go
│       │   ├── email.go
│       │   ├── push.go
│       │   └── sms.go
│       ├── task.go // 任务抽象
│       └── worker.go // 消费者worker定义
├── notify_go.sql
├── pkg
│   ├── item
│   │   └── item.go
│   ├── logger
│   │   ├── logger.go
│   │   ├── logger_test.go
│   │   └── type.go
│   ├── mq
│   │   ├── balancer.go // 发送消息topic负载均衡
│   │   ├── consumer.go // kafka消费者抽象
│   │   ├── kafka_config.go // topic配置热加载
│   │   └── producer.go // kafka生产者抽象
│   └── tool
│       └── tool.go
└── xorm.yaml
```