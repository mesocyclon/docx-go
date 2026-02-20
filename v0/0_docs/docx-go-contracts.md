# DOCX-Go: Контракты модулей

## Как пользоваться этим файлом

Этот файл содержит **контракты** (публичные типы, сигнатуры, зависимости)
всех модулей библиотеки. Загружается вместе с двумя другими файлами.

**Содержимое:**
1. **Публичные типы** каждого модуля (поля структур, но не методы реализации)
2. **Публичные функции** (сигнатуры, но не тела)
3. **Интерфейсы** для связи между модулями (определены в `wml/shared`)
4. **Import paths** для каждого пакета

Реализация модуля = написать код, который **соответствует контракту своего модуля**
и **использует контракты зависимостей**.

### Правила для LLM

```
ПЕРЕД реализацией модуля X:
1. Загрузи contracts.md          — ЧТО реализовать (этот файл)
2. Загрузи reference-appendix.md — XML-примеры, namespaces, порядок элементов
3. Загрузи patterns.md           — КАК реализовать (MarshalXML, RawXML, naming, shared)
4. Найди секцию модуля X в contracts — это твой ВЫХОДНОЙ контракт
5. Найди секции зависимостей X — это твои ВХОДНЫЕ контракты
6. Найди маппинг полей в patterns.md раздел 2 — Go field → XML element name
7. Пиши реализацию, импортируя ТОЛЬКО перечисленные зависимости

ПОСЛЕ реализации модуля X:
1. Убедись что все публичные типы/функции из контракта реализованы
2. Запусти go build — проверь что компилируется
3. Запусти round-trip тест (паттерн в patterns.md раздел 12)
4. НЕ МЕНЯЙ контракт — если нужно изменение, это отдельная задача
```

---

## Go module path

```
module github.com/vortex/docx-go
```

## Дерево пакетов

```
docx-go/
├── units/              ← github.com/vortex/docx-go/units
├── opc/                ← github.com/vortex/docx-go/opc
├── xmltypes/           ← github.com/vortex/docx-go/xmltypes
├── coreprops/          ← github.com/vortex/docx-go/coreprops
├── wml/
│   ├── shared/         ← github.com/vortex/docx-go/wml/shared        ← ИНТЕРФЕЙСЫ
│   ├── rpr/            ← github.com/vortex/docx-go/wml/rpr
│   ├── ppr/            ← github.com/vortex/docx-go/wml/ppr
│   ├── sectpr/         ← github.com/vortex/docx-go/wml/sectpr
│   ├── table/          ← github.com/vortex/docx-go/wml/table
│   ├── tracking/       ← github.com/vortex/docx-go/wml/tracking
│   ├── run/            ← github.com/vortex/docx-go/wml/run
│   ├── para/           ← github.com/vortex/docx-go/wml/para
│   ├── body/           ← github.com/vortex/docx-go/wml/body
│   └── hdft/           ← github.com/vortex/docx-go/wml/hdft
├── dml/                ← github.com/vortex/docx-go/dml
├── parts/
│   ├── document/       ← github.com/vortex/docx-go/parts/document
│   ├── styles/         ← github.com/vortex/docx-go/parts/styles
│   ├── numbering/      ← github.com/vortex/docx-go/parts/numbering
│   ├── settings/       ← github.com/vortex/docx-go/parts/settings
│   ├── fonts/          ← github.com/vortex/docx-go/parts/fonts
│   ├── comments/       ← github.com/vortex/docx-go/parts/comments
│   ├── footnotes/      ← github.com/vortex/docx-go/parts/footnotes
│   ├── headers/        ← github.com/vortex/docx-go/parts/headers
│   ├── theme/          ← github.com/vortex/docx-go/parts/theme
│   └── websettings/    ← github.com/vortex/docx-go/parts/websettings
├── packaging/          ← github.com/vortex/docx-go/packaging
├── validator/          ← github.com/vortex/docx-go/validator
└── docx/               ← github.com/vortex/docx-go/docx              ← публичный API
```

---

## Граф зависимостей (import graph)

```
Кто что импортирует (→ = imports):

=== Уровень 0: Фундамент (нет зависимостей) ===
units         → (nothing)
opc           → (nothing)
xmltypes      → (nothing)
coreprops     → (nothing)

=== Уровень 0.5: Общие интерфейсы ===
wml/shared    → xmltypes                              ← КЛЮЧЕВОЙ ПАКЕТ: интерфейсы контента

=== Уровень 1: WML-типы ===
wml/rpr       → xmltypes, wml/shared                   ← rpr использует shared.RawXML
wml/ppr       → xmltypes, wml/rpr, wml/shared          ← ppr ИМПОРТИРУЕТ rpr (CT_PPr.RPr)
wml/sectpr    → xmltypes, wml/shared
wml/table     → xmltypes, wml/shared                   ← CT_Tc.Content = []shared.BlockLevelElement
wml/tracking  → xmltypes, wml/shared
dml           → xmltypes, wml/shared                    ← dml использует shared.RawXML

=== Уровень 2: WML-составные ===
wml/run       → xmltypes, wml/rpr, wml/shared, dml
wml/para      → xmltypes, wml/rpr, wml/ppr, wml/run, wml/tracking, wml/shared
wml/body      → xmltypes, wml/para, wml/table, wml/sectpr, wml/shared
wml/hdft      → wml/shared

=== Уровень 3: Parts ===
parts/document    → wml/body, opc
parts/styles      → xmltypes, wml/rpr, wml/ppr, wml/table, wml/shared
parts/numbering   → xmltypes, wml/rpr, wml/ppr, wml/shared
parts/settings    → xmltypes, wml/shared
parts/fonts       → xmltypes, wml/shared
parts/comments    → wml/shared
parts/footnotes   → wml/shared
parts/headers     → wml/shared
parts/theme       → (nothing — raw passthrough)
parts/websettings → (nothing — raw passthrough)

=== Уровень 4-5: Интеграция ===
packaging  → opc, coreprops, parts/*, units
validator  → packaging
docx       → packaging, validator, units
```

---

# КОНТРАКТЫ МОДУЛЕЙ

---

## C-00: `units`

**Импортирует**: ничего

```go
package units

// === Типы ===
type EMU int64
type DXA int
type HalfPoint int
type EighthPoint int

// === Конверсии в DXA ===
func InchToDXA(in float64) DXA
func CmToDXA(cm float64) DXA
func MmToDXA(mm float64) DXA
func PtToDXA(pt float64) DXA

// === Конверсии в EMU ===
func InchToEMU(in float64) EMU
func CmToEMU(cm float64) EMU
func MmToEMU(mm float64) EMU
func PtToEMU(pt float64) EMU
func DXAToEMU(d DXA) EMU

// === Обратные конверсии ===
func EMUToDXA(e EMU) DXA
func (d DXA) Inches() float64
func (d DXA) Cm() float64
func (d DXA) Mm() float64
func (d DXA) Pt() float64
func (e EMU) Inches() float64
func (e EMU) Cm() float64

// === Размер шрифта ===
func PtToHalfPoint(pt float64) HalfPoint
func HalfPointToPt(hp HalfPoint) float64
func PtToEighthPoint(pt float64) EighthPoint

// === Парсинг строк из XML ===
func ParseUniversalMeasure(s string) (DXA, error)  // "2.54cm" → 1440

// === Константы страниц (DXA) ===
const (
    LetterW DXA = 12240  // 8.5 inch
    LetterH DXA = 15840  // 11 inch
    A4W     DXA = 11906  // 210 mm
    A4H     DXA = 16838  // 297 mm
)
```

---

## C-01: `opc`

**Импортирует**: ничего (только stdlib)

```go
package opc

import "io"

// === Типы ===
type Package struct { /* private */ }

type Part struct {
    Name        string         // "/word/document.xml"
    ContentType string
    Data        []byte
    Rels        []Relationship // part-level relationships
}

type Relationship struct {
    ID         string // "rId1"
    Type       string // полный URI
    Target     string // "styles.xml" или URL
    TargetMode string // "" (Internal) | "External"
}

// === Открытие/создание ===
func OpenFile(path string) (*Package, error)
func OpenReader(r io.ReaderAt, size int64) (*Package, error)
func New() *Package

// === Сохранение ===
func (p *Package) SaveFile(path string) error
func (p *Package) SaveWriter(w io.Writer) error

// === Операции с частями ===
func (p *Package) Part(name string) (*Part, bool)
func (p *Package) AddPart(name, contentType string, data []byte) *Part
func (p *Package) RemovePart(name string) bool
func (p *Package) Parts() []*Part

// === Package-level relationships ===
func (p *Package) PackageRels() []Relationship
func (p *Package) AddPackageRel(relType, target string) string       // → rId
func (p *Package) PackageRelsByType(relType string) []Relationship

// === Part-level relationships ===
func (pt *Part) AddRel(relType, target string) string                // → rId
func (pt *Part) AddExternalRel(relType, target string) string        // → rId
func (pt *Part) RelsByType(relType string) []Relationship
func (pt *Part) RelByID(id string) (Relationship, bool)

// === Константы Relationship Type URI ===
const (
    RelOfficeDocument  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
    RelCoreProperties  = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties"
    RelExtProperties   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties"
    RelStyles          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
    RelSettings        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings"
    RelWebSettings     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/webSettings"
    RelFontTable       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/fontTable"
    RelNumbering       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering"
    RelFootnotes       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes"
    RelEndnotes        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/endnotes"
    RelComments        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments"
    RelHeader          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/header"
    RelFooter          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer"
    RelImage           = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
    RelHyperlink       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
    RelTheme           = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
)
```

---

## C-02: `xmltypes`

**Импортирует**: ничего (только `encoding/xml`)

