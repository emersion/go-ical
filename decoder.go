package ical

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type lineDecoder struct {
	s string
}

func (ld *lineDecoder) decodeName() (string, error) {
	//name begin with A-Z
	i := strings.IndexAny(ld.s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if i != 0 {
		return "", fmt.Errorf("name not start with A-Z")
	}
	i = strings.IndexAny(ld.s, ";:")
	if i < 0 {
		return "", fmt.Errorf("ical: malformed content line: missing colon")
	} else if i == 0 {
		return "", fmt.Errorf("ical: malformed content line: empty property name")
	}

	name := strings.ToUpper(ld.s[:i])
	ld.s = ld.s[i:]
	return name, nil
}

func (ld *lineDecoder) empty() bool {
	return len(ld.s) == 0
}

func (ld *lineDecoder) peek() byte {
	return ld.s[0]
}

func (ld *lineDecoder) consume(c byte) bool {
	if ld.empty() || ld.peek() != c {
		return false
	}
	ld.s = ld.s[1:]
	return true
}

func (ld *lineDecoder) decodeParamValue() (string, error) {
	var v string
	if ld.consume('"') {
		for !ld.empty() && ld.peek() != '"' {
			v += ld.s[:1]
			ld.s = ld.s[1:]
		}

		if !ld.consume('"') {
			return "", fmt.Errorf("ical: malformed param value: unterminated quoted string")
		}
	} else {
	Loop:
		for !ld.empty() {
			switch c := ld.peek(); c {
			case '"':
				return "", fmt.Errorf("ical: malformed param value: illegal double-quote")
			case ';', ',', ':':
				break Loop
			default:
				v += ld.s[:1]
				ld.s = ld.s[1:]
			}
		}
	}

	return v, nil
}

func (ld *lineDecoder) decodeParam() (string, []string, error) {
	i := strings.IndexByte(ld.s, '=')
	if i < 0 {
		return "", nil, fmt.Errorf("ical: malformed param: missing equal sign")
	} else if i == 0 {
		return "", nil, fmt.Errorf("ical: malformed param: empty param name")
	}

	name := strings.ToUpper(ld.s[:i])
	ld.s = ld.s[i+1:]

	var values []string
Loop:
	for {
		value, err := ld.decodeParamValue()
		if err != nil {
			return "", nil, err
		}
		values = append(values, value)

		switch c := ld.peek(); c {
		case ',':
			ld.s = ld.s[1:]
		case ';', ':':
			break Loop
		default:
			panic(fmt.Errorf("ical: unexpected character %q after decoding param value", c))
		}
	}

	return name, values, nil
}

func (ld *lineDecoder) decodeContentLine() (*Prop, error) {
	name, err := ld.decodeName()
	if err != nil {
		return nil, err
	}

	params := make(map[string][]string)
	for ld.consume(';') {
		paramName, paramValues, err := ld.decodeParam()
		if err != nil {
			return nil, err
		}
		params[paramName] = append(params[paramName], paramValues...)
	}

	if !ld.consume(':') {
		return nil, fmt.Errorf("ical: malformed property: expected colon")
	}

	return &Prop{
		Name:   name,
		Params: params,
		Value:  ld.s,
	}, nil
}

type Decoder struct {
	br *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}
}

func (dec *Decoder) readLine() ([]byte, error) {
	var buf []byte
	for {
		line, isPrefix, err := dec.br.ReadLine()
		if err != nil {
			return nil, err
		}

		if !isPrefix && len(buf) == 0 {
			return line, err
		}

		buf = append(buf, line...)
		if !isPrefix {
			break
		}
	}
	return buf, nil
}

func (dec *Decoder) readContinuedLine() (string, error) {
	var sb strings.Builder

	l, err := dec.readLine()
	if err != nil {
		return "", err
	}
	sb.Write(l)

	for {
		r, _, err := dec.br.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if r != ' ' && r != '\t' {
			dec.br.UnreadRune()
			break
		}

		l, err := dec.readLine()
		if err != nil {
			return "", err
		}
		sb.Write(l)
	}

	return sb.String(), nil
}

func (dec *Decoder) decodeContentLine(lp *Prop) (*Prop, error) {
	for {
		l, err := dec.readContinuedLine()
		if err != nil {
			return nil, err
		}
		if len(l) == 0 {
			continue
		}

		ld := lineDecoder{l}
		p, err := ld.decodeContentLine()
		if err != nil && lp != nil {
			lp.Value += l
			return nil, nil
		}
		return p, err
	}
}

func (dec *Decoder) decodeComponentBody(name string) (*Component, error) {
	var prop *Prop
	props := make(Props)
	var children []*Component
	var lastprop *Prop
Loop:
	for {
		var err error
		//prop maybe multiline, add content to pre prop value and return nil nil
		prop, err = dec.decodeContentLine(lastprop)
		if err != nil {
			return nil, err
		}
		if prop == nil {
			continue
		}
		switch prop.Name {
		case "BEGIN":
			child, err := dec.decodeComponentBody(strings.ToUpper(prop.Value))
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		case "END":
			break Loop
		default:
			props[prop.Name] = append(props[prop.Name], *prop)
		}
		lastprop = prop
	}

	if prop.Name != "END" {
		panic("ical: expected END property")
	}
	if !strings.EqualFold(prop.Value, name) {
		return nil, fmt.Errorf("ical: malformed component: expected END property for %q, got %q", name, prop.Value)
	}

	return &Component{
		Name:     name,
		Props:    props,
		Children: children,
	}, nil
}

func (dec *Decoder) decodeComponent() (*Component, error) {
	prop, err := dec.decodeContentLine(nil)
	if err != nil {
		return nil, err
	}
	if prop.Name != "BEGIN" {
		return nil, fmt.Errorf("ical: malformed component: expected BEGIN property, got %q", prop.Name)
	}

	return dec.decodeComponentBody(strings.ToUpper(prop.Value))
}

func (dec *Decoder) Decode() (*Calendar, error) {
	comp, err := dec.decodeComponent()
	if err != nil {
		return nil, err
	} else if comp.Name != CompCalendar {
		return nil, fmt.Errorf("ical: invalid toplevel component name: expected %q, got %q", CompCalendar, comp.Name)
	}

	return &Calendar{comp}, nil
}
