package ical

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (enc *Encoder) encodeProperty(prop *Property) error {
	var buf bytes.Buffer
	buf.WriteString(prop.Name)

	paramNames := make([]string, 0, len(prop.Params))
	for name := range prop.Params {
		paramNames = append(paramNames, name)
	}
	sort.Strings(paramNames)

	for _, name := range paramNames {
		buf.WriteString(";")
		buf.WriteString(name)
		buf.WriteString("=")

		for i, v := range prop.Params[name] {
			if i > 0 {
				buf.WriteString(",")
			}
			if strings.ContainsRune(v, '"') {
				return fmt.Errorf("ical: failed to encode param value: contains a double-quote")
			}
			if strings.ContainsAny(v, ";:,") {
				buf.WriteString(`"` + v + `"`)
			} else {
				buf.WriteString(v)
			}
		}
	}

	buf.WriteString(":")
	if strings.ContainsAny(prop.Value, "\r\n") {
		return fmt.Errorf("ical: failed to encode property value: contains a CR or LF")
	}
	buf.WriteString(prop.Value)
	buf.WriteString("\r\n")

	_, err := enc.w.Write(buf.Bytes())
	return err
}

func (enc *Encoder) encodeComponent(comp *Component) error {
	err := enc.encodeProperty(&Property{Name: "BEGIN", Value: comp.Name})
	if err != nil {
		return err
	}

	propNames := make([]string, 0, len(comp.Properties))
	for name := range comp.Properties {
		propNames = append(propNames, name)
	}
	sort.Strings(propNames)

	for _, name := range propNames {
		for _, prop := range comp.Properties[name] {
			if err := enc.encodeProperty(&prop); err != nil {
				return err
			}
		}
	}

	for _, child := range comp.Children {
		if err := enc.encodeComponent(&child); err != nil {
			return err
		}
	}

	return enc.encodeProperty(&Property{Name: "END", Value: comp.Name})
}

func (enc *Encoder) EncodeCalendar(cal *Calendar) error {
	return enc.encodeComponent(&cal.Component)
}
