package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
)

func main() {

	// Channel to signal errors
	alarm := make(chan struct{})

	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		fmt.Println("GOMAXPROCS =", runtime.NumCPU())
	} else {
		fmt.Println("GOMAXPROCS is default!")
	}

	argClientID := flag.String("client_id", "", "Arduino IoT Cloud clientID")
	argClientSecret := flag.String("client_secret", "", "Arduino IoT Cloud client secret")
	argThingID := flag.String("thing_id", "", "Arduino IoT Cloud Thing ID")
	argPropertyID := flag.String("property_id", "", "Arduino IoT Cloud Property ID")
	argSendArduino := flag.String("send_arduino", "", "Send contacts to Arduino Cloud")

	flag.Parse()

	if *argClientID == "" {
		log.Fatalln("client_id parameter is mandatory")
	}
	if *argClientSecret == "" {
		log.Fatalln("client_secret parameter is mandatory")
	}
	if *argThingID == "" {
		log.Fatalln("thing_id parameter is mandatory")
	}
	if *argPropertyID == "" {
		log.Fatalln("property_id parameter is mandatory")
	}

	if *argSendArduino != "" {
		ArduinoAPIEndpoint = *argSendArduino
	}
	ClientID = *argClientID
	ClientSecret = *argClientSecret
	ThingID = *argThingID
	PropertyID = *argPropertyID

	tknSync := make(chan string)
	go tokenManager(tknSync)

	// Wait for token synchornization before starting the configuration
	tkn := <-tknSync
	fmt.Println(tkn)
	go confManager()

	// If something go wrong exit
	msg := <-alarm
	fmt.Println(msg)
}
