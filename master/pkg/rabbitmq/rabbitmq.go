package rabbitmq

import (
	"fmt"
	"master/pkg/constant"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

// RabbitMQ represents config for connecting to RMQ
type RmqConfig struct {
	RabbitMQURL    string
	ExchangeName   string
	QueueName      string
	DelayQueueName string
	RoutingKey     string
	Message        string
	MessageID      string
	MessageTTL     string // e.g., "60000" (60 seconds in milliseconds)
}

func setRabbitMQURL(rmq *RmqConfig) *RmqConfig {
	rmq.RabbitMQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	// log.Info("setRabbitMqURL >>> ", rmq.RabbitMQURL)
	return rmq
}

// publishes a message with a specific exchange & routing key
func PublishMessage(rmq *RmqConfig) error {
	// Set RabbitMQ connection URL
	rmq = setRabbitMQURL(rmq)

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rmq.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare an exchange with x-delayed-message type
	err = ch.ExchangeDeclare(
		rmq.ExchangeName,    // Exchange name
		"x-delayed-message", // Exchange type
		true,                // Durable
		false,               // Auto-deleted
		false,               // Internal
		false,               // No-wait
		amqp.Table{
			"x-delayed-type": "direct", // Message routing strategy
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		rmq.QueueName, // Queue name
		true,          // Durable
		false,         // Auto-deleted
		false,         // Exclusive
		false,         // No-wait
		amqp.Table{
			"x-queue-type": constant.RMQ_DEFAULT_QUEUE_TYPE,
		}, // Arguments,

	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind the queue to the delayed exchange
	err = ch.QueueBind(
		q.Name,           // Queue name
		rmq.RoutingKey,   // Routing key
		rmq.ExchangeName, // Exchange name
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind the queue: %w", err)
	}

	// Convert MessageTTL string to integer (milliseconds)
	var delay int
	if rmq.MessageTTL != "" {
		delay, err = strconv.Atoi(rmq.MessageTTL)
		if err != nil {
			return fmt.Errorf("invalid MessageTTL format: %w", err)
		}
	}

	// Publish message with delay
	amqpPub := amqp.Publishing{
		MessageId:    rmq.MessageID,
		ContentType:  "application/json",
		Body:         []byte(rmq.Message),
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Headers: amqp.Table{
			"x-delay": delay, // Set delay in milliseconds
		},
	}

	err = ch.Publish(
		rmq.ExchangeName, // Exchange name
		rmq.RoutingKey,   // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqpPub,          // Message
	)
	if err != nil {
		return fmt.Errorf("failed to publish delayed message: %w", err)
	}

	// Log message information
	log.Infof("[RMQ] Sent Message: %s", rmq.Message)
	log.Infof("[RMQ] Queue Name: %s", rmq.QueueName)
	log.Infof("[RMQ] Delay (ms): %d", delay)

	return nil
}

// Subscribe connects to a RabbitMQ queue and processes messages using the provided callback function.
func Subscribe(queueName string, processMessage func(amqp.Delivery)) error {
	defer func() {
		if err := recover(); err != nil {
			time.Sleep(10 * time.Second)
			log.Info("Queue " + queueName + ", rest for trying subscribe for 10 seconds")
			Subscribe(queueName, processMessage)
		}
	}()
	// Construct the RabbitMQ URL
	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USERNAME"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-queue-type": constant.RMQ_DEFAULT_QUEUE_TYPE,
		}, // Arguments,
	)
	if err != nil {
		fmt.Println("Failed to declare a queue : ", err)
	}

	// Register a consumer
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Infof("Subscribed to messages with queueName: %s", queueName)

	// Channel to keep the subscriber running
	closeChan := make(chan *amqp.Error, 1)
	notifyClose := ch.NotifyClose(closeChan) //Once the consumer's channel has an error, an amqp.Error is generated, and the channel monitors and catches this error
	closeFlag := false

	for {
		select {
		case e := <-notifyClose:
			log.Error("chan channel error : ", e.Error())
			close(closeChan)
			time.Sleep(10 * time.Second)
			Subscribe(queueName, processMessage)
			closeFlag = true
		case msg := <-msgs:
			log.Infof("[RMQ] Received message: %s with routing key: %s\n", string(msg.Body), msg.RoutingKey)

			if len(msg.Body) == 0 {
				close(closeChan)
				time.Sleep(10 * time.Second)
				Subscribe(queueName, processMessage)
				closeFlag = true
			}

			var errProcess error
			chanErr := make(chan error, 1)
			processMessage(msg)
			chanErr <- errProcess

		}
		if closeFlag {
			break
		}
	}

	return nil
}
