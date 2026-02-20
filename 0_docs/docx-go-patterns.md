# DOCX-Go: Паттерны реализации

> Этот файл решает конкретные инженерные проблемы, обнаруженные при аудите.
> Загружать вместе с `contracts.md` и `reference-appendix.md`.

---

# 1. Архитектура: разрыв циклических зависимостей

## 1.1 Проблема

XSD определяет рекурсивные структуры:
- `CT_Tc` (ячейка таблицы) содержит `EG_BlockLevelElts` → `CT_P` и `CT_Tbl`
- `CT_Tbl` содержит `CT_Tc`
- `CT_PPr` содержит `CT_ParaRPr` (тип из rpr)

В Go циклический импорт запрещён.

## 1.2 Решение: пакет `wml/shared`

```
docx-go/
├── wml/
│   ├── shared/       ← НОВЫЙ: интерфейсы и общие типы
│   ├── rpr/          ← импортирует: xmltypes
│   ├── ppr/          ← импортирует: xmltypes, wml/rpr, wml/shared
│   ├── sectpr/       ← импортирует: xmltypes, wml/shared
│   ├── table/        ← импортирует: xmltypes, wml/shared
│   ├── tracking/     ← импортирует: xmltypes, wml/shared
│   ├── run/          ← импортирует: xmltypes, wml/rpr, wml/shared, dml
│   ├── para/         ← импортирует: xmltypes, wml/rpr, wml/ppr, wml/run, wml/tracking, wml/shared
│   ├── body/         ← импортирует: xmltypes, wml/para, wml/table, wml/sectpr, wml/shared
│   └── hdft/         ← импортирует: wml/shared
```

## 1.3 Содержимое `wml/shared`

```go
package shared

import "encoding/xml"

// ==========================================
// ИНТЕРФЕЙСЫ КОНТЕНТА
// ==========================================

// BlockLevelElement — параграф, таблица, SDT, или неизвестный элемент.
// Используется в: body, tc, header, footer, comment, footnote.
type BlockLevelElement interface {
	blockLevelElement()
}

// ParagraphContent — run, hyperlink, bookmark, ins/del, или неизвестный.
// Используется в: para, hyperlink, fldSimple.
type ParagraphContent interface {
	paragraphContent()
}

// RunContent — текст, br, drawing, fldChar, tab, или неизвестный.
// Используется в: run.
type RunContent interface {
	runContent()
}

// ==========================================
// RAW XML ДЛЯ ROUND-TRIP
// ==========================================

// RawXML хранит неизвестный XML-элемент целиком для round-trip.
// При unmarshal: если элемент не распознан → сохранить как RawXML.
// При marshal: восстановить на том же месте.
type RawXML struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Inner   []byte     `xml:",innerxml"`
}

// Реализации интерфейсов для RawXML (попадает куда угодно при round-trip)
func (RawXML) blockLevelElement() {}
func (RawXML) paragraphContent()  {}
func (RawXML) runContent()        {}

// ==========================================
// РЕГИСТРАЦИЯ — чтобы unmarshal знал какой тип создавать
// ==========================================

// BlockLevelFactory создаёт типизированный элемент по имени XML.
// Регистрируется пакетами body/para/table при init().
// Если имя неизвестно → вернуть nil (вызывающий сохранит как RawXML).
type BlockLevelFactory func(name xml.Name) BlockLevelElement

var blockFactories []BlockLevelFactory

func RegisterBlockFactory(f BlockLevelFactory) {
	blockFactories = append(blockFactories, f)
}

func CreateBlockElement(name xml.Name) BlockLevelElement {
	for _, f := range blockFactories {
		if el := f(name); el != nil {
			return el
		}
	}
	return nil
}

// Аналогично для ParagraphContent и RunContent
```

## 1.4 Исправленный граф зависимостей

```
wml/shared    → xmltypes                      (только типы, нет логики)
wml/rpr       → xmltypes, wml/shared
wml/ppr       → xmltypes, wml/rpr, wml/shared ← ppr ИМПОРТИРУЕТ rpr
wml/sectpr    → xmltypes, wml/shared
wml/table     → xmltypes, wml/shared           (CT_Tc.Content = []shared.BlockLevelElement)
wml/tracking  → xmltypes, wml/shared
dml           → xmltypes, wml/shared
wml/run       → xmltypes, wml/rpr, wml/shared, dml
wml/para      → xmltypes, wml/rpr, wml/ppr, wml/run, wml/tracking, wml/shared
wml/body      → xmltypes, wml/para, wml/table, wml/sectpr, wml/shared
wml/hdft      → wml/shared
```

**Почему `ppr` → `rpr` НЕ цикл**: rpr НЕ импортирует ppr. Зависимость односторонняя.

**Почему `table` ↔ `body` НЕ цикл**: оба импортируют `shared`, но не друг друга.
`CT_Tc.Content = []shared.BlockLevelElement` — это интерфейс из shared, а не тип из body.

---

# 2. Конвенция именования: Go field → XML element

## 2.1 Правило

```
Go PascalCase  →  XML camelCase  в пространстве имён w:

Примеры:
  Go: KeepNext       →  XML: <w:keepNext/>
  Go: PageBreakBefore → XML: <w:pageBreakBefore/>
  Go: BCs            →  XML: <w:bCs/>
  Go: SzCs           →  XML: <w:szCs/>
  Go: NumPr          →  XML: <w:numPr>
  Go: OutlineLvl     →  XML: <w:outlineLvl>
```

**Функция конверсии:**

```go
// pascalToXMLName конвертирует Go PascalCase → XML camelCase.
// Первая буква → lowercase, остальное без изменений.
// ИСКЛЮЧЕНИЯ: аббревиатуры в начале (BCs, ICs) уже в правильном регистре.
func pascalToXMLName(s string) string {
    if len(s) == 0 {
        return s
    }
    // Специальные случаи
    switch s {
    case "BCs":
        return "bCs"
    case "ICs":
        return "iCs"
    case "SzCs":
        return "szCs"
    }
    // Общее правило: первая буква → lowercase
    r := []rune(s)
    r[0] = unicode.ToLower(r[0])
    return string(r)
}
```

## 2.2 Полный маппинг CT_RPrBase

```go
// fieldName → xml element local name (в namespace w:)
var rprBaseFieldMap = []fieldMapping{
    {"RStyle",          "rStyle"},
    {"RFonts",          "rFonts"},
    {"B",               "b"},
    {"BCs",             "bCs"},
    {"I",               "i"},
    {"ICs",             "iCs"},
    {"Caps",            "caps"},
    {"SmallCaps",       "smallCaps"},
    {"Strike",          "strike"},
    {"Dstrike",         "dstrike"},
    {"Outline",         "outline"},
    {"Shadow",          "shadow"},
    {"Emboss",          "emboss"},
    {"Imprint",         "imprint"},
    {"NoProof",         "noProof"},
    {"SnapToGrid",      "snapToGrid"},
    {"Vanish",          "vanish"},
    {"WebHidden",       "webHidden"},
    {"Color",           "color"},
    {"Spacing",         "spacing"},
    {"W",               "w"},
    {"Kern",            "kern"},
    {"Position",        "position"},
    {"Sz",              "sz"},
    {"SzCs",            "szCs"},
    {"Highlight",       "highlight"},
    {"U",               "u"},
    {"Effect",          "effect"},
    {"Bdr",             "bdr"},
    {"Shd",             "shd"},
    {"FitText",         "fitText"},
    {"VertAlign",       "vertAlign"},
    {"Rtl",             "rtl"},
    {"Cs",              "cs"},
    {"Em",              "em"},
    {"Lang",            "lang"},
    {"EastAsianLayout", "eastAsianLayout"},
    {"SpecVanish",      "specVanish"},
    {"OMath",           "oMath"},
}
```

