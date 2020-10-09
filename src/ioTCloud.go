package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	iot "github.com/arduino/iot-client-go"
	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
)

// Init the context
var ctx = context.WithValue(context.Background(), iot.ContextAccessToken, "")

func getToken() (*oauth2.Token, error) {
	// We need to pass the additional "audience" var to request an access token
	additionalValues := url.Values{}
	additionalValues.Add("audience", ArduinoAPIEndpoint)
	// Set up OAuth2 configuration
	config := cc.Config{
		ClientID:       ClientID,
		ClientSecret:   ClientSecret,
		TokenURL:       ArduinoAPIEndpoint + "/v1/clients/token",
		EndpointParams: additionalValues,
	}
	// Get the access token in exchange of client_id and client_secret
	tok, err := config.Token(context.Background())
	if err != nil {
		fmt.Printf("Error retrieving access token, %v", err)
		time.Sleep(10 * time.Second)
		return getToken()
	}
	// Confirm we got the token and print expiration time
	log.Printf("Got an access token, will expire on %s", tok.Expiry)
	return tok, nil
}

func tokenManager(tknSync chan string) {
	tok, err := getToken()
	if err != nil {
		fmt.Println("Cannot get a new token ", err)
	} else {
		// We use the token to create a context that will be passed to any API call
		ctx = context.WithValue(context.Background(), iot.ContextAccessToken, tok.AccessToken)
		tknSync <- tok.AccessToken
	}
	ticker := time.NewTicker(200 * time.Second)
	for range ticker.C {
		fmt.Println("Refreshing token")
		tok, err := getToken()
		if err != nil {
			fmt.Println("Cannot get a new token ", err)
		} else {
			// We use the token to create a context that will be passed to any API call
			ctx = context.WithValue(context.Background(), iot.ContextAccessToken, tok.AccessToken)
		}
	}
}