```go
package xmltypes

import "encoding/xml"

// === Namespace URI (полный список) ===
const (
    NSw   = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
    NSr   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
    NSa   = "http://schemas.openxmlformats.org/drawingml/2006/main"
    NSwp  = "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
    NSpic = "http://schemas.openxmlformats.org/drawingml/2006/picture"
    NSmc  = "http://schemas.openxmlformats.org/markup-compatibility/2006"
    NSm   = "http://schemas.openxmlformats.org/officeDocument/2006/math"
    NSv   = "urn:schemas-microsoft-com:vml"
    NSo   = "urn:schemas-microsoft-com:office:office"
    NSw10 = "urn:schemas-microsoft-com:office:word"
    NSw14 = "http://schemas.microsoft.com/office/word/2010/wordml"
    NSw15 = "http://schemas.microsoft.com/office/word/2012/wordml"
    NSw16se = "http://schemas.microsoft.com/office/word/2015/wordml/symex"
    // OPC / package
    NSContentTypes  = "http://schemas.openxmlformats.org/package/2006/content-types"
    NSRelationships = "http://schemas.openxmlformats.org/package/2006/relationships"
    NScp      = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
    NSdc      = "http://purl.org/dc/elements/1.1/"
    NSdcterms = "http://purl.org/dc/terms/"
    NSdcmitype = "http://purl.org/dc/dcmitype/"
    NSxsi     = "http://www.w3.org/2001/XMLSchema-instance"
    // Extended properties
    NSvt = "http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"
    NSep = "http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"
)

// === RawXML — ПЕРЕНЕСЁН в wml/shared ===
// RawXML определён в wml/shared.RawXML (см. C-05).
// НЕ дублировать здесь — все модули используют shared.RawXML.

// === CT_OnOff — тройная логика WML ===
// nil = not set (inherit), Val==nil → true, Val="false" → false
// Подробная реализация с MarshalXML/UnmarshalXML: см. patterns.md раздел 9
type CT_OnOff struct {
    Val *string `xml:"val,attr,omitempty"` // nil → element present = true
}
func NewOnOff(v bool) *CT_OnOff             // true → &CT_OnOff{Val:nil}, false → &CT_OnOff{Val:&"0"}
func (o *CT_OnOff) Bool(defaultVal bool) bool // nil receiver → defaultVal; Val==nil → true; "0"/"false"/"off" → false
func (o *CT_OnOff) IsExplicitlySet() bool
// Реализует xml.Marshaler и xml.Unmarshaler

// === Простые обёрточные типы ===
type CT_String struct {
    Val string `xml:"val,attr"`
}
type CT_DecimalNumber struct {
    Val int `xml:"val,attr"`
}
type CT_UnsignedDecimalNumber struct {
    Val uint `xml:"val,attr"`
}
type CT_TwipsMeasure struct {
    Val int `xml:"val,attr"` // DXA (twips)
}
type CT_SignedTwipsMeasure struct {
    Val int `xml:"val,attr"` // signed DXA
}
type CT_HpsMeasure struct {
    Val int `xml:"val,attr"` // half-points
}
type CT_SignedHpsMeasure struct {
    Val int `xml:"val,attr"` // signed half-points
}
type CT_TextScale struct {
    Val int `xml:"val,attr"` // процент (100 = 100%)
}
type CT_LongHexNumber struct {
    Val string `xml:"val,attr"` // 8-hex e.g. "00A4B3C2"
}
type CT_ShortHexNumber struct {
    Val string `xml:"val,attr"` // 4-hex e.g. "A4B3"
}
type CT_Guid struct {
    Val string `xml:"val,attr"` // {XXXXXXXX-XXXX-...}
}
type CT_Lang struct {
    Val string `xml:"val,attr"` // "en-US", "ru-RU"
}
type CT_Empty struct{} // самозакрывающийся элемент

// === Цвет ===
type CT_Color struct {
    Val        string  `xml:"val,attr"`                   // hex "FF0000" или "auto"
    ThemeColor *string `xml:"themeColor,attr,omitempty"`   // "accent1", "dark1", ...
    ThemeTint  *string `xml:"themeTint,attr,omitempty"`    // hex "BF"
    ThemeShade *string `xml:"themeShade,attr,omitempty"`   // hex "80"
}

// === Подчёркивание ===
type CT_Underline struct {
    Val        *string `xml:"val,attr,omitempty"`          // "single", "double", ...
    Color      *string `xml:"color,attr,omitempty"`        // hex
    ThemeColor *string `xml:"themeColor,attr,omitempty"`
}

// === Highlight ===
type CT_Highlight struct {
    Val string `xml:"val,attr"` // "yellow", "green", ...
}

// === Граница (одна линия) ===
type CT_Border struct {
    Val        string  `xml:"val,attr"`                    // "single", "double", ...
    Sz         *int    `xml:"sz,attr,omitempty"`            // eighth-points
    Space      *int    `xml:"space,attr,omitempty"`         // points
    Color      *string `xml:"color,attr,omitempty"`         // hex
    ThemeColor *string `xml:"themeColor,attr,omitempty"`
    ThemeTint  *string `xml:"themeTint,attr,omitempty"`
    ThemeShade *string `xml:"themeShade,attr,omitempty"`
    Shadow     *bool   `xml:"shadow,attr,omitempty"`
    Frame      *bool   `xml:"frame,attr,omitempty"`
}

// === Заливка ===
type CT_Shd struct {
    Val       string  `xml:"val,attr"`                     // "clear", "solid", "pct10", ...
    Color     *string `xml:"color,attr,omitempty"`
    Fill      *string `xml:"fill,attr,omitempty"`
    ThemeColor     *string `xml:"themeColor,attr,omitempty"`
    ThemeFill      *string `xml:"themeFill,attr,omitempty"`
    ThemeFillTint  *string `xml:"themeFillTint,attr,omitempty"`
    ThemeFillShade *string `xml:"themeFillShade,attr,omitempty"`
}

// === Шрифты ===
type CT_Fonts struct {
    Ascii          *string `xml:"ascii,attr,omitempty"`
    HAnsi          *string `xml:"hAnsi,attr,omitempty"`
    EastAsia       *string `xml:"eastAsia,attr,omitempty"`
    CS             *string `xml:"cs,attr,omitempty"`
    AsciiTheme     *string `xml:"asciiTheme,attr,omitempty"`
    HAnsiTheme     *string `xml:"hAnsiTheme,attr,omitempty"`
    EastAsiaTheme  *string `xml:"eastAsiaTheme,attr,omitempty"`
    CSTheme        *string `xml:"cstheme,attr,omitempty"`
    Hint           *string `xml:"hint,attr,omitempty"`
}

// === Язык ===
type CT_Language struct {
    Val      *string `xml:"val,attr,omitempty"`
    EastAsia *string `xml:"eastAsia,attr,omitempty"`
    Bidi     *string `xml:"bidi,attr,omitempty"`
}
```

---

## C-03: `coreprops`

**Импортирует**: ничего (только stdlib)

```go
package coreprops

import "time"

type CoreProperties struct {
    Title          string
    Subject        string
    Creator        string
    Keywords       string
    Description    string
    LastModifiedBy string
    Revision       string
    Created        time.Time
    Modified       time.Time
    Category       string
    ContentStatus  string
}

type AppProperties struct {
    Template             string
    TotalTime            int
    Pages                int
    Words                int
    Characters           int
    Application          string
    DocSecurity          int
    Lines                int
    Paragraphs           int
    Company              string
    AppVersion           string
    CharactersWithSpaces int
}

func ParseCore(data []byte) (*CoreProperties, error)
func SerializeCore(cp *CoreProperties) ([]byte, error)
func ParseApp(data []byte) (*AppProperties, error)
func SerializeApp(ap *AppProperties) ([]byte, error)
func DefaultCore(creator string) *CoreProperties
func DefaultApp() *AppProperties
```

---

## C-05: `wml/shared`

**Импортирует**: `xmltypes` (только `encoding/xml`)

> **КЛЮЧЕВОЙ ПАКЕТ**. Содержит интерфейсы контента и RawXML.
> Решает циклическую зависимость: table нужен BlockLevelElement (параграфы),
> body нужен table → оба импортируют shared, но не друг друга.
> Подробности архитектуры: см. `patterns.md` раздел 1.

```go
package shared

import "encoding/xml"

// ==========================================
// ИНТЕРФЕЙСЫ КОНТЕНТА
// ==========================================

// BlockLevelElement — параграф, таблица, SDT, или неизвестный элемент.
// Используется в: CT_Body, CT_Tc, CT_HdrFtr, CT_Comment, CT_FtnEdn.
type BlockLevelElement interface {
    blockLevelElement()
}

// ParagraphContent — run, hyperlink, bookmark, ins/del, или неизвестный.
// Используется в: CT_P, CT_Hyperlink, CT_SimpleField.
type ParagraphContent interface {
    paragraphContent()
}

// RunContent — текст, br, drawing, fldChar, tab, или неизвестный.
// Используется в: CT_R.
type RunContent interface {
    runContent()
}

// ==========================================
// RAW XML ДЛЯ ROUND-TRIP
// ==========================================

// RawXML хранит неизвестный XML-элемент целиком.
// При unmarshal: неизвестный элемент → RawXML.
// При marshal: восстановить на том же месте.
// Подробности: см. patterns.md раздел 4.
type RawXML struct {
    XMLName xml.Name
    Attrs   []xml.Attr `xml:",any,attr"`
    Inner   []byte     `xml:",innerxml"`
}

// RawXML реализует все три интерфейса (может попасть на любой уровень)
func (RawXML) blockLevelElement() {}
func (RawXML) paragraphContent()  {}
func (RawXML) runContent()        {}

// ==========================================
// ФАБРИКИ — для unmarshal: XML-имя → типизированный элемент
// Подробности: см. patterns.md раздел 1.3
// ==========================================

// BlockLevelFactory создаёт типизированный BlockLevelElement по XML-имени.
// Регистрируется пакетами body/para/table при init().
// Если имя неизвестно → вернуть nil (вызывающий сохранит как RawXML).
type BlockLevelFactory func(name xml.Name) BlockLevelElement

func RegisterBlockFactory(f BlockLevelFactory)
func CreateBlockElement(name xml.Name) BlockLevelElement

// Аналогично для ParagraphContent и RunContent:
type ParagraphContentFactory func(name xml.Name) ParagraphContent
func RegisterParagraphContentFactory(f ParagraphContentFactory)
func CreateParagraphContent(name xml.Name) ParagraphContent

type RunContentFactory func(name xml.Name) RunContent
func RegisterRunContentFactory(f RunContentFactory)
func CreateRunContent(name xml.Name) RunContent
```

---

## C-10: `wml/rpr`

**Импортирует**: `xmltypes`, `wml/shared`

