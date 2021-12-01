package axml

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type StringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

//MarshalXML
func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

//UnmarshalXML
func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

//Encode
func Encode(data map[string]string) ([]byte, error) {
	result, err := xml.Marshal(StringMap(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Encode
func Decode(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	if err := xml.Unmarshal(data, (*StringMap)(&result)); err != nil {
		return nil, err
	}
	return result, nil
}

//DecodeTo
func DecodeTo(data []byte, result map[string]string) error {
	if err := xml.Unmarshal(data, (*StringMap)(&result)); err != nil {
		return err
	}
	return nil
}

//ToJson
func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		return json.Marshal(r)
	}
}
