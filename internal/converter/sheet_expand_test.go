package converter

import (
	"errors"
	"strings"
	"testing"

	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

func TestConvertSheetExpand(t *testing.T) {
	blocks := []*larkdocx.Block{
		{
			BlockId:   strPtr("sheet1"),
			BlockType: intPtr(int(BlockTypeSheet)),
			Sheet: &larkdocx.Sheet{
				Token: strPtr("sheet_token456"),
			},
		},
	}

	var gotToken, gotSheetID, gotUserToken string
	conv := NewBlockToMarkdown(blocks, ConvertOptions{
		ExpandSheets:    true,
		UserAccessToken: "u-test",
		SheetDataProvider: func(spreadsheetToken, sheetID, userAccessToken string) ([]*SheetData, error) {
			gotToken = spreadsheetToken
			gotSheetID = sheetID
			gotUserToken = userAccessToken
			return []*SheetData{
				{
					Title: "Sheet1",
					Values: [][]any{
						{"名称", "数量"},
						{"苹果", float64(3)},
					},
				},
			}, nil
		},
	})

	got, err := conv.Convert()
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if gotToken != "sheet" || gotSheetID != "token456" || gotUserToken != "u-test" {
		t.Fatalf("provider args = (%q, %q, %q)", gotToken, gotSheetID, gotUserToken)
	}
	if strings.Contains(got, "<sheet") {
		t.Fatalf("expanded output should not contain raw sheet tag:\n%s", got)
	}
	for _, want := range []string{
		`<!-- sheet token="sheet" id="token456" -->`,
		"## Sheet1",
		"| 名称 | 数量 |",
		"| 苹果 | 3 |",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expanded output missing %q:\n%s", want, got)
		}
	}
}

func TestConvertSheetExpandFailureReturnsError(t *testing.T) {
	blocks := []*larkdocx.Block{
		{
			BlockId:   strPtr("sheet1"),
			BlockType: intPtr(int(BlockTypeSheet)),
			Sheet: &larkdocx.Sheet{
				Token: strPtr("sheet_token456"),
			},
		},
	}

	conv := NewBlockToMarkdown(blocks, ConvertOptions{
		ExpandSheets: true,
		SheetDataProvider: func(spreadsheetToken, sheetID, userAccessToken string) ([]*SheetData, error) {
			return nil, errors.New("boom")
		},
	})

	_, err := conv.Convert()
	if err == nil {
		t.Fatal("Convert() error = nil, want sheet expansion error")
	}
	if !strings.Contains(err.Error(), "电子表格自动展开失败") || !strings.Contains(err.Error(), "--expand-sheets=false") {
		t.Fatalf("unexpected error: %v", err)
	}
}
