---
name: feishu-cli-chat
description: >-
  飞书会话浏览、消息互动与群聊管理。查看聊天记录（单聊/群聊）、搜索群聊、获取消息详情、
  Reaction 表情回应、Pin 置顶/取消置顶、删除消息，以及群聊信息管理（获取/更新/解散/成员）。
  支持普通群和话题群（可按 thread_id 获取话题回复，不会自动递归展开所有线程）。
  大多数命令需要 User Token；msg delete 默认 App Token 用于 Bot 自撤回，可显式 User Token 给管理员撤回场景。
argument-hint: <chat_id|群名|用户名>
user-invocable: true
allowed-tools: Bash(feishu-cli auth:*), Bash(feishu-cli msg:*), Bash(feishu-cli chat:*), Bash(feishu-cli search:*), Bash(feishu-cli user:*), Read, Write
---

# 飞书会话浏览与管理

本技能处理“读聊天记录 / 管群 / 消息互动”。发送新消息走 `feishu-cli-msg`；构造卡片走 `feishu-cli-card`。

## 认证

多数命令需要 User Token：

```bash
feishu-cli auth check --scope "im:message:readonly im:message.group_msg:get_as_user"
feishu-cli auth login --domain chat --recommend
```

`msg delete` 例外：默认 App Token，用于 Bot 撤回自己 24 小时内发送的消息；传 `--user-access-token` 或环境变量时可走管理员撤回场景。

## 读消息

```bash
# 群聊历史
feishu-cli msg history --container-id oc_xxx --container-id-type chat --page-size 50 -o json

# 私聊：按邮箱或 open_id 自动反查 P2P chat_id
feishu-cli msg history --user-email user@example.com --page-size 50 -o json
feishu-cli msg history --user-id ou_xxx --page-size 50 -o json

# 消息详情 / 批量详情
feishu-cli msg get <message_id> -o json
feishu-cli msg mget --message-ids <message_id1,message_id2>

# 话题回复：需要明确 thread_id
feishu-cli msg thread-messages <thread_id> --page-size 50
```

JSON 字段使用 snake_case：

```json
{
  "items": [],
  "has_more": true,
  "page_token": "next_page",
  "sender_names": {"ou_xxx": "张三"},
  "merge_forward_sub_messages": {}
}
```

分页时检查 `has_more` 和 `page_token`。合并转发子消息会自动展开到 `merge_forward_sub_messages`，通常不需要二次调用 API。

## 搜索与定位

```bash
# 搜索群聊
feishu-cli msg search-chats --query "项目群" -o json

# 搜索消息内容，P2P 用 search messages 而不是 search-chats
feishu-cli search messages "关键词" --chat-type p2p_chat -o json
feishu-cli search messages "关键词" --chat-ids oc_xxx -o json
```

搜索消息属于 `feishu-cli-search` 的能力；本技能只在聊天阅读任务里顺带调用。

## 消息互动

```bash
# Reaction
feishu-cli msg reaction add <message_id> --emoji-type THUMBSUP
feishu-cli msg reaction remove <message_id> --reaction-id <reaction_id>
feishu-cli msg reaction list <message_id>

# Pin
feishu-cli msg pin <message_id>
feishu-cli msg unpin <message_id>
feishu-cli msg pins --chat-id <chat_id>

# 删除消息
feishu-cli msg delete <message_id>
feishu-cli msg delete <message_id> --user-access-token u-xxx
```

## 群聊管理

```bash
feishu-cli chat get oc_xxx -o json
feishu-cli chat update oc_xxx --name "新群名"
feishu-cli chat member list oc_xxx -o json
feishu-cli chat member add oc_xxx --id-list ou_xxx,ou_yyy
feishu-cli chat member remove oc_xxx --id-list ou_xxx
feishu-cli chat create --name "项目群" --user-ids ou_xxx,ou_yyy
feishu-cli chat delete oc_xxx
```

## 常见决策

| 需求 | 命令 |
|---|---|
| 看某个群最近消息 | `msg history --container-id <chat_id> --container-id-type chat` |
| 看和某人的私聊 | `msg history --user-email <email>` 或 `--user-id <open_id>` |
| 找群 | `msg search-chats --query "<name>"` |
| 按关键词搜消息 | `search messages "<keyword>" ...` |
| 读合并转发 | `msg get -o json` / `msg mget --message-ids ...` / `msg history -o json`，查看 `merge_forward_sub_messages` |
| 看话题回复 | `msg thread-messages <thread_id>` |
| 发送新消息 | 转到 `feishu-cli-msg` |

## 输出处理

1. 保存 JSON 到临时文件再分析，避免长消息刷屏。
2. 文本内容通常在 `body.content` 里，需按 `msg_type` 解析 JSON 字符串。
3. `sender_names` 已自动给出 open_id 到姓名的映射，优先用它展示发送者。
