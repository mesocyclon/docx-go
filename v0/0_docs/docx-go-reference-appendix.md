# DOCX Reference Appendix — Справочник для реализации

Этот документ содержит всё, что нужно для написания кода без обращения к внешним источникам.

---

# 1. Namespaces — Полный справочник

## 1.1 Transitional namespaces (99.9% реальных документов)

```go
const (
    // === WordprocessingML (main) ===
    NSw  = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

    // === Relationships ===
    NSr  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"

    // === DrawingML ===
    NSa  = "http://schemas.openxmlformats.org/drawingml/2006/main"
    NSwp = "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
    NSpic = "http://schemas.openxmlformats.org/drawingml/2006/picture"

    // === Math ===
    NSm  = "http://schemas.openxmlformats.org/officeDocument/2006/math"

    // === VML (legacy) ===
    NSv  = "urn:schemas-microsoft-com:vml"
    NSo  = "urn:schemas-microsoft-com:office:office"
    NSw10 = "urn:schemas-microsoft-com:office:word"

    // === Markup Compatibility ===
    NSmc = "http://schemas.openxmlformats.org/markup-compatibility/2006"

    // === Package ===
    NSContentTypes = "http://schemas.openxmlformats.org/package/2006/content-types"
    NSRelationships = "http://schemas.openxmlformats.org/package/2006/relationships"

    // === Document Properties ===
    NScp      = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
    NSdc      = "http://purl.org/dc/elements/1.1/"
    NSdcterms = "http://purl.org/dc/terms/"
    NSdcmitype = "http://purl.org/dc/dcmitype/"
    NSextProps = "http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"
    NSvt      = "http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"
    NScustomProps = "http://schemas.openxmlformats.org/officeDocument/2006/custom-properties"

    // === Schema Library ===
    NSsl = "http://schemas.openxmlformats.org/schemaLibrary/2006/main"

    // === Microsoft Extensions (Word 2010+) ===
    NSw14  = "http://schemas.microsoft.com/office/word/2010/wordml"
    NSwp14 = "http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"
    NSwpc  = "http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
    NSwpg  = "http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"
    NSwpi  = "http://schemas.microsoft.com/office/word/2010/wordprocessingInk"
    NSwps  = "http://schemas.microsoft.com/office/word/2010/wordprocessingShape"

    // === Microsoft Extensions (Word 2013+) ===
    NSw15 = "http://schemas.microsoft.com/office/word/2012/wordml"

    // === Microsoft Extensions (Word 2016+) ===
    NSw16se  = "http://schemas.microsoft.com/office/word/2015/wordml/symex"
    NSw16cid = "http://schemas.microsoft.com/office/word/2016/wordml/cid"
    NSw16    = "http://schemas.microsoft.com/office/word/2018/wordml"
    NSw16sdtdh = "http://schemas.microsoft.com/office/word/2020/wordml/sdtdatahash"

    // === Word 2006 Extensions ===
    NSwne = "http://schemas.microsoft.com/office/word/2006/wordml"

    // === XSI ===
    NSxsi = "http://www.w3.org/2001/XMLSchema-instance"
)
```

## 1.2 Strict namespaces (ECMA-376 Strict, ISO/IEC 29500 Strict)

```go
const (
    NSw_Strict  = "http://purl.oclc.org/ooxml/wordprocessingml/main"
    NSr_Strict  = "http://purl.oclc.org/ooxml/officeDocument/relationships"
    NSa_Strict  = "http://purl.oclc.org/ooxml/drawingml/main"
    NSwp_Strict = "http://purl.oclc.org/ooxml/drawingml/wordprocessingDrawing"
    NSm_Strict  = "http://purl.oclc.org/ooxml/officeDocument/math"
    NSsl_Strict = "http://purl.oclc.org/ooxml/schemaLibrary/main"
    NSpic_Strict = "http://purl.oclc.org/ooxml/drawingml/picture"
)
```

**Критически важно**: при unmarshal принимать ОБА варианта namespace. При marshal — писать Transitional (для совместимости).

## 1.3 Relationship Types — полный список

```go
const (
    // === Package-level ===
    RelOfficeDocument  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
    RelCoreProperties  = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties"
    RelExtProperties   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties"
    RelCustomProperties = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/custom-properties"
    RelThumbnail       = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"

    // === Part-level (от document.xml) ===
    RelStyles      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
    RelSettings    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings"
    RelFontTable   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/fontTable"
    RelNumbering   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering"
    RelFootnotes   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes"
    RelEndnotes    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/endnotes"
    RelComments    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments"
    RelHeader      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/header"
    RelFooter      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer"
    RelImage       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
    RelHyperlink   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
    RelTheme       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
    RelWebSettings = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/webSettings"
    RelGlossary    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/glossaryDocument"
    RelCustomXml   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml"
    RelChart       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
    RelOleObject   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/oleObject"
    RelPackage     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package"

    // === Microsoft Extensions ===
    RelCommentsExtended = "http://schemas.microsoft.com/office/2011/relationships/commentsExtended"
    RelCommentsIds      = "http://schemas.microsoft.com/office/2016/09/relationships/commentsIds"
    RelPeople           = "http://schemas.microsoft.com/office/2011/relationships/people"
)
```

## 1.4 Content Types — полный список

```go
const (
    CTDocument   = "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"
    CTStyles     = "application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"
    CTSettings   = "application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"
    CTFontTable  = "application/vnd.openxmlformats-officedocument.wordprocessingml.fontTable+xml"
    CTNumbering  = "application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"
    CTFootnotes  = "application/vnd.openxmlformats-officedocument.wordprocessingml.footnotes+xml"
    CTEndnotes   = "application/vnd.openxmlformats-officedocument.wordprocessingml.endnotes+xml"
    CTComments   = "application/vnd.openxmlformats-officedocument.wordprocessingml.comments+xml"
    CTHeader     = "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"
    CTFooter     = "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"
    CTWebSettings = "application/vnd.openxmlformats-officedocument.wordprocessingml.webSettings+xml"
    CTTheme      = "application/vnd.openxmlformats-officedocument.theme+xml"
    CTCore       = "application/vnd.openxmlformats-package.core-properties+xml"
    CTExtended   = "application/vnd.openxmlformats-officedocument.extended-properties+xml"
    CTCustom     = "application/vnd.openxmlformats-officedocument.custom-properties+xml"
    CTRels       = "application/vnd.openxmlformats-package.relationships+xml"
    CTGlossary   = "application/vnd.openxmlformats-officedocument.wordprocessingml.document.glossary+xml"
    CTCommentsExt = "application/vnd.openxmlformats-officedocument.wordprocessingml.commentsExtended+xml"
    CTCommentsIds = "application/vnd.openxmlformats-officedocument.wordprocessingml.commentsIds+xml"
    CTPeople     = "application/vnd.openxmlformats-officedocument.wordprocessingml.people+xml"

    // Media defaults (по расширению)
    CTPng  = "image/png"
    CTJpeg = "image/jpeg"
    CTGif  = "image/gif"
    CTBmp  = "image/bmp"
    CTTiff = "image/tiff"
    CTSvg  = "image/svg+xml"
    CTEmf  = "image/x-emf"
    CTWmf  = "image/x-wmf"
)
```

---

