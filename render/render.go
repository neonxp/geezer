package render

import (
	"io"
)

type Renderer interface {
	Render(io.Writer) error
	ContentType() string
}
