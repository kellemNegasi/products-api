package data

import "testing"

func TestValidator(t *testing.T) {
	p := &Product{
		Name:  "test product",
		Price: 1.0,
		SKU:   "abc-abc-abcd-",
	}
	err := p.validate()
	if err != nil {
		t.Fatal(err)
	}
}
