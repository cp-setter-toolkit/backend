package memory

import "strconv"

type Amount int

const (
	Byte Amount = 1
	KB   Amount = 1000 * Byte
	KiB  Amount = 1024 * Byte
	MB   Amount = 1000 * KB
	MiB  Amount = 1024 * KiB
	GB   Amount = 1000 * MB
	GiB  Amount = 1024 * MiB
)

func (a Amount) String() string {
	switch {
	case a >= GiB:
		return strconv.FormatFloat(float64(a)/float64(GiB), 'f', 3, 64) + "GiB"
	case a >= MiB:
		return strconv.FormatFloat(float64(a)/float64(MiB), 'f', 3, 64) + "MiB"
	case a >= KiB:
		return strconv.FormatFloat(float64(a)/float64(KiB), 'f', 3, 64) + "KiB"
	default:
		return strconv.Itoa(int(a)) + "B"
	}
}
