package render

import "io"

type TextRender struct {
	contentType string
	Body        string
}

func (t *TextRender) Render(writer io.Writer) error {
	_, err := writer.Write([]byte(t.Body))
	return err
}

func (t *TextRender) ContentType() string {
	return t.contentType
}

func Text(contentType string, body string) *TextRender {
	return &TextRender{contentType: contentType, Body: body}
}
