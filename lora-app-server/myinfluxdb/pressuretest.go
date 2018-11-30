package myinfluxdb

import (
	"crypto/tls"
	"flag"
	"fmt" //"log"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var chanOK = make(chan string, 1000)

// MQTTRX ..
func MQTTRX() {
	// MQTT.DEBUG = log.New(os.Stdout, "", 0)
	// MQTT.ERROR = log.New(os.Stdout, "", 0)
	// stdin := bufio.NewReader(os.Stdin)
	hostname, _ := os.Hostname()

	// server := flag.String("server", "tcp://120.78.76.139:1883", "The full URL of the MQTT server to connect to")
	server := flag.String("server", "tcp://127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	// topic := flag.String("topic", "gateway/fffb02426f96c917/rx", "Topic to publish the messages on")
	topic := flag.String("topic", "pressure_test", "Topic to publish the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "entre", "A username to authenticate to the MQTT server")
	password := flag.String("password", "wcf123123", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Printf("Connected to %s\n", *server)
	for OK := range chanOK {
		client.Publish(*topic, byte(*qos), *retained, OK)
	}
}

// func main() {
// 	go MQTTRX()
// 	for {
// 		chanOK <- "0k"
// 		time.Sleep(time.Second * 1)
// 	}
// }
