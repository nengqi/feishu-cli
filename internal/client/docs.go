package client

import (
	"fmt"

	larkdocs "github.com/larksuite/oapi-sdk-go/v3/service/docs/v1"
)

// GetDocMarkdown 获取文档内容（Markdown 格式），仅支持新版 docx 文档
// 旧版 doc 类型不支持，docs/v1/content API 只接受 doc_type=docx
func GetDocMarkdown(docToken string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	req := larkdocs.NewGetContentReqBuilder().
		DocToken(docToken).
		DocType("docx").
		ContentType("markdown").
		Build()

	resp, err := client.Docs.V1.Content.Get(Context(), req)
	if err != nil {
		return "", fmt.Errorf("获取文档 Markdown 内容失败: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("获取文档 Markdown 内容失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data.Content == nil {
		return "", nil
	}

	return *resp.Data.Content, nil
}
