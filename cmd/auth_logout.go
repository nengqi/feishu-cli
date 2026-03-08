package cmd

import (
	"fmt"

	"github.com/riba2534/feishu-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录（清除本地 token）",
	Long: `清除本地存储的 OAuth token 信息。

删除文件: ~/.feishu-cli/token.json

示例:
  feishu-cli auth logout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.DeleteToken(); err != nil {
			return err
		}

		path, _ := auth.TokenPath()
		fmt.Printf("已清除本地授权信息 (%s)\n", path)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLogoutCmd)
}
