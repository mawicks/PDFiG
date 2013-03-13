/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

const (
	regularCharacter = iota
	whiteSpaceCharacter
	delimiterCharacter
)

var (
	whiteSpaceCharacters = "\000\t\n\f\r "
	delimiterCharacters = "%()/<>[]{}"
	characterTypes [256]byte
)

func init() {
	for i := range characterTypes {
		characterTypes[i] = regularCharacter
	}

	for _,w := range whiteSpaceCharacters {
		characterTypes[w] = whiteSpaceCharacter
	}

	for _,d := range delimiterCharacters {
		characterTypes[d] = delimiterCharacter
	}
}

func  IsWhiteSpace (b byte) bool {
	return characterTypes[b] == whiteSpaceCharacter
}

func IsDelimiter (b byte) bool {
	return characterTypes[b] == delimiterCharacter
}

func IsRegular (b byte) bool {
	return characterTypes[b] == regularCharacter
}

