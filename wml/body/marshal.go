package body

import (
	"encoding/xml"

	"github.com/vortex/docx-go/xmltypes"
)

// ---------------------------------------------------------------------------
// CT_Document — MarshalXML
// ---------------------------------------------------------------------------

// MarshalXML writes the <w:document> root element with preserved namespaces.
func (doc *CT_Document) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "w:document"}

	if len(doc.Namespaces) > 0 {
		start.Attr = doc.Namespaces
	} else {
		start.Attr = defaultDocumentNamespaces()
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Body
	if doc.Body != nil {
		bodyStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "body"}}
		if err := e.EncodeElement(doc.Body, bodyStart); err != nil {
			return err
		}
	}

	// Extra (unknown children of <w:document> besides <w:body>)
	for i := range doc.Extra {
		if err := e.EncodeElement(&doc.Extra[i], xml.StartElement{Name: doc.Extra[i].XMLName}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// ---------------------------------------------------------------------------
// CT_Body — MarshalXML
// ---------------------------------------------------------------------------

// MarshalXML writes the <w:body> element in document order:
// block-level content first, then the trailing sectPr.
func (b *CT_Body) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Block-level content in order
	for _, el := range b.Content {
		if err := marshalBlockElement(e, el); err != nil {
			return err
		}
	}

	// Trailing section properties
	if b.SectPr != nil {
		sectStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "sectPr"}}
		if err := e.EncodeElement(b.SectPr, sectStart); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// ---------------------------------------------------------------------------
// CT_SdtBlock — MarshalXML
// ---------------------------------------------------------------------------

// MarshalXML writes the <w:sdt> element.
func (sdt *CT_SdtBlock) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if sdt.SdtPr != nil {
		prStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "sdtPr"}}
		if err := e.EncodeElement(sdt.SdtPr, prStart); err != nil {
			return err
		}
	}
	if sdt.SdtEndPr != nil {
		eprStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "sdtEndPr"}}
		if err := e.EncodeElement(sdt.SdtEndPr, eprStart); err != nil {
			return err
		}
	}

	// sdtContent wrapper
	if len(sdt.SdtContent) > 0 {
		contentStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "sdtContent"}}
		if err := e.EncodeToken(contentStart); err != nil {
			return err
		}
		for _, el := range sdt.SdtContent {
			if err := marshalBlockElement(e, el); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(contentStart.End()); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// ---------------------------------------------------------------------------
// marshalBlockElement dispatches a BlockLevelElement to the encoder.
// ---------------------------------------------------------------------------

func marshalBlockElement(e *xml.Encoder, el interface{}) error {
	switch v := el.(type) {
	case ParagraphElement:
		pStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "p"}}
		return e.EncodeElement(v.P, pStart)
	case TableElement:
		tStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "tbl"}}
		return e.EncodeElement(v.T, tStart)
	case SdtBlockElement:
		sdtStart := xml.StartElement{Name: xml.Name{Space: xmltypes.NSw, Local: "sdt"}}
		return e.EncodeElement(v.Sdt, sdtStart)
	case RawBlockElement:
		return e.EncodeElement(&v.Raw, xml.StartElement{Name: v.Raw.XMLName})
	}
	return nil
}

// ---------------------------------------------------------------------------
// Default namespaces for new documents
// ---------------------------------------------------------------------------

func defaultDocumentNamespaces() []xml.Attr {
	return []xml.Attr{
		{Name: xml.Name{Local: "xmlns:wpc"}, Value: "http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"},
		{Name: xml.Name{Local: "xmlns:mc"}, Value: xmltypes.NSmc},
		{Name: xml.Name{Local: "xmlns:o"}, Value: xmltypes.NSo},
		{Name: xml.Name{Local: "xmlns:r"}, Value: xmltypes.NSr},
		{Name: xml.Name{Local: "xmlns:m"}, Value: xmltypes.NSm},
		{Name: xml.Name{Local: "xmlns:v"}, Value: xmltypes.NSv},
		{Name: xml.Name{Local: "xmlns:wp14"}, Value: "http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"},
		{Name: xml.Name{Local: "xmlns:wp"}, Value: xmltypes.NSwp},
		{Name: xml.Name{Local: "xmlns:w10"}, Value: xmltypes.NSw10},
		{Name: xml.Name{Local: "xmlns:w"}, Value: xmltypes.NSw},
		{Name: xml.Name{Local: "xmlns:w14"}, Value: xmltypes.NSw14},
		{Name: xml.Name{Local: "xmlns:w15"}, Value: xmltypes.NSw15},
		{Name: xml.Name{Local: "xmlns:wpg"}, Value: "http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"},
		{Name: xml.Name{Local: "xmlns:wpi"}, Value: "http://schemas.microsoft.com/office/word/2010/wordprocessingInk"},
		{Name: xml.Name{Local: "xmlns:wne"}, Value: "http://schemas.microsoft.com/office/word/2006/wordml"},
		{Name: xml.Name{Local: "xmlns:wps"}, Value: "http://schemas.microsoft.com/office/word/2010/wordprocessingShape"},
		{Name: xml.Name{Local: "mc:Ignorable"}, Value: "w14 w15 wp14"},
	}
}
