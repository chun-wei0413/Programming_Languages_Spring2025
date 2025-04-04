package documentgenerator

type BaseGenerator struct {
	Doc DocumentGenerator
}

func (b *BaseGenerator) Generate() string {
	if b.Doc == nil {
		return "Error: no document generator provided"
	}
	data := b.Doc.PrepareData()
	formatted := b.Doc.FormatContent(data)
	return b.Doc.Save(formatted)
}
