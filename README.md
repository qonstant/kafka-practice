
# Kafka Producer and Consumer with Docker

This documentation provides instructions on setting up a Kafka producer and consumer using the Confluent Kafka Go library and Docker.

## Prerequisites

- Docker and Docker Compose installed on your machine.
- Go programming language installed.

## Project Structure

```
├── producer/main.go
├── consumer1/main.go
├── consumer2/main.go
├── docker-compose.yml
```


## Producer

The `producer/main.go` file contains the code for producing messages to a Kafka topic.

### Code

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderPlacer struct {
    producer     *kafka.Producer
    topic        string
    deliveryChan chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
    return &OrderPlacer{
        producer:     p,
        topic:        topic,
        deliveryChan: make(chan kafka.Event, 10000),
    }
}

func (op *OrderPlacer) placeOrder(orderType string, size int) error {
    var (
        format  = fmt.Sprintf("%s - %d", orderType, size)
        payload = []byte(format)
    )
    err := op.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &op.topic,
            Partition: kafka.PartitionAny,
        },
        Value: payload,
    },
    op.deliveryChan,
    )
    if err != nil {
        log.Fatal(err)
    }

    <-op.deliveryChan

    fmt.Printf("Placed order on the queue: %s\n", format)
    return nil
}

func main() {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "client.id":         "producer",
        "acks":              "all",
    })

    if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
    }

    op := NewOrderPlacer(p, "HVSE")
    for i := 0; i < 1000; i++ {
        if err := op.placeOrder("order", i+1); err != nil {
            log.Fatal(err)
        }

        time.Sleep(time.Second * 3)
    }
}

The `producer/main.go` file contains the code for producing messages to a Kafka topic.
```

## Consumer (Data Team)

The `consumer_data_team.go` file contains the code for The consumer1/main.go file contains the code for consuming messages from a Kafka topic by the Data Team.

### Code

```go
package main

import (
    "fmt"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    topic := "HVSE"

    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  "localhost:9092",
        "group.id":           "data-team",
        "auto.offset.reset":  "smallest",
        "enable.auto.commit": "true",
        "auto.commit.interval.ms": "1000",
    })  
    if err != nil {
        log.Fatal(err)
    }
    err = consumer.Subscribe(topic, nil)
    if err != nil {
        log.Fatal(err)
    }

    for {
        ev := consumer.Poll(100)
        switch e := ev.(type) {
        case *kafka.Message:
            fmt.Printf("Consumed order by Data Team: %s\n", string(e.Value))
            _, err := consumer.CommitMessage(e)
            if err != nil {
                log.Printf("Failed to commit message: %s\n", err)
            }
        case *kafka.Error:
            fmt.Printf("%s\n", e)
        }
    }
}
```

## Consumer (Nothing Group)

The `consumer2/main.go` file contains the code for consuming messages from a Kafka topic by the Nothing Group.

### Code

```go
package main

import (
    "fmt"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    topic := "HVSE"

    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  "localhost:9092",
        "group.id":           "nothing-group",
        "auto.offset.reset":  "smallest",
        "enable.auto.commit": "true",
        "auto.commit.interval.ms": "1000",
    })  
    if err != nil {
        log.Fatal(err)
    }
    err = consumer.Subscribe(topic, nil)
    if err != nil {
        log.Fatal(err)
    }

    for {
        ev := consumer.Poll(100)
        switch e := ev.(type) {
        case *kafka.Message:
            fmt.Printf("Consumed order by Nothing Group: %s\n", string(e.Value))
            _, err := consumer.CommitMessage(e)
            if err != nil {
                log.Printf("Failed to commit message: %s\n", err)
            }
        case *kafka.Error:
            fmt.Printf("%s\n", e)
        }
    }
}
```

## Docker Compose

The `docker-compose.yml` file sets up the Kafka and Zookeeper services using Docker.

### Configuration

```yaml
version: '3.5'
services:
    zookeeper:
        image: confluentinc/cp-zookeeper:latest
        environment:
            ZOOKEEPER_CLIENT_PORT: 2181
            ZOOKEEPER_TICK_TIME: 2000
        logging:
            driver: "json-file"
            options:
              max-file: "10"
              max-size: "10m"
    
    kafka:
        image: confluentinc/cp-kafka:latest
        depends_on:
            - zookeeper
        ports:
            - 9092:9092
            - 39092:39092
        environment:
            KAFKA_BROKER_ID: 1
            KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
            KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,DOCKER_HOST://host.docker.internal:39092,LOCALHOST://localhost:9092
            KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,DOCKER_HOST:PLAINTEXT,LOCALHOST:PLAINTEXT
            KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
            KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

## How to Run

1. **Start Docker services**:

    ```bash
    docker-compose up -d
    ```

2. **Run the Producer**:

    ```bash
    go run producer/main.go
    ```

3. **Run the Consumers**:

    ```bash
    go run consumer1/main.go
    go run consumer2/main.go
    ```
