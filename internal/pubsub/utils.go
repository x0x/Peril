package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Printf("Unable to marshal value to jsonData: %v", err)
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonData),
	}

	ctx := context.Background()

	err = ch.PublishWithContext(ctx, exchange, key, false, false, msg)
	if err != nil {
		log.Printf("Unable to publish with context :%v", err)
		return err
	}

	return nil
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("conn.Channel : %v", err)
		return nil, amqp.Queue{}, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueType.Durable,
		queueType.AutoDelete,
		queueType.Exclusive,
		queueType.NoWait,
		queueType.Args,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName,
		key,
		exchange,
		false,
		nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	log.Printf("Queue Bound with queue name %v", queue)
	return ch, queue, nil
}

func SubscribeJSON[T any](conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType, 
) error {
	_, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		log.Printf("Error in declaring and binding the queue : %v", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			var msg T
			err := json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				log.Printf("Unable to unmarshal the data :%v", err)
				continue
			}
			fmt.Println("Received message from the exchange %v", msg)
			ackType := handler(msg)
			if ackType == Ack {
				err = delivery.Ack(false)
			} else if ackType == NackReque {
				err = delivery.Nack(false, true)
			} else {
				err = delivery.Nack(false, false)
			}
			if err != nil {
				log.Printf("Ack error: %+v", err)
			}
		}
	}()

	return nil
}
