package ical

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const exampleCalendarStr = `BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
DTSTAMP:19960704T120000Z
UID:uid1@example.com
ORGANIZER:mailto:jsmith@example.com
DTSTART:19960918T143000Z
DTEND:19960920T220000Z
STATUS:CONFIRMED
CATEGORIES:CONFERENCE
SUMMARY;FOO=bar,"baz":Networld+Interop Conference
DESCRIPTION;ALTREP="cid:part1.0001@example.org":Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta\, Georgia
END:VEVENT
END:VCALENDAR`

var exampleCalendar = &Calendar{Component{
	Name: "VCALENDAR",
	Properties: Properties{
		"PRODID": []Property{{
			Name:   "PRODID",
			Params: Params{},
			Value:  "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN",
		}},
		"VERSION": []Property{{
			Name:   "VERSION",
			Params: Params{},
			Value:  "2.0",
		}},
	},
	Children: []Component{
		Component{
			Name: "VEVENT",
			Properties: Properties{
				"DTSTAMP": []Property{{
					Name:   "DTSTAMP",
					Params: Params{},
					Value:  "19960704T120000Z",
				}},
				"UID": []Property{{
					Name:   "UID",
					Params: Params{},
					Value:  "uid1@example.com",
				}},
				"ORGANIZER": []Property{{
					Name:   "ORGANIZER",
					Params: Params{},
					Value:  "mailto:jsmith@example.com",
				}},
				"DTSTART": []Property{{
					Name:   "DTSTART",
					Params: Params{},
					Value:  "19960918T143000Z",
				}},
				"DTEND": []Property{{
					Name:   "DTEND",
					Params: Params{},
					Value:  "19960920T220000Z",
				}},
				"STATUS": []Property{{
					Name:   "STATUS",
					Params: Params{},
					Value:  "CONFIRMED",
				}},
				"CATEGORIES": []Property{{
					Name:   "CATEGORIES",
					Params: Params{},
					Value:  "CONFERENCE",
				}},
				"SUMMARY": []Property{{
					Name: "SUMMARY",
					Params: Params{
						"FOO": []string{"bar", "baz"},
					},
					Value: "Networld+Interop Conference",
				}},
				"DESCRIPTION": []Property{{
					Name: "DESCRIPTION",
					Params: Params{
						"ALTREP": []string{"cid:part1.0001@example.org"},
					},
					Value: `Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta\, Georgia`,
				}},
			},
		},
	},
}}

func TestDecoder(t *testing.T) {
	dec := NewDecoder(strings.NewReader(exampleCalendarStr))

	cal, err := dec.DecodeCalendar()
	if err != nil {
		t.Fatalf("DecodeCalendar() = %v", err)
	}
	if !reflect.DeepEqual(cal, exampleCalendar) {
		t.Errorf("DecodeCalendar() = \n%#v\n, want \n%#v", cal, exampleCalendar)
	}

	if _, err := dec.DecodeCalendar(); err != io.EOF {
		t.Errorf("DecodeCalendar() = %v, want io.EOF", err)
	}
}