# 2. Минимальный валидный .docx — Скелет для `docx.New()`

## 2.1 Структура файлов

> **ПРОВЕРЕНО**: этот набор файлов открывается в MS Word без ошибок и предупреждений.
> Версия без theme/webSettings/clrSchemeMapping/shapeDefaults/rsids вызывает
> ошибку «содержимое не удалось прочитать» и принудительное восстановление.

```
minimal.docx (ZIP)
├── [Content_Types].xml            ← ОБЯЗАТЕЛЬНО
├── _rels/
│   └── .rels                      ← ОБЯЗАТЕЛЬНО
├── word/
│   ├── _rels/
│   │   └── document.xml.rels      ← ОБЯЗАТЕЛЬНО
│   ├── document.xml               ← ОБЯЗАТЕЛЬНО
│   ├── styles.xml                 ← ОБЯЗАТЕЛЬНО (без него Word пересоздаёт)
│   ├── settings.xml               ← ОБЯЗАТЕЛЬНО (без clrSchemeMapping → ошибка)
│   ├── webSettings.xml            ← ОБЯЗАТЕЛЬНО (может быть пустым, но должен быть)
│   ├── fontTable.xml              ← ОБЯЗАТЕЛЬНО (без него Word пересоздаёт)
│   └── theme/
│       └── theme1.xml             ← ОБЯЗАТЕЛЬНО (стили ссылаются на themeColor/themeFonts)
└── docProps/
    ├── core.xml                   ← рекомендуется
    └── app.xml                    ← рекомендуется
```

## 2.2 `[Content_Types].xml`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
  <Override PartName="/word/settings.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>
  <Override PartName="/word/webSettings.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.webSettings+xml"/>
  <Override PartName="/word/fontTable.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.fontTable+xml"/>
  <Override PartName="/word/theme/theme1.xml"
    ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
  <Override PartName="/docProps/core.xml"
    ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
  <Override PartName="/docProps/app.xml"
    ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
</Types>
```

## 2.3 `_rels/.rels`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>
```

## 2.4 `word/_rels/document.xml.rels`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/webSettings" Target="webSettings.xml"/>
  <Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/fontTable" Target="fontTable.xml"/>
  <Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>
</Relationships>
```

## 2.5 `word/document.xml`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
            xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
            xmlns:o="urn:schemas-microsoft-com:office:office"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
            xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
            xmlns:v="urn:schemas-microsoft-com:vml"
            xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"
            xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
            xmlns:w10="urn:schemas-microsoft-com:office:word"
            xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
            xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
            xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"
            xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk"
            xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml"
            xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape"
            mc:Ignorable="w14 w15 wp14">
  <w:body>
    <w:p w14:paraId="00000001" w14:textId="77777777"
         w:rsidR="00000001" w:rsidRDefault="00000001">
      <w:r>
        <w:t>Hello World</w:t>
      </w:r>
    </w:p>
    <w:sectPr w:rsidR="00000001">
      <w:pgSz w:w="12240" w:h="15840"/>
      <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"
               w:header="720" w:footer="720" w:gutter="0"/>
      <w:cols w:space="720"/>
      <w:docGrid w:linePitch="360"/>
    </w:sectPr>
  </w:body>
</w:document>
```

## 2.6 `word/styles.xml` — Минимальный

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
          xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
          mc:Ignorable="w14 w15">
  <w:docDefaults>
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:asciiTheme="minorHAnsi" w:eastAsiaTheme="minorHAnsi"
                  w:hAnsiTheme="minorHAnsi" w:cstheme="minorBidi"/>
        <w:sz w:val="24"/>
        <w:szCs w:val="24"/>
        <w:lang w:val="en-US" w:eastAsia="en-US" w:bidi="ar-SA"/>
      </w:rPr>
    </w:rPrDefault>
    <w:pPrDefault>
      <w:pPr>
        <w:spacing w:after="160" w:line="259" w:lineRule="auto"/>
      </w:pPr>
    </w:pPrDefault>
  </w:docDefaults>
  <w:latentStyles w:defLockedState="0" w:defUIPriority="99"
    w:defSemiHidden="0" w:defUnhideWhenUsed="0" w:defQFormat="0" w:count="376"/>
  <!-- Четыре обязательных default-стиля -->
  <w:style w:type="paragraph" w:default="1" w:styleId="Normal">
    <w:name w:val="Normal"/>
    <w:qFormat/>
  </w:style>
  <w:style w:type="character" w:default="1" w:styleId="DefaultParagraphFont">
    <w:name w:val="Default Paragraph Font"/>
    <w:uiPriority w:val="1"/>
    <w:semiHidden/>
    <w:unhideWhenUsed/>
  </w:style>
  <w:style w:type="table" w:default="1" w:styleId="TableNormal">
    <w:name w:val="Normal Table"/>
    <w:uiPriority w:val="99"/>
    <w:semiHidden/>
    <w:unhideWhenUsed/>
    <w:tblPr>
      <w:tblInd w:w="0" w:type="dxa"/>
      <w:tblCellMar>
        <w:top w:w="0" w:type="dxa"/>
        <w:left w:w="108" w:type="dxa"/>
        <w:bottom w:w="0" w:type="dxa"/>
        <w:right w:w="108" w:type="dxa"/>
      </w:tblCellMar>
    </w:tblPr>
  </w:style>
  <w:style w:type="numbering" w:default="1" w:styleId="NoList">
    <w:name w:val="No List"/>
    <w:uiPriority w:val="99"/>
    <w:semiHidden/>
    <w:unhideWhenUsed/>
  </w:style>
