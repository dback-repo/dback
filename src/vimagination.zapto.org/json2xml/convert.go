// Package json2xml converts a JSON structure to XML.
//
// json2xml wraps each type within xml tags named after the type. For example:-
//
// An object is wrapped in `<object></object>`<br />
// An array is wrapped in `<array></array>`<br />
// A boolean is wrapped in `<boolean></boolean>` , with either "true" or "false" as chardata<br />
// A number is wrapped in `<number></number>`<br />
// A string is wrapped in `<string></string>`<br />
// A null becomes `<null></null>` , with no chardata
//
// When a type is a member of an object, the name of the key becomes an
// attribute on the type tag, for example: -
//
// {
// 	"Location": {
// 		"Longitude": -1.8262,
// 		"Latitude": 51.1789
// 	}
// }
//
// ...becomes...
//
// `<object>
//	<object name="Location">
//		<number name="Longitude">-1.8262</number>
// 		<number name="Latitude">51.1789</number>
// 	</object>
// </object>`
package json2xml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

type ttype byte

const (
	typObject ttype = iota
	typArray
	typBool
	typNumber
	typString
	typNull
)

var ttypeNames = [...]string{"object", "array", "boolean", "number", "string", "null"}

// JSONDecoder represents a type that gives out JSON tokens, usually
// implemented by *json.Decoder
// It is encouraged for implementers of this interface to output numbers using
// the json.Number type, as it reduces needless conversions.
// Users of the json.Decoder implementation should call the UseNumber method to
// achieve this
type JSONDecoder interface {
	Token() (json.Token, error)
}

// XMLEncoder represents a type that takes XML tokens, usually implemented by
// *xml.Encoder
type XMLEncoder interface {
	EncodeToken(xml.Token) error
}

// Converter represents the ongoing conversion from JSON to XML
type Converter struct {
	decoder JSONDecoder
	types   []ttype
	data    *string
}

// Tokens provides a JSON converter that implements the xml.TokenReader
// interface
func Tokens(j JSONDecoder) *Converter {
	return &Converter{
		decoder: j,
	}
}

// Token gets a xml.Token from the Converter, as per the xml.TokenReader
// interface
func (c *Converter) Token() (xml.Token, error) {
	if c.data != nil {
		token := xml.CharData(*c.data)
		c.data = nil
		return token, nil
	}
	if len(c.types) > 0 {
		switch c.types[len(c.types)-1] {
		case typObject, typArray:
		default:
			return c.outputEnd(), nil
		}
	}
	var keyName *string
	token, err := c.decoder.Token()
	if err != nil {
		return nil, err
	}
	if len(c.types) > 0 && c.types[len(c.types)-1] == typObject && token != json.Delim('}') {
		tokenStr, ok := token.(string)
		if !ok {
			return nil, ErrInvalidKey
		}
		keyName = &tokenStr
		token, err = c.decoder.Token()
		if err != nil {
			return nil, err
		}
	}
	switch token := token.(type) {
	case json.Delim:
		switch token {
		case '{':
			return c.outputStart(typObject, keyName), nil
		case '[':
			return c.outputStart(typArray, keyName), nil
		case '}':
			if len(c.types) == 0 || c.types[len(c.types)-1] != typObject {
				return nil, ErrInvalidToken
			}
			return c.outputEnd(), nil
		case ']':
			if len(c.types) == 0 || c.types[len(c.types)-1] != typArray {
				return nil, ErrInvalidToken
			}
			return c.outputEnd(), nil
		default:
			return nil, ErrUnknownToken
		}
	case bool:
		if token {
			return c.outputType(typBool, &cTrue, keyName), nil
		}
		return c.outputType(typBool, &cFalse, keyName), nil
	case float64:
		number := strconv.FormatFloat(token, 'f', -1, 64)
		return c.outputType(typNumber, &number, keyName), nil
	case json.Number:
		return c.outputType(typNumber, (*string)(&token), keyName), nil
	case string:
		return c.outputType(typString, &token, keyName), nil
	case nil:
		return c.outputType(typNull, nil, keyName), nil
	default:
		return nil, ErrUnknownToken
	}
}

func (c *Converter) outputType(typ ttype, data *string, keyName *string) xml.Token {
	c.data = data
	return c.outputStart(typ, keyName)
}

func (c *Converter) outputStart(typ ttype, keyName *string) xml.Token {
	c.types = append(c.types, typ)
	var attr []xml.Attr
	if keyName != nil {
		attr = []xml.Attr{
			xml.Attr{
				Name: xml.Name{
					Local: "name",
				},
				Value: *keyName,
			},
		}
	}
	return xml.StartElement{
		Name: xml.Name{
			Local: ttypeNames[typ],
		},
		Attr: attr,
	}
}

func (c *Converter) outputEnd() xml.Token {
	typ := c.types[len(c.types)-1]
	c.types = c.types[:len(c.types)-1]
	return xml.EndElement{
		Name: xml.Name{
			Local: ttypeNames[typ],
		},
	}
}

// Convert converts JSON and sends it to the given XML encoder
func Convert(j JSONDecoder, x XMLEncoder) error {
	c := Converter{
		decoder: j,
	}
	for {
		tk, err := c.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err = x.EncodeToken(tk); err != nil {
			return err
		}
	}
}

var (
	cTrue  = "true"
	cFalse = "false"
)

// Errors
var (
	ErrInvalidKey   = errors.New("invalid key type")
	ErrUnknownToken = errors.New("unknown token type")
	ErrInvalidToken = errors.New("invalid token")
)
