package documentgenerator

type HTMLDocument struct {
	BaseGenerator
}

func NewHTMLDocument() *HTMLDocument {
	doc := &HTMLDocument{}
	doc.Doc = doc //field promotion, promote the doc field, then it can be used in BaseGenerator. :D
	return doc
}

func (h *HTMLDocument) PrepareData() string {
	return "<html><body>This is raw HTML data.</body></html>"
}

func (h *HTMLDocument) FormatContent(data string) string {
	return "<div>" + data + "</div>"
}

func (h *HTMLDocument) Save(content string) string {
	return "Saving HTML document: " + content
}
