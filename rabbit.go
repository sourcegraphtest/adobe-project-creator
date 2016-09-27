package main

import (
    "fmt"
    "github.com/streadway/amqp"
    "log"
)

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    tag     string
    done    chan error
}

func NewConsumer(uri, key, ctag string) (*Consumer, error) {
    c := &Consumer{
        conn:    nil,
        channel: nil,
        tag:     ctag,
        done:    make(chan error),
    }

    var err error

    c.conn, err = amqp.Dial(uri)
    if err != nil {
        return nil, fmt.Errorf("Dial: %s", err)
    }

    go func() {
        fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
    }()

    log.Printf("got Connection, getting Channel")
    c.channel, err = c.conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("Channel: %s", err)
    }

    log.Print("got Channel, declaring Exchange (ft.project-creator.exchange)")
    if err = c.channel.ExchangeDeclare(
        "ft.project-creator.exchange",     // name of the exchange
        "direct",                          // type
        true,                              // durable
        false,                             // delete when complete
        false,                             // internal
        false,                             // noWait
        nil,                               // arguments
    ); err != nil {
        return nil, fmt.Errorf("Exchange Declare: %s", err)
    }

    log.Print("declared Exchange, declaring Queue ft.project-creator.create")
    queue, err := c.channel.QueueDeclare(
        "ft.project-creator.create",   // name of the queue
        true,                          // durable
        false,                         // delete when usused
        false,                         // exclusive
        false,                         // noWait
        nil,                           // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Declare: %s", err)
    }

    log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
        queue.Name, queue.Messages, queue.Consumers, key)

    if err = c.channel.QueueBind(
        queue.Name,                     // name of the queue
        key,                            // bindingKey
        "ft.project-creator.exchange",  // sourceExchange
        false,                          // noWait
        nil,                            // arguments
    ); err != nil {
        return nil, fmt.Errorf("Queue Bind: %s", err)
    }

    log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
    deliveries, err := c.channel.Consume(
        queue.Name, // name
        c.tag,      // consumerTag,
        false,      // noAck
        false,      // exclusive
        false,      // noLocal
        false,      // noWait
        nil,        // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Consume: %s", err)
    }

    go handle(deliveries, c.done)
    select{}
}

func (c *Consumer) Shutdown() error {
    // will close() the deliveries channel
    if err := c.channel.Cancel(c.tag, true); err != nil {
        return fmt.Errorf("Consumer cancel failed: %s", err)
    }

    if err := c.conn.Close(); err != nil {
        return fmt.Errorf("AMQP connection close error: %s", err)
    }

    defer log.Printf("AMQP shutdown OK")

    // wait for handle() to exit
    return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
    for d := range deliveries {
        log.Printf("[%v] : Received %q", d.DeliveryTag, d.Body)

        status,errors := NewProject(string(d.Body))
        if errors {
            log.Printf("[%v] : errors %s", d.DeliveryTag, status)
            d.Ack(false)
        } else {
            log.Printf("[%v] : %s", d.DeliveryTag, status)
            d.Ack(true)
        }
    }
    log.Printf("handle: deliveries channel closed")
    done <- nil
}
