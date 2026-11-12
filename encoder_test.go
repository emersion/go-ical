package ical

import (
	"bytes"
	"strings"
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

func TestEncoderInvalidPropError(t *testing.T) {
	for _, tc := range []struct {
		prop Prop
		want string
	}{
		{Prop{Name: PropSummary, Params: make(Params), Value: "a\rb"}, `"SUMMARY"`},
		{Prop{Name: PropAttendee, Params: Params{ParamCommonName: {`a"b`}}, Value: "mailto:a@example.org"}, `"ATTENDEE": param "CN"`},
	} {
		var buf bytes.Buffer
		err := NewEncoder(&buf).encodeProp(&tc.prop)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("encodeProp(%v) = %v, want an error naming %s", tc.prop.Name, err, tc.want)
		}
	}
}
