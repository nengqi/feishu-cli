package cmd

import (
	"fmt"

	"github.com/riba2534/feishu-cli/internal/auth"
	"github.com/riba2534/feishu-cli/internal/client"
	"github.com/riba2534/feishu-cli/internal/config"
	"github.com/spf13/cobra"
)

var searchDocsCmd = &cobra.Command{
	Use:   "docs <query>",
	Short: "搜索云文档",
	Long: `搜索当前用户可见的飞书云文档。

注意：此功能需要 User Access Token（用户授权令牌），推荐通过 auth login 获取。

参数:
  query           搜索关键词（必需）

选项:
  --count         返回数量（0-50，默认 20）
  --offset        偏移量（offset + count < 200）
  --owner-ids     文件所有者 Open ID 列表（逗号分隔）
  --chat-ids      文件所在群 ID 列表（逗号分隔）
  --docs-types    文档类型列表（逗号分隔，可选值：doc/docx/sheet/slides/bitable/mindnote/file/wiki/shortcut）

示例:
  # 先登录获取 Token（推荐）
  feishu-cli auth login

  # 搜索包含"产品需求"的文档
  feishu-cli search docs "产品需求"

  # 搜索特定类型的文档
  feishu-cli search docs "季度报告" --docs-types doc,sheet

  # 指定返回数量和偏移
  feishu-cli search docs "技术方案" --count 10 --offset 0

  # 也可以手动指定 Token
  feishu-cli search docs "产品需求" --user-access-token <token>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Validate(); err != nil {
			return err
		}

		query := args[0]

		// 获取 user access token
		flagToken, _ := cmd.Flags().GetString("user-access-token")
		cfg := config.Get()
		userAccessToken, err := auth.ResolveUserAccessToken(flagToken, cfg.UserAccessToken, cfg.AppID, cfg.AppSecret, cfg.BaseURL)
		if err != nil {
			return err
		}

		// 获取其他参数
		count, _ := cmd.Flags().GetInt("count")
		offset, _ := cmd.Flags().GetInt("offset")
		ownerIDsStr, _ := cmd.Flags().GetString("owner-ids")
		chatIDsStr, _ := cmd.Flags().GetString("chat-ids")
		docsTypesStr, _ := cmd.Flags().GetString("docs-types")
		output, _ := cmd.Flags().GetString("output")

		// 解析逗号分隔的列表
		var ownerIDs, chatIDs, docsTypes []string
		if ownerIDsStr != "" {
			ownerIDs = splitAndTrim(ownerIDsStr)
		}
		if chatIDsStr != "" {
			chatIDs = splitAndTrim(chatIDsStr)
		}
		if docsTypesStr != "" {
			docsTypes = splitAndTrim(docsTypesStr)
		}

		opts := client.SearchDocWikiOptions{
			Query:    query,
			Count:    count,
			Offset:   offset,
			OwnerIDs: ownerIDs,
			ChatIDs:  chatIDs,
			DocTypes: docsTypes,
		}

		result, err := client.SearchDocWiki(opts, userAccessToken)
		if err != nil {
			return err
		}

		if output == "json" {
			if err := printJSON(result); err != nil {
				return err
			}
		} else {
			if len(result.ResUnits) == 0 {
				fmt.Println("未找到匹配的文档")
				return nil
			}

			fmt.Printf("搜索结果（共 %d 条）:\n\n", result.Total)
			for i, unit := range result.ResUnits {
				fmt.Printf("[%d] %s\n", i+1, unit.Title)
				if unit.DocsType != "" {
					fmt.Printf("    类型: %s\n", unit.DocsType)
				}
				if unit.URL != "" {
					fmt.Printf("    链接: %s\n", unit.URL)
				}
				if unit.OwnerID != "" {
					fmt.Printf("    所有者: %s\n", unit.OwnerID)
				}
				fmt.Println()
			}

			if result.HasMore {
				fmt.Printf("\n还有更多结果，使用 --offset %d 获取下一页\n", offset+count)
			}
		}

		return nil
	},
}

func init() {
	searchCmd.AddCommand(searchDocsCmd)

	searchDocsCmd.Flags().String("user-access-token", "", "User Access Token（用户授权令牌）")
	searchDocsCmd.Flags().Int("count", 20, "返回数量（0-50）")
	searchDocsCmd.Flags().Int("offset", 0, "偏移量（offset + count < 200）")
	searchDocsCmd.Flags().String("owner-ids", "", "文件所有者 Open ID 列表（逗号分隔）")
	searchDocsCmd.Flags().String("chat-ids", "", "文件所在群 ID 列表（逗号分隔）")
	searchDocsCmd.Flags().String("docs-types", "", "文档类型列表（逗号分隔，可选值：doc/docx/sheet/slides/bitable/mindnote/file/wiki/shortcut）")
	searchDocsCmd.Flags().StringP("output", "o", "", "输出格式（json）")
}
