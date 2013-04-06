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
	return map[StandardFont] string {
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

