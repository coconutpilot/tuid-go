package tuid

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {

	fmt.Printf("# Testing invalid tuid spec")
	tuider, err := New("N54C8R2☺")
	if tuider != nil || err == nil {
		t.Errorf("Decoding invalid tuid spec succeeded")
	}

	fmt.Printf("# Testing a valid tuid spec")
	tuider, err = New("N54C8R2")
	if tuider == nil || err != nil {
		t.Errorf("Decoding valid tuid spec failed, err: %v", err)
	}
	fmt.Printf("# tuider: %v", tuider)

}