</w:styles>
```

## 2.7 `word/settings.xml` — Минимальный

> **ВАЖНО**: settings.xml без `clrSchemeMapping`, `shapeDefaults` и `rsids`
> вызывает ошибку «содержимое не удалось прочитать» в MS Word.

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:settings xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
            xmlns:o="urn:schemas-microsoft-com:office:office"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
            xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
            xmlns:v="urn:schemas-microsoft-com:vml"
            xmlns:w10="urn:schemas-microsoft-com:office:word"
            xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
            xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
            xmlns:sl="http://schemas.openxmlformats.org/schemaLibrary/2006/main"
            mc:Ignorable="w14 w15">
  <w:zoom w:percent="100"/>
  <w:proofState w:spelling="clean" w:grammar="clean"/>
  <w:defaultTabStop w:val="720"/>
  <w:characterSpacingControl w:val="doNotCompress"/>
  <w:compat>
    <w:compatSetting w:name="compatibilityMode"
      w:uri="http://schemas.microsoft.com/office/word" w:val="15"/>
    <w:compatSetting w:name="overrideTableStyleFontSizeAndJustification"
      w:uri="http://schemas.microsoft.com/office/word" w:val="1"/>
    <w:compatSetting w:name="enableOpenTypeFeatures"
      w:uri="http://schemas.microsoft.com/office/word" w:val="1"/>
    <w:compatSetting w:name="doNotFlipMirrorIndents"
      w:uri="http://schemas.microsoft.com/office/word" w:val="1"/>
    <w:compatSetting w:name="differentiateMultirowTableHeaders"
      w:uri="http://schemas.microsoft.com/office/word" w:val="1"/>
    <w:compatSetting w:name="useWord2013TrackBottomHyphenation"
      w:uri="http://schemas.microsoft.com/office/word" w:val="0"/>
  </w:compat>
  <!-- ОБЯЗАТЕЛЬНО: rsids — без них Word показывает ошибку -->
  <w:rsids>
    <w:rsidRoot w:val="00000001"/>
    <w:rsid w:val="00000001"/>
  </w:rsids>
  <m:mathPr>
    <m:mathFont m:val="Cambria Math"/>
    <m:brkBin m:val="before"/>
    <m:brkBinSub m:val="--"/>
    <m:smallFrac m:val="0"/>
    <m:dispDef/>
    <m:lMargin m:val="0"/>
    <m:rMargin m:val="0"/>
    <m:defJc m:val="centerGroup"/>
    <m:wrapIndent m:val="1440"/>
    <m:intLim m:val="subSup"/>
    <m:naryLim m:val="undOvr"/>
  </m:mathPr>
  <w:themeFontLang w:val="en-US"/>
  <!-- ОБЯЗАТЕЛЬНО: clrSchemeMapping — маппинг цветовой схемы из theme1.xml -->
  <w:clrSchemeMapping w:bg1="light1" w:t1="dark1" w:bg2="light2" w:t2="dark2"
    w:accent1="accent1" w:accent2="accent2" w:accent3="accent3" w:accent4="accent4"
    w:accent5="accent5" w:accent6="accent6" w:hyperlink="hyperlink"
    w:followedHyperlink="followedHyperlink"/>
  <!-- ОБЯЗАТЕЛЬНО: shapeDefaults — значения по умолчанию для VML-фигур -->
  <w:shapeDefaults>
    <o:shapedefaults v:ext="edit" spidmax="1026"/>
    <o:shapelayout v:ext="edit">
      <o:idmap v:ext="edit" data="1"/>
    </o:shapelayout>
  </w:shapeDefaults>
  <w:decimalSymbol w:val="."/>
  <w:listSeparator w:val=","/>
  <w14:docId w14:val="00000001"/>
  <w15:chartTrackingRefBased/>
  <w15:docId w15:val="{00000000-0000-0000-0000-000000000001}"/>
</w:settings>
```

## 2.8 `word/webSettings.xml` — Минимальный

> Может быть пустым (самозакрывающийся тег), но part и relationship обязательны.

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:webSettings xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
               xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
               xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
               xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
               xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
               mc:Ignorable="w14 w15"/>
```

## 2.9 `word/theme/theme1.xml` — Минимальный

> Тема ОБЯЗАТЕЛЬНА, потому что styles.xml ссылается на `asciiTheme="minorHAnsi"`,
> `w:themeColor="accent1"` и т.д. Без темы эти ссылки не могут разрешиться.
>
> Полная тема — ~8 КБ XML (цветовая схема, шрифты, форматирование).
> Ниже — минимальная тема, которую принимает Word:

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme">
  <a:themeElements>
    <a:clrScheme name="Office">
      <a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1>
      <a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1>
      <a:dk2><a:srgbClr val="44546A"/></a:dk2>
      <a:lt2><a:srgbClr val="E7E6E6"/></a:lt2>
      <a:accent1><a:srgbClr val="4472C4"/></a:accent1>
      <a:accent2><a:srgbClr val="ED7D31"/></a:accent2>
      <a:accent3><a:srgbClr val="A5A5A5"/></a:accent3>
      <a:accent4><a:srgbClr val="FFC000"/></a:accent4>
      <a:accent5><a:srgbClr val="5B9BD5"/></a:accent5>
      <a:accent6><a:srgbClr val="70AD47"/></a:accent6>
      <a:hlink><a:srgbClr val="0563C1"/></a:hlink>
      <a:folHlink><a:srgbClr val="954F72"/></a:folHlink>
    </a:clrScheme>
    <a:fontScheme name="Office">
      <a:majorFont>
        <a:latin typeface="Calibri Light" panose="020F0302020204030204"/>
        <a:ea typeface=""/>
        <a:cs typeface=""/>
      </a:majorFont>
      <a:minorFont>
        <a:latin typeface="Calibri" panose="020F0502020204030204"/>
        <a:ea typeface=""/>
        <a:cs typeface=""/>
      </a:minorFont>
    </a:fontScheme>
    <a:fmtScheme name="Office">
      <a:fillStyleLst>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
      </a:fillStyleLst>
      <a:lnStyleLst>
        <a:ln w="6350"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
        <a:ln w="12700"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
        <a:ln w="19050"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
      </a:lnStyleLst>
      <a:effectStyleLst>
        <a:effectStyle><a:effectLst/></a:effectStyle>
        <a:effectStyle><a:effectLst/></a:effectStyle>
        <a:effectStyle><a:effectLst/></a:effectStyle>
      </a:effectStyleLst>
      <a:bgFillStyleLst>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
      </a:bgFillStyleLst>
    </a:fmtScheme>
  </a:themeElements>
  <a:objectDefaults/>
  <a:extraClrSchemeLst/>
</a:theme>
```

## 2.10 `word/fontTable.xml` — Минимальный

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:fonts xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:font w:name="Calibri">
    <w:panose1 w:val="020F0502020204030204"/>
    <w:charset w:val="00"/>
    <w:family w:val="swiss"/>
    <w:pitch w:val="variable"/>
  </w:font>
  <w:font w:name="Times New Roman">
    <w:panose1 w:val="02020603050405020304"/>
    <w:charset w:val="00"/>
    <w:family w:val="roman"/>
    <w:pitch w:val="variable"/>
  </w:font>
  <w:font w:name="Calibri Light">
    <w:panose1 w:val="020F0302020204030204"/>
    <w:charset w:val="00"/>
    <w:family w:val="swiss"/>
    <w:pitch w:val="variable"/>
  </w:font>
</w:fonts>
```

## 2.11 `docProps/core.xml`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
                   xmlns:dc="http://purl.org/dc/elements/1.1/"
                   xmlns:dcterms="http://purl.org/dc/terms/"
                   xmlns:dcmitype="http://purl.org/dc/dcmitype/"
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <dc:title/>
  <dc:subject/>
  <dc:creator>Author</dc:creator>
  <cp:keywords/>
  <dc:description/>
  <cp:lastModifiedBy>Author</cp:lastModifiedBy>
  <cp:revision>1</cp:revision>
  <dcterms:created xsi:type="dcterms:W3CDTF">2025-01-01T00:00:00Z</dcterms:created>
  <dcterms:modified xsi:type="dcterms:W3CDTF">2025-01-01T00:00:00Z</dcterms:modified>
</cp:coreProperties>
```

## 2.12 `docProps/app.xml`

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"
            xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
  <Template>Normal</Template>
  <TotalTime>0</TotalTime>
  <Pages>1</Pages>
  <Words>0</Words>
  <Characters>0</Characters>
  <Application>docx-go</Application>
  <DocSecurity>0</DocSecurity>
  <Lines>1</Lines>
  <Paragraphs>1</Paragraphs>
  <ScaleCrop>false</ScaleCrop>
  <Company/>
  <LinksUpToDate>false</LinksUpToDate>
  <CharactersWithSpaces>0</CharactersWithSpaces>
  <SharedDoc>false</SharedDoc>
  <HyperlinksChanged>false</HyperlinksChanged>
  <AppVersion>16.0000</AppVersion>