## 2.3 Полный маппинг CT_PPrBase (СТРОГИЙ ПОРЯДОК xsd:sequence!)

```go
var pprBaseFieldMap = []fieldMapping{
    {"PStyle",              "pStyle"},              //  1
    {"KeepNext",            "keepNext"},             //  2
    {"KeepLines",           "keepLines"},             //  3
    {"PageBreakBefore",     "pageBreakBefore"},       //  4
    {"FramePr",             "framePr"},               //  5
    {"WidowControl",        "widowControl"},          //  6
    {"NumPr",               "numPr"},                 //  7
    {"SuppressLineNumbers", "suppressLineNumbers"},   //  8
    {"PBdr",                "pBdr"},                  //  9
    {"Shd",                 "shd"},                   // 10
    {"Tabs",                "tabs"},                  // 11
    {"SuppressAutoHyphens", "suppressAutoHyphens"},   // 12
    {"Kinsoku",             "kinsoku"},               // 13
    {"WordWrap",            "wordWrap"},               // 14
    {"OverflowPunct",       "overflowPunct"},          // 15
    {"TopLinePunct",        "topLinePunct"},            // 16
    {"AutoSpaceDE",         "autoSpaceDE"},             // 17
    {"AutoSpaceDN",         "autoSpaceDN"},             // 18
    {"Bidi",                "bidi"},                    // 19
    {"AdjustRightInd",      "adjustRightInd"},          // 20
    {"SnapToGrid",          "snapToGrid"},              // 21
    {"Spacing",             "spacing"},                 // 22
    {"Ind",                 "ind"},                     // 23
    {"ContextualSpacing",   "contextualSpacing"},       // 24
    {"MirrorIndents",       "mirrorIndents"},            // 25
    {"SuppressOverlap",     "suppressOverlap"},          // 26
    {"Jc",                  "jc"},                       // 27
    {"TextDirection",       "textDirection"},             // 28
    {"TextAlignment",       "textAlignment"},             // 29
    {"TextboxTightWrap",    "textboxTightWrap"},          // 30
    {"OutlineLvl",          "outlineLvl"},                // 31
    {"DivId",               "divId"},                     // 32
    {"CnfStyle",            "cnfStyle"},                   // 33
}
```

## 2.4 Полный маппинг CT_TblPrBase (СТРОГИЙ ПОРЯДОК xsd:sequence!)

```go
var tblPrBaseFieldMap = []fieldMapping{
    {"TblStyle",            "tblStyle"},
    {"TblpPr",              "tblpPr"},
    {"TblOverlap",          "tblOverlap"},
    {"BidiVisual",          "bidiVisual"},
    {"TblStyleRowBandSize", "tblStyleRowBandSize"},
    {"TblStyleColBandSize", "tblStyleColBandSize"},
    {"TblW",                "tblW"},
    {"Jc",                  "jc"},
    {"TblCellSpacing",      "tblCellSpacing"},
    {"TblInd",              "tblInd"},
    {"TblBorders",          "tblBorders"},
    {"Shd",                 "shd"},
    {"TblLayout",           "tblLayout"},
    {"TblCellMar",          "tblCellMar"},
    {"TblLook",             "tblLook"},
    {"TblCaption",          "tblCaption"},
    {"TblDescription",      "tblDescription"},
}
```

## 2.5 Полный маппинг CT_TcPrBase (СТРОГИЙ ПОРЯДОК xsd:sequence!)

```go
var tcPrBaseFieldMap = []fieldMapping{
    {"CnfStyle",       "cnfStyle"},
    {"TcW",            "tcW"},
    {"GridSpan",       "gridSpan"},
    {"VMerge",         "vMerge"},
    {"TcBorders",      "tcBorders"},
    {"Shd",            "shd"},
    {"NoWrap",         "noWrap"},
    {"TcMar",          "tcMar"},
    {"TextDirection",  "textDirection"},
    {"TcFitText",      "tcFitText"},
    {"VAlign",         "vAlign"},
    {"HideMark",       "hideMark"},
    {"Headers",        "headers"},
}
```

## 2.6 Маппинг CT_TrPrBase (xsd:choice — порядок не строгий, но рекомендуется)

```go
var trPrBaseFieldMap = []fieldMapping{
    {"CnfStyle",       "cnfStyle"},
    {"DivId",          "divId"},
    {"GridBefore",     "gridBefore"},
    {"GridAfter",      "gridAfter"},
    {"WBefore",        "wBefore"},
    {"WAfter",         "wAfter"},
    {"CantSplit",      "cantSplit"},
    {"TrHeight",       "trHeight"},
    {"TblHeader",      "tblHeader"},
    {"TblCellSpacing", "tblCellSpacing"},
    {"Jc",             "jc"},
    {"Hidden",         "hidden"},
}
```

## 2.7 Маппинг EG_SectPrContents (СТРОГИЙ ПОРЯДОК xsd:sequence!)

```go
var sectPrFieldMap = []fieldMapping{
    // EG_HdrFtrReferences (max 6: default/first/even × header/footer)
    {"HeaderReference", "headerReference"},  // w:type + r:id attrs
    {"FooterReference", "footerReference"},
    // EG_SectPrContents
    {"FootnotePr",      "footnotePr"},
    {"EndnotePr",       "endnotePr"},
    {"Type",            "type"},
    {"PgSz",            "pgSz"},
    {"PgMar",           "pgMar"},
    {"PaperSrc",        "paperSrc"},
    {"PgBorders",       "pgBorders"},
    {"LnNumType",       "lnNumType"},
    {"PgNumType",       "pgNumType"},
    {"Cols",            "cols"},
    {"FormProt",        "formProt"},
    {"VAlign",          "vAlign"},
    {"NoEndnote",       "noEndnote"},
    {"TitlePg",         "titlePg"},
    {"TextDirection",   "textDirection"},
    {"Bidi",            "bidi"},
    {"RtlGutter",       "rtlGutter"},
    {"DocGrid",         "docGrid"},
    {"PrinterSettings", "printerSettings"},
    // extension
    {"SectPrChange",    "sectPrChange"},
}
```

## 2.8 Маппинг CT_Style (xsd:sequence)

```go
var styleFieldMap = []fieldMapping{
    {"Name",            "name"},
    {"Aliases",         "aliases"},
    {"BasedOn",         "basedOn"},
    {"Next",            "next"},
    {"Link",            "link"},
    {"AutoRedefine",    "autoRedefine"},
    {"Hidden",          "hidden"},
    {"UIpriority",      "uiPriority"},
    {"SemiHidden",      "semiHidden"},
    {"UnhideWhenUsed",  "unhideWhenUsed"},
    {"QFormat",         "qFormat"},
    {"Locked",          "locked"},
    {"Personal",        "personal"},
    {"PersonalCompose", "personalCompose"},
    {"PersonalReply",   "personalReply"},
    {"Rsid",            "rsid"},
    {"PPr",             "pPr"},
    {"RPr",             "rPr"},
    {"TblPr",           "tblPr"},
    {"TrPr",            "trPr"},
    {"TcPr",            "tcPr"},
    {"TblStylePr",      "tblStylePr"},
}
```

## 2.9 Маппинг CT_Lvl (xsd:sequence)

```go
var lvlFieldMap = []fieldMapping{
    {"Start",           "start"},
    {"NumFmt",          "numFmt"},
    {"LvlRestart",      "lvlRestart"},
    {"PStyle",          "pStyle"},
    {"IsLgl",           "isLgl"},
    {"Suff",            "suff"},
    {"LvlText",         "lvlText"},
    {"LvlPicBulletId",  "lvlPicBulletId"},
    {"LvlJc",           "lvlJc"},
    {"PPr",             "pPr"},
    {"RPr",             "rPr"},
}
```

