package xsens

// UTCValidity represents the validity of an Xsens UTC timestamp.
type UTCValidity uint8

const (
	UTCDateValidFlag              = 0x01
	UTCTimeOfDayValidFlag         = 0x02
	UTCTimeOfDayFullyResolvedFlag = 0x04
)

func (u UTCValidity) IsDateValid() bool {
	return u&UTCDateValidFlag > 0
}

func (u UTCValidity) IsTimeOfDayValid() bool {
	return u&UTCTimeOfDayValidFlag > 0
}

func (u UTCValidity) IsTimeOfDayFullyResolved() bool {
	return u&UTCTimeOfDayFullyResolvedFlag > 0
}