</Properties>
```

---

# 3. Примеры реального XML

## 3.1 Форматированный параграф

```xml
<!-- Параграф: Heading 1, выравнивание по центру -->
<w:p w:rsidR="00A77B3E" w:rsidRDefault="00A77B3E" w:rsidP="00A77B3E">
  <w:pPr>
    <w:pStyle w:val="Heading1"/>
    <w:jc w:val="center"/>
  </w:pPr>
  <w:r w:rsidRPr="00C83215">
    <w:rPr>
      <w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:cs="Arial"/>
      <w:b/>
      <w:color w:val="2F5496" w:themeColor="accent1" w:themeShade="BF"/>
      <w:sz w:val="32"/>
      <w:szCs w:val="32"/>
    </w:rPr>
    <w:t>Document Title</w:t>
  </w:r>
</w:p>

<!-- Параграф: обычный текст с пробелами, жирный + курсив -->
<w:p w:rsidR="00B22C47" w:rsidRDefault="00B22C47">
  <w:r>
    <w:t xml:space="preserve">This is </w:t>
  </w:r>
  <w:r w:rsidRPr="00D714A3">
    <w:rPr>
      <w:b/>
      <w:bCs/>
    </w:rPr>
    <w:t>bold</w:t>
  </w:r>
  <w:r>
    <w:t xml:space="preserve"> and </w:t>
  </w:r>
  <w:r>
    <w:rPr>
      <w:i/>
      <w:iCs/>
    </w:rPr>
    <w:t>italic</w:t>
  </w:r>
  <w:r>
    <w:t xml:space="preserve"> text.</w:t>
  </w:r>
</w:p>
```

## 3.2 Таблица 2×2

```xml
<w:tbl>
  <w:tblPr>
    <w:tblStyle w:val="TableGrid"/>
    <w:tblW w:w="0" w:type="auto"/>
    <w:tblLook w:firstRow="1" w:lastRow="0" w:firstColumn="1"
               w:lastColumn="0" w:noHBand="0" w:noVBand="1"/>
  </w:tblPr>
  <w:tblGrid>
    <w:gridCol w:w="4675"/>
    <w:gridCol w:w="4675"/>
  </w:tblGrid>
  <!-- Row 1 (header) -->
  <w:tr w:rsidR="009A2C41" w:rsidTr="009A2C41">
    <w:tc>
      <w:tcPr>
        <w:tcW w:w="4675" w:type="dxa"/>
        <w:shd w:val="clear" w:color="auto" w:fill="D9E2F3" w:themeFill="accent1" w:themeFillTint="33"/>
      </w:tcPr>
      <w:p>
        <w:pPr><w:jc w:val="center"/></w:pPr>
        <w:r>
          <w:rPr><w:b/></w:rPr>
          <w:t>Header 1</w:t>
        </w:r>
      </w:p>
    </w:tc>
    <w:tc>
      <w:tcPr>
        <w:tcW w:w="4675" w:type="dxa"/>
        <w:shd w:val="clear" w:color="auto" w:fill="D9E2F3" w:themeFill="accent1" w:themeFillTint="33"/>
      </w:tcPr>
      <w:p>
        <w:pPr><w:jc w:val="center"/></w:pPr>
        <w:r>
          <w:rPr><w:b/></w:rPr>
          <w:t>Header 2</w:t>
        </w:r>
      </w:p>
    </w:tc>
  </w:tr>
  <!-- Row 2 -->
  <w:tr w:rsidR="009A2C41" w:rsidTr="009A2C41">
    <w:tc>
      <w:tcPr>
        <w:tcW w:w="4675" w:type="dxa"/>
      </w:tcPr>
      <w:p>
        <w:r><w:t>Cell A</w:t></w:r>
      </w:p>
    </w:tc>
    <w:tc>
      <w:tcPr>
        <w:tcW w:w="4675" w:type="dxa"/>
      </w:tcPr>
      <w:p>
        <w:r><w:t>Cell B</w:t></w:r>
      </w:p>
    </w:tc>
  </w:tr>
</w:tbl>
```

## 3.3 Нумерованный список

```xml
<!-- numbering.xml (фрагмент) -->
<w:abstractNum w:abstractNumId="0">
  <w:nsid w:val="3A5C117E"/>
  <w:multiLevelType w:val="hybridMultilevel"/>
  <w:tmpl w:val="E6A2FD28"/>
  <w:lvl w:ilvl="0" w:tplc="04090001">
    <w:start w:val="1"/>
    <w:numFmt w:val="bullet"/>
    <w:lvlText w:val="&#xF0B7;"/>
    <w:lvlJc w:val="left"/>
    <w:pPr>
      <w:ind w:left="720" w:hanging="360"/>
    </w:pPr>
    <w:rPr>
      <w:rFonts w:ascii="Symbol" w:hAnsi="Symbol" w:hint="default"/>
    </w:rPr>
  </w:lvl>
  <w:lvl w:ilvl="1" w:tplc="04090003">
    <w:start w:val="1"/>
    <w:numFmt w:val="bullet"/>
    <w:lvlText w:val="o"/>
    <w:lvlJc w:val="left"/>
    <w:pPr>
      <w:ind w:left="1440" w:hanging="360"/>
    </w:pPr>
    <w:rPr>
      <w:rFonts w:ascii="Courier New" w:hAnsi="Courier New" w:cs="Courier New" w:hint="default"/>
    </w:rPr>
  </w:lvl>
</w:abstractNum>
<w:num w:numId="1">
  <w:abstractNumId w:val="0"/>
</w:num>

<!-- document.xml: ссылка на список -->
<w:p>
  <w:pPr>
    <w:pStyle w:val="ListParagraph"/>
    <w:numPr>
      <w:ilvl w:val="0"/>
      <w:numId w:val="1"/>
    </w:numPr>
  </w:pPr>
  <w:r>
    <w:t>First bullet item</w:t>
  </w:r>
</w:p>
<w:p>
  <w:pPr>
    <w:pStyle w:val="ListParagraph"/>
    <w:numPr>
      <w:ilvl w:val="1"/>
      <w:numId w:val="1"/>
    </w:numPr>
  </w:pPr>
  <w:r>
    <w:t>Nested sub-item</w:t>
  </w:r>
</w:p>
```

## 3.4 Track Changes — вставка и удаление

```xml
<w:p>
  <w:r>
    <w:t xml:space="preserve">The contract term is </w:t>
  </w:r>
  <w:del w:id="1" w:author="Jane Smith" w:date="2025-01-15T10:30:00Z">
    <w:r w:rsidDel="00F12AB3">
      <w:rPr>
        <w:b/>
      </w:rPr>
      <w:delText>30</w:delText>
    </w:r>
  </w:del>
  <w:ins w:id="2" w:author="Jane Smith" w:date="2025-01-15T10:30:00Z">
    <w:r w:rsidR="00F12AB3">
      <w:rPr>
        <w:b/>
      </w:rPr>
      <w:t>60</w:t>
    </w:r>
  </w:ins>
  <w:r>
    <w:t xml:space="preserve"> days.</w:t>
  </w:r>
