package cell

import (
	"image"
	"reflect"
	"testing"
)

const twentyFourBits = 0xffffff

var cellTests = []struct {
	in        Field
	expectErr bool
}{
	{Field{Width: 0, Height: 0}, false},
	{Field{Width: 1, Height: 1}, false},
	{Field{Width: 2, Height: 4}, false},
}

func TestNewFieldReturnsField(t *testing.T) {
	for _, td := range cellTests {
		w := td.in.Width
		h := td.in.Height
		got := NewField(w, h)
		if len(got.State) != w*h {
			t.Errorf(
				"NewField(%d, %d) => len(c.State) = %d, want %d",
				w,
				h,
				len(got.State),
				w*h)
		}
		if got.Width != w {
			t.Errorf("NewField(%d, %d) => c.Width = %d, want %d",
				w,
				h,
				got.Width,
				w)
		}
		if got.Height != h {
			t.Errorf(
				"NewField(%d, %d) => c.Height = %d, want %d",
				w,
				h,
				got.Height,
				h)
		}
	}
}

func TestRandomFieldReturns24BitField(t *testing.T) {
	got := RandomField(128, 128)
	for i, s := range got.State {
		if s > twentyFourBits {
			t.Errorf(
				"RandomField(128, 128)[%d] => %x, want <= %x",
				i,
				s,
				twentyFourBits)

		}
	}
}

var imgTests = []struct {
	in   Field
	want image.NRGBA
}{
	{
		in:   Field{},
		want: *image.NewNRGBA(image.Rect(0, 0, 0, 0)),
	},
	{
		in: Field{
			State:  []uint32{0xff000000},
			Width:  1,
			Height: 1,
		},
		want: image.NRGBA{
			Pix:    []uint8{0, 0, 0, 0xff},
			Stride: 4,
			Rect:   image.Rect(0, 0, 1, 1),
		},
	},
	{
		in: Field{
			State:  []uint32{0xabcd0000},
			Width:  1,
			Height: 1,
		},
		want: image.NRGBA{
			Pix:    []uint8{0xcd, 0, 0, 0xff},
			Stride: 4,
			Rect:   image.Rect(0, 0, 1, 1),
		},
	},
	{
		in: Field{
			State:  []uint32{0xf000da00},
			Width:  1,
			Height: 1,
		},
		want: image.NRGBA{
			Pix:    []uint8{0, 0xda, 0, 0xff},
			Stride: 4,
			Rect:   image.Rect(0, 0, 1, 1),
		},
	},
	{
		in: Field{
			State:  []uint32{0xf00000da},
			Width:  1,
			Height: 1,
		},
		want: image.NRGBA{
			Pix:    []uint8{0, 0, 0xda, 0xff},
			Stride: 4,
			Rect:   image.Rect(0, 0, 1, 1),
		},
	},
	{
		in: Field{
			State:  []uint32{0xf00000da},
			Width:  1,
			Height: 1,
		},
		want: image.NRGBA{
			Pix:    []uint8{0, 0, 0xda, 0xff},
			Stride: 4,
			Rect:   image.Rect(0, 0, 1, 1),
		},
	},
	{
		in: Field{
			State:  []uint32{0x11223344, 0x55667788, 0x99001122, 0x33445566},
			Width:  2,
			Height: 2,
		},
		want: image.NRGBA{
			Pix: []uint8{
				0x22, 0x33, 0x44, 0xff,
				0x66, 0x77, 0x88, 0xff,
				0x00, 0x11, 0x22, 0xff,
				0x44, 0x55, 0x66, 0xff,
			},
			Stride: 8,
			Rect:   image.Rect(0, 0, 2, 2),
		},
	},
}

func TestToImageValidFieldReturnsImage(t *testing.T) {
	for _, td := range imgTests {
		got := ToImage(&td.in)
		if got == nil {
			t.Errorf("ToImage(&%q) => img: nil, want %v", td.in, td.want)
			continue
		}
		if !reflect.DeepEqual(got, &td.want) {
			t.Errorf("ToImage(&%q) => %#v, want %#v", td.in, *got, td.want)
			continue
		}
	}
}
