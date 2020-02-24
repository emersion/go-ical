package ical

import (
	"bytes"
	"testing"
)

func TestEncoder(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).EncodeCalendar(exampleCalendar); err != nil {
		t.Fatalf("Encode() = %v", err)
	}
	s := buf.String()

	if s != exampleCalendarStr {
		t.Errorf("Encode() = \n%v\nbut want:\n%v", s, exampleCalendarStr)
	}
}
