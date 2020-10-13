package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	iot "github.com/arduino/iot-client-go"
)

// Parameters - List of parameters to be sent
var Parameters = make(map[string]int)

func confManager() {

	// Read parameters from configuration file
	csvfile, err := os.Open("../configuration.csv")
	if err != nil {
		log.Println("Couldn't open the configuration file", err)
	}
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		i, err := strconv.Atoi(record[1])
		if err == nil {
			Parameters[record[0]] = i
		}
	}

	things := getThings()
	count := 0
	for _, thing := range things {
		if strings.HasPrefix(thing.Name, "Contact-tracing-GW") {
			count++
		}
	}

	//Create every the number of thing specified in the file
	to_be_created := Parameters["THINGS_TO_BE_CREATED"] - count
	if to_be_created > 0 {
		for i := 1; i == to_be_created; i++ {
			thing := createThing("Contact-tracing-GW-" + strconv.Itoa(count+i))
			//Add the new thing to the things array
			things = append(things, thing)
		}
	}

	for i, thing := range things {
		if strings.HasPrefix(thing.Name, "Contact-tracing") {
			fmt.Println(i, thing.Id)
			properties := getProperties(thing.Id)

			//Create a property to handle contacts
			findOrCreate("Contacts", thing.Id, properties)

			//Create a list of properties to handle parameters
			for k, v := range Parameters {
				if strings.HasSuffix(k, "_PARAM") {
					property := findOrCreate(k, thing.Id, properties)
					publishProperty(thing.Id, property.Id, v)
				}
			}
		}
	}
}

// The function retrieve a list of things from a user
func getThings() []iot.ArduinoThing {

	// Create an instance of the iot-api Go client, we pass an empty config
	// because defaults are ok
	client := iot.NewAPIClient(iot.NewConfiguration())

	// Get the list of things for the current user
	things, _, err := client.ThingsV2Api.ThingsV2List(ctx, nil)
	if err != nil {
		log.Printf("Error retrieving things, %v", err)
		return nil
	}

	return things
}

// The function retrieve a list of properties from a thing
func getProperties(thingid string) []iot.ArduinoProperty {

	// Create an instance of the iot-api Go client, we pass an empty config
	// because defaults are ok
	client := iot.NewAPIClient(iot.NewConfiguration())

	// Get the list of things for the current user
	properties, _, err := client.PropertiesV2Api.PropertiesV2List(ctx, thingid, nil)
	if err != nil {
		log.Printf("Error retrieving properties, %v", err)
		return nil
	}

	return properties
}

func createProperty(thingid string, name string) iot.ArduinoProperty {
	// Create an instance of the iot-api Go client, we pass an empty config
	// because defaults are ok
	client := iot.NewAPIClient(iot.NewConfiguration())

	var prop iot.Property
	prop.Name = name
	prop.Permission = "READ_WRITE"
	prop.Type = "INT"
	prop.MinValue = -128
	prop.MaxValue = 127
	prop.Persist = true
	prop.UpdateParameter = 0
	prop.UpdateStrategy = "ON_CHANGE"
	prop.VariableName = name

	// Get the list of things for the current user
	property, _, err := client.PropertiesV2Api.PropertiesV2Create(ctx, thingid, prop)
	if err != nil {
		log.Printf("Error creating the property, %v", err)
	}
	return property

}

func createThing(name string) iot.ArduinoThing {
	// Create an instance of the iot-api Go client, we pass an empty config
	// because defaults are ok
	client := iot.NewAPIClient(iot.NewConfiguration())

	var t iot.CreateThingsV2Payload
	t.Name = name

	// Get the list of things for the current user
	thing, _, err := client.ThingsV2Api.ThingsV2Create(ctx, t, nil)
	if err != nil {
		log.Printf("Error creating the thing, %v", err)
	}
	return thing

}

// The function set a new value for the property
func publishProperty(thingid string, propertyid string, v int) {

	// Create an instance of the iot-api Go client, we pass an empty config
	// because defaults are ok
	client := iot.NewAPIClient(iot.NewConfiguration())

	value := iot.PropertyValue{
		DeviceId: "",
		Value:    v,
	}

	// Get the list of devices for the current user
	_, err := client.PropertiesV2Api.PropertiesV2Publish(ctx, thingid, propertyid, value)
	if err != nil {
		log.Fatalf("Error publishing the property, %v", err)
	}
}

func findOrCreate(name string, thingid string, properties []iot.ArduinoProperty) iot.ArduinoProperty {
	for _, property := range properties {
		if property.Name == name {
			fmt.Println("Property found", property.Id)
			return property
		}
	}
	return createProperty(thingid, name)
}
