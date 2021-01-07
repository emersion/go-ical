package ical

import (
	"bytes"
	"testing"
)

func TestEncoder(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(exampleCalendar); err != nil {
		t.Fatalf("Encode() = %v", err)
	}
	s := buf.String()

	if s != exampleCalendarStr {
		t.Errorf("Encode() = \n%v\nbut want:\n%v", s, exampleCalendarStr)
	}
}

func TestEncoder_encodeProp(t *testing.T) {
	enc := &Encoder{
		maxLineLength: 13,
	}

	tests := []struct {
		name string
		prop Prop
		want string
	}{
		{name: "short", prop: Prop{Name: "FOO", Value: "BAR"}, want: "FOO:BAR"},
		{name: "multibyte", prop: Prop{Name: "A", Value: "Ḽơᶉëᶆ ȋṕšᶙṁ"}, want: "A:Ḽơᶉë\r\n ᶆ ȋṕš\r\n ᶙṁ"},
		{name: "exact length", prop: Prop{Name: "FOUR", Value: "+ eight!"}, want: "FOUR:+ eight!"},
		{
			name: "exceeding line limit",
			prop: Prop{
				Name:  "FOO",
				Value: "Exceeding line limit",
			},
			want: "FOO:Exceeding\r\n  line limit",
		},
		{
			name: "with params",
			prop: Prop{
				Name: "A",
				Params: Params{
					ParamEncoding: []string{"8bit"},
				},
				Value: "B",
			},
			want: "A;ENCODING=8b\r\n it:B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := enc
			enc.w = &buf

			err := enc.encodeProp(&tt.prop)
			if err != nil {
				t.Errorf("encodeProp() error = %v", err)
			}

			got := buf.String()
			want := tt.want + "\r\n"
			if got != want {
				t.Errorf("expected %q, but got %q", want, got)
			}
		})
	}
}
