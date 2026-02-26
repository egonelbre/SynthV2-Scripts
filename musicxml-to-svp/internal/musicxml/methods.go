package musicxml

func (v YesNo) Bool() bool {
	switch v {
	case Yes:
		return true
	case No:
		return false
	default:
		panic("invalid value")
	}
}

func (v AboveBelow) IsAbove() bool {
	switch v {
	case Above:
		return true
	case Below:
		return false
	default:
		panic("invalid value")
	}
}

func (v AboveBelow) IsBelow() bool { return !v.IsAbove() }