```go
package rpr

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/shared"
)

// CT_RPr — полные свойства run (включая rPrChange для track changes)
type CT_RPr struct {
    Base    CT_RPrBase
    RPrChange *CT_RPrChange          // track changes
    Extra   []shared.RawXML        // w14:*, w15:*, неизвестные элементы
}

// CT_RPrBase — базовые свойства форматирования символов
// Порядок полей = порядок в EG_RPrBase (choice, но Word пишет в фиксированном порядке)
type CT_RPrBase struct {
    RStyle          *xmltypes.CT_String            // ссылка на character style
    RFonts          *xmltypes.CT_Fonts
    B               *xmltypes.CT_OnOff
    BCs             *xmltypes.CT_OnOff
    I               *xmltypes.CT_OnOff
    ICs             *xmltypes.CT_OnOff
    Caps            *xmltypes.CT_OnOff
    SmallCaps       *xmltypes.CT_OnOff
    Strike          *xmltypes.CT_OnOff
    Dstrike         *xmltypes.CT_OnOff
    Outline         *xmltypes.CT_OnOff
    Shadow          *xmltypes.CT_OnOff
    Emboss          *xmltypes.CT_OnOff
    Imprint         *xmltypes.CT_OnOff
    NoProof         *xmltypes.CT_OnOff
    SnapToGrid      *xmltypes.CT_OnOff
    Vanish          *xmltypes.CT_OnOff
    WebHidden       *xmltypes.CT_OnOff
    Color           *xmltypes.CT_Color
    Spacing         *xmltypes.CT_SignedTwipsMeasure
    W               *xmltypes.CT_TextScale
    Kern            *xmltypes.CT_HpsMeasure
    Position        *xmltypes.CT_SignedHpsMeasure
    Sz              *xmltypes.CT_HpsMeasure
    SzCs            *xmltypes.CT_HpsMeasure
    Highlight       *xmltypes.CT_Highlight
    U               *xmltypes.CT_Underline
    Effect          *CT_TextEffect
    Bdr             *xmltypes.CT_Border
    Shd             *xmltypes.CT_Shd
    FitText         *CT_FitText
    VertAlign       *CT_VerticalAlignRun
    Rtl             *xmltypes.CT_OnOff
    Cs              *xmltypes.CT_OnOff
    Em              *CT_Em
    Lang            *xmltypes.CT_Language
    EastAsianLayout *CT_EastAsianLayout
    SpecVanish      *xmltypes.CT_OnOff
    OMath           *xmltypes.CT_OnOff
    Extra           []shared.RawXML
}

// CT_RPrChange — отслеживание изменений форматирования run
type CT_RPrChange struct {
    ID     int    `xml:"id,attr"`
    Author string `xml:"author,attr"`
    Date   string `xml:"date,attr,omitempty"`
    RPr    *CT_RPrBase
}

// CT_ParaRPr — default run properties параграфа (pPr/rPr)
// Отличается от CT_RPr наличием ins/del
type CT_ParaRPr struct {
    Base  CT_RPrBase
    Ins   *CT_TrackChangeRef  // если rPr был вставлен
    Del   *CT_TrackChangeRef  // если rPr был удалён
    Extra []shared.RawXML
}

type CT_TrackChangeRef struct {
    ID     int    `xml:"id,attr"`
    Author string `xml:"author,attr"`
    Date   string `xml:"date,attr,omitempty"`
}

// Мелкие типы, специфичные для rpr
type CT_TextEffect struct {
    Val string `xml:"val,attr"` // "blinkBackground", "lights", ...
}
type CT_FitText struct {
    Val int    `xml:"val,attr"` // DXA
    ID  *int   `xml:"id,attr,omitempty"`
}
type CT_VerticalAlignRun struct {
    Val string `xml:"val,attr"` // "baseline", "superscript", "subscript"
}
type CT_Em struct {
    Val string `xml:"val,attr"` // "none", "dot", "comma", "circle", "underDot"
}
type CT_EastAsianLayout struct {
    ID        *int  `xml:"id,attr,omitempty"`
    Combine   *bool `xml:"combine,attr,omitempty"`
    CombineBrackets *string `xml:"combineBrackets,attr,omitempty"`
    Vert      *bool `xml:"vert,attr,omitempty"`
    VertCompress *bool `xml:"vertCompress,attr,omitempty"`
}

// Marshal/Unmarshal
func (r *CT_RPr) MarshalXML(e *xml.Encoder, start xml.StartElement) error
func (r *CT_RPr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
```

---

## C-11: `wml/ppr`

**Импортирует**: `xmltypes`, `wml/rpr`, `wml/shared`

```go
package ppr

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/rpr"
    "github.com/vortex/docx-go/wml/shared"
)

// CT_PPr — полные свойства параграфа
type CT_PPr struct {
    Base      CT_PPrBase
    RPr       *rpr.CT_ParaRPr           // из пакета wml/rpr (НЕ цикл: rpr не импортирует ppr)
    SectPr    *CT_SectPrRef             // raw XML ссылка (чтобы не тянуть sectpr)
    PPrChange *CT_PPrChange
    Extra     []shared.RawXML
}

// CT_PPrBase — свойства параграфа в СТРОГОМ ПОРЯДКЕ (xsd:sequence!)
// Нарушение порядка → Word показывает "файл повреждён"
type CT_PPrBase struct {
    PStyle              *xmltypes.CT_String        //  1. pStyle
    KeepNext            *xmltypes.CT_OnOff         //  2. keepNext
    KeepLines           *xmltypes.CT_OnOff         //  3. keepLines
    PageBreakBefore     *xmltypes.CT_OnOff         //  4. pageBreakBefore
    FramePr             *CT_FramePr                //  5. framePr
    WidowControl        *xmltypes.CT_OnOff         //  6. widowControl
    NumPr               *CT_NumPr                  //  7. numPr
    SuppressLineNumbers *xmltypes.CT_OnOff         //  8
    PBdr                *CT_PBdr                   //  9. pBdr
    Shd                 *xmltypes.CT_Shd           // 10. shd
    Tabs                *CT_Tabs                   // 11. tabs
    SuppressAutoHyphens *xmltypes.CT_OnOff         // 12
    Kinsoku             *xmltypes.CT_OnOff         // 13
    WordWrap            *xmltypes.CT_OnOff         // 14
    OverflowPunct       *xmltypes.CT_OnOff         // 15
    TopLinePunct        *xmltypes.CT_OnOff         // 16
    AutoSpaceDE         *xmltypes.CT_OnOff         // 17
    AutoSpaceDN         *xmltypes.CT_OnOff         // 18
    Bidi                *xmltypes.CT_OnOff         // 19
    AdjustRightInd      *xmltypes.CT_OnOff         // 20
    SnapToGrid          *xmltypes.CT_OnOff         // 21
    Spacing             *CT_Spacing                // 22. spacing
    Ind                 *CT_Ind                    // 23. ind
    ContextualSpacing   *xmltypes.CT_OnOff         // 24
    MirrorIndents       *xmltypes.CT_OnOff         // 25
    SuppressOverlap     *xmltypes.CT_OnOff         // 26
    Jc                  *CT_Jc                     // 27. jc
    TextDirection       *CT_TextDirection          // 28
    TextAlignment       *CT_TextAlignment          // 29
    TextboxTightWrap    *CT_TextboxTightWrap       // 30
    OutlineLvl          *xmltypes.CT_DecimalNumber // 31
    DivId               *xmltypes.CT_DecimalNumber // 32
    CnfStyle            *CT_Cnf                    // 33
    Extra               []shared.RawXML
}

// Дочерние типы
type CT_Spacing struct {
    Before             *int    `xml:"before,attr,omitempty"`      // DXA
    BeforeLines        *int    `xml:"beforeLines,attr,omitempty"` // 1/100 строки
    BeforeAutospacing  *bool   `xml:"beforeAutospacing,attr,omitempty"`
    After              *int    `xml:"after,attr,omitempty"`       // DXA
    AfterLines         *int    `xml:"afterLines,attr,omitempty"`
    AfterAutospacing   *bool   `xml:"afterAutospacing,attr,omitempty"`
    Line               *int    `xml:"line,attr,omitempty"`        // DXA или 1/240 pt
    LineRule           *string `xml:"lineRule,attr,omitempty"`    // "auto"|"exact"|"atLeast"
}

type CT_Ind struct {
    Start          *int `xml:"start,attr,omitempty"`         // DXA (или "left" в transitional)
    StartChars     *int `xml:"startChars,attr,omitempty"`
    End            *int `xml:"end,attr,omitempty"`           // DXA (или "right" в transitional)
    EndChars       *int `xml:"endChars,attr,omitempty"`
    Hanging        *int `xml:"hanging,attr,omitempty"`       // DXA
    HangingChars   *int `xml:"hangingChars,attr,omitempty"`
    FirstLine      *int `xml:"firstLine,attr,omitempty"`     // DXA
    FirstLineChars *int `xml:"firstLineChars,attr,omitempty"`
}

type CT_Jc struct {
    Val string `xml:"val,attr"` // "start"|"center"|"end"|"both"|"left"|"right"
}

type CT_NumPr struct {
    Ilvl  *xmltypes.CT_DecimalNumber // уровень вложенности (0-8)
    NumId *xmltypes.CT_DecimalNumber // ссылка на numbering.xml
}

type CT_Tabs struct {
    Tab []CT_TabStop
}
type CT_TabStop struct {
    Val    string  `xml:"val,attr"`    // "start"|"center"|"end"|"decimal"|"bar"|"clear"
    Pos    int     `xml:"pos,attr"`    // DXA
    Leader *string `xml:"leader,attr,omitempty"` // "none"|"dot"|"hyphen"|"underscore"|"heavy"
}

type CT_PBdr struct {
    Top     *xmltypes.CT_Border
    Bottom  *xmltypes.CT_Border
    Left    *xmltypes.CT_Border
    Right   *xmltypes.CT_Border
    Between *xmltypes.CT_Border
    Bar     *xmltypes.CT_Border
}

type CT_FramePr struct { /* dropCap, lines, w, h, hSpace, vSpace, wrap, anchors */ }
type CT_TextDirection struct { Val string `xml:"val,attr"` }
type CT_TextAlignment struct { Val string `xml:"val,attr"` }
type CT_TextboxTightWrap struct { Val string `xml:"val,attr"` }
type CT_Cnf struct { /* 12 boolean атрибутов conditional formatting */ }
type CT_SectPrRef struct { InnerXML []byte `xml:",innerxml"` } // raw reference

type CT_PPrChange struct {
    ID     int    `xml:"id,attr"`
    Author string `xml:"author,attr"`
    Date   string `xml:"date,attr,omitempty"`
    PPr    *CT_PPrBase
}

// КРИТИЧНО: кастомный MarshalXML для соблюдения порядка xsd:sequence
func (p *CT_PPrBase) MarshalXML(e *xml.Encoder, start xml.StartElement) error
func (p *CT_PPrBase) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
```

