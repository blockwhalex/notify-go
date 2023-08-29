package mq

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	hosts                []string
	topicBalanceMappings map[string]Balancer
}

func NewProducer(mqCfg *Config) *Producer {
	var tbMappings map[string]Balancer
	tbMappings = make(map[string]Balancer, len(mqCfg.TopicMappings))

	//	为每一个channel创建好balancer
	for channel, topics := range mqCfg.TopicMappings {
		bala := NewBalanceBuilder(channel, topics.Topics).Build(topics.Strategy)
		tbMappings[channel] = bala
	}

	go func() {
		for {
			mu.Lock()
			for !topicChange {
				changeSignal.Wait()
			}
			log.Println("producer update topics")
			for channel, topics := range mqCfg.TopicMappings {
				bala := NewBalanceBuilder(channel, topics.Topics).Build(topics.Strategy)
				tbMappings[channel] = bala
			}
			topicChange = false
			mu.Unlock()
		}
	}()

	return &Producer{
		hosts:                mqCfg.Host,
		topicBalanceMappings: tbMappings,
	}
}

// Send data里面应该包含channel信息，方法找到topic
func (p *Producer) Send(channel string, data []byte) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true    // 设定是否需要返回错误信息
	config.Producer.Return.Successes = true // 设定是否需要返回成功信息
	producer, err := sarama.NewAsyncProducer(p.hosts, config)
	if err != nil {
		log.Fatal("NewSyncProducer err:", err)
	}
	var (
		wg                                   sync.WaitGroup
		enqueued, timeout, successes, errors int
	)
	// [!important] 异步生产者发送后必须把返回值从 Errors 或者 Successes 中读出来 不然会阻塞 sarama 内部处理逻辑 导致只能发出去一条消息
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range producer.Errors() {
			log.Printf("[Producer] Errors：err:%v msg:%+v \n", e.Msg, e.Err)
			errors++
		}
	}()

	// 根据channel类型，和路由策略选取发送的topic
	topic, err := p.topicBalanceMappings[channel].GetNext()
	if err != nil {
		log.Printf("[Producer] choose topic fail: channel:%v error:%+v \n", channel, err)
	}
	msg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(data)}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	select {
	case producer.Input() <- msg:
		enqueued++
	case <-ctx.Done():
		timeout++
	}
	cancel()

	producer.AsyncClose()
	wg.Wait()
	log.Printf("发送完毕[%s] enqueued:%d timeout:%d successes: %d errors: %d\n", topic, enqueued,
		timeout, successes, errors)
}
