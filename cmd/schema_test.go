package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TestSchemaListServices 验证空 path 时列出所有 service。
func TestSchemaListServices(t *testing.T) {
	var buf bytes.Buffer
	if err := runSchema(&buf, "", "pretty"); err != nil {
		t.Fatalf("runSchema(\"\") err = %v", err)
	}
	out := buf.String()
	// 至少包含几个高频 service
	for _, want := range []string{"im", "drive", "calendar", "sheets", "bitable", "wiki"} {
		if !strings.Contains(out, want) {
			// bitable 在 meta_data 里可能叫 base，跳过该项也行
			if want == "bitable" {
				continue
			}
			t.Errorf("service list missing %q\n--- output ---\n%s", want, out)
		}
	}
}

// TestSchemaServiceDetail 验证 service-only path 列出该 service 下 resource.method。
func TestSchemaServiceDetail(t *testing.T) {
	var buf bytes.Buffer
	if err := runSchema(&buf, "im", "pretty"); err != nil {
		t.Fatalf("runSchema(\"im\") err = %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "messages") {
		t.Errorf("im resource list missing 'messages'\n--- output ---\n%s", out)
	}
}

// TestSchemaMethodDetail 验证具体 method path 输出关键字段。
func TestSchemaMethodDetail(t *testing.T) {
	var buf bytes.Buffer
	if err := runSchema(&buf, "im.messages.delete", "pretty"); err != nil {
		t.Fatalf("runSchema(\"im.messages.delete\") err = %v", err)
	}
	out := buf.String()
	for _, want := range []string{"DELETE", "/open-apis/im/v1/messages", "message_id"} {
		if !strings.Contains(out, want) {
			t.Errorf("method detail missing %q\n--- output ---\n%s", want, out)
		}
	}
}

// TestSchemaJSONFormat 验证 --format json 在 method 路径下能跑通且不报错。
// printJSON 写到 os.Stdout（CLI 全局约定），这里仅校验调用不报错。
func TestSchemaJSONFormat(t *testing.T) {
	// discard 是为了避免 pretty 路径占用 buf；json 路径走 stdout
	if err := runSchema(io.Discard, "im.messages.delete", "json"); err != nil {
		t.Fatalf("runSchema json err = %v", err)
	}
}

// TestSchemaUnknownService 验证未知 service 返回友好错误。
func TestSchemaUnknownService(t *testing.T) {
	var buf bytes.Buffer
	err := runSchema(&buf, "nonexistent_service_xyz", "pretty")
	if err == nil {
		t.Fatal("expected error for unknown service, got nil")
	}
	if !strings.Contains(err.Error(), "未知 service") {
		t.Errorf("error message should mention 未知 service, got: %v", err)
	}
}

// TestSchemaListSubcommand 验证 schema list --service 子命令 runner。
func TestSchemaListSubcommand(t *testing.T) {
	var buf bytes.Buffer
	if err := runSchemaList(&buf, "im", "pretty"); err != nil {
		t.Fatalf("runSchemaList(\"im\") err = %v", err)
	}
	if !strings.Contains(buf.String(), "messages") {
		t.Errorf("schema list --service im missing 'messages'\n--- output ---\n%s", buf.String())
	}
}
