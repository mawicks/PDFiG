package pdf

type Font interface {
	Indirect(f File) *Indirect
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
	return map[StandardFont] string {
		TimesRoman: "Times-Roman",
		Helvetica: "Helvetica",
		Courier: "Courier",
		Symbol: "Symbol",
		TimesBold: "Times-Bold",
		HelveticaBold: "Helvetica-Bold",
		CourierBold: "Courier-Bold",
		ZapfDingbats: "ZapfDingbats",
		TimesItalic: "Times-Italic",
		HelveticaOblique: "Helvetica-Oblique",
		CourierOblique: "Courier-Oblique",
		TimesBoldItalic: "Times-BoldItalic",
		HelveticaBoldOblique: "Helvetica-BoldOblique",
		CourierBoldOblique: "Courier-BoldOblique" } [font]
}

type standardFont struct {
	fileBindings map[File]*Indirect
	dictionary *Dictionary
}

func NewStandardFont(font StandardFont) Font {
	result := new(standardFont)
	result.fileBindings = make(map[File]*Indirect,5)
	result.dictionary = NewDictionary()
	result.dictionary.Add ("Type", NewName("Font"))
	result.dictionary.Add ("Subtype", NewName("Type1"))
	result.dictionary.Add ("BaseFont", NewName(StandardFontToName(font)))

	// Note: Internal encoding is used; no encoding is specified.
	// Note: Use of a /Name entry (not included here) is required
	// in PDF 1.0 and deprecated in later versions.
	return result
}

func (font *standardFont) Indirect(file File) *Indirect {
	i,exists := font.fileBindings[file]
	if (!exists) {
		i = file.AddObject (font.dictionary)
		font.fileBindings[file] = i
	}
	return i
}

