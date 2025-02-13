package scan

// WhiteSpace
const (
	NULL            byte = byte(0)
	HORIZONTAL_TAB  byte = byte(9)
	LINE_FEED       byte = byte(10)
	FORM_FEED       byte = byte(12)
	CARRIAGE_RETURN byte = byte(13)
	SPACE           byte = byte(32)
)

// Delimiter
const (
	LEFT_PARENTHESIS     byte = byte(40)
	RIGHT_PARENTHESIS    byte = byte(41)
	LESS_THAN            byte = byte(60)
	GREATER_THAN         byte = byte(62)
	LEFT_SQUARE_BRACKET  byte = byte(91)
	RIGTH_SQUARE_BRACKET byte = byte(93)
	LEFT_CURLY_BRACKET   byte = byte(123)
	RIGHT_CURLY_BRACKET  byte = byte(125)
	SOLIDUS              byte = byte(47)
	PERCENT_SIGN         byte = byte(37)
)

func isWhitespace(b byte) bool {
	return b == SPACE ||
		b == CARRIAGE_RETURN ||
		b == LINE_FEED ||
		b == FORM_FEED ||
		b == HORIZONTAL_TAB ||
		b == NULL
}
