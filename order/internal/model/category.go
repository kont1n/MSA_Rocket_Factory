package model

type Category int

const (
	UNKNOWN  Category = iota // 0
	ENGINE                   // 1
	FUEL                     // 2
	PORTHOLE                 // 3
	WING                     // 4
)

func (c Category) String() string {
	switch c {
	case ENGINE:
		return "ENGINE"
	case FUEL:
		return "FUEL"
	case PORTHOLE:
		return "PORTHOLE"
	case WING:
		return "WING"
	default:
		return "UNKNOWN"
	}
}

func CategoryName(value int) string {
	return Category(value).String()
}

func ToCategory(value int) Category {
	switch value {
	case int(ENGINE), int(FUEL), int(PORTHOLE), int(WING):
		return Category(value)
	default:
		return UNKNOWN
	}
}
