package pdf

type Font interface {
	Indirect(f File) *Indirect
	Name() string
}

type StandardFont uint8

const ( TimesRoman StandardFont = iota
	Helvetica
	Courier
	Symbol
	TimesBold
	HelveticaBold
	CourierBold
	ZapfDingbats
	TimesItalic
	HelveticaOblique
	CourierOblique
	TimesBoldItalic
	HelveticaBoldOblique
	CourierBoldOblique )

func StandardFontToName (font StandardFont)  (result string) {
	switch (font) {
	case Helvetica:
		result = "Helvetica"
	case Courier:
		result = "Courier"
	case Symbol:
		result = "Symbol"
	case TimesBold:
		result = "Times-Bold"
	case HelveticaBold:
		result = "Helvetica-Bold"
	case CourierBold:
		result = "Courier-Bold"
	case ZapfDingbats:
		result = "ZapfDingbats"
	case TimesItalic:
		result = "Times-Italic"
	case HelveticaOblique:
		result = "Helvetica-Oblique"
	case CourierOblique:
		result = "Courier-Oblique"
	case TimesBoldItalic:
		result = "Times-BoldItalic"
	case HelveticaBoldOblique:
		result = "Helvetica-BoldOblique"
	case CourierBoldOblique:
		result = "Courier-BoldOblique"
	}
	return result
}

type standardFont struct {
	fileBindings map[File]*Indirect
	name string
	descriptor *Dictionary
}

func NewStandardFont(font StandardFont, name string) Font {
	result := new(standardFont)
	result.fileBindings = make(map[File]*Indirect,5)
	result.name = name
	result.descriptor = NewDictionary()
	result.descriptor.Add ("Type", NewName("Font"))
	result.descriptor.Add ("Subtype", NewName("Type1"))
	result.descriptor.Add ("Name", NewName(result.name))
	result.descriptor.Add ("BaseFont", NewName(StandardFontToName(font)))
	result.descriptor.Add ("Encoding", NewName("MacRomanEncoding"))
	return result
}

func (font *standardFont) Indirect(file File) *Indirect {
	i,ok := font.fileBindings[file]
	if (!ok) {
		i = file.AddObject (font.descriptor)
		font.fileBindings[file] = i
	}
	return i
}

func (font *standardFont) Name() string {
	return font.name
}

