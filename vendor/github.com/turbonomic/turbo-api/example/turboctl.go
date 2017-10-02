package main

import (
	"fmt"
	"net/url"

	"github.com/turbonomic/turbo-api/pkg/api"
	"github.com/turbonomic/turbo-api/pkg/client"
)

func main() {
	discoverTargetExample()
}

func addTarget() {
	serverAddress, err := url.Parse("<Server_Address>")
	if err != nil {
		fmt.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		APIPath("/vmturbo/rest").
		BasicAuthentication("<UI-username>", "UI-password").
		Create()
	client, err := client.NewAPIClientWithBA(config)
	if err != nil {
		fmt.Errorf("Error creating client: %s", err)
	}

	target := &api.Target{
		Category: "Hypervisor",
		Type:     "vCenter",
		InputFields: []*api.InputField{
			{
				Value:           "<VC_Address>",
				Name:            "nameOrAddress",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<VC_Username>",
				Name:            "username",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<VC_Password>",
				Name:            "password",
				GroupProperties: []*api.List{},
			},
		},
	}
	resp, err := client.AddTarget(target)
	if err != nil {
		fmt.Errorf("Error adding target: %s", err)
		return
	}
	fmt.Printf("Response is %++v", resp)
}

// Add an external target. This type of type is registered through SDK.
// Here we use Kubernetes target for example.
func addExternalTarget() {
	// Get Turbonomic server address.
	serverAddress, err := url.Parse("<SERVER_ADDRESS>")
	if err != nil {
		fmt.Printf("Incorrect URL: %s\n", err)
	}

	// Create API client config.
	config := client.NewConfigBuilder(serverAddress).
		APIPath("/vmturbo/rest").
		BasicAuthentication("<UI_USERNAME>", "<UI_PASSWORD>").
		Create()
	client, err := client.NewAPIClientWithBA(config)
	if err != nil {
		fmt.Printf("Error creating client: %s\n", err)
	}

	// Configure target data.
	target := &api.Target{
		Category: "Custom",
		Type:     "Kubernetes",
		InputFields: []*api.InputField{
			{
				Value:           "<Kubernetes_TargetID>",
				Name:            "targetIdentifier",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<Kubernetes_Target_Username>",
				Name:            "username",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<Kubernetes_Target_Password>",
				Name:            "password",
				GroupProperties: []*api.List{},
			},
		},
	}

	// Make API calls.
	resp, err := client.AddTarget(target)
	if err != nil {
		fmt.Printf("Error adding target: %s\n", err)
		return
	}
	fmt.Printf("Response is %++v\n", resp)
}

func discoverTargetExample() {
	serverAddress, err := url.Parse("<SERVER_ADDRESS>")
	if err != nil {
		fmt.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		APIPath("/vmturbo/rest").
		BasicAuthentication("<UI_USERNAME>", "<UI_PASSWORD>").
		Create()
	client, err := client.NewAPIClientWithBA(config)
	if err != nil {
		fmt.Errorf("Error creating client: %s", err)
	}
	uuid := "<TARGET_UUID>"
	resp, err := client.DiscoverTarget(uuid)
	if err != nil {
		fmt.Errorf("Error adding target: %s", err)
		return
	}
	fmt.Printf("Response is %++v", resp)
}
