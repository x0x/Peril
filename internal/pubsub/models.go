package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type SimpleQueueType struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

var DurableQueueType = SimpleQueueType{
	Durable:    true,
	AutoDelete: false,
	Exclusive:  false,
	NoWait:     false,
	Args:       nil,
}

var TransientQueueType = SimpleQueueType{
	Durable:    false,
	AutoDelete: true,
	Exclusive:  true,
	NoWait:     false,
	Args:       nil,
}

type AckType int

const (
	Ack AckType = iota
	NackReque
	NackDiscard
)