## 2.10 Маппинг EG_RunInnerContent (xsd:choice)

```go
var runInnerContentMap = map[string]string{
    "br":                     "br",
    "t":                      "t",
    "contentPart":            "contentPart",
    "delText":                "delText",
    "instrText":              "instrText",
    "delInstrText":           "delInstrText",
    "noBreakHyphen":          "noBreakHyphen",
    "softHyphen":             "softHyphen",
    "dayShort":               "dayShort",
    "monthShort":             "monthShort",
    "yearShort":              "yearShort",
    "dayLong":                "dayLong",
    "monthLong":              "monthLong",
    "yearLong":               "yearLong",
    "annotationRef":          "annotationRef",
    "footnoteRef":            "footnoteRef",
    "endnoteRef":             "endnoteRef",
    "separator":              "separator",
    "continuationSeparator":  "continuationSeparator",
    "sym":                    "sym",
    "pgNum":                  "pgNum",
    "cr":                     "cr",
    "tab":                    "tab",
    "object":                 "object",
    "fldChar":                "fldChar",
    "ruby":                   "ruby",
    "footnoteReference":      "footnoteReference",
    "endnoteReference":       "endnoteReference",
    "commentReference":       "commentReference",
    "drawing":                "drawing",
    "ptab":                   "ptab",
    "lastRenderedPageBreak":  "lastRenderedPageBreak",
}
```

---

# 3. Паттерн: Кастомный MarshalXML для xsd:sequence

## 3.1 Проблема

Go `encoding/xml` с тегами `xml:"name"` не гарантирует порядок полей.
Поля сериализуются в порядке объявления в структуре, но:
- nil-поля пропускаются (это ОК)
- `Extra []RawXML` нужно вставить на правильные позиции

## 3.2 Решение: таблица полей + reflect

```go
package ppr

import (
    "encoding/xml"
    "reflect"
    "github.com/vortex/docx-go/xmltypes"
)

// fieldMapping описывает одно поле структуры
type fieldMapping struct {
    GoField  string // имя поля в Go-структуре
    XMLLocal string // локальное имя XML-элемента
}

// orderedFields — порядок элементов CT_PPrBase из XSD
var orderedFields = []fieldMapping{
    {"PStyle", "pStyle"},
    {"KeepNext", "keepNext"},
    {"KeepLines", "keepLines"},
    // ... (полный список из раздела 2.3)
}

func (p *CT_PPrBase) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    start.Name = xml.Name{Space: xmltypes.NSw, Local: "pPr"}
    e.EncodeToken(start)

    v := reflect.ValueOf(p).Elem()

    // Сериализуем типизированные поля в правильном порядке
    for _, fm := range orderedFields {
        fv := v.FieldByName(fm.GoField)
        if fv.IsNil() {
            continue
        }
        elemStart := xml.StartElement{
            Name: xml.Name{Space: xmltypes.NSw, Local: fm.XMLLocal},
        }
        if err := e.EncodeElement(fv.Interface(), elemStart); err != nil {
            return err
        }
    }

    // Сериализуем Extra (неизвестные элементы) — В КОНЦЕ
    for _, raw := range p.Extra {
        if err := e.EncodeElement(raw, xml.StartElement{Name: raw.XMLName}); err != nil {
            return err
        }
    }

    return e.EncodeToken(start.End())
}
```

## 3.3 Альтернатива БЕЗ reflect (рекомендуется для hot path)

```go
func (p *CT_PPrBase) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    e.EncodeToken(start)

    // Жёстко закодированный порядок — быстрее reflect
    if p.PStyle != nil {
        encodeChild(e, "pStyle", p.PStyle)
    }
    if p.KeepNext != nil {
        encodeOnOff(e, "keepNext", p.KeepNext)
    }
    if p.KeepLines != nil {
        encodeOnOff(e, "keepLines", p.KeepLines)
    }
    if p.PageBreakBefore != nil {
        encodeOnOff(e, "pageBreakBefore", p.PageBreakBefore)
    }
    // ... все 33 поля ...

    // Extra в конце
    for _, raw := range p.Extra {
        e.EncodeToken(xml.StartElement{Name: raw.XMLName, Attr: raw.Attrs})
        e.Flush()
        e.EncodeToken(xml.CharData(raw.Inner))
        e.EncodeToken(xml.EndElement{Name: raw.XMLName})
    }

    return e.EncodeToken(start.End())
}

// Хелперы
func encodeChild(e *xml.Encoder, local string, v interface{}) error {
    return e.EncodeElement(v, xml.StartElement{
        Name: xml.Name{Space: xmltypes.NSw, Local: local},
    })
}

func encodeOnOff(e *xml.Encoder, local string, o *xmltypes.CT_OnOff) error {
    start := xml.StartElement{
        Name: xml.Name{Space: xmltypes.NSw, Local: local},
    }
    if o.Val != nil {
        start.Attr = append(start.Attr, xml.Attr{
            Name: xml.Name{Space: xmltypes.NSw, Local: "val"},
            Value: *o.Val,
        })
    }
    e.EncodeToken(start)
    return e.EncodeToken(start.End())
}
```

---

# 4. Паттерн: RawXML round-trip

## 4.1 Unmarshal: перехват неизвестных элементов

```go
func (p *CT_PPrBase) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        tok, err := d.Token()
        if err != nil {
            return err
        }

        switch t := tok.(type) {
        case xml.StartElement:
            switch t.Name.Local {
            case "pStyle":
                p.PStyle = &xmltypes.CT_String{}
                if err := d.DecodeElement(p.PStyle, &t); err != nil {
                    return err
                }
            case "keepNext":
                p.KeepNext = &xmltypes.CT_OnOff{}
                if err := d.DecodeElement(p.KeepNext, &t); err != nil {
                    return err
                }
            // ... все известные элементы ...

            default:
                // НЕИЗВЕСТНЫЙ ЭЛЕМЕНТ → сохранить как RawXML для round-trip
                var raw shared.RawXML
                if err := d.DecodeElement(&raw, &t); err != nil {
                    return err
                }
                p.Extra = append(p.Extra, raw)
            }

        case xml.EndElement:
            return nil // конец pPr
        }
    }
}
```

## 4.2 Marshal: восстановление RawXML

```go
// В MarshalXML — после всех типизированных полей:
for _, raw := range p.Extra {
    // Восстанавливаем элемент как есть
    start := xml.StartElement{Name: raw.XMLName, Attr: raw.Attrs}
    e.EncodeToken(start)
    // Inner содержит всё между <start> и </start> включая детей
    e.Flush() // важно: flush перед записью raw bytes
    // Для encoding/xml нельзя писать raw bytes напрямую.
    // Вместо этого используем CharData или рекурсивный подход:
    if len(raw.Inner) > 0 {
        // Перепарсить и переиграть токены
        innerDec := xml.NewDecoder(bytes.NewReader(raw.Inner))
        for {
            innerTok, err := innerDec.Token()
            if err != nil {
                break
            }
            e.EncodeToken(xml.CopyToken(innerTok))
        }
    }
    e.EncodeToken(start.End())
}
```

## 4.3 Упрощённый вариант: храним весь элемент как []byte

