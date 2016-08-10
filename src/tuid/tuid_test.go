package tuid

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func Test_New(t *testing.T) {

	ctx := &tuid{}

	ctx.Init("N54C8R2")

	return
}
