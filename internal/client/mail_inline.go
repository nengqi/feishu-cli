package client

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MailInlineImageRef 描述 body 中一个本地图片引用的扫描结果
// RawSrc:     原 HTML 中 <img src="..."> 里的字面值
// LocalPath:  解析后的绝对路径
// CID:        生成的唯一 CID（不含 "<>" 包裹，不含 "cid:" 前缀）
// FileToken:  上传到 drive 后返回的 file_token（用于附件回写 EML body part）
// FileName:   原始文件名
// Bytes:      文件内容（构造 multipart/related 时使用）
// MIME:       内容类型（image/png 等，按扩展名兜底为 application/octet-stream）
type MailInlineImageRef struct {
	RawSrc    string
	LocalPath string
	CID       string
	FileToken string
	FileName  string
	Bytes     []byte
	MIME      string
}

// imgSrcRegexp 抓 <img ...src="...">，捕获 src 的字面值
// 不区分大小写，兼容 <IMG SRC='...'>，单/双引号皆可
var imgSrcRegexp = regexp.MustCompile(`(?i)<img\b[^>]*\bsrc\s*=\s*(?:"([^"]+)"|'([^']+)')[^>]*>`)

// inlineURISchemeRegexp 检测带 scheme 的 URL（已经是 cid:/http(s):/data: 不再扫描）
var inlineURISchemeRegexp = regexp.MustCompile(`(?i)^[a-z][a-z0-9+.\-]*:`)

// ScanInlineImagePaths 扫描 HTML body 中的 <img src="local-path">，仅返回本地路径
// 已经是 cid:/http:/https:/data: 等 scheme 的会跳过
// 同一文件路径只返回一次（去重）
func ScanInlineImagePaths(body string) []string {
	matches := imgSrcRegexp.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var out []string
	for _, m := range matches {
		var src string
		if len(m) >= 2 && m[1] != "" {
			src = m[1]
		} else if len(m) >= 3 && m[2] != "" {
			src = m[2]
		}
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		// 协议无关 URL（//cdn.com/x.png）跳过
		if strings.HasPrefix(src, "//") {
			continue
		}
		// 有 scheme（http:/https:/data:/cid: ...）跳过
		if inlineURISchemeRegexp.MatchString(src) {
			continue
		}
		if seen[src] {
			continue
		}
		seen[src] = true
		out = append(out, src)
	}
	return out
}

// ReplaceInlineImageSrc 把 body 中所有匹配 rawSrc 的 <img src="rawSrc"> 替换为 cid:cid
// 使用倒序替换以保持索引稳定
func ReplaceInlineImageSrc(body string, refs []MailInlineImageRef) string {
	if len(refs) == 0 {
		return body
	}
	matches := imgSrcRegexp.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return body
	}
	srcToCID := make(map[string]string, len(refs))
	for _, r := range refs {
		srcToCID[r.RawSrc] = r.CID
	}
	// 倒序，避免 index 漂移
	out := body
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		// m[2]:m[3] 为双引号 src；m[4]:m[5] 为单引号 src
		var s, e int
		if m[2] >= 0 {
			s, e = m[2], m[3]
		} else if m[4] >= 0 {
			s, e = m[4], m[5]
		} else {
			continue
		}
		srcVal := strings.TrimSpace(out[s:e])
		cid, ok := srcToCID[srcVal]
		if !ok {
			continue
		}
		out = out[:s] + "cid:" + cid + out[e:]
	}
	return out
}

// GenerateMailCID 生成一个 20-hex CID（与 lark-cli 风格一致）
func GenerateMailCID() (string, error) {
	var b [10]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("生成 CID 失败: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

// guessMIMEByExt 按扩展名兜底 MIME（避免 multipart/related 缺 Content-Type）
func guessMIMEByExt(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// LoadInlineImageBytes 读取本地文件并填充 FileName/Bytes/MIME
// 用于 EML builder 构造 multipart/related part；调用方可选择只用 FileToken
// 而不读盘（如已经走 drive 上传完毕、只需 cid 替换）
func LoadInlineImageBytes(ref *MailInlineImageRef) error {
	abs, err := filepath.Abs(ref.LocalPath)
	if err != nil {
		return fmt.Errorf("解析路径失败 %s: %w", ref.LocalPath, err)
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return fmt.Errorf("读取本地图片失败 %s: %w", abs, err)
	}
	ref.Bytes = data
	if ref.FileName == "" {
		ref.FileName = filepath.Base(abs)
	}
	if ref.MIME == "" {
		ref.MIME = guessMIMEByExt(ref.FileName)
	}
	return nil
}