```go
// Альтернативный подход — хранить полный XML элемента
type RawElement struct {
    FullXML []byte // "<w14:shadow ... >...</w14:shadow>"
}

// При unmarshal неизвестного элемента:
func captureRawElement(d *xml.Decoder, start xml.StartElement) ([]byte, error) {
    var buf bytes.Buffer
    enc := xml.NewEncoder(&buf)
    enc.EncodeToken(start)

    depth := 1
    for depth > 0 {
        tok, err := d.Token()
        if err != nil {
            return nil, err
        }
        enc.EncodeToken(xml.CopyToken(tok))
        switch tok.(type) {
        case xml.StartElement:
            depth++
        case xml.EndElement:
            depth--
        }
    }
    enc.Flush()
    return buf.Bytes(), nil
}

// При marshal — просто пишем bytes
// (через промежуточный decoder → encoder, т.к. xml.Encoder не принимает raw bytes)
```

---

# 5. Паттерн: Dual namespace (Transitional + Strict)

## 5.1 Проблема

99.9% документов используют Transitional namespaces, но Strict допустим.
Парсер должен принимать оба варианта. Генератор всегда пишет Transitional.

## 5.2 Решение: маппинг при чтении

```go
package xmltypes

// Strict → Transitional маппинг
var strictToTransitional = map[string]string{
    "http://purl.oclc.org/ooxml/wordprocessingml/main":              NSw,
    "http://purl.oclc.org/ooxml/officeDocument/relationships":       NSr,
    "http://purl.oclc.org/ooxml/drawingml/main":                     NSa,
    "http://purl.oclc.org/ooxml/drawingml/wordprocessingDrawing":    NSwp,
    "http://purl.oclc.org/ooxml/officeDocument/math":                NSm,
    "http://purl.oclc.org/ooxml/drawingml/picture":                  NSpic,
}

// NormalizeNamespace конвертирует Strict → Transitional.
// Вызывается при unmarshal для каждого элемента/атрибута.
func NormalizeNamespace(ns string) string {
    if mapped, ok := strictToTransitional[ns]; ok {
        return mapped
    }
    return ns
}

// WrapDecoder оборачивает xml.Decoder для автоматической нормализации
type NormalizingDecoder struct {
    inner *xml.Decoder
}

func NewNormalizingDecoder(r io.Reader) *NormalizingDecoder {
    return &NormalizingDecoder{inner: xml.NewDecoder(r)}
}

func (d *NormalizingDecoder) Token() (xml.Token, error) {
    tok, err := d.inner.Token()
    if err != nil {
        return tok, err
    }
    switch t := tok.(type) {
    case xml.StartElement:
        t.Name.Space = NormalizeNamespace(t.Name.Space)
        for i := range t.Attr {
            t.Attr[i].Name.Space = NormalizeNamespace(t.Attr[i].Name.Space)
        }
        return t, nil
    case xml.EndElement:
        t.Name.Space = NormalizeNamespace(t.Name.Space)
        return t, nil
    }
    return tok, nil
}
```

## 5.3 Использование

```go
// При открытии файла:
func parseDocumentXML(data []byte) (*CT_Document, error) {
    dec := xmltypes.NewNormalizingDecoder(bytes.NewReader(data))
    // Теперь все namespace'ы нормализованы к Transitional
    // Обычный unmarshal работает только с Transitional-именами
}
```

---

# 6. Паттерн: Сохранение namespace declarations

## 6.1 Проблема

`<w:document>` содержит 15-20 xmlns деклараций. `encoding/xml` в Go при
marshal добавляет namespace заново и может потерять неиспользуемые.

## 6.2 Решение

```go
type CT_Document struct {
    Body       *CT_Body
    Extra      []shared.RawXML
    // Сохранённые namespace declarations из исходного документа
    Namespaces []xml.Attr
}

func (doc *CT_Document) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    // Сохраняем ВСЕ атрибуты корневого элемента (включая xmlns:*)
    doc.Namespaces = make([]xml.Attr, len(start.Attr))
    copy(doc.Namespaces, start.Attr)

    // Далее парсим содержимое...
    for {
        tok, err := d.Token()
        // ...
    }
}

func (doc *CT_Document) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    start.Name = xml.Name{Local: "w:document"}

    // Восстанавливаем все namespace declarations
    if len(doc.Namespaces) > 0 {
        start.Attr = doc.Namespaces
    } else {
        // Для нового документа: минимальный набор
        start.Attr = defaultDocumentNamespaces()
    }

    e.EncodeToken(start)
    // ... тело ...
    return e.EncodeToken(start.End())
}

func defaultDocumentNamespaces() []xml.Attr {
    return []xml.Attr{
        {Name: xml.Name{Local: "xmlns:w"}, Value: xmltypes.NSw},
        {Name: xml.Name{Local: "xmlns:r"}, Value: xmltypes.NSr},
        {Name: xml.Name{Local: "xmlns:mc"}, Value: xmltypes.NSmc},
        {Name: xml.Name{Local: "xmlns:wp"}, Value: xmltypes.NSwp},
        {Name: xml.Name{Local: "xmlns:w14"}, Value: xmltypes.NSw14},
        {Name: xml.Name{Local: "xmlns:w15"}, Value: xmltypes.NSw15},
        // ... полный набор см. reference-appendix раздел 2.5
    }
}
```

---

# 7. Эталонная реализация: модуль `units`

Полный готовый код модуля `units` — используйте как образец стиля.

```go
// файл: units/units.go
package units

import (
    "fmt"
    "math"
    "strconv"
    "strings"
)

// EMU — English Metric Unit. 1 inch = 914400 EMU.
type EMU int64

// DXA — единица измерения в twips. 1 inch = 1440 DXA. 1 pt = 20 DXA.
type DXA int

// HalfPoint — половина пункта. Используется для размера шрифта. 12pt = 24 half-points.
type HalfPoint int

// EighthPoint — 1/8 пункта. Используется для толщины границ.
type EighthPoint int

// ========================================
// Конверсии в DXA
// ========================================

func InchToDXA(in float64) DXA  { return DXA(math.Round(in * 1440)) }
func CmToDXA(cm float64) DXA    { return DXA(math.Round(cm * 567.0)) }
func MmToDXA(mm float64) DXA    { return DXA(math.Round(mm * 56.7)) }
func PtToDXA(pt float64) DXA    { return DXA(math.Round(pt * 20)) }

// ========================================
// Конверсии в EMU
// ========================================

func InchToEMU(in float64) EMU  { return EMU(math.Round(in * 914400)) }
func CmToEMU(cm float64) EMU    { return EMU(math.Round(cm * 360000)) }
func MmToEMU(mm float64) EMU    { return EMU(math.Round(mm * 36000)) }
func PtToEMU(pt float64) EMU    { return EMU(math.Round(pt * 12700)) }
func DXAToEMU(d DXA) EMU        { return EMU(d) * 635 }

// ========================================
// Обратные конверсии
// ========================================

func EMUToDXA(e EMU) DXA        { return DXA(e / 635) }

func (d DXA) Inches() float64   { return float64(d) / 1440.0 }
func (d DXA) Cm() float64       { return float64(d) / 567.0 }
func (d DXA) Mm() float64       { return float64(d) / 56.7 }
func (d DXA) Pt() float64       { return float64(d) / 20.0 }

func (e EMU) Inches() float64   { return float64(e) / 914400.0 }
func (e EMU) Cm() float64       { return float64(e) / 360000.0 }

// ========================================
// Размер шрифта
// ========================================

func PtToHalfPoint(pt float64) HalfPoint       { return HalfPoint(math.Round(pt * 2)) }
func HalfPointToPt(hp HalfPoint) float64        { return float64(hp) / 2.0 }
func PtToEighthPoint(pt float64) EighthPoint    { return EighthPoint(math.Round(pt * 8)) }

// ========================================
// Парсинг
// ========================================

// ParseUniversalMeasure парсит строку вида "2.54cm", "1in", "72pt", "20mm".
func ParseUniversalMeasure(s string) (DXA, error) {
    s = strings.TrimSpace(s)
    if len(s) < 3 {
        return 0, fmt.Errorf("units: invalid measure %q", s)
    }
    unit := s[len(s)-2:]
    numStr := s[:len(s)-2]
    num, err := strconv.ParseFloat(numStr, 64)
    if err != nil {
        return 0, fmt.Errorf("units: invalid number in %q: %w", s, err)
    }
    switch unit {
    case "in":
        return InchToDXA(num), nil
    case "cm":
        return CmToDXA(num), nil
    case "mm":
        return MmToDXA(num), nil
    case "pt":
        return PtToDXA(num), nil
    default:
        return 0, fmt.Errorf("units: unknown unit %q in %q", unit, s)
    }
}

// ========================================
// Константы
// ========================================

const (
    LetterW DXA = 12240 // 8.5 inch
    LetterH DXA = 15840 // 11 inch
    A4W     DXA = 11906 // 210 mm
    A4H     DXA = 16838 // 297 mm
    LegalW  DXA = 12240 // 8.5 inch
    LegalH  DXA = 20163 // 14 inch

    DefaultMargin DXA = 1440 // 1 inch
    DefaultHeader DXA = 720  // 0.5 inch
    DefaultFooter DXA = 720  // 0.5 inch
    DefaultGutter DXA = 0

    DefaultTabStop DXA = 720 // 0.5 inch
)
```

