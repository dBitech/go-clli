package main

import (
	"fmt"
	"log"

	"github.com/dbitech/go-clli/pkg/clli"
)

func main() {
	result, err := clli.Parse("LSANCA12")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Input: LSANCA12\n")
	fmt.Printf("Place: %s\n", result.Place)
	fmt.Printf("Region: %s\n", result.Region)
	fmt.Printf("NetworkSite: %s\n", result.NetworkSite)
	fmt.Printf("EntityCode: %s\n", result.EntityCode)
	fmt.Printf("LocationCode: %s\n", result.LocationCode)
	fmt.Printf("LocationID: %s\n", result.LocationID)
	fmt.Printf("CustomerCode: %s\n", result.CustomerCode)
	fmt.Printf("CustomerID: %s\n", result.CustomerID)
	fmt.Printf("Type: %d\n", result.Type())
	fmt.Printf("IsEntityCLLI: %t\n", result.IsEntityCLLI())
	fmt.Printf("IsNonBuildingCLLI: %t\n", result.IsNonBuildingCLLI())
	fmt.Printf("IsCustomerCLLI: %t\n", result.IsCustomerCLLI())
}
