package documentgenerator

type TextDocument struct {
	BaseGenerator
}

func NewTextDocument() *TextDocument {
	doc := &TextDocument{}
	doc.Doc = doc //field promotion
	return doc
}

func (t *TextDocument) PrepareData() string {
	return "This is the raw text data."
}

func (t *TextDocument) FormatContent(data string) string {
	return "Formatted Text: " + data
}

func (t *TextDocument) Save(content string) string {
	return "Saving text document: " + content
}
