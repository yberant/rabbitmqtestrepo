package main

import (
	"log"
	"fmt"
	//"strconv"
	"bytes"
  	"encoding/json"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message map[string]interface{}

func deserialize(b []byte) (Message, error) {
    var msg Message
    buf := bytes.NewBuffer(b)
    decoder := json.NewDecoder(buf)
    err := decoder.Decode(&msg)
    return msg, err
}


func main() {

	fmt.Println("ingrese direccion de ip de logistica")
	var addrS string
	fmt.Scanln(&addrS)
		
	/*portnum:
		fmt.Println("ingrese puerto en el que ip de logistica esta escuchando a serv. financiero")
		var port string
		fmt.Scanln(&port)
		if p,err:=strconv.Atoi(port);err!=nil{
			goto portnum
		} else {
			_=p
		}*/

	// /5672
	conn, err := amqp.Dial("amqp://user:pass@"+addrS+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}