package display

import (
	"io"
)

type Displayer interface {
	Display(interface{}) error
}

type View interface {
	Header(io.Writer) error
	Body(io.Writer, interface{}) error
}

type Renderer struct {
	w io.Writer
}

func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

func (r *Renderer) Render(view View, v interface{}) error {
	if err := view.Header(r.w); err != nil {
		return err
	}

	if err := view.Body(r.w, v); err != nil {
		return err
	}
	return nil
}
