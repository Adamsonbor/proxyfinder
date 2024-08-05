package rabbit

import (
	"context"
	"fmt"
	"proxyfinder/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitService struct {
	cfg   *config.Config
	conn  *amqp.Connection
	ch    *amqp.Channel
	q     *amqp.Queue
	qname string
}

func (r *RabbitService) Publish(ctx context.Context, body []byte) error {
	err := r.ch.PublishWithContext(
		ctx,
		"",
		r.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitService) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *RabbitService) Close() error {
	err := r.ch.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitService) Connect() error {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		r.cfg.Rabbit.User,
		r.cfg.Rabbit.Pass,
		r.cfg.Rabbit.Host,
		r.cfg.Rabbit.Port,
	))
	if err != nil {
		return err
	}
	r.conn = conn
	r.ch, err = conn.Channel()
	if err != nil {
		return err
	}
	q, err := r.ch.QueueDeclare(
		r.qname,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	r.q = &q

	return nil
}

func NewRabbit(cfg *config.Config, qname string) *RabbitService {
	service := &RabbitService{
		cfg:   cfg,
		qname: qname,
	}

	err := service.Connect()
	if err != nil {
		panic(err)
	}

	return service
}
