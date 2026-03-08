package cmd

import (
	"fmt"
	"time"

	"github.com/riba2534/feishu-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看当前授权状态",
	Long: `查看本地存储的 OAuth token 状态。

显示内容:
  - Access Token（脱敏）和有效期
  - Refresh Token 状态和有效期
  - 授权范围

示例:
  feishu-cli auth status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return fmt.Errorf("读取 token 失败: %w", err)
		}

		if token == nil {
			fmt.Println("授权状态: 未登录")
			fmt.Println("  使用 feishu-cli auth login 进行授权")
			return nil
		}

		fmt.Println("授权状态: 已登录")
		fmt.Printf("  Access Token:   %s\n", auth.MaskToken(token.AccessToken))

		if token.IsAccessTokenValid() {
			remaining := time.Until(token.ExpiresAt)
			fmt.Printf("  有效期至:        %s（剩余 %s）\n",
				token.ExpiresAt.Format("2006-01-02 15:04:05"),
				formatDuration(remaining))
		} else {
			fmt.Printf("  有效期至:        %s（已过期）\n", token.ExpiresAt.Format("2006-01-02 15:04:05"))
		}

		if token.RefreshToken != "" {
			if token.IsRefreshTokenValid() {
				if token.RefreshExpiresAt.IsZero() {
					fmt.Println("  Refresh Token:  有效（过期时间未知）")
				} else {
					remaining := time.Until(token.RefreshExpiresAt)
					fmt.Printf("  Refresh Token:  有效（剩余 %s）\n", formatDuration(remaining))
				}
			} else {
				fmt.Println("  Refresh Token:  已过期")
			}
		}

		if token.Scope != "" {
			fmt.Printf("  授权范围:        %s\n", token.Scope)
		}

		return nil
	},
}

// formatDuration 格式化时间间隔为友好显示
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "已过期"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d 天 %d 小时", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%d 小时 %d 分", hours, minutes)
	}
	return fmt.Sprintf("%d 分钟", minutes)
}

func init() {
	authCmd.AddCommand(authStatusCmd)
}
