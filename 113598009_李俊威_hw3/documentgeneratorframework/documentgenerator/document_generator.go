package documentgenerator

type DocumentGenerator interface {
	PrepareData() string
	FormatContent(data string) string
	Save(content string) string
}
