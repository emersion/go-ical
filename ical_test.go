package ical

import (
	"strings"
	"testing"
	"time"
)

func toCRLF(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

var exampleCalendarStr = toCRLF(`BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
CATEGORIES:CONFERENCE
DESCRIPTION;ALTREP="cid:part1.0001@example.org":Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta\, Georgia
DTEND:19960920T220000Z
DTSTAMP:19960704T120000Z
DTSTART:19960918T143000Z
ORGANIZER:mailto:jsmith@example.com
STATUS:CONFIRMED
SUMMARY;FOO=bar,"b:az":Networld+Interop Conference
UID:uid1@example.com
END:VEVENT
END:VCALENDAR
`)

var exampleCalendar = &Calendar{&Component{
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
	Children: []*Component{
		{
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
						"FOO": []string{"bar", "b:az"},
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

func TestCalendar(t *testing.T) {
	events := exampleCalendar.Events()
	if len(events) != 1 {
		t.Fatalf("len(Calendar.Events()) = %v, want 1", len(events))
	}
	event := events[0]

	wantSummary := "Networld+Interop Conference"
	if summary, err := event.Properties.Text(PropSummary); err != nil {
		t.Errorf("Event.Properties.Text(PropSummary) = %v", err)
	} else if summary != wantSummary {
		t.Errorf("Event.Properties.Text(PropSummary) = %v, want %v", summary, wantSummary)
	}

	wantDesc := "Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta, Georgia"
	if desc, err := event.Properties.Text(PropDescription); err != nil {
		t.Errorf("Event.Properties.Text(PropDescription) = %v", err)
	} else if desc != wantDesc {
		t.Errorf("Event.Properties.Text(PropDescription) = %v, want %v", desc, wantDesc)
	}

	wantDTStamp := time.Date(1996, 07, 04, 12, 0, 0, 0, time.UTC)
	if dtStamp, err := event.Properties.DateTime(PropDateTimeStamp, nil); err != nil {
		t.Errorf("Event.Properties.DateTime(PropDateTimeStamp) = %v", err)
	} else if dtStamp != wantDTStamp {
		t.Errorf("Event.Properties.DateTime(PropDateTimeStamp) = %v, want %v", dtStamp, wantDTStamp)
	}
}
