package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

const subscriptionID = "5073fd4c-3a1b-4559-8371-21e034f70820"

func main() {
	azlog.SetListener(func(event azlog.Event, s string) {
		fmt.Println(s)
	})

	// include only azidentity credential logs
	azlog.SetEvents(azidentity.EventAuthentication)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	rcFactory, err := armresources.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	rgClient := rcFactory.NewResourceGroupsClient()

	ctx := context.Background()
	resultPager := rgClient.NewListPager(nil)

	resourceGroups := make([]*armresources.ResourceGroup, 0)
	for resultPager.More() {
		pageResp, err := resultPager.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		resourceGroups = append(resourceGroups, pageResp.ResourceGroupListResult.Value...)
	}

	jsonData, err := json.MarshalIndent(resourceGroups, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))

}