---

## C-12: `wml/sectpr`

**Импортирует**: `xmltypes`, `wml/shared`

```go
package sectpr

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/shared"
)

type CT_SectPr struct {
    // EG_HdrFtrReferences
    HeaderRefs []CT_HdrFtrRef
    FooterRefs []CT_HdrFtrRef
    // EG_SectPrContents (в СТРОГОМ ПОРЯДКЕ)
    FootnotePr *CT_FtnProps
    EndnotePr  *CT_EdnProps
    Type       *CT_SectType
    PgSz       *CT_PageSz
    PgMar      *CT_PageMar
    PaperSrc   *CT_PaperSource
    PgBorders  *CT_PageBorders
    LnNumType  *CT_LineNumber
    PgNumType  *CT_PageNumber
    Cols       *CT_Columns
    FormProt   *xmltypes.CT_OnOff
    VAlign     *CT_VerticalJc
    NoEndnote  *xmltypes.CT_OnOff
    TitlePg    *xmltypes.CT_OnOff
    TextDirection *CT_TextDirection
    Bidi       *xmltypes.CT_OnOff
    RtlGutter  *xmltypes.CT_OnOff
    DocGrid    *CT_DocGrid
    // Атрибуты
    RsidR      *string `xml:"rsidR,attr,omitempty"`
    RsidSect   *string `xml:"rsidSect,attr,omitempty"`
    Extra      []shared.RawXML
}

type CT_HdrFtrRef struct {
    Type string `xml:"type,attr"` // "default"|"first"|"even"
    RID  string `xml:"id,attr"`   // r:id → relationship
}

type CT_PageSz struct {
    W      int     `xml:"w,attr"`                      // DXA
    H      int     `xml:"h,attr"`                      // DXA
    Orient *string `xml:"orient,attr,omitempty"`        // "portrait"|"landscape"
    Code   *int    `xml:"code,attr,omitempty"`          // paper size code
}

type CT_PageMar struct {
    Top    int `xml:"top,attr"`    // signed DXA
    Right  int `xml:"right,attr"`
    Bottom int `xml:"bottom,attr"` // signed DXA
    Left   int `xml:"left,attr"`
    Header int `xml:"header,attr"`
    Footer int `xml:"footer,attr"`
    Gutter int `xml:"gutter,attr"`
}

type CT_Columns struct {
    EqualWidth *bool `xml:"equalWidth,attr,omitempty"`
    Space      *int  `xml:"space,attr,omitempty"`       // DXA
    Num        *int  `xml:"num,attr,omitempty"`
    Sep        *bool `xml:"sep,attr,omitempty"`
    Col        []CT_Column
}
type CT_Column struct {
    W     *int `xml:"w,attr,omitempty"`
    Space *int `xml:"space,attr,omitempty"`
}

type CT_DocGrid struct {
    Type      *string `xml:"type,attr,omitempty"`
    LinePitch *int    `xml:"linePitch,attr,omitempty"`
    CharSpace *int    `xml:"charSpace,attr,omitempty"`
}

type CT_SectType struct { Val string `xml:"val,attr"` }
type CT_VerticalJc struct { Val string `xml:"val,attr"` }
type CT_PageNumber struct { /* fmt, start, chapStyle, chapSep */ }
type CT_PageBorders struct { /* top, bottom, left, right CT_PageBorder */ }
type CT_LineNumber struct { /* countBy, start, restart, distance */ }
type CT_PaperSource struct { /* first, other */ }
type CT_FtnProps struct { /* pos, numFmt, numStart, numRestart */ }
type CT_EdnProps struct { /* pos, numFmt, numStart, numRestart */ }
type CT_TextDirection struct { Val string `xml:"val,attr"` }
```

---

## C-13: `wml/table`

**Импортирует**: `xmltypes`, `wml/shared`

```go
package table

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/shared"
)

type CT_Tbl struct {
    TblPr   *CT_TblPr
    TblGrid *CT_TblGrid
    // Содержимое: rows + bookmarks + track changes
    Content []TblContent // interface: CT_Row | RawXML (для EG_RunLevelElts)
    Extra   []shared.RawXML
}

// TblContent — интерфейс для содержимого таблицы
type TblContent interface {
    isTblContent()
}

type CT_TblPr struct {
    // СТРОГИЙ ПОРЯДОК (xsd:sequence)
    TblStyle              *xmltypes.CT_String
    TblpPr                *CT_TblPPr
    TblOverlap            *CT_TblOverlap
    BidiVisual            *xmltypes.CT_OnOff
    TblStyleRowBandSize   *xmltypes.CT_DecimalNumber
    TblStyleColBandSize   *xmltypes.CT_DecimalNumber
    TblW                  *CT_TblWidth
    Jc                    *CT_JcTable
    TblCellSpacing        *CT_TblWidth
    TblInd                *CT_TblWidth
    TblBorders            *CT_TblBorders
    Shd                   *xmltypes.CT_Shd
    TblLayout             *CT_TblLayoutType
    TblCellMar            *CT_TblCellMar
    TblLook               *CT_TblLook
    TblCaption            *xmltypes.CT_String
    TblDescription        *xmltypes.CT_String
    TblPrChange           *CT_TblPrChange
    Extra                 []shared.RawXML
}

type CT_TblGrid struct {
    GridCol []CT_TblGridCol
}
type CT_TblGridCol struct {
    W int `xml:"w,attr"` // DXA
}

type CT_Row struct {
    TblPrEx *CT_TblPrEx
    TrPr    *CT_TrPr
    Content []RowContent // CT_Tc | CT_CustomXmlCell | CT_SdtCell
    // Атрибуты
    RsidR   *string `xml:"rsidR,attr,omitempty"`
    RsidTr  *string `xml:"rsidTr,attr,omitempty"`
    Extra   []shared.RawXML
}

type RowContent interface {
    isRowContent()
}

type CT_Tc struct {
    TcPr    *CT_TcPr
    Content []shared.BlockLevelElement // параграфы и таблицы внутри ячейки
    // ВАЖНО: всегда содержит ≥1 элемент (минимум пустой <w:p/>)
}

type CT_TcPr struct {
    CnfStyle       *CT_Cnf
    TcW            *CT_TblWidth
    GridSpan       *xmltypes.CT_DecimalNumber
    HMerge         *CT_HMerge
    VMerge         *CT_VMerge
    TcBorders      *CT_TcBorders
    Shd            *xmltypes.CT_Shd
    NoWrap         *xmltypes.CT_OnOff
    TcMar          *CT_TblCellMar
    TextDirection  *CT_TextDirection
    TcFitText      *xmltypes.CT_OnOff
    VAlign         *CT_VerticalJc
    HideMark       *xmltypes.CT_OnOff
    TcPrChange     *CT_TcPrChange
    Extra          []shared.RawXML
}

type CT_TrPr struct {
    CnfStyle       *CT_Cnf
    GridBefore     *xmltypes.CT_DecimalNumber
    GridAfter      *xmltypes.CT_DecimalNumber
    WBefore        *CT_TblWidth
    WAfter         *CT_TblWidth
    CantSplit      *xmltypes.CT_OnOff
    TrHeight       *CT_Height
    TblHeader      *xmltypes.CT_OnOff
    TblCellSpacing *CT_TblWidth
    Jc             *CT_JcTable
    Hidden         *xmltypes.CT_OnOff
    TrPrChange     *CT_TrPrChange
    Extra          []shared.RawXML
}

type CT_TblWidth struct {
    W    int    `xml:"w,attr"`
    Type string `xml:"type,attr"` // "nil"|"pct"|"dxa"|"auto"
}

type CT_TblBorders struct {
    Top     *xmltypes.CT_Border
    Start   *xmltypes.CT_Border // или Left в transitional
    Bottom  *xmltypes.CT_Border
    End     *xmltypes.CT_Border // или Right в transitional
    InsideH *xmltypes.CT_Border
    InsideV *xmltypes.CT_Border
}

type CT_TcBorders struct {
    Top     *xmltypes.CT_Border
    Start   *xmltypes.CT_Border
    Bottom  *xmltypes.CT_Border
    End     *xmltypes.CT_Border
    InsideH *xmltypes.CT_Border
    InsideV *xmltypes.CT_Border
    Tl2br   *xmltypes.CT_Border
    Tr2bl   *xmltypes.CT_Border
}

type CT_TblCellMar struct {
    Top    *CT_TblWidth
    Start  *CT_TblWidth
    Bottom *CT_TblWidth
    End    *CT_TblWidth
}

type CT_TblLook struct {
    FirstRow    *bool `xml:"firstRow,attr,omitempty"`
    LastRow     *bool `xml:"lastRow,attr,omitempty"`
    FirstColumn *bool `xml:"firstColumn,attr,omitempty"`
    LastColumn  *bool `xml:"lastColumn,attr,omitempty"`
    NoHBand     *bool `xml:"noHBand,attr,omitempty"`
    NoVBand     *bool `xml:"noVBand,attr,omitempty"`
}

type CT_Height struct {
    Val   *int    `xml:"val,attr,omitempty"`
    HRule *string `xml:"hRule,attr,omitempty"` // "auto"|"exact"|"atLeast"
}
type CT_HMerge struct { Val *string `xml:"val,attr,omitempty"` } // "restart"|"continue"
type CT_VMerge struct { Val *string `xml:"val,attr,omitempty"` }
type CT_TblLayoutType struct { Type string `xml:"type,attr"` } // "fixed"|"autofit"
type CT_JcTable struct { Val string `xml:"val,attr"` }
type CT_TblPPr struct { /* floating table positioning */ }
type CT_TblOverlap struct { Val string `xml:"val,attr"` }
type CT_VerticalJc struct { Val string `xml:"val,attr"` }
type CT_TextDirection struct { Val string `xml:"val,attr"` }
type CT_Cnf struct { /* 12 booleans */ }
type CT_TblPrChange struct { ID int `xml:"id,attr"`; Author string `xml:"author,attr"` }
type CT_TrPrChange struct { ID int `xml:"id,attr"`; Author string `xml:"author,attr"` }
type CT_TcPrChange struct { ID int `xml:"id,attr"`; Author string `xml:"author,attr"` }
type CT_TblPrEx struct { /* row-level table property overrides */ }
```

