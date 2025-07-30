package queue

import (
	"context"
	"encoding/json"
	"ledger/internal/domain"
	"log"

	"github.com/streadway/amqp"
)

type TransactionPublisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewTransactionPublisher(amqpURL, queueName string) (*TransactionPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		err := ch.Close()
		if err != nil {
			return nil, err
		}
		err = conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &TransactionPublisher{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (p *TransactionPublisher) Publish(ctx context.Context, msg domain.Transaction) error {
	if p.channel == nil {
		return amqp.ErrClosed
	}
	// Ensure the context is not cancelled before publishing
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish transaction: %v", err)
		return err
	}

	log.Printf("Transaction published: %+v", msg)
	return nil
}

func (p *TransactionPublisher) Close() {
	if p.channel != nil {
		err := p.channel.Close()
		if err != nil {
			return
		}
	}
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return
		}
	}
}
