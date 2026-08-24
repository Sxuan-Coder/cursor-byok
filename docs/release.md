# 发布规则

流程：**改版本号 → 更新 release-notes.md → 推送 tag → CI 自动打包发布**。

## 1. 修改版本号

`build/config.yml` 中的 `info.version`：

```yaml
info:
  version: "0.0.50"   # 与即将推送的 tag 一致（不带 v 前缀）
```

CI 会以 tag 为准自动同步版本号，构建资产（Info.plist / info.json / nfpm.yaml）无需手动重新生成。

## 2. 更新发布说明

重写 `release-notes.md`（整份替换为本版变更条目，即 GitHub Release 正文，也是 `update.json` 的 `release_notes`）：

```markdown
- 修复xxx
- 支持xxx
```

## 3. 提交并推送 tag

```bash
git add build/config.yml release-notes.md
git commit -m "release: 0.0.50"
git push origin main
git tag v0.0.50
git push origin v0.0.50    # 推送 tag 即触发 CI
```

## 4. CI 自动完成（`.github/workflows/release.yml`）

- 四平台构建：`windows-amd64.zip`、`macos-arm64.tar.gz`、`macos-amd64.tar.gz`、`linux-amd64.tar.gz`
- 生成 `update.json`（应用内更新清单：版本、说明、下载地址、sha256，指向本仓库）
- 创建 GitHub Release：标题为 tag，正文为 release-notes.md，附全部产物

全流程约 10~15 分钟，进度见 [Actions 页面](https://github.com/Sxuan-Coder/cursor-byok/actions)。

## 注意事项

- tag 必须为 `v主.次.修订` 格式（如 `v0.0.50`）才触发；`v0.1.0-beta` 等预发布 tag 不触发。
- 重发同一版本：`git tag -f v0.0.50 main && git push --force origin v0.0.50`（会覆盖同名 Release 的产物）。
- 发布后访问 [Releases 页面](https://github.com/Sxuan-Coder/cursor-byok/releases) 确认产物与 update.json 齐全。