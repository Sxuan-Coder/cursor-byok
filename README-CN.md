<div align="center">

# cursor-byok（Go 版）

Cursor 的本地模型网关 —— 基于 **Go 实现**的 Fork，持续二次开发。

[下载最新版](https://github.com/Sxuan-Coder/cursor-byok/releases/latest) · [发布列表](https://github.com/Sxuan-Coder/cursor-byok/releases) · [问题反馈](https://github.com/Sxuan-Coder/cursor-byok/issues) · [English](./README.md)

[![Release](https://img.shields.io/github/v/release/Sxuan-Coder/cursor-byok?style=flat-square)](https://github.com/Sxuan-Coder/cursor-byok/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/Sxuan-Coder/cursor-byok/total?style=flat-square)](https://github.com/Sxuan-Coder/cursor-byok/releases)
[![License](https://img.shields.io/github/license/Sxuan-Coder/cursor-byok?style=flat-square)](./LICENSE)
[![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-lightgrey?style=flat-square)](https://github.com/Sxuan-Coder/cursor-byok/releases/latest)

</div>

![cursor-byok 支持接入多种模型 API](./images/cn-brand.png)

![cursor-byok 主界面](./images/cn-home.png)

## 项目介绍

cursor-byok 是一个开源的 Cursor 本地模型接入工具。它通过运行在本机的服务连接 Cursor 与你配置的模型 API，让模型请求使用自己的渠道处理，同时保留 Cursor Agent 的工具调用、Skills 和 MCP 等能力。

你可以接入 OpenAI、Anthropic 及其兼容服务，自由配置接口地址、模型、密钥和请求参数，不再局限于平台预设的模型渠道。

原仓库主线此后转向了 **Rust 重构**（`0.1.0-beta`）。本仓库**不跟随**该重构路线——这里的新功能和修复均构建在原有 Go 实现之上，两个项目会随时间逐渐分化。

## 核心能力

- **自定义模型渠道**：配置自己的 API 地址、访问密钥和模型标识。
- **多种接口协议**：支持 OpenAI、Anthropic 兼容接口及自定义端点。
- **模型管理**：添加、复制、编辑、排序和批量测试多个模型配置。
- **连接性能测试**：查看首字延迟、生成速度与模型服务的原始响应。
- **Agent 工作流**：支持工具调用、Skills、MCP 和多轮会话。
- **会话统计**：查看 Token 消耗、缓存命中率、对话轮次和价值估算。
- **跨平台运行**：支持 macOS、Windows 和 Linux。

## 快速开始

1. 从 [Releases](https://github.com/Sxuan-Coder/cursor-byok/releases/latest) 下载对应平台的最新版本。
2. 启动应用，打开“模型配置”，填写接口地址、API Key 和模型标识。
3. 测试模型配置；测试通过后返回主界面启动服务。
4. 打开 Cursor，选择已配置的模型并开始使用 Agent。

## 工作原理

```text
Cursor 客户端
    │
    │ Agent 请求与工具结果
    ▼
cursor-byok 本地服务
    │
    │ OpenAI / Anthropic 兼容请求
    ▼
你配置的模型 API
```

cursor-byok 在本机负责协议适配、模型请求转发、工具调用衔接与会话状态管理。模型 API Key 和应用配置保存在本机；实际请求仍会发送到你所配置的模型服务商。

## 开发

需要 Go 1.25、Node.js/Yarn、[Task](https://taskfile.dev) 和 Wails v3 CLI。

```bash
task dev    # 开发模式运行
task build  # 构建当前系统分发包
```

推送版本标签即可自动发布，详见 [docs/release.md](./docs/release.md)。

## 致谢

基于原项目 [leookun/cursor-byok](https://github.com/leookun/cursor-byok) 及其贡献者。

## 许可证

本项目基于 [MIT License](./LICENSE) 开源。