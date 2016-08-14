package tuid

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func Test_New(t *testing.T) {

	//tuid.Init("N54C8R2☺")

	log.Printf("Testing a valid tuid spec")
	tuid, err := New("N54C8R2")
	if tuid == nil || err != nil {
		t.Errorf("Decoding valid tuid failed, err: %v, tuid: %v", err, tuid)
	}
	log.Printf("tuid: %v", tuid)

	log.Printf("Testing invalid tuid spec")
	tuid, err = New("N54C8R2☺")
	if err != nil && tuid != nil {
		t.Errorf("Decoding invalid tuid succeeded, err: %v, tuid: %v", err, tuid)
	}
}
