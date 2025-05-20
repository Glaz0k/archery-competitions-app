package enums

type Bow string

const (
	BowClassic            Bow = "classic"
	BowBlock              Bow = "block"
	BowClassicNewbie      Bow = "classic_newbie"
	Bow3DClassic          Bow = "3D_classic"
	Bow3DCompound         Bow = "3D_compound"
	Bow3DLong             Bow = "3D_long"
	BowPeripheral         Bow = "peripheral"
	BowPeripheralWithRing Bow = "peripheral_with_ring"
)

func (s Bow) IsValid() bool {
	switch s {
	case BowClassic, BowBlock, BowClassicNewbie, Bow3DClassic, Bow3DCompound, Bow3DLong, BowPeripheral, BowPeripheralWithRing:
		return true
	}
	return false
}