---

## C-14: `wml/tracking`

**Импортирует**: `xmltypes`, `wml/shared`

```go
package tracking

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/shared"
)

// CT_RunTrackChange — обёртка для <w:ins>, <w:del>, <w:moveFrom>, <w:moveTo>
type CT_RunTrackChange struct {
    ID       int            `xml:"id,attr"`
    Author   string         `xml:"author,attr"`
    Date     *string        `xml:"date,attr,omitempty"`
    Content  []interface{}  // RunContent внутри ins/del
}

// CT_Markup — базовый тип для commentReference, bookmarkStart и т.д.
type CT_Markup struct {
    ID int `xml:"id,attr"`
}

type CT_Bookmark struct {
    ID       int    `xml:"id,attr"`
    Name     string `xml:"name,attr"`
    ColFirst *int   `xml:"colFirst,attr,omitempty"`
    ColLast  *int   `xml:"colLast,attr,omitempty"`
}

type CT_MarkupRange struct {
    ID int `xml:"id,attr"`
}

type CT_MoveBookmark struct {
    ID     int    `xml:"id,attr"`
    Author string `xml:"author,attr"`
    Name   string `xml:"name,attr"`
}
```

---

## C-15: `wml/run`

**Импортирует**: `xmltypes`, `wml/rpr`, `wml/shared`, `dml`

```go
package run

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/rpr"
    "github.com/vortex/docx-go/wml/shared"
)

type CT_R struct {
    RPr     *rpr.CT_RPr
    Content []shared.RunContent // интерфейс из wml/shared
    // Атрибуты
    RsidR    *string `xml:"rsidR,attr,omitempty"`
    RsidRPr  *string `xml:"rsidRPr,attr,omitempty"`
    RsidDel  *string `xml:"rsidDel,attr,omitempty"`
}

// Реализации shared.RunContent (все реализуют runContent()):
type CT_Text struct {
    Space *string `xml:"space,attr,omitempty"` // "preserve"
    Value string  `xml:",chardata"`
}
func (CT_Text) runContent() {}

type CT_Br struct {
    Type  *string `xml:"type,attr,omitempty"`  // "page"|"column"|"textWrapping"
    Clear *string `xml:"clear,attr,omitempty"` // "none"|"left"|"right"|"all"
}
func (CT_Br) runContent() {}

type CT_Drawing struct {
    Inline []interface{} // dml.WP_Inline
    Anchor []interface{} // dml.WP_Anchor
}
func (CT_Drawing) runContent() {}

type CT_FldChar struct {
    FldCharType string `xml:"fldCharType,attr"` // "begin"|"separate"|"end"
    FldLock     *bool  `xml:"fldLock,attr,omitempty"`
    Dirty       *bool  `xml:"dirty,attr,omitempty"`
}
func (CT_FldChar) runContent() {}

type CT_InstrText struct {
    Space *string `xml:"space,attr,omitempty"`
    Value string  `xml:",chardata"`
}
func (CT_InstrText) runContent() {}

type CT_Sym struct {
    Font string `xml:"font,attr"`
    Char string `xml:"char,attr"`
}
func (CT_Sym) runContent() {}

type CT_FtnEdnRef struct {
    ID int `xml:"id,attr"`
}
func (CT_FtnEdnRef) runContent() {}

// Пустые элементы: tab, cr, pgNum, noBreakHyphen, softHyphen,
// footnoteRef, endnoteRef, annotationRef, separator, continuationSeparator
type CT_EmptyRunContent struct {
    XMLName xml.Name
}
func (CT_EmptyRunContent) runContent() {}

// Неизвестный контент (для round-trip)
type CT_RawRunContent struct {
    Raw shared.RawXML
}
func (CT_RawRunContent) runContent() {}
```

---

## C-16: `wml/para`

**Импортирует**: `xmltypes`, `wml/rpr`, `wml/ppr`, `wml/run`, `wml/tracking`, `wml/shared`

```go
package para

import (
    "github.com/vortex/docx-go/wml/ppr"
    "github.com/vortex/docx-go/wml/run"
    "github.com/vortex/docx-go/wml/tracking"
    "github.com/vortex/docx-go/wml/shared"
    "github.com/vortex/docx-go/xmltypes"
)

type CT_P struct {
    PPr     *ppr.CT_PPr
    Content []shared.ParagraphContent // интерфейс из wml/shared
    // Атрибуты
    RsidR        *string `xml:"rsidR,attr,omitempty"`
    RsidRDefault *string `xml:"rsidRDefault,attr,omitempty"`
    RsidP        *string `xml:"rsidP,attr,omitempty"`
    RsidRPr      *string `xml:"rsidRPr,attr,omitempty"`
    RsidDel      *string `xml:"rsidDel,attr,omitempty"`
    ParaId       *string // w14:paraId
    TextId       *string // w14:textId
}

// Реализации shared.ParagraphContent (все реализуют paragraphContent()):
type RunItem struct { R *run.CT_R }
func (RunItem) paragraphContent() {}

type HyperlinkItem struct { H *CT_Hyperlink }
func (HyperlinkItem) paragraphContent() {}

type SimpleFieldItem struct { F *CT_SimpleField }
func (SimpleFieldItem) paragraphContent() {}

type InsItem struct { Ins *tracking.CT_RunTrackChange }
func (InsItem) paragraphContent() {}

type DelItem struct { Del *tracking.CT_RunTrackChange }
func (DelItem) paragraphContent() {}

type BookmarkStartItem struct { B *tracking.CT_Bookmark }
func (BookmarkStartItem) paragraphContent() {}

type BookmarkEndItem struct { B *tracking.CT_MarkupRange }
func (BookmarkEndItem) paragraphContent() {}

type CommentRangeStartItem struct { C *tracking.CT_Markup }
func (CommentRangeStartItem) paragraphContent() {}

type CommentRangeEndItem struct { C *tracking.CT_Markup }
func (CommentRangeEndItem) paragraphContent() {}

type SdtRunItem struct { Sdt *CT_SdtRun }
func (SdtRunItem) paragraphContent() {}

type RawParagraphContent struct { Raw shared.RawXML }
func (RawParagraphContent) paragraphContent() {}

// Hyperlink
type CT_Hyperlink struct {
    RID        *string `xml:"id,attr,omitempty"`     // r:id для external
    Anchor     *string `xml:"anchor,attr,omitempty"` // для internal
    Content    []shared.ParagraphContent // runs, fields внутри гиперссылки
}

// Простое поле
type CT_SimpleField struct {
    Instr   string `xml:"instr,attr"`
    Content []shared.ParagraphContent
}

// SDT (content control)
type CT_SdtRun struct {
    SdtPr      *shared.RawXML // сложная структура — raw для начала
    SdtEndPr   *shared.RawXML
    SdtContent []shared.ParagraphContent
}
```

---

## C-17: `wml/body`

**Импортирует**: `xmltypes`, `wml/para`, `wml/table`, `wml/sectpr`, `wml/shared`

```go
package body

import (
    "github.com/vortex/docx-go/wml/para"
    "github.com/vortex/docx-go/wml/table"
    "github.com/vortex/docx-go/wml/sectpr"
    "github.com/vortex/docx-go/wml/shared"
    "github.com/vortex/docx-go/xmltypes"
)

type CT_Document struct {
    Body  *CT_Body
    Extra []shared.RawXML
    // Сохранённые namespace declarations корневого элемента
    // Подробности: см. patterns.md раздел 6
    Namespaces []xml.Attr
}

type CT_Body struct {
    Content []shared.BlockLevelElement // интерфейс из wml/shared
    SectPr  *sectpr.CT_SectPr          // последний sectPr (body-level)
}

// Реализации shared.BlockLevelElement (все реализуют blockLevelElement()):
type ParagraphElement struct { P *para.CT_P }
func (ParagraphElement) blockLevelElement() {}

type TableElement struct { T *table.CT_Tbl }
func (TableElement) blockLevelElement() {}

type SdtBlockElement struct { Sdt *CT_SdtBlock }
func (SdtBlockElement) blockLevelElement() {}

type RawBlockElement struct { Raw shared.RawXML }
func (RawBlockElement) blockLevelElement() {}

type CT_SdtBlock struct {
    SdtPr      *shared.RawXML
    SdtEndPr   *shared.RawXML
    SdtContent []shared.BlockLevelElement
}
```

---

## C-18: `wml/hdft`

**Импортирует**: `wml/shared`

```go
package hdft

import "github.com/vortex/docx-go/wml/shared"

type CT_HdrFtr struct {
    Content []shared.BlockLevelElement
    // ВАЖНО: всегда содержит ≥1 элемент (минимум <w:p/>)
}

func Parse(data []byte) (*CT_HdrFtr, error)
func Serialize(hf *CT_HdrFtr, rootName string) ([]byte, error) // rootName: "w:hdr" | "w:ftr"
```

---

## C-19: `dml`

**Импортирует**: `xmltypes`, `wml/shared`

