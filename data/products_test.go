package data

import "testing"

// Test Case to Validate Struct
func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "Shubham Dhage",
		Price: 20,
		SKU: "abc-abcd-abcde",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
