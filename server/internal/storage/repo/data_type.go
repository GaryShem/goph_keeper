package repo

type DataType string

var (
	TextType   DataType = "text"
	CardType   DataType = "card"
	BinaryType DataType = "binary"
)

func (d DataType) IsValid() bool {
	switch d {
	case TextType, CardType, BinaryType:
		return true
	default:
		return false
	}
}
