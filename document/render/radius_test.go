package render

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/document/layout"
	"github.com/gpdf-dev/gpdf/pdf"
)

func TestRoundedRectClampsRadius(t *testing.T) {
	path := roundedRect(document.Rectangle{X: 10, Y: 20, Width: 40, Height: 10}, 100)
	start := path.Segments[0].Points[0]
	if start != (document.Point{X: 15, Y: 20}) {
		t.Errorf("start = %+v, want radius clamped to half the height", start)
	}
	if last := path.Segments[len(path.Segments)-2].Points[2]; last != start {
		t.Errorf("path ends at %+v, want %+v", last, start)
	}
}

func TestBorderRadius(t *testing.T) {
	bg := pdf.RGBHex(0xF5F7F7)
	var img bytes.Buffer
	if err := png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name       string
		node       document.DocumentNode
		want, deny string
	}{
		{"square box", &document.Box{BoxStyle: document.BoxStyle{Background: &bg}}, " re\n", " c\n"},
		{"rounded box", &document.Box{BoxStyle: document.BoxStyle{Background: &bg, BorderRadius: document.Pt(6)}}, " c\n", " re\n"},
		{"rounded image", &document.Image{
			Source:   document.ImageSource{Data: img.Bytes(), Format: document.ImagePNG, Width: 2, Height: 2},
			ImgStyle: document.Style{BorderRadius: document.Pt(6)},
		}, "W n\n", ""},
	} {
		r, _ := newTestRenderer(t)
		if err := r.BeginPage(document.Size{Width: 595, Height: 842}); err != nil {
			t.Fatal(err)
		}
		pn := layout.PlacedNode{Node: tc.node, Position: document.Point{X: 10, Y: 10}, Size: document.Size{Width: 100, Height: 50}}
		if err := r.renderPlacedNode(pn, 0, 0); err != nil {
			t.Fatal(err)
		}
		content := string(r.pageContent)
		if !strings.Contains(content, tc.want) || tc.deny != "" && strings.Contains(content, tc.deny) {
			t.Errorf("%s: content %q, want %q and not %q", tc.name, content, tc.want, tc.deny)
		}
	}
}
