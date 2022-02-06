package render

import (
	"encoding/json"
	"io"
)

type JsonRender struct {
	Data any
}

func JSON(data any) *JsonRender {
	return &JsonRender{Data: data}
}

func (j *JsonRender) Render(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(j.Data)
}

func (JsonRender) ContentType() string {
	return "application/json; charset=utf-8"
}
