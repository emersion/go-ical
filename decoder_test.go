package ical

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecoder(t *testing.T) {
	dec := NewDecoder(strings.NewReader(exampleCalendarStr))

	cal, err := dec.DecodeCalendar()
	if err != nil {
		t.Fatalf("DecodeCalendar() = %v", err)
	}
	if !reflect.DeepEqual(cal, exampleCalendar) {
		t.Errorf("DecodeCalendar() = \n%#v\nbut want:\n%#v", cal, exampleCalendar)
	}

	if _, err := dec.DecodeCalendar(); err != io.EOF {
		t.Errorf("DecodeCalendar() = %v, want io.EOF", err)
	}
}
