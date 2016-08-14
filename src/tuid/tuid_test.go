package tuid

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {

	fmt.Printf("# Testing invalid TUID spec\n")
	tuider, err := New("N54C8R2☺")
	if tuider != nil || err == nil {
		t.Errorf("Decoding invalid TUID spec succeeded")
	}

	fmt.Printf("# Testing a valid TUID spec\n")
	tuider, err = New("N54C8R2")
	if tuider == nil || err != nil {
		t.Errorf("Decoding valid TUID spec failed, err: %v", err)
	}
	fmt.Printf("# tuider: %v\n", tuider)
}

func Test_Gen1(t *testing.T) {

	tuider, err := New("N5I1")
	if err != nil {
		t.Errorf("Decoding valid TUID spec failed, err: %v", err)
		return
	}

	last := uint64(1) // first TUID will be 0
	for i := 0; i < 100; i++ {
		id := tuider.Gen()
		fmt.Printf("# TUID: %b\n", id)
		if id == last {
			t.Errorf("Duplicated TUID: %v\n", id)
			return
		}
		last = id
	}
}

func Test_Gen2(t *testing.T) {

	tuider, err := New("N60I1")
	if err != nil {
		t.Errorf("Decoding valid TUID spec failed, err: %v", err)
		return
	}

	var last uint64
	for i := 0; i < 100; i++ {
		id := tuider.Gen()
		fmt.Printf("# TUID: %X\n", id)
		if id == last {
			t.Errorf("Duplicated TUID: %v\n", id)
			return
		}
		last = id
	}
}

func Test_Overflow_Epoch1(t *testing.T) {
	_, err := New("N54C8R2E12345678901234567890")
	if err == nil {
		t.Errorf("Verifying Epoch overflow failed\n")
	}
}

func Test_Overflow_Nanosec1(t *testing.T) {
	_, err := New("N65")
	if err == nil {
		t.Errorf("Verifying Nanosec overflow failed\n")
	}
}

func Test_Overflow_Nanosec2(t *testing.T) {
	_, err := New("C10N55")
	if err == nil {
		t.Errorf("Verifying Nanosec overflow failed\n")
	}
}

func Test_Overflow_Counter1(t *testing.T) {
	_, err := New("C65")
	if err == nil {
		t.Errorf("Verifying Counter overflow failed\n")
	}
}

func Test_Overflow_Counter2(t *testing.T) {
	_, err := New("N54C11")
	if err == nil {
		t.Errorf("Verifying Counter overflow failed\n")
	}
}

func Test_Overflow_Random1(t *testing.T) {
	_, err := New("R65")
	if err == nil {
		t.Errorf("Verifying Random overflow failed\n")
	}
}

func Test_Overflow_Random2(t *testing.T) {
	_, err := New("N54R11")
	if err == nil {
		t.Errorf("Verifying Random overflow failed\n")
	}
}

func Test_Overflow_ID(t *testing.T) {
	_, err := New("N54C8R2I12345678901234567890")
	if err == nil {
		t.Errorf("Verifying ID overflow failed\n")
	}
}
