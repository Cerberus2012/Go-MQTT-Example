package main

import (
	"log"
	"fmt"
	"time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    log.Printf("Connect lost: %v", err)
}

func publish (client mqtt.Client) {
	token := client.Publish("test", 0, false, time.Now().String())
	token.Wait()
}

func main() {
	log.Println("Configuring MQTT options...")
    var broker = "192.168.1.102"
    var port = 1883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("emqx")
    opts.SetPassword("public")
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
	log.Println("Connecting to MQTT broker...")
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
  	}

	msgCount := 50000
	start := time.Now()
	for i := 0; i < msgCount; i ++ {
		publish(client)
  	}
	duration := time.Since(start)

	log.Printf("Total duration for %d messages: %d ms", msgCount, duration.Milliseconds())
}