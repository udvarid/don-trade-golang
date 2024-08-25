package user

import (
	"slices"
	"testing"
)

func TestGetRandomUniqueName(t *testing.T) {
	existingNames := []string{
		"Avaricious Merchant", "Covetous Dealer", "Rapacious Vendor",
		"Grasping Broker", "Voracious Retailer", "Acquisitive Wholesaler",
	}

	newName := GetRandomUniqueName(existingNames)

	if slices.Contains(existingNames, newName) {
		t.Errorf("Expected a unique name, but got an existing name: %s", newName)
	}

	if newName == "" {
		t.Errorf("Expected a non-empty name, but got an empty string")
	}
}
