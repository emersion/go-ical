package ical_test

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/emersion/go-ical"
)

func ExampleDecoder() {
	// Let's assume r is an io.Reader containing iCal data
	var r io.Reader

	dec := ical.NewDecoder(r)
	for {
		cal, err := dec.DecodeCalendar()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for _, event := range cal.Events() {
			summary, err := event.Properties.Text(ical.PropSummary)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Found event: %v", summary)
		}
	}
}

func ExampleEncoder() {
	event := ical.NewEvent()
	event.Properties.SetText(ical.PropUID, "uid@example.org")
	event.Properties.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	event.Properties.SetText(ical.PropSummary, "My awesome event")
	event.Properties.SetDateTime(ical.PropDateTimeStart, time.Now())

	cal := ical.NewCalendar()
	cal.Children = append(cal.Children, event.Component)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf); err != nil {
		log.Fatal(err)
	}

	log.Print(buf.String())
}
