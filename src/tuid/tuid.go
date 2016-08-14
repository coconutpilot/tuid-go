package tuid

import (
	"fmt"
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

func New(spec string) (*tuid, error) {
	// XXX: temp logging
	fmt.Printf("Decoding TUID spec: %s\n", spec)

	ctx := &tuid{random: 0x5248c8561600f46d}

	bitpos := uint8(64)

	max := len(spec)
	i := 0

	// this string parser is a direct translation from C which means
	// it isn't idiomatic Go.  Undecided if I should rewrite it or
	// keep the libraries as similar as possible ...
	for i < max {
		identifier := spec[i]
		i++

		var value uint64
		for i < max {
			v := spec[i]
			if v < '0' || v > '9' {
				break
			}
			i++
			value = value*10 + uint64(v-'0')
		}
		// XXX: temp logging
		fmt.Printf("Key: %v Val: %v\n", string(identifier), value)

		switch identifier {
		case 'E':
			ctx.nsec_offset = value

		case 'N':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			ctx.nsec_mask = (^uint64(0)) >> (64 - value)
			bitpos -= uint8(value)
			ctx.nsec_shift = bitpos

		case 'C':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			ctx.counter_max = (^uint64(0)) >> (64 - value)
			ctx.counter = ctx.counter_max /* this forces initialization */
			bitpos -= uint8(value)
			ctx.counter_shift = uint8(bitpos)

		case 'R':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			ctx.random_mask = (^uint64(0)) >> (64 - value)
			bitpos -= uint8(value)
			ctx.random_shift = bitpos

		case 'I':
			ctx.id = value

		default:
			return nil, fmt.Errorf("tuid spec error at: %v%v", identifier, value)
		}

	}
	return ctx, nil
}