```go
// файл: units/units_test.go
package units

import (
    "math"
    "testing"
)

func TestInchToDXA(t *testing.T) {
    tests := []struct {
        in   float64
        want DXA
    }{
        {1.0, 1440},
        {8.5, 12240},
        {0.5, 720},
        {0, 0},
    }
    for _, tt := range tests {
        got := InchToDXA(tt.in)
        if got != tt.want {
            t.Errorf("InchToDXA(%v) = %v, want %v", tt.in, got, tt.want)
        }
    }
}

func TestCmToDXA(t *testing.T) {
    got := CmToDXA(2.54)
    if got != 1440 {
        t.Errorf("CmToDXA(2.54) = %v, want 1440", got)
    }
}

func TestDXARoundTrip(t *testing.T) {
    original := DXA(1440)
    inches := original.Inches()
    back := InchToDXA(inches)
    if back != original {
        t.Errorf("round trip failed: %v → %v → %v", original, inches, back)
    }
}

func TestEMUConversions(t *testing.T) {
    emu := InchToEMU(2.0)
    if emu != 1828800 {
        t.Errorf("InchToEMU(2.0) = %v, want 1828800", emu)
    }
    dxa := EMUToDXA(emu)
    if dxa != 2880 {
        t.Errorf("EMUToDXA(%v) = %v, want 2880", emu, dxa)
    }
}

func TestParseUniversalMeasure(t *testing.T) {
    tests := []struct {
        input string
        want  DXA
    }{
        {"2.54cm", 1440},
        {"1in", 1440},
        {"72pt", 1440},
        {"25.4mm", 1440},
    }
    for _, tt := range tests {
        got, err := ParseUniversalMeasure(tt.input)
        if err != nil {
            t.Errorf("ParseUniversalMeasure(%q) error: %v", tt.input, err)
            continue
        }
        if math.Abs(float64(got-tt.want)) > 1 {
            t.Errorf("ParseUniversalMeasure(%q) = %v, want %v", tt.input, got, tt.want)
        }
    }
}

func TestConstants(t *testing.T) {
    // Letter = 8.5" × 11"
    if LetterW != 12240 { t.Errorf("LetterW = %v", LetterW) }
    if LetterH != 15840 { t.Errorf("LetterH = %v", LetterH) }
    // A4 = 210mm × 297mm
    if A4W != 11906 { t.Errorf("A4W = %v", A4W) }
    if A4H != 16838 { t.Errorf("A4H = %v", A4H) }
}
```

---

# 8. Недостающие типы (дополнение к contracts.md)

## 8.1 CT_Settings — стратегия partial parsing

CT_Settings содержит 253 элемента. Типизируем ~30 часто используемых,
остальное храним как RawXML.

```go
package settings

type CT_Settings struct {
    // Типизированные (часто используемые)
    WriteProtection       *CT_WriteProtection
    Zoom                  *CT_Zoom
    ProofState            *CT_Proof
    DefaultTabStop        *xmltypes.CT_TwipsMeasure
    CharacterSpacingControl *xmltypes.CT_String
    Compat                *CT_Compat
    Rsids                 *CT_DocRsids
    MathPr                *shared.RawXML  // сложная структура Math — raw для MVP
    ThemeFontLang         *CT_ThemeFontLang
    ClrSchemeMapping      *CT_ClrSchemeMapping  // ОБЯЗАТЕЛЬНЫЙ
    ShapeDefaults         *shared.RawXML        // ОБЯЗАТЕЛЬНЫЙ (VML — raw)
    DecimalSymbol         *xmltypes.CT_String
    ListSeparator         *xmltypes.CT_String
    DocId14               *xmltypes.CT_LongHexNumber  // w14:docId
    DocId15               *xmltypes.CT_Guid            // w15:docId
    TrackRevisions        *xmltypes.CT_OnOff
    DoNotTrackMoves       *xmltypes.CT_OnOff
    DoNotTrackFormatting  *xmltypes.CT_OnOff
    DocumentProtection    *CT_DocProtect
    EvenAndOddHeaders     *xmltypes.CT_OnOff
    MirrorMargins         *xmltypes.CT_OnOff

    // ВСЕ остальные 230+ элементов
    Extra []shared.RawXML

    // Порядок элементов для round-trip
    // (при unmarshal запоминаем порядок, при marshal восстанавливаем)
    elementOrder []string // ["zoom", "proofState", <raw:3>, "defaultTabStop", ...]
}

type CT_Zoom struct {
    Percent int    `xml:"percent,attr"`
    Val     string `xml:"val,attr,omitempty"` // "none"|"fullPage"|"bestFit"|...
}

type CT_Proof struct {
    Spelling *string `xml:"spelling,attr,omitempty"` // "clean"|"dirty"
    Grammar  *string `xml:"grammar,attr,omitempty"`
}

type CT_Compat struct {
    CompatSetting []CT_CompatSetting
    Extra         []shared.RawXML
}

type CT_CompatSetting struct {
    Name string `xml:"name,attr"`
    URI  string `xml:"uri,attr"`
    Val  string `xml:"val,attr"`
}

type CT_DocRsids struct {
    RsidRoot *xmltypes.CT_LongHexNumber
    Rsid     []xmltypes.CT_LongHexNumber
}

type CT_ThemeFontLang struct {
    Val      *string `xml:"val,attr,omitempty"`
    EastAsia *string `xml:"eastAsia,attr,omitempty"`
    Bidi     *string `xml:"bidi,attr,omitempty"`
}

type CT_ClrSchemeMapping struct {
    Bg1               string `xml:"bg1,attr"`
    T1                string `xml:"t1,attr"`
    Bg2               string `xml:"bg2,attr"`
    T2                string `xml:"t2,attr"`
    Accent1           string `xml:"accent1,attr"`
    Accent2           string `xml:"accent2,attr"`
    Accent3           string `xml:"accent3,attr"`
    Accent4           string `xml:"accent4,attr"`
    Accent5           string `xml:"accent5,attr"`
    Accent6           string `xml:"accent6,attr"`
    Hyperlink         string `xml:"hyperlink,attr"`
    FollowedHyperlink string `xml:"followedHyperlink,attr"`
}

type CT_WriteProtection struct { /* recommended, algorithm, hash, salt, ... */ }
type CT_DocProtect struct { /* edit, enforcement, algorithm, hash, salt, ... */ }
```