```go
package dml

import (
    "github.com/vortex/docx-go/xmltypes"
    "github.com/vortex/docx-go/wml/shared"
)

// Inline image
type WP_Inline struct {
    DistT        int           `xml:"distT,attr"`
    DistB        int           `xml:"distB,attr"`
    DistL        int           `xml:"distL,attr"`
    DistR        int           `xml:"distR,attr"`
    Extent       WP_Extent
    EffectExtent *WP_EffectExtent
    DocPr        WP_DocPr
    Graphic      A_Graphic
}

// Floating image
type WP_Anchor struct {
    BehindDoc      bool
    DistT, DistB   int
    DistL, DistR   int
    RelativeHeight int
    SimplePos      bool
    Locked         bool
    LayoutInCell   bool
    AllowOverlap   bool
    SimplePosXY    *WP_Point
    PositionH      WP_PosH
    PositionV      WP_PosV
    Extent         WP_Extent
    EffectExtent   *WP_EffectExtent
    WrapType       interface{} // WP_WrapNone | WP_WrapSquare | ...
    DocPr          WP_DocPr
    Graphic        A_Graphic
}

type WP_Extent struct {
    CX int64 `xml:"cx,attr"` // EMU
    CY int64 `xml:"cy,attr"` // EMU
}
type WP_EffectExtent struct {
    L, T, R, B int64 // EMU
}
type WP_DocPr struct {
    ID    int    `xml:"id,attr"`
    Name  string `xml:"name,attr"`
    Descr string `xml:"descr,attr,omitempty"`
}
type WP_PosH struct {
    RelativeFrom string // "column"|"page"|"margin"|...
    PosOffset    *int64  // EMU
    Align        *string // "left"|"center"|"right"
}
type WP_PosV struct {
    RelativeFrom string
    PosOffset    *int64
    Align        *string
}
type WP_Point struct { X, Y int64 }
type WP_WrapNone struct{}
type WP_WrapSquare struct { WrapText string }
type WP_WrapTight struct { WrapText string }
type WP_WrapTopAndBottom struct{}

type A_Graphic struct {
    GraphicData A_GraphicData
}
type A_GraphicData struct {
    URI     string // "http://schemas.openxmlformats.org/drawingml/2006/picture"
    Pic     *PIC_Pic
    RawData *shared.RawXML // fallback для не-picture
}
type PIC_Pic struct {
    NvPicPr  PIC_NvPicPr
    BlipFill PIC_BlipFill
    SpPr     A_SpPr
}
type PIC_NvPicPr struct {
    CNvPr    WP_DocPr
}
type PIC_BlipFill struct {
    Blip    A_Blip
    Stretch *A_Stretch
}
type A_Blip struct {
    Embed string `xml:"embed,attr"` // r:embed → rId
    Link  string `xml:"link,attr,omitempty"`
}
type A_Stretch struct{}
type A_SpPr struct {
    Xfrm     *A_Xfrm
    PrstGeom *A_PrstGeom
}
type A_Xfrm struct {
    Off A_Off
    Ext A_Ext
}
type A_Off struct { X, Y int64 }
type A_Ext struct { CX, CY int64 }
type A_PrstGeom struct {
    Prst string // "rect", "ellipse", ...
}
```

---

## C-20..C-29: `parts/*`

Все part-модули имеют одинаковый контракт:

```go
package <partname>

// Parse десериализует XML-байты в типизированную структуру
func Parse(data []byte) (*<RootType>, error)

// Serialize сериализует структуру обратно в XML-байты
func Serialize(obj *<RootType>) ([]byte, error)
```

| Пакет | RootType | Зависимости |
|-------|----------|-------------|
| `parts/document` | `body.CT_Document` | wml/body, opc |
| `parts/styles` | `CT_Styles` | xmltypes, wml/rpr, wml/ppr, wml/table, wml/shared |
| `parts/numbering` | `CT_Numbering` | xmltypes, wml/rpr, wml/ppr, wml/shared |
| `parts/settings` | `CT_Settings` (определён ниже) | xmltypes, wml/shared |
| `parts/fonts` | `CT_FontsList` (определён ниже) | xmltypes, wml/shared |
| `parts/comments` | `CT_Comments` | wml/shared |
| `parts/footnotes` | `CT_Footnotes` | wml/shared |
| `parts/headers` | `hdft.CT_HdrFtr` | wml/shared |
| `parts/theme` | `[]byte` (raw passthrough) | — |
| `parts/websettings` | `[]byte` (raw passthrough) | — |

### parts/styles — ключевые типы:

```go
type CT_Styles struct {
    DocDefaults   *CT_DocDefaults
    LatentStyles  *CT_LatentStyles
    Style         []CT_Style
    Extra         []shared.RawXML
}

type CT_DocDefaults struct {
    RPrDefault *CT_RPrDefault
    PPrDefault *CT_PPrDefault
}
type CT_RPrDefault struct { RPr *rpr.CT_RPr }
type CT_PPrDefault struct { PPr *ppr.CT_PPrBase }

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

type CT_Style struct {
    // Атрибуты
    Type         string  `xml:"type,attr"`
    Default      *bool   `xml:"default,attr,omitempty"`
    CustomStyle  *bool   `xml:"customStyle,attr,omitempty"`
    StyleID      string  `xml:"styleId,attr"`
    // Элементы (СТРОГИЙ ПОРЯДОК xsd:sequence — маппинг в patterns.md 2.8)
    Name         *xmltypes.CT_String
    Aliases      *xmltypes.CT_String
    BasedOn      *xmltypes.CT_String
    Next         *xmltypes.CT_String
    Link         *xmltypes.CT_String
    AutoRedefine *xmltypes.CT_OnOff
    Hidden       *xmltypes.CT_OnOff
    UIpriority   *xmltypes.CT_DecimalNumber
    SemiHidden   *xmltypes.CT_OnOff
    UnhideWhenUsed *xmltypes.CT_OnOff
    QFormat      *xmltypes.CT_OnOff
    Locked       *xmltypes.CT_OnOff
    Personal     *xmltypes.CT_OnOff
    PersonalCompose *xmltypes.CT_OnOff
    PersonalReply *xmltypes.CT_OnOff
    Rsid         *xmltypes.CT_LongHexNumber
    PPr          *ppr.CT_PPrBase
    RPr          *rpr.CT_RPrBase
    TblPr        *table.CT_TblPr
    TrPr         *table.CT_TrPr
    TcPr         *table.CT_TcPr
    TblStylePr   []CT_TblStylePr
    Extra        []shared.RawXML
}

type CT_TblStylePr struct {
    Type string          `xml:"type,attr"` // "firstRow"|"lastRow"|"band1Vert"|...
    PPr  *ppr.CT_PPrBase
    RPr  *rpr.CT_RPrBase
    TblPr *table.CT_TblPr
    TrPr  *table.CT_TrPr
    TcPr  *table.CT_TcPr
}
```

### parts/numbering — ключевые типы:

```go
type CT_Numbering struct {
    AbstractNum []CT_AbstractNum
    Num         []CT_Num
    Extra       []shared.RawXML
}
type CT_AbstractNum struct {
    AbstractNumID int
    Nsid          *xmltypes.CT_LongHexNumber
    MultiLevelType *xmltypes.CT_String
    Tmpl          *xmltypes.CT_LongHexNumber
    Name          *xmltypes.CT_String
    StyleLink     *xmltypes.CT_String
    NumStyleLink  *xmltypes.CT_String
    Lvl           []CT_Lvl
}
type CT_Lvl struct {
    // Атрибуты
    Ilvl       int     `xml:"ilvl,attr"`
    Tplc       *string `xml:"tplc,attr,omitempty"`
    Tentative  *bool   `xml:"tentative,attr,omitempty"`
    // Элементы (СТРОГИЙ ПОРЯДОК — маппинг в patterns.md 2.9)
    Start      *xmltypes.CT_DecimalNumber
    NumFmt     *xmltypes.CT_String
    LvlRestart *xmltypes.CT_DecimalNumber
    PStyle     *xmltypes.CT_String
    IsLgl      *xmltypes.CT_OnOff
    Suff       *xmltypes.CT_String
    LvlText    *xmltypes.CT_String
    LvlPicBulletId *xmltypes.CT_DecimalNumber
    LvlJc      *ppr.CT_Jc
    PPr        *ppr.CT_PPrBase
    RPr        *rpr.CT_RPrBase
    Extra      []shared.RawXML
}
type CT_Num struct {
    NumID          int `xml:"numId,attr"`
    AbstractNumID  xmltypes.CT_DecimalNumber
    LvlOverride    []CT_NumLvl
}
type CT_NumLvl struct {
    Ilvl          int `xml:"ilvl,attr"`
    StartOverride *xmltypes.CT_DecimalNumber
    Lvl           *CT_Lvl // полное переопределение уровня
}
```

### parts/settings — CT_Settings (partial parsing, 253 элемента в XSD):

> Стратегия: ~20 часто используемых полей типизированы, остальные 230+ → RawXML.
> Порядок элементов сохраняется через elementOrder для round-trip.

```go
type CT_Settings struct {
    // Типизированные (часто используемые)
    WriteProtection         *CT_WriteProtection
    Zoom                    *CT_Zoom
    ProofState              *CT_Proof
    DefaultTabStop          *xmltypes.CT_TwipsMeasure
    CharacterSpacingControl *xmltypes.CT_String
    EvenAndOddHeaders       *xmltypes.CT_OnOff
    MirrorMargins           *xmltypes.CT_OnOff
    TrackRevisions          *xmltypes.CT_OnOff
    DoNotTrackMoves         *xmltypes.CT_OnOff
    DoNotTrackFormatting    *xmltypes.CT_OnOff
    DocumentProtection      *CT_DocProtect
    Compat                  *CT_Compat
    Rsids                   *CT_DocRsids
    MathPr                  *shared.RawXML         // Math ML — raw для MVP
    ThemeFontLang           *CT_ThemeFontLang
    ClrSchemeMapping        *CT_ClrSchemeMapping   // ⚠️ ОБЯЗАТЕЛЬНЫЙ (см. reference-appendix 5.11)
    ShapeDefaults           *shared.RawXML         // ⚠️ ОБЯЗАТЕЛЬНЫЙ (VML — raw)
    DecimalSymbol           *xmltypes.CT_String
    ListSeparator           *xmltypes.CT_String
    DocId14                 *xmltypes.CT_LongHexNumber // w14:docId
    DocId15                 *xmltypes.CT_Guid          // w15:docId

    // ВСЕ остальные 230+ элементов
    Extra []shared.RawXML

    // Порядок элементов для round-trip (private)
    // elementOrder []string
}

type CT_Zoom struct {
    Percent int    `xml:"percent,attr"`
    Val     string `xml:"val,attr,omitempty"`
}
type CT_Proof struct {
    Spelling *string `xml:"spelling,attr,omitempty"`
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
type CT_WriteProtection struct { /* recommended, algorithm, hash, salt */ }
type CT_DocProtect struct { /* edit, enforcement, algorithm, hash, salt */ }
```

### parts/fonts — CT_FontsList:

```go
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
    EmbedRegular    *CT_FontRel
    EmbedBold       *CT_FontRel
    EmbedItalic     *CT_FontRel
    EmbedBoldItalic *CT_FontRel
    Extra           []shared.RawXML
}
type CT_Panose struct { Val string `xml:"val,attr"` }
type CT_Charset struct { Val string `xml:"val,attr"` }
type CT_FontFamily struct { Val string `xml:"val,attr"` }
type CT_Pitch struct { Val string `xml:"val,attr"` }
type CT_FontSig struct {
    Usb0 string `xml:"usb0,attr"`
    Usb1 string `xml:"usb1,attr"`
    Usb2 string `xml:"usb2,attr"`
    Usb3 string `xml:"usb3,attr"`
    Csb0 string `xml:"csb0,attr"`
    Csb1 string `xml:"csb1,attr"`
}
type CT_FontRel struct {
    FontKey    *string `xml:"fontKey,attr,omitempty"`
    SubsetInfo *string `xml:"subsetted,attr,omitempty"`
    ID         string  `xml:"id,attr"` // r:id
}
```