</w:p>
```

## 3.5 Комментарий

```xml
<!-- document.xml -->
<w:p>
  <w:commentRangeStart w:id="0"/>
  <w:r>
    <w:t>This text has a comment.</w:t>
  </w:r>
  <w:commentRangeEnd w:id="0"/>
  <w:r>
    <w:rPr>
      <w:rStyle w:val="CommentReference"/>
    </w:rPr>
    <w:commentReference w:id="0"/>
  </w:r>
</w:p>

<!-- comments.xml -->
<w:comments xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:comment w:id="0" w:author="Reviewer" w:date="2025-01-20T14:00:00Z" w:initials="R">
    <w:p>
      <w:pPr>
        <w:pStyle w:val="CommentText"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:rStyle w:val="CommentReference"/>
        </w:rPr>
        <w:annotationRef/>
      </w:r>
      <w:r>
        <w:t>Please verify this statement.</w:t>
      </w:r>
    </w:p>
  </w:comment>
</w:comments>
```

## 3.6 Inline-изображение

```xml
<!-- document.xml.rels -->
<Relationship Id="rId10" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.png"/>

<!-- [Content_Types].xml -->
<Default Extension="png" ContentType="image/png"/>

<!-- document.xml -->
<w:r>
  <w:rPr>
    <w:noProof/>
  </w:rPr>
  <w:drawing>
    <wp:inline distT="0" distB="0" distL="0" distR="0">
      <wp:extent cx="1828800" cy="1371600"/>
      <wp:effectExtent l="0" t="0" r="0" b="0"/>
      <wp:docPr id="1" name="Picture 1" descr="Logo"/>
      <wp:cNvGraphicFramePr>
        <a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/>
      </wp:cNvGraphicFramePr>
      <a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
        <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
          <pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
            <pic:nvPicPr>
              <pic:cNvPr id="1" name="image1.png"/>
              <pic:cNvPicPr/>
            </pic:nvPicPr>
            <pic:blipFill>
              <a:blip r:embed="rId10"/>
              <a:stretch>
                <a:fillRect/>
              </a:stretch>
            </pic:blipFill>
            <pic:spPr>
              <a:xfrm>
                <a:off x="0" y="0"/>
                <a:ext cx="1828800" cy="1371600"/>
              </a:xfrm>
              <a:prstGeom prst="rect">
                <a:avLst/>
              </a:prstGeom>
            </pic:spPr>
          </pic:pic>
        </a:graphicData>
      </a:graphic>
    </wp:inline>
  </w:drawing>
</w:r>
```

## 3.7 Колонтитулы

```xml
<!-- document.xml: sectPr с колонтитулами -->
<w:sectPr>
  <w:headerReference w:type="default" r:id="rId8"/>
  <w:headerReference w:type="first" r:id="rId9"/>
  <w:footerReference w:type="default" r:id="rId10"/>
  <w:pgSz w:w="12240" w:h="15840"/>
  <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"
           w:header="720" w:footer="720" w:gutter="0"/>
  <w:titlePg/>   <!-- включает отдельную первую страницу -->
</w:sectPr>

<!-- header1.xml (default) -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <w:p>
    <w:pPr>
      <w:pStyle w:val="Header"/>
      <w:jc w:val="right"/>
    </w:pPr>
    <w:r>
      <w:t xml:space="preserve">Page </w:t>
    </w:r>
    <w:r>
      <w:fldChar w:fldCharType="begin"/>
    </w:r>
    <w:r>
      <w:instrText xml:space="preserve"> PAGE </w:instrText>
    </w:r>
    <w:r>
      <w:fldChar w:fldCharType="separate"/>
    </w:r>
    <w:r>
      <w:t>1</w:t>
    </w:r>
    <w:r>
      <w:fldChar w:fldCharType="end"/>
    </w:r>
  </w:p>
</w:hdr>

<!-- footer1.xml (default) -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:p>
    <w:pPr>
      <w:pStyle w:val="Footer"/>
      <w:jc w:val="center"/>
    </w:pPr>
    <w:r>
      <w:t>Confidential</w:t>
    </w:r>
  </w:p>
</w:ftr>
```

## 3.8 Сноски

```xml
<!-- footnotes.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:footnotes xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <!-- ОБЯЗАТЕЛЬНО: separator (id=0) -->
  <w:footnote w:type="separator" w:id="-1">
    <w:p>
      <w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr>
      <w:r><w:separator/></w:r>
    </w:p>
  </w:footnote>
  <!-- ОБЯЗАТЕЛЬНО: continuationSeparator (id=1) -->
  <w:footnote w:type="continuationSeparator" w:id="0">
    <w:p>
      <w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr>
      <w:r><w:continuationSeparator/></w:r>
    </w:p>
  </w:footnote>
  <!-- Пользовательская сноска -->
  <w:footnote w:id="1">
    <w:p>
      <w:pPr><w:pStyle w:val="FootnoteText"/></w:pPr>
      <w:r>
        <w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr>
        <w:footnoteRef/>
      </w:r>
      <w:r>
        <w:t xml:space="preserve"> See the original source for details.</w:t>
      </w:r>
    </w:p>
  </w:footnote>
</w:footnotes>

<!-- document.xml: ссылка на сноску -->
<w:r>
  <w:rPr>
    <w:rStyle w:val="FootnoteReference"/>
  </w:rPr>
  <w:footnoteReference w:id="1"/>
</w:r>
```

## 3.9 Разрыв секции (landscape + 2 колонки)

```xml
<!-- Последний параграф перед разрывом секции -->
<w:p>
  <w:pPr>
    <w:sectPr>
      <!-- Свойства ПРЕДЫДУЩЕЙ секции -->
      <w:pgSz w:w="12240" w:h="15840"/>
      <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"
               w:header="720" w:footer="720" w:gutter="0"/>
    </w:sectPr>
  </w:pPr>
  <w:r><w:t>End of portrait section.</w:t></w:r>
</w:p>
<!-- Следующая секция — landscape, 2 колонки -->
<w:p>
  <w:r><w:t>This is in landscape with two columns.</w:t></w:r>
</w:p>
<!-- sectPr в body — для ПОСЛЕДНЕЙ секции -->
<w:sectPr>
  <w:type w:val="continuous"/>
  <w:pgSz w:w="15840" w:h="12240" w:orient="landscape"/>
  <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"
           w:header="720" w:footer="720" w:gutter="0"/>
  <w:cols w:num="2" w:space="720"/>
</w:sectPr>
```

## 3.10 mc:AlternateContent (Markup Compatibility)

```xml
<!-- Word вставляет AlternateContent для фич, которые не все потребители поддерживают -->
<w:r>
  <mc:AlternateContent>
    <mc:Choice Requires="wps">
      <!-- Word 2010+ рендеринг -->
      <w:drawing>
        <wp:anchor><!-- ... modern shape ... --></wp:anchor>
      </w:drawing>
    </mc:Choice>
    <mc:Fallback>
      <!-- Fallback для старых потребителей -->
      <w:pict>
        <v:rect><!-- ... VML shape ... --></v:rect>
      </w:pict>
    </mc:Fallback>
  </mc:AlternateContent>
</w:r>
```

**Для round-trip**: хранить весь `mc:AlternateContent` как RawXML. Не пытаться парсить внутренности.

---

# 4. Порядок элементов (обязательные xsd:sequence)

## 4.1 CT_PPrBase — порядок элементов

Нарушение этого порядка → Word выдаёт "файл повреждён".

```
 1. pStyle
 2. keepNext
 3. keepLines
 4. pageBreakBefore
 5. framePr
 6. widowControl
 7. numPr
 8. suppressLineNumbers
 9. pBdr