## 8.2 CT_FontsList

```go
package fonts

type CT_FontsList struct {
    Font []CT_Font
}

type CT_Font struct {
    Name    string `xml:"name,attr"`
    Panose1 *CT_Panose
    Charset *CT_Charset
    Family  *CT_FontFamily
    Pitch   *CT_Pitch
    Sig     *CT_FontSig
    // Для embedded fonts
    EmbedRegular    *CT_FontRel
    EmbedBold       *CT_FontRel
    EmbedItalic     *CT_FontRel
    EmbedBoldItalic *CT_FontRel
    Extra           []shared.RawXML
}

type CT_Panose struct { Val string `xml:"val,attr"` }
type CT_Charset struct { Val string `xml:"val,attr"` }
type CT_FontFamily struct { Val string `xml:"val,attr"` } // "roman"|"swiss"|"modern"|...
type CT_Pitch struct { Val string `xml:"val,attr"` }       // "fixed"|"variable"|"default"
type CT_FontSig struct {
    Usb0 string `xml:"usb0,attr"`
    Usb1 string `xml:"usb1,attr"`
    Usb2 string `xml:"usb2,attr"`
    Usb3 string `xml:"usb3,attr"`
    Csb0 string `xml:"csb0,attr"`
    Csb1 string `xml:"csb1,attr"`
}
type CT_FontRel struct {
    FontKey  *string `xml:"fontKey,attr,omitempty"`
    SubsetInfo *string `xml:"subsetted,attr,omitempty"`
    ID string `xml:"id,attr"` // r:id
}
```

## 8.3 CT_DocDefaults

```go
// parts/styles
type CT_DocDefaults struct {
    RPrDefault *CT_RPrDefault
    PPrDefault *CT_PPrDefault
}

type CT_RPrDefault struct {
    RPr *rpr.CT_RPr
}

type CT_PPrDefault struct {
    PPr *ppr.CT_PPrBase
}
```

## 8.4 CT_LatentStyles

```go
type CT_LatentStyles struct {
    DefLockedState    *bool `xml:"defLockedState,attr,omitempty"`
    DefUIPriority     *int  `xml:"defUIPriority,attr,omitempty"`
    DefSemiHidden     *bool `xml:"defSemiHidden,attr,omitempty"`
    DefUnhideWhenUsed *bool `xml:"defUnhideWhenUsed,attr,omitempty"`
    DefQFormat        *bool `xml:"defQFormat,attr,omitempty"`
    Count             *int  `xml:"count,attr,omitempty"`
    LsdException      []CT_LsdException
}

type CT_LsdException struct {
    Name           string `xml:"name,attr"`
    Locked         *bool  `xml:"locked,attr,omitempty"`
    UIPriority     *int   `xml:"uiPriority,attr,omitempty"`
    SemiHidden     *bool  `xml:"semiHidden,attr,omitempty"`
    UnhideWhenUsed *bool  `xml:"unhideWhenUsed,attr,omitempty"`
    QFormat        *bool  `xml:"qFormat,attr,omitempty"`
}
```

## 8.5 CT_TblStylePr

```go
type CT_TblStylePr struct {
    Type string          `xml:"type,attr"` // "firstRow"|"lastRow"|"band1Vert"|...
    PPr  *ppr.CT_PPrBase
    RPr  *rpr.CT_RPrBase
    TblPr *CT_TblPrBase  // подмножество table props
    TrPr  *CT_TrPrBase
    TcPr  *CT_TcPrBase
}
```

## 8.6 CT_NumLvl (level override)

```go
type CT_NumLvl struct {
    Ilvl       int `xml:"ilvl,attr"`
    StartOverride *xmltypes.CT_DecimalNumber
    Lvl        *CT_Lvl // полное переопределение уровня
}
```

## 8.7 parts/theme и parts/websettings

```go
// parts/theme — хранить как raw XML (DrawingML слишком сложен для MVP)
package theme

func Parse(data []byte) ([]byte, error) { return data, nil }  // passthrough
func Serialize(data []byte) ([]byte, error) { return data, nil }

// parts/websettings — аналогично
package websettings

func Parse(data []byte) ([]byte, error) { return data, nil }
func Serialize(data []byte) ([]byte, error) { return data, nil }
```

---

# 9. Паттерн: CT_OnOff полная реализация

```go
package xmltypes

import "encoding/xml"

type CT_OnOff struct {
    Val *string // nil → element present means true
}

// NewOnOff создаёт CT_OnOff.
// true  → element без атрибута val (самая компактная форма)
// false → element с val="0"
func NewOnOff(v bool) *CT_OnOff {
    if v {
        return &CT_OnOff{Val: nil}
    }
    s := "0"
    return &CT_OnOff{Val: &s}
}

// Bool возвращает логическое значение.
// Если CT_OnOff == nil → возвращает defaultVal (наследование от стиля).
func (o *CT_OnOff) Bool(defaultVal bool) bool {
    if o == nil {
        return defaultVal
    }
    if o.Val == nil {
        return true // <w:b/> без val = true
    }
    switch *o.Val {
    case "true", "1", "on":
        return true
    case "false", "0", "off":
        return false
    default:
        return true // неизвестное значение → true (Word behavior)
    }
}

func (o *CT_OnOff) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if o.Val != nil {
        start.Attr = append(start.Attr, xml.Attr{
            Name:  xml.Name{Space: NSw, Local: "val"},
            Value: *o.Val,
        })
    }
    // Самозакрывающийся элемент: <w:b/> или <w:b w:val="false"/>
    e.EncodeToken(start)
    return e.EncodeToken(start.End())
}

func (o *CT_OnOff) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for _, attr := range start.Attr {
        if attr.Name.Local == "val" {
            s := attr.Value
            o.Val = &s
        }
    }
    // Пропустить содержимое (обычно пустое)
    d.Skip()
    return nil
}
```

---

# 10. Паттерн: CT_Text с xml:space

```go
package run

import "encoding/xml"

type CT_Text struct {
    Space *string `xml:"space,attr,omitempty"` // xml:space
    Value string  `xml:",chardata"`
}

func NewText(s string) CT_Text {
    t := CT_Text{Value: s}
    // Авто-добавление xml:space="preserve" если есть пробелы по краям
    if len(s) > 0 && (s[0] == ' ' || s[len(s)-1] == ' ' || s[0] == '\t') {
        preserve := "preserve"
        t.Space = &preserve
    }
    return t
}

func (t *CT_Text) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if t.Space != nil {
        start.Attr = append(start.Attr, xml.Attr{
            Name:  xml.Name{Space: "http://www.w3.org/XML/1998/namespace", Local: "space"},
            Value: *t.Space,
        })
    }
    e.EncodeToken(start)
    e.EncodeToken(xml.CharData(t.Value))
    return e.EncodeToken(start.End())
}
```

---

# 11. Обновлённое дерево пакетов