### parts/comments:

```go
type CT_Comments struct {
    Comment []CT_Comment
}
type CT_Comment struct {
    ID       int    `xml:"id,attr"`
    Author   string `xml:"author,attr"`
    Date     string `xml:"date,attr,omitempty"`
    Initials string `xml:"initials,attr,omitempty"`
    Content  []shared.BlockLevelElement
}
```

### parts/footnotes (и endnotes — аналогично):

```go
type CT_Footnotes struct {
    Footnote []CT_FtnEdn
}
type CT_FtnEdn struct {
    Type    *string `xml:"type,attr,omitempty"` // "normal"|"separator"|"continuationSeparator"
    ID      int     `xml:"id,attr"`
    Content []shared.BlockLevelElement
}
```

### parts/theme и parts/websettings (raw passthrough):

```go
// parts/theme — хранить как raw XML (DrawingML слишком сложен для MVP)
// При open → сохранить bytes, при save → записать те же bytes
func Parse(data []byte) ([]byte, error)    { return data, nil }
func Serialize(data []byte) ([]byte, error) { return data, nil }

// parts/websettings — аналогично
// Файл обязателен (reference-appendix 5.13), но содержимое можно не парсить
```

---

## C-30: `packaging`

**Импортирует**: `opc`, `coreprops`, все `parts/*`, `units`

```go
package packaging

import (
    "github.com/vortex/docx-go/opc"
    "github.com/vortex/docx-go/coreprops"
    "github.com/vortex/docx-go/wml/body"
    // ... и другие parts
)

type Document struct {
    // Типизированные parts
    Document  *body.CT_Document
    Styles    *styles.CT_Styles
    Numbering *numbering.CT_Numbering    // nil если нет
    Settings  *settings.CT_Settings
    Fonts     *fonts.CT_FontsList
    Comments  *comments.CT_Comments       // nil если нет
    Footnotes *footnotes.CT_Footnotes     // nil если нет
    Endnotes  *footnotes.CT_Footnotes     // nil если нет
    CoreProps *coreprops.CoreProperties
    AppProps  *coreprops.AppProperties

    // Колонтитулы: rId → parsed
    Headers   map[string]*hdft.CT_HdrFtr
    Footers   map[string]*hdft.CT_HdrFtr

    // Raw parts (round-trip)
    Theme        []byte
    WebSettings  []byte
    Media        map[string][]byte   // "image1.png" → bytes
    UnknownParts map[string][]byte   // path → raw bytes
    UnknownRels  []opc.Relationship  // неизвестные relationships

    // Internal
    pkg *opc.Package
}

func Open(path string) (*Document, error)
func OpenReader(r io.ReaderAt, size int64) (*Document, error)
func (d *Document) Save(path string) error
func (d *Document) SaveWriter(w io.Writer) error

// Вспомогательные
func (d *Document) NextRelID() string     // генерирует следующий rId
func (d *Document) NextBookmarkID() int   // уникальный ID для bookmark/comment/ins/del
func (d *Document) AddMedia(filename string, data []byte) string // → rId
```

---

## C-31: `validator`

**Импортирует**: `packaging`

```go
package validator

import "github.com/vortex/docx-go/packaging"

type Severity int
const (
    Warn  Severity = iota
    Error
    Fatal
)

type Issue struct {
    Severity Severity
    Part     string // "word/document.xml"
    Path     string // "body/p[3]/r[1]/t"
    Code     string // "EMPTY_TC", "MISSING_REL", ...
    Message  string
}

func Validate(doc *packaging.Document) []Issue
func AutoFix(doc *packaging.Document) []string // описания исправлений
```

---

## C-32: `docx` (Public API)

**Импортирует**: `packaging`, `validator`, `units`,
`wml/body`, `wml/para`, `wml/run`, `wml/table`, `wml/sectpr`, `wml/hdft`, `wml/shared`,
`wml/rpr`, `wml/ppr`, `coreprops`, `parts/styles`

> C-32 — тонкая обёртка поверх WML-типов. Каждый wrapper предоставляет
> convenience-методы (Add/Remove/Set/Get) **и** escape hatch `X()` для
> прямого доступа к нижележащему типу (см. patterns.md раздел 14).
>
> **Как Body.Content хранит элементы (факт реализации C-17 + C-13):**
> - `*table.CT_Tbl` — указатель (CT_Tbl содержит embedded `shared.BlockLevelMarker`)
> - `*para.CT_P` — указатель (аналогично)
> - `shared.RawXML` — value
>
> **Как tbl.Content хранит строки (факт реализации C-13):**
> - `table.CT_Row` — **value** (не указатель), реализует `TblContent` через value receiver
> - `table.RawTblContent` — value
>
> **Как row.Content хранит ячейки (факт реализации C-13):**
> - `table.CT_Tc` — **value** (не указатель), реализует `RowContent` через value receiver
> - `table.RawRowContent` — value
>
> Следствие: Row и Cell wrapper'ы не могут хранить указатель на свой WML-тип
> (значение живёт внутри interface slice). Вместо этого хранят координаты
> (tbl + idx). Подробности: patterns.md раздел 15.

