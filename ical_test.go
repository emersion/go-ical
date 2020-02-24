package ical

import (
	"strings"
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
