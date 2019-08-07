package xsens

// UTCValidity represents the validity of an Xsens UTC timestamp.
type UTCValidity uint8

func (u UTCValidity) IsDateValid() bool {
	return u&0x01 > 0
}

func (u UTCValidity) IsTimeOfDayValid() bool {
	return u&0x02 > 0
}

func (u UTCValidity) IsTimeOfDayFullyResolved() bool {
	return u&0x04 > 0
}