10. shd
11. tabs
12. suppressAutoHyphens
13. kinsoku
14. wordWrap
15. overflowPunct
16. topLinePunct
17. autoSpaceDE
18. autoSpaceDN
19. bidi
20. adjustRightInd
21. snapToGrid
22. spacing
23. ind
24. contextualSpacing
25. mirrorIndents
26. suppressOverlap
27. jc
28. textDirection
29. textAlignment
30. textboxTightWrap
31. outlineLvl
32. divId
33. cnfStyle
(затем в CT_PPr: rPr, sectPr, pPrChange)
```

## 4.2 EG_SectPrContents — порядок элементов

```
 1. footnotePr
 2. endnotePr
 3. type
 4. pgSz
 5. pgMar
 6. paperSrc
 7. pgBorders
 8. lnNumType
 9. pgNumType
10. cols
11. formProt
12. vAlign
13. noEndnote
14. titlePg
15. textDirection
16. bidi
17. rtlGutter
18. docGrid
19. printerSettings
```

## 4.3 CT_TblPr — порядок элементов

```
 1. tblStyle
 2. tblpPr
 3. tblOverlap
 4. bidiVisual
 5. tblStyleRowBandSize
 6. tblStyleColBandSize
 7. tblW
 8. jc
 9. tblCellSpacing
10. tblInd
11. tblBorders
12. shd
13. tblLayout
14. tblCellMar
15. tblLook
16. tblCaption
17. tblDescription
(затем tblPrChange)
```

## 4.4 EG_RPrBase — НЕ sequence, а choice

EG_RPrBase — это `xsd:choice`, поэтому формально порядок не обязателен. Однако Word всегда записывает в определённом порядке. Для максимальной совместимости рекомендуется:

```
rStyle, rFonts, b, bCs, i, iCs, caps, smallCaps, strike, dstrike,
outline, shadow, emboss, imprint, noProof, snapToGrid, vanish,
webHidden, color, spacing, w, kern, position, sz, szCs, highlight,
u, effect, bdr, shd, fitText, vertAlign, rtl, cs, em, lang,
eastAsianLayout, specVanish, oMath
```

---

# 5. Ловушки и подводные камни

## 5.1 CT_OnOff — полная семантика

```go
// Варианты в XML              Go-значение             Bool
// <w:b/>                      &CT_OnOff{Val: nil}     true
// <w:b w:val="true"/>         &CT_OnOff{Val: "true"}  true
// <w:b w:val="1"/>            &CT_OnOff{Val: "1"}     true
// <w:b w:val="on"/>           &CT_OnOff{Val: "on"}    true
// <w:b w:val="false"/>        &CT_OnOff{Val: "false"} false
// <w:b w:val="0"/>            &CT_OnOff{Val: "0"}     false
// <w:b w:val="off"/>          &CT_OnOff{Val: "off"}   false
// элемент отсутствует         nil                     inherit (from style)

// При маршализации для "true":
//   Если Val == nil → <w:b/>
//   Если Val == "true" → <w:b w:val="true"/>
//   (сохранять оригинальную форму для round-trip)
```

## 5.2 xml:space="preserve"

```go
// ОБЯЗАТЕЛЬНО для <w:t> с leading/trailing whitespace
// <w:t xml:space="preserve"> text </w:t>   ← правильно
// <w:t> text </w:t>                        ← НЕПРАВИЛЬНО: пробелы потеряются

// Автодобавление при маршализации:
func marshalText(t CT_Text) {
    if strings.HasPrefix(t.Value, " ") || strings.HasSuffix(t.Value, " ") ||
       strings.Contains(t.Value, "\t") {
        // Добавить xml:space="preserve"
    }
}
```

## 5.3 RSID — Revision Save ID

```
- Формат: 8-символьный hex (4 байта): "00AB1234"
- Диапазон: 0x00000001 — 0x7FFFFFFE
- 0x00000000 и >= 0x7FFFFFFF — невалидны
- Генерация: crypto/rand → hex → обрезать до 8 символов
- При создании нового документа: генерировать один rsid для всего документа
- Word генерирует новый rsid при каждом сохранении
```

## 5.4 Уникальность w:id

```
- Все w:id в ins/del/moveFrom/moveTo/bookmarkStart/bookmarkEnd/
  commentRangeStart/commentRangeEnd ДОЛЖНЫ быть уникальны
  В ПРЕДЕЛАХ ВСЕГО ДОКУМЕНТА (включая headers, footers, comments, footnotes)
- Тип: int (не hex)
- При генерации: использовать автоинкремент начиная с 0
```

## 5.5 Минимум сноски

```
footnotes.xml и endnotes.xml ОБЯЗАНЫ содержать как минимум два элемента:
- footnote/endnote с type="separator" (обычно id=-1)
- footnote/endnote с type="continuationSeparator" (обычно id=0)
Без них Word показывает ошибку при открытии.
```

## 5.6 Минимум в ячейке таблицы

```
Каждая <w:tc> ОБЯЗАНА содержать хотя бы один <w:p>.
Пустая ячейка:
<w:tc>
  <w:tcPr><w:tcW w:w="4675" w:type="dxa"/></w:tcPr>
  <w:p/>
