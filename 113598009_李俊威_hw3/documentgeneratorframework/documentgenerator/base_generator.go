package documentgenerator

type BaseGenerator struct {
	Doc DocumentGenerator
}

func (b *BaseGenerator) Generate() string {
	data := b.Doc.PrepareData()
	formatted := b.Doc.FormatContent(data)
	return b.Doc.Save(formatted)
}
