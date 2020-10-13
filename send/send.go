package main

import (
  "log"
  "fmt"
  "net"
  "github.com/streadway/amqp"
  "bytes"
  "encoding/json"
  //"strconv"
)

type Message map[string]interface{}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func getIPAddr() string{
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func serialize(msg Message) ([]byte, error) {
    var b bytes.Buffer
    encoder := json.NewEncoder(&b)
    err := encoder.Encode(msg)
    return b.Bytes(), err
}

func main(){

	IPAddr:=getIPAddr()
	//log.Println(IpAddr)
	fmt.Println("dirección IP de logística: ",IPAddr)
	
	/*
	portnum:
		fmt.Println("ingrese puerto para escuchar a a serv. financiero")
		var port string
		fmt.Scanln(&port)
		if p,err:=strconv.Atoi(port);err!=nil{
			goto portnum
		} else {
			_=p
		}
	*/

	conn, err := amqp.Dial("amqp://user:pass@"+IPAddr+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//declare queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  
	failOnError(err, "Failed to declare a queue")

	body:=Message{
		"name":"raul",
		"age":19,
	}

	ser,err:=serialize(body)
	  
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
		  Body:        ser,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}