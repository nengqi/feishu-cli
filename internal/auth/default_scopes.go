package auth

import (
	"sort"
	"strings"
)

var defaultLoginScopes = []string{
	// 用于 /authen/v1/user_info 和“当前登录用户是谁”这类最基础能力。
	"auth:user.id:read",
}

// NormalizeScopeList 规范化 scope 列表：按空白切分、去重多余空格、重新用单空格拼接。
func NormalizeScopeList(scope string) string {
	return strings.Join(UniqueScopeList(scope), " ")
}

// DefaultLoginScopes 返回 auth login 默认申请的最小核心 user scopes。
//
// 设计目的：
//   - 让首次 auth login 稳定成功，不再因默认 scope 过多触发飞书数量上限
//   - 让后续缺权限场景通过 `auth check` + `auth login --scope "..."` 按需补授权
func DefaultLoginScopes() string {
	return strings.Join(defaultLoginScopes, " ")
}

// DefaultLoginScopeList returns auth login 默认申请的最小核心 user scopes 切片。
func DefaultLoginScopeList() []string {
	return append([]string(nil), defaultLoginScopes...)
}

// UniqueScopeList splits a scope string into a de-duplicated ordered slice.
func UniqueScopeList(scope string) []string {
	seen := make(map[string]struct{}, 64)
	parts := make([]string, 0, 64)
	for _, item := range strings.Fields(scope) {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		parts = append(parts, item)
	}
	return parts
}

// MergeScopeLists merges multiple scope slices while preserving first-seen order.
func MergeScopeLists(groups ...[]string) []string {
	seen := make(map[string]struct{}, 64)
	parts := make([]string, 0, 64)
	for _, group := range groups {
		for _, item := range group {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			parts = append(parts, item)
		}
	}
	return parts
}

// JoinScopes joins multiple scope slices into a normalized space-separated string.
func JoinScopes(groups ...[]string) string {
	return strings.Join(MergeScopeLists(groups...), " ")
}

// SortScopeList returns a sorted copy, useful in tests/help output.
func SortScopeList(scopes []string) []string {
	out := append([]string(nil), scopes...)
	sort.Strings(out)
	return out
}
