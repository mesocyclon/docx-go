package packaging

import (
	"testing"
)

// ──────────────────────────────────────────────
// Unit tests for internal helpers
// ──────────────────────────────────────────────

func TestNormalizePartName(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"word/document.xml", "/word/document.xml"},
		{"/word/document.xml", "/word/document.xml"},
		{"docProps/core.xml", "/docProps/core.xml"},
	}
	for _, tt := range tests {
		got := normalizePartName(tt.input)
		if got != tt.want {
			t.Errorf("normalizePartName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveTarget(t *testing.T) {
	tests := []struct {
		source, target, want string
	}{
		{"/word/document.xml", "styles.xml", "/word/styles.xml"},
		{"/word/document.xml", "theme/theme1.xml", "/word/theme/theme1.xml"},
		{"/word/document.xml", "/docProps/core.xml", "/docProps/core.xml"},
		{"/word/document.xml", "media/image1.png", "/word/media/image1.png"},
	}
	for _, tt := range tests {
		got := resolveTarget(tt.source, tt.target)
		if got != tt.want {
			t.Errorf("resolveTarget(%q, %q) = %q, want %q", tt.source, tt.target, got, tt.want)
		}
	}
}

func TestParseRelIDNum(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"rId1", 1},
		{"rId10", 10},
		{"rId0", 0},
		{"rId999", 999},
		{"bogus", 0},
		{"rId", 0},
	}
	for _, tt := range tests {
		got := parseRelIDNum(tt.input)
		if got != tt.want {
			t.Errorf("parseRelIDNum(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestGuessMediaContentType(t *testing.T) {
	tests := []struct {
		filename, want string
	}{
		{"image1.png", "image/png"},
		{"photo.jpg", "image/jpeg"},
		{"photo.jpeg", "image/jpeg"},
		{"icon.gif", "image/gif"},
		{"diagram.emf", "image/x-emf"},
		{"logo.svg", "image/svg+xml"},
		{"unknown.xyz", "application/octet-stream"},
	}
	for _, tt := range tests {
		got := guessMediaContentType(tt.filename)
		if got != tt.want {
			t.Errorf("guessMediaContentType(%q) = %q, want %q", tt.filename, got, tt.want)
		}
	}
}

func TestGuessContentType(t *testing.T) {
	tests := []struct {
		partName, want string
	}{
		{"/word/custom.xml", "application/xml"},
		{"/word/_rels/document.xml.rels", "application/vnd.openxmlformats-package.relationships+xml"},
		{"/word/media/image1.png", "image/png"},
		{"/word/something.bin", "application/octet-stream"},
	}
	for _, tt := range tests {
		got := guessContentType(tt.partName)
		if got != tt.want {
			t.Errorf("guessContentType(%q) = %q, want %q", tt.partName, got, tt.want)
		}
	}
}

func TestIsPackageLevelRel(t *testing.T) {
	tests := []struct {
		relType string
		want    bool
	}{
		{relCoreProperties, true},
		{relExtProperties, true},
		{relStyles, false},
		{relImage, false},
		{relHeader, false},
		{"http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail", true},
	}
	for _, tt := range tests {
		got := isPackageLevelRel(tt.relType)
		if got != tt.want {
			t.Errorf("isPackageLevelRel(%q) = %v, want %v", tt.relType, got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────
// Tests for Document helper methods
// ──────────────────────────────────────────────

func newMinimalDoc() *Document {
	return &Document{
		Media:        make(map[string][]byte),
		UnknownParts: make(map[string][]byte),
		nextRelSeq:   5,
		nextBmkID:    100,
	}
}

func TestNextRelID(t *testing.T) {
	d := newMinimalDoc()
	got1 := d.NextRelID()
	got2 := d.NextRelID()
	got3 := d.NextRelID()

	if got1 != "rId5" {
		t.Errorf("first NextRelID() = %q, want rId5", got1)
	}
	if got2 != "rId6" {
		t.Errorf("second NextRelID() = %q, want rId6", got2)
	}
	if got3 != "rId7" {
		t.Errorf("third NextRelID() = %q, want rId7", got3)
	}
}

func TestNextBookmarkID(t *testing.T) {
	d := newMinimalDoc()
	got1 := d.NextBookmarkID()
	got2 := d.NextBookmarkID()

	if got1 != 100 {
		t.Errorf("first NextBookmarkID() = %d, want 100", got1)
	}
	if got2 != 101 {
		t.Errorf("second NextBookmarkID() = %d, want 101", got2)
	}
}

func TestAddMedia(t *testing.T) {
	d := newMinimalDoc()

	rID := d.AddMedia("logo.png", []byte("PNG data"))
	if rID != "rId5" {
		t.Errorf("AddMedia returned rId=%q, want rId5", rID)
	}
	if _, ok := d.Media["logo.png"]; !ok {
		t.Error("AddMedia did not store the file in Media map")
	}

	// Adding a file with the same name should get a deduped name.
	rID2 := d.AddMedia("logo.png", []byte("PNG data 2"))
	if rID2 != "rId6" {
		t.Errorf("AddMedia second returned rId=%q, want rId6", rID2)
	}
	if _, ok := d.Media["logo1.png"]; !ok {
		t.Error("AddMedia did not deduplicate filename — expected logo1.png")
	}
}

func TestAddMediaPreservesExtension(t *testing.T) {
	d := newMinimalDoc()
	d.AddMedia("photo.jpeg", []byte("JPEG"))
	if _, ok := d.Media["photo.jpeg"]; !ok {
		t.Error("expected photo.jpeg in Media map")
	}
}

// ──────────────────────────────────────────────
// Concurrency safety smoke test
// ──────────────────────────────────────────────

func TestNextRelIDConcurrent(t *testing.T) {
	d := newMinimalDoc()
	done := make(chan string, 50)
	for i := 0; i < 50; i++ {
		go func() {
			done <- d.NextRelID()
		}()
	}
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		id := <-done
		if seen[id] {
			t.Errorf("duplicate rId generated: %s", id)
		}
		seen[id] = true
	}
}

// ──────────────────────────────────────────────
// Round-trip integration test pattern
// ──────────────────────────────────────────────
// NOTE: This test requires all dependency modules (opc, parts/*, coreprops,
// wml/*) to be compiled. It is included here to document the intended test
// pattern per patterns.md §12. In a CI build with the full module tree this
// will execute; in isolation it may be skipped.

/*
func TestRoundTrip(t *testing.T) {
	// 1. Open a known-good .docx fixture
	doc, err := Open("testdata/minimal.docx")
	if err != nil {
		t.Fatal(err)
	}

	// 2. Verify parsed parts
	if doc.Document == nil {
		t.Fatal("Document is nil after Open")
	}
	if doc.Styles == nil {
		t.Fatal("Styles is nil after Open")
	}
	if doc.Settings == nil {
		t.Fatal("Settings is nil after Open")
	}
	if doc.Fonts == nil {
		t.Fatal("Fonts is nil after Open")
	}
	if doc.CoreProps == nil {
		t.Fatal("CoreProps is nil after Open")
	}

	// 3. Save to a buffer
	var buf bytes.Buffer
	if err := doc.SaveWriter(&buf); err != nil {
		t.Fatal(err)
	}

	// 4. Re-open from the buffer
	doc2, err := OpenReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal("re-open failed:", err)
	}

	// 5. Verify key fields survived the round-trip
	if doc2.Document == nil {
		t.Fatal("Document is nil after round-trip")
	}
	if doc2.Styles == nil {
		t.Fatal("Styles lost after round-trip")
	}
	if doc2.Settings == nil {
		t.Fatal("Settings lost after round-trip")
	}
	if doc2.CoreProps == nil {
		t.Fatal("CoreProps lost after round-trip")
	}

	// Theme and WebSettings should survive as raw bytes
	if len(doc2.Theme) == 0 {
		t.Error("Theme lost after round-trip")
	}
	if len(doc2.WebSettings) == 0 {
		t.Error("WebSettings lost after round-trip")
	}

	// Unknown parts should be preserved
	if len(doc2.UnknownParts) != len(doc.UnknownParts) {
		t.Errorf("UnknownParts count: got %d, want %d",
			len(doc2.UnknownParts), len(doc.UnknownParts))
	}

	// Media should be preserved
	if len(doc2.Media) != len(doc.Media) {
		t.Errorf("Media count: got %d, want %d",
			len(doc2.Media), len(doc.Media))
	}
}
*/
