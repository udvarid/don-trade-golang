package user

import (
	"testing"
)

func TestGetRandomUniqueName(t *testing.T) {
	existingNames := []string{
		"Avaricious Merchant", "Covetous Dealer", "Rapacious Vendor",
		"Grasping Broker", "Voracious Retailer", "Acquisitive Wholesaler",
	}

	newName := GetRandomUniqueName(existingNames)

	if newName == "" {
		t.Errorf("Expected a non-empty name, but got an empty string")
	}
}
