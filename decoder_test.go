package ical

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecoder(t *testing.T) {
	dec := NewDecoder(strings.NewReader(exampleCalendarStr))

	cal, err := dec.Decode()
	if err != nil {
		t.Fatalf("DecodeCalendar() = %v", err)
	}
	if !reflect.DeepEqual(cal, exampleCalendar) {
		t.Errorf("DecodeCalendar() = \n%#v\nbut want:\n%#v", cal, exampleCalendar)
	}

	if _, err := dec.Decode(); err != io.EOF {
		t.Errorf("DecodeCalendar() = %v, want io.EOF", err)
	}
}

func TestDecoderLongLine(t *testing.T) {
	template := `BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
DESCRIPTION:%v
DTSTAMP:19960704T120000Z
DTSTART:19960918T143000Z
UID:uid1@example.com
END:VEVENT
END:VCALENDAR
`
	description := strings.Repeat("Networld+Interop Conference", 500)
	calendarStr := toCRLF(fmt.Sprintf(template, description))
	calendar := &Calendar{&Component{
		Name: "VCALENDAR",
		Props: Props{
			"PRODID": []Prop{{
				Name:   "PRODID",
				Params: Params{},
				Value:  "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN",
			}},
			"VERSION": []Prop{{
				Name:   "VERSION",
				Params: Params{},
				Value:  "2.0",
			}},
		},
		Children: []*Component{
			{
				Name: "VEVENT",
				Props: Props{
					"DTSTAMP": []Prop{{
						Name:   "DTSTAMP",
						Params: Params{},
						Value:  "19960704T120000Z",
					}},
					"UID": []Prop{{
						Name:   "UID",
						Params: Params{},
						Value:  "uid1@example.com",
					}},
					"DTSTART": []Prop{{
						Name:   "DTSTART",
						Params: Params{},
						Value:  "19960918T143000Z",
					}},
					"DESCRIPTION": []Prop{{
						Name:   "DESCRIPTION",
						Params: Params{},
						Value:  description,
					}},
				},
			},
		},
	}}

	dec := NewDecoder(strings.NewReader(calendarStr))

	cal, err := dec.Decode()
	if err != nil {
		t.Fatalf("DecodeCalendar() = %v", err)
	}
	if !reflect.DeepEqual(cal, calendar) {
		t.Errorf("DecodeCalendar() = \n%#v\nbut want:\n%#v", cal, calendar)
	}
}
