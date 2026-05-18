package converter

import "testing"

func TestBuildSheetReadChunksSplitsLargeGrid(t *testing.T) {
	chunks, err := buildSheetReadChunks("sheet1", 200, 40)
	if err != nil {
		t.Fatalf("buildSheetReadChunks() error = %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("len(chunks) = %d, want 2", len(chunks))
	}
	if chunks[0].Range != "sheet1!A1:AN125" {
		t.Fatalf("chunk[0].Range = %q", chunks[0].Range)
	}
	if chunks[1].Range != "sheet1!A126:AN200" {
		t.Fatalf("chunk[1].Range = %q", chunks[1].Range)
	}
}

func TestBuildSheetReadChunksRejectsHugeGrid(t *testing.T) {
	_, err := buildSheetReadChunks("sheet1", 1001, 100)
	if err == nil {
		t.Fatal("expected huge grid to be rejected")
	}
}

func TestMergeSheetChunkValues(t *testing.T) {
	dest := make([][]any, 3)
	mergeSheetChunkValues(dest, 2, 3, [][]any{
		{"c2", "d2"},
		{"c3"},
	})

	if got := dest[1][2]; got != "c2" {
		t.Fatalf("dest row2 col3 = %v", got)
	}
	if got := dest[1][3]; got != "d2" {
		t.Fatalf("dest row2 col4 = %v", got)
	}
	if got := dest[2][2]; got != "c3" {
		t.Fatalf("dest row3 col3 = %v", got)
	}
}
