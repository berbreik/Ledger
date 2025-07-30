package queue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"ledger/internal/domain"
)

type TransactionConsumer struct {
	conn               *amqp.Connection
	channel            *amqp.Channel
	queueName          string
	transactionService domain.TransactionService
}

// StartTransactionConsumer wrapper for the transaction consumer that listens to a RabbitMQ queue and processes transactions
func StartTransactionConsumer(amqpURL, queueName string, service domain.TransactionService) error {
	consumer, err := NewTransactionConsumer(amqpURL, queueName, service)
	if err != nil {
		return err
	}
	go func() {
		if err := consumer.StartConsuming(); err != nil {
			log.Fatalf("consumer failed: %v", err)
		}
	}()
	return nil
}

func NewTransactionConsumer(amqpURL, queueName string, service domain.TransactionService) (*TransactionConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &TransactionConsumer{
		conn:               conn,
		channel:            ch,
		queueName:          queueName,
		transactionService: service,
	}, nil
}

func (c *TransactionConsumer) StartConsuming() error {
	ctx := context.Background()
	if c.channel == nil {
		return amqp.ErrClosed
	}
	if c.conn == nil {
		return amqp.ErrClosed
	}
	// Start consuming messages from the queue
	log.Printf("Starting consumer for queue: %s", c.queueName)

	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg domain.Transaction
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Invalid transaction message: %v", err)
				continue
			}

			log.Printf("Processing transaction: %+v", msg)

			err := c.transactionService.ProcessTransaction(ctx, &msg)
			if err != nil {
				log.Printf("Failed to process transaction: %v", err)
			} else {
				log.Printf("Transaction processed successfully")
			}
		}
	}()

	log.Printf("Consumer started. Waiting for messages...")
	return nil
}

func (c *TransactionConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