```
docx-go/
├── go.mod
├── units/              C-00  (0 deps)
├── opc/                C-01  (0 deps)
├── xmltypes/           C-02  (0 deps)
├── coreprops/          C-03  (0 deps)
├── wml/
│   ├── shared/         C-NEW (← xmltypes)      ← НОВЫЙ ПАКЕТ
│   ├── rpr/            C-10  (← xmltypes, shared)
│   ├── ppr/            C-11  (← xmltypes, rpr, shared)  ← ppr ИМПОРТИРУЕТ rpr
│   ├── sectpr/         C-12  (← xmltypes, shared)
│   ├── table/          C-13  (← xmltypes, shared)
│   ├── tracking/       C-14  (← xmltypes, shared)
│   ├── run/            C-15  (← xmltypes, rpr, shared, dml)
│   ├── para/           C-16  (← xmltypes, rpr, ppr, run, tracking, shared)
│   ├── body/           C-17  (← xmltypes, para, table, sectpr, shared)
│   └── hdft/           C-18  (← shared)
├── dml/                C-19  (← xmltypes, shared)
├── parts/
│   ├── document/       C-20  (← body, opc)
│   ├── styles/         C-21  (← xmltypes, rpr, ppr, table, shared)
│   ├── numbering/      C-22  (← xmltypes, rpr, ppr, shared)
│   ├── settings/       C-23  (← xmltypes, shared)
│   ├── fonts/          C-24  (← xmltypes, shared)
│   ├── comments/       C-25  (← shared)
│   ├── footnotes/      C-26  (← shared)
│   ├── headers/        C-27  (← shared)
│   ├── theme/          C-28  (0 deps — raw passthrough)
│   └── websettings/    C-29  (0 deps — raw passthrough)
├── packaging/          C-30  (← opc, coreprops, parts/*, units)
├── validator/          C-31  (← packaging)
└── docx/               C-32  (← packaging, validator, units, wml/*, coreprops)
```

---

# 12. Тестовый паттерн: round-trip

```go
func TestPPrBaseRoundTrip(t *testing.T) {
    input := `<w:pPr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml">` +
        `<w:pStyle w:val="Heading1"/>` +
        `<w:keepNext/>` +
        `<w:keepLines/>` +
        `<w:spacing w:before="240" w:after="0"/>` +
        `<w:jc w:val="center"/>` +
        `<w:outlineLvl w:val="0"/>` +
        `<w14:someExtension w14:val="test"/>` + // неизвестный элемент
        `</w:pPr>`

    // Unmarshal
    var ppr CT_PPrBase
    err := xml.Unmarshal([]byte(input), &ppr)
    if err != nil {
        t.Fatal(err)
    }

    // Проверки
    if ppr.PStyle == nil || ppr.PStyle.Val != "Heading1" {
        t.Error("PStyle not parsed")
    }
    if !ppr.KeepNext.Bool(false) {
        t.Error("KeepNext not true")
    }
    if ppr.Jc == nil || ppr.Jc.Val != "center" {
        t.Error("Jc not parsed")
    }
    if len(ppr.Extra) != 1 {
        t.Errorf("Expected 1 Extra element, got %d", len(ppr.Extra))
    }

    // Marshal
    output, err := xml.Marshal(&ppr)
    if err != nil {
        t.Fatal(err)
    }

    // Re-unmarshal и сравнить
    var ppr2 CT_PPrBase
    err = xml.Unmarshal(output, &ppr2)
    if err != nil {
        t.Fatal(err)
    }

    // Сравнить ключевые поля
    if ppr2.PStyle.Val != ppr.PStyle.Val {
        t.Error("round-trip lost PStyle")
    }
    if len(ppr2.Extra) != len(ppr.Extra) {
        t.Error("round-trip lost Extra elements")
    }
    if ppr2.Extra[0].XMLName.Local != "someExtension" {
        t.Error("round-trip lost extension element name")
    }
}
```

---

# 13. Чеклист: перед реализацией каждого модуля

```
□ Прочитать контракт модуля (contracts.md, секция C-XX)
□ Прочитать контракты зависимостей (только публичные типы)
□ Прочитать маппинг полей (patterns.md, раздел 2)
□ Проверить: нужен ли кастомный MarshalXML? (xsd:sequence → ДА)
□ Реализовать Extra []shared.RawXML для round-trip
□ XML namespace: использовать xmltypes.NSw для всех элементов w:
□ CT_OnOff: использовать xmltypes.CT_OnOff, не *bool
□ CT_Text: автоматический xml:space="preserve"
□ Написать round-trip тест (unmarshal → marshal → unmarshal → сравнить)
□ go vet && go test ./...
```

---

# 14. Паттерн: X() escape hatch

## 14.1 Проблема

Высокоуровневый API (пакет `docx`, C-32) предоставляет convenience-методы
`Add*`/`Remove*`/`Set*` для типичных сценариев. Однако покрыть **все**
возможные операции с OOXML невозможно — спецификация слишком обширна.
Пользователю нужен способ «выйти за рамки» обёртки к нижнему уровню.

## 14.2 Решение

Каждый wrapper-тип в пакете `docx` имеет метод `X()`, который возвращает
доступ к нижележащему WML-типу:

```go
// Pointer-based wrappers (хранят указатель → X() возвращает указатель):
func (d *Document) X() *packaging.Document
func (b *Body) X() *body.CT_Body
func (p *Paragraph) X() *para.CT_P       // Body.Content хранит *para.CT_P
func (r *Run) X() *run.CT_R              // para.RunItem хранит *run.CT_R
func (t *Table) X() *table.CT_Tbl        // Body.Content хранит *table.CT_Tbl
func (s *Section) X() *sectpr.CT_SectPr
func (h *Header) X() *hdft.CT_HdrFtr
func (f *Footer) X() *hdft.CT_HdrFtr

// Value-based wrappers (CT_Row/CT_Tc — value в interface slice → X() возвращает КОПИЮ):
func (r *Row) X() table.CT_Row           // value, не pointer!
func (c *Cell) X() table.CT_Tc           // value, не pointer!
```

**Почему Row и Cell возвращают value, а не pointer:**

`CT_Row` хранится как value-type в `tbl.Content []TblContent` (interface slice).
Go не предоставляет способа получить указатель на элемент внутри interface slice.
При извлечении `tbl.Content[i].(table.CT_Row)` мы получаем **копию**.
Мутации копии не отражаются на оригинале. Подробности: раздел 15.

## 14.3 Конвенция

| Правило | Описание |
|---------|----------|
| Имя | Всегда `X()` — одна буква, легко запомнить |
| Возвращаемый тип | Указатель на WML-структуру для pointer-based; value-copy для Row/Cell (см. раздел 15) |
| Семантика | Прямой доступ без копирования — мутации через `X()` отражаются в документе |
| Гарантии | Никаких — пользователь берёт ответственность за инварианты (≥1 `<w:p>` в ячейке и пр.) |
| Документирование | Каждый `X()` документирован комментарием с указанием типа и ограничений |

## 14.4 Пример использования

```go
doc := docx.New()
p := doc.Body().AddParagraph()

// Через высокоуровневый API:
p.SetStyle("Heading1")

// Через escape hatch — прямой доступ к PPr:
raw := p.X()
raw.PPr.KeepNext = xmltypes.NewOnOff(true)

// Прямая манипуляция Body.Content:
bd := doc.Body().X()
bd.Content = append(bd.Content[:2], bd.Content[3:]...) // удалить элемент [2]
```

## 14.5 Реализация в wrapper-структуре

Каждый wrapper хранит **указатель** на нижележащий тип. `X()` — это
просто геттер:

```go
// Body — пример реализации wrapper.
type Body struct {
    raw *body.CT_Body
}

func (b *Body) X() *body.CT_Body {
    return b.raw
}
```

Wrapper **не копирует** данные — он всего лишь хранит указатель.
Следствие: два Wrapper, полученные из одного Document, указывают на
одни и те же данные (share state).

---

# 15. Паттерн: Value-type в interface slice (CT_Row, CT_Tc)

## 15.1 Проблема

В пакете `wml/table` строки и ячейки хранятся как **значения** в
interface-слайсах:

```go
type CT_Tbl struct {
    Content []TblContent    // TblContent — interface
}