```go
package docx

import (
    "github.com/vortex/docx-go/units"
    "github.com/vortex/docx-go/packaging"
    "github.com/vortex/docx-go/validator"
    "github.com/vortex/docx-go/coreprops"
    "github.com/vortex/docx-go/wml/body"
    "github.com/vortex/docx-go/wml/para"
    "github.com/vortex/docx-go/wml/run"
    "github.com/vortex/docx-go/wml/table"
    "github.com/vortex/docx-go/wml/sectpr"
    "github.com/vortex/docx-go/wml/hdft"
    "github.com/vortex/docx-go/wml/shared"
    "github.com/vortex/docx-go/parts/styles"
)

// ================================================================
// Жизненный цикл
// ================================================================

func Open(path string) (*Document, error)
func New() *Document
func (d *Document) Save(path string) error
func (d *Document) Validate() []validator.Issue

// ================================================================
// Document
// ================================================================

type Document struct { /* wraps *packaging.Document */ }
func (d *Document) Body() *Body
func (d *Document) Styles() *Styles
func (d *Document) AddHeader(hdrType string) *Header    // "default"|"first"|"even"
func (d *Document) AddFooter(ftrType string) *Footer
func (d *Document) CoreProperties() *coreprops.CoreProperties
func (d *Document) SetCoreProperties(cp *coreprops.CoreProperties)
func (d *Document) X() *packaging.Document              // escape hatch

// ================================================================
// Styles — обёртка над styles.CT_Styles
// ================================================================

type Styles struct { /* wraps *styles.CT_Styles */ }
func (s *Styles) X() *styles.CT_Styles                  // escape hatch
// Convenience-методы для Styles — v2 (пока пользователь работает через X())

// ================================================================
// Body — обёртка над body.CT_Body
// ================================================================

// Body хранит *body.CT_Body (указатель). Мутации через wrapper
// отражаются в документе напрямую.
type Body struct {
    raw *body.CT_Body
}

// --- Добавление ---
func (b *Body) AddParagraph() *Paragraph
func (b *Body) AddHeading(text string, level int) *Paragraph
func (b *Body) AddTable(rows, cols int) *Table
func (b *Body) AddPageBreak()
func (b *Body) AddSectionBreak(breakType string) *Section

// --- Вставка в позицию ---
// index = позиция в Body.Content (среди ВСЕХ block-level элементов:
// *para.CT_P, *table.CT_Tbl, shared.RawXML, etc.).
// Элемент вставляется ПЕРЕД указанной позицией.
// index == ElementCount() → вставка в конец (аналогично Add).
// Возвращает error если index < 0 || index > ElementCount().
func (b *Body) InsertParagraphAt(index int) (*Paragraph, error)
func (b *Body) InsertTableAt(index int, rows, cols int) (*Table, error)

// --- Чтение ---
func (b *Body) Paragraphs() []*Paragraph              // только *para.CT_P (фильтр по type assert)
func (b *Body) Tables() []*Table                       // только *table.CT_Tbl (фильтр по type assert)
func (b *Body) ElementCount() int                      // len(raw.Content)
func (b *Body) Section() *Section                      // body-level raw.SectPr

// --- Удаление ---
// index = позиция в Body.Content (как для InsertAt).
// Возвращает error если index вне диапазона.
func (b *Body) RemoveElement(index int) error          // удалить любой block element
func (b *Body) Clear()                                 // удалить всё содержимое (SectPr сохраняется)

// --- Поиск и замена ---
// Работает ТОЛЬКО по параграфам верхнего уровня Body.Content.
// НЕ ищет в таблицах, headers, footers, footnotes (v1 ограничение).
// Работает в пределах одного Run (текст, разбитый между ранами, не находится).
// Подробности и ограничения: см. patterns.md раздел 16.
func (b *Body) FindText(needle string) []TextLocation  // поиск по body-level параграфам
func (b *Body) ReplaceText(old, new string) int         // замена, возвращает кол-во замен

// --- Escape hatch ---
func (b *Body) X() *body.CT_Body                       // прямой доступ к нижнему слою

// ================================================================
// TextLocation — результат FindText
// ================================================================

type TextLocation struct {
    // BlockIndex — позиция элемента в Body.Content[]
    // (по этому индексу лежит *para.CT_P).
    BlockIndex int
    // RunIndex — позиция RunItem в CT_P.Content[] (среди ВСЕХ ParagraphContent,
    // не только RunItem; для получения конкретного рана используй RunIndex
    // как индекс в полном слайсе CT_P.Content).
    RunIndex  int
    Paragraph *Paragraph // wrapper
    Run       *Run       // wrapper
}

// ================================================================
// Paragraph — обёртка над para.CT_P
// ================================================================

// Paragraph хранит *para.CT_P (указатель, т.к. Body.Content хранит
// *para.CT_P). Мутации отражаются в документе.
type Paragraph struct {
    raw *para.CT_P
}

// --- Добавление ---
func (p *Paragraph) AddRun(text string) *Run
func (p *Paragraph) AddHyperlink(text, url string) *Run

// --- Вставка ---
// index — позиция в CT_P.Content[] (все ParagraphContent).
// Возвращает error при index вне диапазона.
func (p *Paragraph) InsertRunAt(index int, text string) (*Run, error)

// --- Свойства ---
func (p *Paragraph) SetStyle(styleID string)
func (p *Paragraph) SetAlignment(jc string)
func (p *Paragraph) SetSpacing(before, after units.DXA, line units.DXA, lineRule string)
func (p *Paragraph) SetIndent(left, right, firstLine units.DXA)
func (p *Paragraph) SetNumbering(numID, level int)

// --- Чтение ---
func (p *Paragraph) Runs() []*Run                     // только RunItem из Content (фильтр)
func (p *Paragraph) RunCount() int                     // len(Runs())
func (p *Paragraph) Text() string                      // конкатенация текста всех RunItem
func (p *Paragraph) Style() string                     // PPr.PStyle.Val или ""

// --- Удаление ---
// index — позиция в CT_P.Content[] (все ParagraphContent, не только RunItem).
// Возвращает error при index вне диапазона.
func (p *Paragraph) RemoveRun(index int) error
func (p *Paragraph) Clear()                            // удалить Content (PPr сохраняется)

// --- Escape hatch ---
func (p *Paragraph) X() *para.CT_P

// ================================================================
// Run — обёртка над run.CT_R
// ================================================================

// Run хранит *run.CT_R (указатель, т.к. para.RunItem хранит *run.CT_R).
type Run struct {
    raw *run.CT_R
}

// --- Контент ---
func (r *Run) SetText(text string)
func (r *Run) Text() string
func (r *Run) AddBreak(breakType string)               // "page"|"column"|"textWrapping"
func (r *Run) AddImage(imgData []byte, ext string, width, height units.EMU) error
func (r *Run) AddTab()
func (r *Run) Clear()                                  // удалить CT_R.Content (RPr сохраняется)

// --- Форматирование (character properties) ---
func (r *Run) SetBold(v bool)
func (r *Run) SetItalic(v bool)
func (r *Run) SetUnderline(style string)
func (r *Run) SetFontFamily(name string)
func (r *Run) SetFontSize(pt float64)
func (r *Run) SetColor(hex string)
func (r *Run) SetHighlight(color string)
func (r *Run) SetStrikethrough(v bool)
func (r *Run) SetSuperscript()
func (r *Run) SetSubscript()

// --- Чтение форматирования ---
func (r *Run) Bold() bool
func (r *Run) Italic() bool
func (r *Run) FontSize() float64                       // pt, 0 если не задан
func (r *Run) FontFamily() string                      // "" если не задан
func (r *Run) Color() string                           // hex без #, "" если не задан

// --- Escape hatch ---
func (r *Run) X() *run.CT_R

// ================================================================
// Table — обёртка над table.CT_Tbl
// ================================================================

// Table хранит *table.CT_Tbl (указатель, т.к. Body.Content хранит
// *table.CT_Tbl через embedded shared.BlockLevelMarker).
type Table struct {
    raw *table.CT_Tbl
}

// --- Добавление ---
func (t *Table) AddRow() *Row

// --- Вставка ---
// index — позиция среди строк в CT_Tbl.Content[] (только CT_Row, не RawTblContent).
// Возвращает error при index вне диапазона.
func (t *Table) InsertRowAt(index int) (*Row, error)

// --- Свойства ---
func (t *Table) SetStyle(styleID string)
func (t *Table) SetWidth(w int, wType string)
func (t *Table) SetBorders(style string)

// --- Чтение ---
func (t *Table) Cell(row, col int) *Cell
func (t *Table) Rows() []*Row
func (t *Table) RowCount() int
func (t *Table) ColCount() int                         // по TblGrid.GridCol или длине первой строки

// --- Удаление ---
// index — позиция строки (как в Rows(), не в raw Content).
// Возвращает error при index вне диапазона.
func (t *Table) RemoveRow(index int) error

// --- Escape hatch ---
func (t *Table) X() *table.CT_Tbl

// ================================================================
// Row — обёртка над table.CT_Row (value type в interface slice!)
// ================================================================

// ⚠ CT_Row хранится как VALUE в tbl.Content []TblContent.
// Row НЕ может хранить указатель на CT_Row — вместо этого хранит
// указатель на родительскую таблицу + индекс строки.
// При каждой операции: извлечь CT_Row из slice, мутировать, записать обратно.
// Подробности: patterns.md раздел 15.
type Row struct {
    tbl *table.CT_Tbl  // родительская таблица
    idx int            // позиция в tbl.Content (среди всех TblContent)
}

// --- Чтение ---
func (r *Row) Cells() []*Cell
func (r *Row) CellCount() int

// --- Добавление ---
func (r *Row) AddCell() *Cell

// --- Удаление ---
// index — позиция ячейки (как в Cells(), не в raw Content).
// Возвращает error при index вне диапазона.
func (r *Row) RemoveCell(index int) error

// --- Escape hatch ---
// Возвращает КОПИЮ CT_Row (value type). Мутации копии НЕ отражаются
// в документе. Для прямых мутаций: Table.X().Content[i].(table.CT_Row),
// мутация, запись обратно в Table.X().Content[i] (см. patterns.md 15).
func (r *Row) X() table.CT_Row

// ================================================================
// Cell — обёртка над table.CT_Tc (value type в interface slice!)
// ================================================================

// ⚠ CT_Tc хранится как VALUE в row.Content []RowContent.
// Cell НЕ может хранить указатель на CT_Tc — вместо этого хранит
// координаты: указатель на таблицу + индекс строки + индекс ячейки.
// Подробности: patterns.md раздел 15.
type Cell struct {
    tbl    *table.CT_Tbl  // корневая таблица
    rowIdx int            // позиция CT_Row в tbl.Content
    colIdx int            // позиция CT_Tc в row.Content
}

// --- Добавление ---
func (c *Cell) AddParagraph() *Paragraph

// --- Свойства ---
func (c *Cell) SetWidth(w int, wType string)
func (c *Cell) SetShading(fill string)
func (c *Cell) MergeHorizontal(span int)
func (c *Cell) MergeVertical(vmerge string)

// --- Чтение ---
func (c *Cell) Paragraphs() []*Paragraph
func (c *Cell) Tables() []*Table                       // вложенные таблицы

// --- Удаление ---
// Clear удаляет весь Content ячейки и вставляет пустой *para.CT_P{}.
// Инвариант OOXML: каждая tc содержит ≥1 <w:p> (reference-appendix 5.6).
// TcPr сохраняется.
func (c *Cell) Clear()

// --- Escape hatch ---
// Возвращает КОПИЮ CT_Tc (value type). Мутации копии НЕ отражаются
// в документе. Для прямых мутаций: извлечь row из Table.X().Content,
// извлечь tc из row.Content, мутировать, записать обратно оба.
func (c *Cell) X() table.CT_Tc

// ================================================================
// Section — обёртка над sectpr.CT_SectPr
// ================================================================

type Section struct { /* wraps *sectpr.CT_SectPr */ }
func (s *Section) SetPageSize(w, h units.DXA)
func (s *Section) SetLandscape()
func (s *Section) SetPortrait()
func (s *Section) SetMargins(top, right, bottom, left units.DXA)
func (s *Section) SetColumns(num int, space units.DXA)
func (s *Section) AddHeader(hdrType string) *Header
func (s *Section) AddFooter(ftrType string) *Footer

// --- Escape hatch ---
func (s *Section) X() *sectpr.CT_SectPr

// ================================================================
// Header / Footer — обёртки над hdft.CT_HdrFtr
// ================================================================

type Header struct { /* wraps *hdft.CT_HdrFtr */ }
func (h *Header) AddParagraph() *Paragraph
func (h *Header) Paragraphs() []*Paragraph
func (h *Header) Clear()                               // удалить Content, вставить *para.CT_P{}
func (h *Header) X() *hdft.CT_HdrFtr

type Footer struct { /* wraps *hdft.CT_HdrFtr */ }
func (f *Footer) AddParagraph() *Paragraph
func (f *Footer) Paragraphs() []*Paragraph
func (f *Footer) Clear()                               // удалить Content, вставить *para.CT_P{}
func (f *Footer) X() *hdft.CT_HdrFtr
```

---

## Промпт для начала реализации модуля

При запуске нового чата для реализации модуля X, используйте этот промпт:

````
Реализуй модуль [MODULE_NAME] библиотеки docx-go на Go.

Прикреплённые файлы (загрузи все 3):
1. contracts.md           — ЧТО реализовать (типы, сигнатуры, зависимости)
2. reference-appendix.md  — XML-примеры, namespaces, порядок элементов, ловушки
3. patterns.md            — КАК реализовать (MarshalXML, RawXML, naming, shared)

Твоя задача:
- Реализовать ТОЛЬКО модуль [MODULE_NAME]
- Контракт модуля: секция C-XX в contracts.md
- Зависимости: импортируй ТОЛЬКО перечисленные пакеты
- Не меняй контракт (публичные типы и сигнатуры)
- Маппинг Go field → XML element: см. patterns.md раздел 2
- Кастомный MarshalXML для xsd:sequence: см. patterns.md раздел 3
- RawXML round-trip для неизвестных элементов: см. patterns.md раздел 4
- Напиши round-trip тест: unmarshal XML из appendix → marshal → сравнить

Структура вывода:
- pkg/[module_path]/*.go — исходный код
- pkg/[module_path]/*_test.go — тесты
- Один файл = один логический блок
````

### Дополнение для C-32 (`docx`)

Модуль `docx` — особый: это тонкая обёртка, а не XML-маршализатор.
Добавь к промпту:

````
Дополнительно для C-32 (`docx`):
- Прочитай patterns.md разделы 14-17 (X(), value-type gotcha, FindText, Remove)
- Прочитай reference-appendix.md раздел 5.15 (инварианты при редактировании)
- Каждый wrapper хранит указатель на WML-тип, НЕ копирует данные
- X() возвращает указатель (кроме Row — см. patterns.md раздел 15)
- Cell.Clear(), Header.Clear(), Footer.Clear() — вставляют пустой <w:p/>
- Все Remove*/InsertAt* — index-based, документируй инвалидацию индексов
- Write-back при мутации CT_Row/CT_Tc (patterns.md раздел 15)
- Тесты: Add→Read→Remove→Read→Validate для каждого уровня (Body/Para/Run/Table)
````