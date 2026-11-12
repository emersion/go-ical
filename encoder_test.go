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

func TestEncoderParamLineBreak(t *testing.T) {
	for _, v := range []string{"a\r\nATTACH:b", "a\rb", "a\nb"} {
		prop := Prop{Name: PropAttendee, Params: Params{ParamCommonName: {v}}, Value: "mailto:a@example.org"}
		if err := NewEncoder(&bytes.Buffer{}).encodeProp(&prop); err == nil {
			t.Errorf("encodeProp() with param value %q = nil, want an error", v)
		}
	}
}