type CT_Row struct {        // CT_Row — value type, не *CT_Row
    Content []RowContent    // RowContent — interface
}

// Unmarshal добавляет value (не pointer):
tbl.Content = append(tbl.Content, row)  // row — CT_Row, не *CT_Row
```

Это означает, что при извлечении строки из `tbl.Content` мы получаем
**копию**. Мутация копии НЕ отражается на оригинале:

```go
row := tbl.Content[0].(table.CT_Row)  // копия!
row.TrPr = &table.CT_TrPr{...}       // изменяем копию
// tbl.Content[0] — НЕ изменилось!
```

## 15.2 Решение в обёртке

### Вариант A: Write-back после мутации

После изменения копии — записать её обратно в слайс:

```go
row := tbl.Content[i].(table.CT_Row)
// ... мутации ...
tbl.Content[i] = row  // записать обратно
```

Это паттерн, уже используемый в `autofix.go` (validator).

### Вариант B: Wrapper хранит индекс + указатель на таблицу

```go
type Row struct {
    tbl *table.CT_Tbl  // родительская таблица
    idx int            // позиция в tbl.Content
}

// При каждом чтении — извлекаем актуальное значение:
func (r *Row) cells() []table.RowContent {
    row := r.tbl.Content[r.idx].(table.CT_Row)
    return row.Content
}

// При мутации — извлекаем, меняем, записываем обратно:
func (r *Row) AddCell() *Cell {
    row := r.tbl.Content[r.idx].(table.CT_Row)
    tc := table.CT_Tc{Content: []shared.BlockLevelElement{&para.CT_P{}}}
    row.Content = append(row.Content, tc)
    r.tbl.Content[r.idx] = row  // write-back!
    // ...
}
```

### Вариант C (рекомендуемый): X() возвращает копию, рядом — XMut()

```go
// X() для чтения — возвращает копию (безопасно).
func (r *Row) X() table.CT_Row {
    return r.tbl.Content[r.idx].(table.CT_Row)
}
```

Для сложных мутаций — пользователь использует `Table.X().Content` напрямую.

## 15.3 Аналогичная ситуация с CT_Tc

`CT_Tc` хранится как value в `CT_Row.Content`. Wrapper `Cell` должен
применять те же паттерны write-back.

```go
type Cell struct {
    tbl    *table.CT_Tbl
    rowIdx int
    colIdx int
}
```

## 15.4 Ключевые инварианты

| Инвариант | Где проверять |
|-----------|--------------|
| Индекс Row/Cell валиден | Перед каждой операцией чтения/записи |
| CT_Tc.Content ≥ 1 элемент | После каждой мутации (Cell.Clear вставляет `<w:p/>`) |
| Write-back после мутации | ВСЕГДА при работе с value-type в interface slice |

---

# 16. Паттерн: FindText / ReplaceText

## 16.1 Область применения

`Body.FindText()` и `Body.ReplaceText()` работают на уровне отдельных
Run-ов (`<w:r>/<w:t>`). Это осознанное ограничение для v1.

## 16.2 Ограничения

| Ограничение | Причина |
|-------------|---------|
| Текст, разбитый между ранами, **не находится** | Word может разбить одно слово на несколько `<w:r>` при редактировании, вставке rsid, проверке орфографии |
| **Не ищет внутри таблиц** | CT_Tc хранит value type, TextLocation не может однозначно адресовать ячейку; v2+ |
| Не ищет в headers/footers/footnotes/comments | Они хранятся в отдельных parts |
| Не ищет внутри SDT (content controls) | SDT вложены как RawXML |
| Case-sensitive | Для v1 достаточно; case-insensitive можно добавить позже |

## 16.3 Алгоритм FindText

```
Для каждого BlockLevelElement в Body.Content:
    Если это *para.CT_P (type assert):
        Для каждого ParagraphContent в CT_P.Content:
            Если это para.RunItem:
                Собрать весь текст из CT_R.Content (CT_Text элементы)
                Если strings.Contains(fullText, needle):
                    Добавить TextLocation{BlockIndex, RunIndex, Paragraph, Run}
    Если это *table.CT_Tbl:
        ПРОПУСТИТЬ (v1 — не ищем в таблицах)
    Иначе (shared.RawXML и пр.):
        ПРОПУСТИТЬ
```

## 16.4 Алгоритм ReplaceText

```
Аналогично FindText, но при нахождении:
    Для каждого CT_Text в CT_R.Content:
        text.Value = strings.ReplaceAll(text.Value, old, new)
    Счётчик += количество замен
```

## 16.5 Пример

```go
doc, _ := docx.Open("report.docx")

// Найти все вхождения:
locs := doc.Body().FindText("PLACEHOLDER")

// Заменить:
count := doc.Body().ReplaceText("PLACEHOLDER", "Actual Value")
fmt.Printf("Replaced %d occurrences\n", count)

doc.Save("report_filled.docx")
```

## 16.6 Будущие расширения (v2+)

- **Search in tables**: рекурсивный обход ячеек. Потребует расширить
  `TextLocation` полями `Table`, `Row`, `Col` для адресации.
- **Cross-run search**: объединять текст соседних ранов для поиска через
  границы `<w:r>`. Сложность: при замене нужно корректно перераспределять
  текст между ранами, сохраняя форматирование.
- **Regex search**: `FindRegex(pattern string)`.
- **Search in headers/footers**: через `Header.FindText()` / `Footer.FindText()`.

---

# 17. Паттерн: Remove / Insert — поддержание инвариантов

## 17.1 Инварианты при удалении элементов

При удалении block-level элементов из Body, строк из Table, ранов из
Paragraph — необходимо учитывать инварианты OOXML:

| Операция | Инварианты |
|----------|-----------|
| `Body.RemoveElement(i)` | Если удалённый параграф содержал section break (PPr.SectPr != nil), секция теряется — это допустимо, но пользователь должен быть осведомлён |
| `Body.Clear()` | Не удаляет Body.SectPr (секция body-level сохраняется) |
| `Table.RemoveRow(i)` | Если таблица становится пустой (0 строк) — допустимо (Word обработает), но validator выдаст предупреждение |
| `Paragraph.RemoveRun(i)` | Параграф без ранов допустим (`<w:p/>` = пустая строка) |
| `Paragraph.Clear()` | Сохраняет PPr (стиль, нумерация) — удаляет только Content |
| `Cell.Clear()` | ОБЯЗАТЕЛЬНО вставляет `&para.CT_P{}` после очистки (инвариант ≥1 `<w:p>`) |
| `Header.Clear()` / `Footer.Clear()` | Аналогично Cell — вставляет `&para.CT_P{}` |

## 17.2 Orphan relationships

При удалении элементов могут остаться «осиротевшие» relationships:

- Удалён параграф с гиперссылкой → rId в document.xml.rels ссылается на
  несуществующий контент
- Удалён ран с картинкой → relationship на image part висит в воздухе

**Решение для v1**: orphan relationships допустимы. Word их игнорирует,
они не вызывают ошибок. Очистка — будущая оптимизация.

## 17.3 Index-based API

Все `Remove*`/`InsertAt*` используют integer index в соответствующем
Content-слайсе. Это осознанный выбор:

- Просто, предсказуемо, zero-alloc
- Соответствует Go-идиоме работы со слайсами
- Пользователь получает индексы через `Paragraphs()`, `Rows()`, `Runs()`

**Правило**: после `Remove` или `InsertAt` все ранее полученные индексы
и wrapper-объекты для элементов **после** точки вставки/удаления —
**инвалидированы**. Это задокументировано в godoc каждого метода.