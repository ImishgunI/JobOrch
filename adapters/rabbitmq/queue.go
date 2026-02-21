package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitQueue struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewRabbitQueue(ctx context.Context, connection, queue string) (*rabbitQueue, error) {
	conn, err := getconn(connection)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = ch.Qos(1, 0, false); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &rabbitQueue{
		conn:  conn,
		ch:    ch,
		queue: q.Name,
	}, nil
}

func getconn(connection string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connection)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (q *rabbitQueue) CloseConnection() {
	q.conn.Close()
}

func (q *rabbitQueue) CloseChannel() {
	q.ch.Close()
}

func (q *rabbitQueue) PublishWithContext(ctx context.Context, jobID string) error {
	err := q.ch.PublishWithContext(ctx, "", q.queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(jobID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (q *rabbitQueue) ConsumeWithContext(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := q.ch.ConsumeWithContext(ctx, q.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