</w:tc>
```

## 5.7 Четыре обязательных стиля по умолчанию

```
styles.xml ДОЛЖЕН содержать четыре стиля с default="1":
1. type="paragraph" styleId="Normal"
2. type="character" styleId="DefaultParagraphFont"
3. type="table"     styleId="TableNormal"
4. type="numbering"  styleId="NoList"
```

## 5.8 Landscape: w и h НЕ меняются местами

```xml
<!-- Landscape: orient="landscape", но w < h! -->
<w:pgSz w:w="15840" w:h="12240" w:orient="landscape"/>
<!-- Здесь w=15840 (11") — это физическая ширина страницы,
     h=12240 (8.5") — физическая высота.
     Word просто поворачивает координатную систему. -->

<!-- НО! Некоторые документы генерируют по-другому (w > h без orient).
     Парсер должен принимать оба варианта.
     При генерации: ставить orient="landscape" + w > h. -->
```

## 5.9 ZIP-порядок

```
Некоторые потребители (особенно Google Docs) чувствительны к порядку
файлов в ZIP. Рекомендуемый порядок:
1.  [Content_Types].xml
2.  _rels/.rels
3.  word/document.xml
4.  word/_rels/document.xml.rels
5.  word/styles.xml
6.  word/settings.xml
7.  word/webSettings.xml            ← обязателен
8.  word/fontTable.xml
9.  word/numbering.xml              (если есть)
10. word/footnotes.xml              (если есть)
11. word/endnotes.xml               (если есть)
12. word/comments.xml               (если есть)
13. word/header*.xml                (если есть)
14. word/footer*.xml                (если есть)
15. word/theme/theme1.xml           ← обязателен
16. word/media/*                    (если есть)
17. docProps/core.xml
18. docProps/app.xml
```

## 5.10 Encoding

```
- Все XML: UTF-8, без BOM
- XML declaration: <?xml version="1.0" encoding="UTF-8" standalone="yes"?>
- Всегда standalone="yes"
- Перевод строк: \r\n или \n (Word использует \r\n, но принимает оба)
```

## 5.11 ⚠️ settings.xml — Обязательные элементы (ПРОВЕРЕНО НА ПРАКТИКЕ)

```
Без следующих элементов в settings.xml MS Word показывает ошибку
«содержимое, которое не удалось прочитать» и запускает восстановление:

1. w:clrSchemeMapping — маппинг абстрактных цветовых слотов
   (bg1→light1, t1→dark1, accent1→accent1, hyperlink→hyperlink, ...)
   на именованные цвета из theme1.xml.
   БЕЗ НЕГО: Word не может разрешить themeColor="accent1" в стилях.

2. w:shapeDefaults — значения по умолчанию для VML-фигур (o:shapedefaults).
   Содержит spidmax (макс. shape ID) и shapelayout/idmap.
   БЕЗ НЕГО: Word считает файл повреждённым.

3. w:rsids — таблица Revision Save IDs.
   Минимум: rsidRoot + хотя бы один rsid.
   Значения rsid в параграфах (w:rsidR, w:rsidRDefault) ДОЛЖНЫ
   присутствовать в таблице rsids.
   БЕЗ НЕГО: Word считает файл повреждённым.

4. m:mathPr — настройки Math ML.
   Без этого блока Word добавляет его при восстановлении.

5. w:compat с compatibilityMode="15" + все 5 стандартных compatSetting.
   Неполный блок compat может привести к legacy-рендерингу.

ВЫВОД: для docx.New() копировать весь блок settings.xml как шаблон,
не пытаться собирать по частям.
```

## 5.12 ⚠️ theme1.xml — ОБЯЗАТЕЛЬНА (ПРОВЕРЕНО НА ПРАКТИКЕ)

```
word/theme/theme1.xml ОБЯЗАТЕЛЕН, если styles.xml содержит ссылки на тему:
- w:asciiTheme="minorHAnsi"  → разрешается через fontScheme/minorFont/latin
- w:cstheme="minorBidi"      → разрешается через fontScheme/minorFont/cs
- w:themeColor="accent1"     → разрешается через clrScheme/accent1
- w:themeFill="accent1"      → аналогично
- w:themeShade="BF"          → модификатор яркости

Поскольку стандартный styles.xml (docDefaults) всегда использует
asciiTheme/hAnsiTheme/cstheme для шрифтов → theme1.xml де-факто обязателен.

Минимальная тема должна содержать:
- clrScheme с 12 цветами (dk1, lt1, dk2, lt2, accent1-6, hlink, folHlink)
- fontScheme с majorFont и minorFont (latin, ea, cs)
- fmtScheme с fillStyleLst, lnStyleLst, effectStyleLst, bgFillStyleLst
  (по 3 элемента в каждом — Word проверяет количество)

Также ОБЯЗАТЕЛЬНЫ:
- relationship в document.xml.rels (Type=".../theme")
- Override в [Content_Types].xml
```

## 5.13 ⚠️ webSettings.xml — ОБЯЗАТЕЛЕН (ПРОВЕРЕНО НА ПРАКТИКЕ)

```
word/webSettings.xml должен присутствовать как part с relationship
и content type, даже если он пустой (самозакрывающийся тег).
Word добавляет его при восстановлении.
```

## 5.14 Полный чеклист «файл открывается без ошибок»

```
Перед выдачей .docx пользователю, проверить:

Parts:
☑ [Content_Types].xml — Override для каждого part
☑ _rels/.rels — officeDocument, core-properties, extended-properties
☑ word/document.xml — с полным набором xmlns
☑ word/_rels/document.xml.rels — styles, settings, webSettings, fontTable, theme
☑ word/styles.xml — 4 default стиля + docDefaults
☑ word/settings.xml — zoom, defaultTabStop, compat, rsids, mathPr,
    clrSchemeMapping, shapeDefaults, decimalSymbol, listSeparator
☑ word/webSettings.xml — может быть пустым
☑ word/fontTable.xml — шрифты, используемые в документе
☑ word/theme/theme1.xml — clrScheme + fontScheme + fmtScheme

Целостность:
☑ Каждый rId в XML → relationship в .rels
☑ Каждый part → Override в [Content_Types].xml
☑ rsidR/rsidRDefault на параграфах → значение есть в w:rsids
☑ Каждая tc → минимум один p
☑ xml:space="preserve" на w:t с пробелами по краям
```

## 5.15 ⚠️ Инварианты при редактировании и удалении элементов

```
При программном удалении/вставке/замене элементов необходимо
учитывать инварианты OOXML. Нарушение приводит к «repair» при
открытии в Word или к потере данных.

1. ЯЧЕЙКА ТАБЛИЦЫ (tc) — минимум один параграф
   Каждая <w:tc> ОБЯЗАНА содержать ≥1 <w:p>.
   При Cell.Clear() → удалить Content, вставить &para.CT_P{}.
   При удалении параграфов из ячейки → проверить что остался хотя бы один.
   Без этого: Word показывает «содержимое повреждено».

2. HEADER / FOOTER — минимум один параграф
   CT_HdrFtr.Content ОБЯЗАН содержать ≥1 элемент (минимум <w:p/>).
   При Header.Clear() / Footer.Clear() → аналогично Cell.Clear().

3. SECTION BREAK в параграфе
   Если удаляемый параграф содержит PPr.SectPr (section break),
   секция теряется. Текст до и после «склеивается» в одну секцию.
   Это допустимо, но может изменить layout.

4. ORPHAN RELATIONSHIPS
   При удалении параграфа с гиперссылкой или рана с картинкой,
   relationship (rId) в document.xml.rels остаётся.
   Word НЕ падает от orphan relationships — они безопасно игнорируются.
   Очистка — необязательна, но уменьшает размер файла.

5. ORPHAN RSID
   Если удалён параграф, его rsidR может остаться в w:rsids таблице
   settings.xml. Это безопасно — Word не падает.
   Обратная ситуация опаснее: если rsidR на параграфе НЕ в w:rsids —
   Word может показать repair dialog.

6. BOOKMARKS и COMMENT RANGES
   bookmarkStart/bookmarkEnd, commentRangeStart/commentRangeEnd
   должны идти парами. Если удалён параграф, содержащий один из
   парных элементов → нарушена целостность.
   Для v1: эта проверка не реализуется (допустимо).

7. НУМЕРАЦИЯ (w:numPr)
   При удалении нумерованных параграфов нумерация автоматически
   пересчитывается Word-ом. Нет необходимости корректировать numbering.xml.

8. TRACK CHANGES (ins/del)
   При удалении параграфа с <w:ins> или <w:del> — tracked change
   теряется. Это допустимо, если пользователь осознанно удаляет.
```

---

# 6. Единицы измерения — Краткая шпаргалка

| Контекст | Единица | 1 дюйм = | Пример |
|----------|---------|---------|--------|
| Размер страницы, поля, отступы | DXA (twips) | 1440 | `w:w="12240"` = 8.5" |
| Размер шрифта | Half-Point | 144 | `w:val="24"` = 12pt |
| Межсимвольный интервал | Twips | 1440 | `w:val="20"` = 1pt |
| Толщина границ | Eighth-Point | 576 | `w:val="4"` = 0.5pt |
| Размер изображений | EMU | 914400 | `cx="1828800"` = 2" |
| Отступы изображений | EMU | 914400 | `distT="114300"` |
| Ширина таблицы (pct) | Fiftieths of % | — | `w:w="5000"` = 100% |

```
Шпаргалка конверсий:
1 inch = 1440 DXA = 914400 EMU = 72 pt = 144 half-pt = 576 eighth-pt
1 cm   = 567 DXA  = 360000 EMU = 28.35 pt
1 mm   = 56.7 DXA = 36000 EMU  = 2.835 pt
1 pt   = 20 DXA   = 12700 EMU
```

---

# 7. Enum-значения из XSD — Полный список ключевых

## ST_Jc (выравнивание параграфа)

```
start, center, end, both, mediumKashida, distribute,
numTab, highKashida, lowKashida, thaiDistribute,
left, right  (legacy синонимы start/end)
```

## ST_LineSpacingRule

```
auto     — line × 1/240 пункта (line=240 → одинарный, line=480 → двойной)
exact    — точно line twips
atLeast  — не менее line twips
```

## ST_Border (типы границ)

```
nil, none, single, thick, double, dotted, dashed, dotDash, dotDotDash,
triple, thinThickSmallGap, thickThinSmallGap, thinThickThinSmallGap,
thinThickMediumGap, thickThinMediumGap, thinThickThinMediumGap,
thinThickLargeGap, thickThinLargeGap, thinThickThinLargeGap,
wave, doubleWave, dashSmallGap, dashDotStroked,
threeDEmboss, threeDEngrave, outset, inset,
+ ~180 декоративных (apples, balloons3Colors, christmasTree, ...)
```

## ST_Shd (заливка)

```
nil, clear, solid, horzStripe, vertStripe, reverseDiagStripe,
diagStripe, horzCross, diagCross, thinHorzStripe, thinVertStripe,
thinReverseDiagStripe, thinDiagStripe, thinHorzCross, thinDiagCross,
pct5, pct10, pct12, pct15, pct20, pct25, pct30, pct35, pct37,
pct40, pct45, pct50, pct55, pct60, pct62, pct65, pct70, pct75,
pct80, pct85, pct87, pct90, pct95
```

## ST_ThemeColor

```
dark1, light1, dark2, light2,
accent1, accent2, accent3, accent4, accent5, accent6,
hyperlink, followedHyperlink,
none, background1, text1, background2, text2
```

## ST_Underline

```
single, words, double, thick, dotted, dottedHeavy, dash, dashedHeavy,
dashLong, dashLongHeavy, dotDash, dashDotHeavy, dotDotDash,
dashDotDotHeavy, wave, wavyHeavy, wavyDouble, none
```

## ST_VerticalAlignRun

```
baseline, superscript, subscript
```

## ST_FldCharType

```
begin, separate, end
```

## ST_SectType

```
nextPage, nextColumn, continuous, evenPage, oddPage
```

## ST_NumFmt (формат нумерации)

```
decimal, upperRoman, lowerRoman, upperLetter, lowerLetter,
ordinal, cardinalText, ordinalText, hex, chicago,
ideographDigital, japaneseCounting, aiueo, iroha,
decimalFullWidth, decimalHalfWidth, japaneseLegal,
japaneseDigitalTenThousand, decimalEnclosedCircle,
decimalFullWidth2, aiueoFullWidth, irohaFullWidth,
decimalZero, bullet, ganada, chosung,
decimalEnclosedFullstop, decimalEnclosedParen,
decimalEnclosedCircleChinese, ideographEnclosedCircle,
ideographTraditional, ideographZodiac,
ideographZodiacTraditional, taiwaneseCounting,
ideographLegalTraditional, taiwaneseCountingThousand,
taiwaneseDigital, chineseCounting, chineseLegalSimplified,
chineseCountingThousand, koreanDigital, koreanCounting,
koreanLegal, koreanDigital2, vietnameseCounting,
russianLower, russianUpper, none,
numberInDash, hebrew1, hebrew2, arabicAlpha, arabicAbjad,
hindiVowels, hindiConsonants, hindiNumbers, hindiCounting,
thaiLetters, thaiNumbers, thaiCounting,
bahtText, dollarText, custom
```

## ST_HdrFtr (тип колонтитула)

```
even, default, first
```

## ST_TblWidth (тип ширины таблицы)

```
nil   — нулевая ширина
pct   — процент (×50, т.е. 5000 = 100%)
dxa   — twips (абсолютная)
auto  — автоматическая
```

---

# 8. Чеклист готовности к разработке

Перед началом реализации модуля убедитесь:

- [ ] Прочитан раздел модуля в `docx-go-implementation-tree.md`
- [ ] Прочитан раздел модуля в `docx-go-dev-plan.md` (Go-структуры)
- [ ] Для типов с `xsd:sequence` — известен порядок элементов (см. раздел 4)
- [ ] Namespace URI взяты из раздела 1 этого документа
- [ ] Для round-trip: реализовано хранение `Extra []RawXML`
- [ ] Для round-trip: реализовано сохранение namespace declarations
- [ ] CT_OnOff: реализована полная семантика (см. раздел 5.1)
- [ ] xml:space="preserve": автоматическое добавление (см. раздел 5.2)
- [ ] Написан round-trip тест: unmarshal → marshal → сравнить
- [ ] Тест на реальном XML из раздела 3

Перед реализацией `docx.New()` (создание документа с нуля) убедитесь:

- [ ] Включены ВСЕ обязательные parts (см. раздел 2.1, 5.14)
- [ ] settings.xml содержит clrSchemeMapping, shapeDefaults, rsids (см. раздел 5.11)
- [ ] theme1.xml включён и содержит clrScheme + fontScheme + fmtScheme (см. раздел 5.12)
- [ ] webSettings.xml включён, хотя бы пустой (см. раздел 5.13)
- [ ] [Content_Types].xml содержит Override для КАЖДОГО part
- [ ] document.xml.rels содержит Relationship для КАЖДОГО part
- [ ] Результат проверен открытием в MS Word БЕЗ ошибок восстановления

Перед реализацией `docx` (C-32) — операции редактирования и удаления:

- [ ] Прочитан patterns.md раздел 14 (X() escape hatch)
- [ ] Прочитан patterns.md раздел 15 (value-type в interface slice)
- [ ] Прочитан patterns.md раздел 16 (FindText / ReplaceText)
- [ ] Прочитан patterns.md раздел 17 (инварианты Remove / Insert)
- [ ] Cell.Clear() вставляет пустой `<w:p/>` (раздел 5.6, 5.15)
- [ ] Header.Clear() / Footer.Clear() вставляют пустой `<w:p/>` (раздел 5.15)
- [ ] Paragraph.Clear() сохраняет PPr (стиль не теряется)
- [ ] Все Remove*/InsertAt* документируют инвалидацию индексов в godoc
- [ ] Write-back при мутации CT_Row/CT_Tc (value в interface slice, раздел 15)
- [ ] Тест: удаление → Validate() → нет новых Fatal/Error
- [ ] Тест: Cell.Clear() → ячейка содержит ровно один `<w:p/>`
- [ ] Тест: ReplaceText round-trip → Open → Replace → Save → Open → проверить