// Package footnotes implements parsing and serialization of footnotes.xml and
// endnotes.xml parts (CT_Footnotes / CT_FtnEdn) in an OOXML word-processing
// document.
//
// Both footnotes and endnotes share the same schema; the only difference is
// the root element name (w:footnotes vs w:endnotes) and the child element
// name (w:footnote vs w:endnote). The exported types handle both cases.
//
// See contracts.md C-26 and reference-appendix.md §footnotes.
package footnotes

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/vortex/docx-go/wml/shared"
)

// Namespace URI for the WordprocessingML main namespace.
const nsW = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

// ---------------------------------------------------------------------------
// CT_Footnotes — root element of footnotes.xml (or endnotes.xml)
// ---------------------------------------------------------------------------

// CT_Footnotes represents the <w:footnotes> (or <w:endnotes>) root element.
type CT_Footnotes struct {
	Footnote []CT_FtnEdn

	// rootLocal is the local name of the root XML element that was parsed.
	// It is "footnotes" for footnotes.xml and "endnotes" for endnotes.xml.
	// When zero-valued, MarshalXML defaults to "footnotes".
	rootLocal string

	// childLocal is the local name of each child element.
	// It is "footnote" or "endnote".  Defaults to "footnote".
	childLocal string

	// namespaces preserves the original xmlns:* declarations on the root
	// element so they survive a round-trip without loss.
	namespaces []xml.Attr
}

// UnmarshalXML implements xml.Unmarshaler for CT_Footnotes.
func (fn *CT_Footnotes) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Remember root element name ("footnotes" or "endnotes").
	fn.rootLocal = start.Name.Local
	switch fn.rootLocal {
	case "footnotes":
		fn.childLocal = "footnote"
	case "endnotes":
		fn.childLocal = "endnote"
	default:
		fn.childLocal = "footnote"
	}

	// Preserve namespace declarations for round-trip fidelity.
	fn.namespaces = make([]xml.Attr, len(start.Attr))
	copy(fn.namespaces, start.Attr)

	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			// We expect <w:footnote> or <w:endnote> children.
			if t.Name.Local == fn.childLocal {
				var fe CT_FtnEdn
				if err := fe.unmarshal(d, t); err != nil {
					return err
				}
				fn.Footnote = append(fn.Footnote, fe)
			} else {
				// Unknown child — skip.
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// MarshalXML implements xml.Marshaler for CT_Footnotes.
func (fn CT_Footnotes) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	local := fn.rootLocal
	if local == "" {
		local = "footnotes"
	}
	childLocal := fn.childLocal
	if childLocal == "" {
		childLocal = "footnote"
	}

	start := xml.StartElement{
		Name: xml.Name{Local: "w:" + local},
	}

	if len(fn.namespaces) > 0 {
		start.Attr = fn.namespaces
	} else {
		start.Attr = defaultNamespaces()
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for i := range fn.Footnote {
		if err := fn.Footnote[i].marshal(e, childLocal); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// defaultNamespaces returns the minimal set of xmlns declarations for a
// newly-created footnotes/endnotes root element.
func defaultNamespaces() []xml.Attr {
	return []xml.Attr{
		{Name: xml.Name{Local: "xmlns:w"}, Value: nsW},
	}
}

// ---------------------------------------------------------------------------
// CT_FtnEdn — a single footnote or endnote entry
// ---------------------------------------------------------------------------

// CT_FtnEdn represents a <w:footnote> or <w:endnote> element.
type CT_FtnEdn struct {
	Type    *string `xml:"type,attr,omitempty"` // "normal"|"separator"|"continuationSeparator"
	ID      int     `xml:"id,attr"`
	Content []shared.BlockLevelElement
}

// unmarshal decodes a <w:footnote>/<w:endnote> element from the decoder.
// It is called from CT_Footnotes.UnmarshalXML after reading the start element.
func (fe *CT_FtnEdn) unmarshal(d *xml.Decoder, start xml.StartElement) error {
	// Parse attributes (w:type, w:id).
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "type":
			v := attr.Value
			fe.Type = &v
		case "id":
			id, err := strconv.Atoi(attr.Value)
			if err != nil {
				return fmt.Errorf("footnotes: invalid id %q: %w", attr.Value, err)
			}
			fe.ID = id
		}
	}

	// Parse child elements — block-level content (paragraphs, tables, …).
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			// Try the registered factory first.
			if el := shared.CreateBlockElement(t.Name); el != nil {
				if err := d.DecodeElement(el, &t); err != nil {
					return err
				}
				fe.Content = append(fe.Content, el)
			} else {
				// Unknown element — preserve as RawXML for round-trip.
				var raw shared.RawXML
				if err := d.DecodeElement(&raw, &t); err != nil {
					return err
				}
				fe.Content = append(fe.Content, raw)
			}
		case xml.EndElement:
			return nil
		}
	}
}

// marshal encodes a <w:footnote>/<w:endnote> element into the encoder.
func (fe *CT_FtnEdn) marshal(e *xml.Encoder, childLocal string) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "w:" + childLocal},
	}

	// Write attributes using w: prefix.
	if fe.Type != nil {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "w:type"},
			Value: *fe.Type,
		})
	}
	start.Attr = append(start.Attr, xml.Attr{
		Name:  xml.Name{Local: "w:id"},
		Value: strconv.Itoa(fe.ID),
	})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode block-level children.
	for _, bl := range fe.Content {
		switch v := bl.(type) {
		case shared.RawXML:
			if err := v.MarshalXML(e, xml.StartElement{}); err != nil {
				return err
			}
		default:
			// Typed elements (CT_P, CT_Tbl, …) implement xml.Marshaler.
			if m, ok := bl.(xml.Marshaler); ok {
				if err := m.MarshalXML(e, xml.StartElement{}); err != nil {
					return err
				}
			}
		}
	}

	return e.EncodeToken(start.End())
}

// ---------------------------------------------------------------------------
// Parse / Serialize — top-level API
// ---------------------------------------------------------------------------

// Parse deserializes XML bytes into a CT_Footnotes structure.
func Parse(data []byte) (*CT_Footnotes, error) {
	var fn CT_Footnotes
	if err := xml.Unmarshal(data, &fn); err != nil {
		return nil, fmt.Errorf("footnotes.Parse: %w", err)
	}
	return &fn, nil
}

// Serialize serializes a CT_Footnotes structure back to XML bytes.
// The output includes the XML declaration.
func Serialize(fn *CT_Footnotes) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	if err := enc.Encode(fn); err != nil {
		return nil, fmt.Errorf("footnotes.Serialize: %w", err)
	}
	if err := enc.Flush(); err != nil {
		return nil, fmt.Errorf("footnotes.Serialize: %w", err)
	}

	return buf.Bytes(), nil
}
