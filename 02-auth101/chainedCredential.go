package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

const subscriptionID = "5073fd4c-3a1b-4559-8371-21e034f70820"

func main() {
	azlog.SetListener(func(event azlog.Event, s string) {
		fmt.Println(s)
	})

	azlog.SetEvents(azidentity.EventAuthentication)

	cli, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	env, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	mic, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	cred, err := azidentity.NewChainedTokenCredential([]azcore.TokenCredential{cli, mic, env}, nil)
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
