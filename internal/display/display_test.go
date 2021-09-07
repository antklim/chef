package display_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/antklim/chef/internal/display"
	"github.com/stretchr/testify/assert"
)

type testData [3]string
type testView struct{}

func (v *testView) Header(w io.Writer) error {
	_, err := fmt.Fprintln(w, "FIRST\tSECOND\tTHIRD")
	return err
}

func (v *testView) Body(w io.Writer, data interface{}) error {
	d, ok := data.(testData)
	if !ok {
		return errors.New("invalid body")
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\n", d[0], d[1], d[2])
	return err
}

var _ display.View = (*testView)(nil)

func TestRenderer(t *testing.T) {
	var buf bytes.Buffer
	renderer := display.NewRenderer(&buf)
	td := testData{"1", "2", "3"}

	err := renderer.Render(&testView{}, td)
	assert.NoError(t, err)

	expected := fmt.Sprintln("FIRST\tSECOND\tTHIRD\n1\t2\t3")
	assert.Equal(t, expected, buf.String())
}

type testDisplay struct {
	w io.Writer

	headerFormat string
	rowFormat    string
}

func newTestDisplay(w io.Writer) *testDisplay {
	return &testDisplay{
		w:            w,
		headerFormat: "FIRST\tSECOND\tTHIRD",
		rowFormat:    "%s\t%s\t%s\n",
	}
}

func (d *testDisplay) Display(v interface{}) error {
	rows, ok := v.([]testData)
	if !ok {
		return errors.New("display error")
	}

	_, err := fmt.Fprintln(d.w, d.headerFormat)
	if err != nil {
		return err
	}

	for _, row := range rows {
		_, err := fmt.Fprintf(d.w, d.rowFormat, row[0], row[1], row[2])
		if err != nil {
			return err
		}
	}
	return nil
}

var _ display.Displayer = (*testDisplay)(nil)

func TestDisplay(t *testing.T) {
	var buf bytes.Buffer
	displayer := newTestDisplay(&buf)
	td := []testData{{"1", "2", "3"}}

	err := displayer.Display(td)
	assert.NoError(t, err)

	expected := fmt.Sprintln("FIRST\tSECOND\tTHIRD\n1\t2\t3")
	assert.Equal(t, expected, buf.String())
}
