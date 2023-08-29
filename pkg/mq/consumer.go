package mq

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

// 这里写consumer的抽象方法，然后具体的消费逻辑写到worker里面，然后路由到sender中不同的处理方法里面

type ConsumerGroup struct {
	handler       sarama.ConsumerGroupHandler
	sCg           sarama.ConsumerGroup
	topics        *[]string
	restartSignal chan struct{}
}

func NewConsumerGroup(mqCfg *KafkaConfig, channel string, handler sarama.ConsumerGroupHandler) *ConsumerGroup {
	sCfg := sarama.NewConfig()
	sCfg.Consumer.Return.Errors = true

	topics := mqCfg.GetTopicsByChannel(channel)
	groupId := mqCfg.GetGroupIdByChannel(channel)

	cg, err := sarama.NewConsumerGroup(mqCfg.Host, groupId, sCfg)
	if err != nil {
		log.Fatal("NewConsumerGroup err: ", err)
	}

	mcg := &ConsumerGroup{
		handler:       handler,
		sCg:           cg,
		topics:        &topics,
		restartSignal: make(chan struct{}),
	}

	go func() {
		for {
			mu.Lock()
			// topic变化后热更新
			for !kafkaTopicChange {
				changeSignal.Wait()
			}
			topics = mqCfg.GetTopicsByChannel(channel)
			log.Println("consumer update topics, ", topics)
			kafkaTopicChange = false
			select {
			case mcg.restartSignal <- struct{}{}:
			default:
			}
			mu.Unlock()
		}
	}()

	return mcg
}

func (c *ConsumerGroup) Start(ctx context.Context) {
	defer func() { _ = c.sCg.Close() }()

	for {
		newCtx, cancel := context.WithCancel(ctx)
		fmt.Println("rerunning consume...", c.topics)

		go func() {
			select {
			// 消费者接收topic配置变化信号，取消当前的消费
			case <-c.restartSignal:
				fmt.Println("new topics...", c.topics)
				cancel()
			}
		}()

		// newCtx被取消之后，这里会在下一轮循环重新执行，似乎要执行re-balance，那么还是不能无损解决消息积压问题
		err := c.sCg.Consume(newCtx, *c.topics, c.handler)
		if err != nil {
			log.Println("Consume err: ", err)
		}
	}
}
