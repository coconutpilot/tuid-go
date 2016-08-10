package tuid

import (
	"log"
)

type tuid struct {
	nsec_offset        uint64 /* nsec offset from unixtime epoch */
	nsec               uint64 /* unixtime in nsec, adjusted by offset */
	nsec_mask          uint64 /* the bits of nsec to be used in the TUID */
	nsec_shift         uint8  /* how far to left shift the nsec bits */
	nsec_min_increment uint64 /* minimum increment required if nsec mask was too coarse */
	id                 uint64 /* static id, ie machine id */
	counter            uint64 /* a counter, current value */
	counter_max        uint64 /* when to reset the counter to 0 */
	counter_shift      uint8  /* how far to left shift the counter bits */
	random             uint64 /* a random #, also used as seed */
	random_mask        uint64 /* the bits of the random # to be used */
	random_shift       uint8  /* XXX: turns out this is not needed */
}

//func init() {
//	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
//}

func (ctx *tuid) Init(spec string) {
	log.Println("tuid.Init()")

	ctx.nsec_offset = 0
	ctx.nsec = 0
	ctx.nsec_mask = 0
	ctx.nsec_shift = 0

	ctx.id = 0

	ctx.counter = 0
	ctx.counter_max = 0
	ctx.counter_shift = 0

	ctx.random = 0x5248c8561600f46d
	ctx.random_mask = 0
	ctx.random_shift = 0

	log.Printf("Decoding TUID spec: %s\n", spec)

}
