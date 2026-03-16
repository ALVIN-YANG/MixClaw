# OpenClaw 功能点全景图与 MixClaw 裁剪计划

为了打造极致轻量级的 **MixClaw**（基于 Go 的单文件 AI Agent），我们全面梳理了 OpenClaw 的功能点，并制定了相应的“断舍离”裁剪计划。

## 1. 核心架构与网关 (Gateway)
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **Local-first Gateway** (本地统一控制平面) | ✅ **保留** | MixClaw 也会有一个本地 HTTP 网关，处理所有请求。 |
| **Web Control UI & WebChat** | ✅ **保留并强化** | 核心卖点！直接将 Web 前端打包进 Go 二进制文件。 |
| **多 Agent 复杂路由规则** | ✂️ **裁剪** | 初期只保留一个主 Agent，加上可通过工具动态衍生 (Spawn) 的子任务代理，去掉复杂的路由分发配置。 |
| **Session Model** (会话与群组隔离) | ⚠️ **简化** | 保留最基础的会话隔离（Session ID），去掉复杂的群组响应规则。 |

## 2. 消息渠道接入 (Channels)
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **支持 20+ 聊天软件** (WhatsApp/Telegram/Slack/Discord 等) | ✂️ **大幅裁剪** | 太过臃肿！MixClaw 核心仅支持自带的 **Web UI** 和 **CLI**。未来将国内主流渠道（如 Feishu/钉钉/QQ）作为独立、可选的扩展件，不打包进核心中。 |
| **群聊提及与路由 (Group Routing)** | ✂️ **裁剪** | 初期不处理群聊逻辑。 |
| **媒体管道 (音频/视频/图片解析)** | ✂️ **裁剪** | 不内置重型多媒体解析和语音转文字引擎，维持极低内存占用。 |

## 3. 大模型与提供商 (Models & Auth)
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **多模型与 Provider 支持** | ✅ **保留** | 重点优化接入如火山引擎 (Volcengine)、DeepSeek、OpenAI 以及本地模型 (Ollama/vLLM)。 |
| **Model Failover (故障转移与轮询)** | ✂️ **裁剪** | 没必要搞这么复杂，调用失败直接反馈给用户即可。 |

## 4. 客户端与周边节点 (Apps & Nodes)
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **macOS 状态栏原生 App** | ✂️ **砍掉** | 专注 Web UI 跨平台体验，不写原生客户端。 |
| **iOS / Android 伴随节点 (调用手机相机/传感器)** | ✂️ **砍掉** | 偏离核心需求，太重了。 |
| **Voice Wake (语音唤醒) & Talk Mode** | ✂️ **砍掉** | 不需要实时语音功能。 |
| **Live Canvas (实时可视化画布)** | ✂️ **砍掉** | 保留 Markdown 渲染就足够了，不需要同步画板。 |

## 5. 工具箱与自动化 (Tools)
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **Browser Control (控制 Chrome 浏览器爬虫)** | ⚠️ **简化** | 抛弃沉重的 Puppeteer/Playwright 依赖，改用轻量级的 HTTP Fetch (类似 `curl` 或 Jina Reader)。 |
| **系统命令与文件操作** | ✅ **保留并加固** | 核心工具，但要在 Go 里实现严格的路径限制（工作区沙箱）。 |
| **Cron 定时任务 & Webhooks** | ✂️ **裁剪/延后** | 初期不引入复杂的定时调度引擎。 |
| **Skills Platform (技能扩展市场)** | ✅ **保留轻量版** | 允许用户通过 Web 界面贴入代码或配置文件，快速扩展新能力。 |

## 6. 部署与运维
| OpenClaw 功能点 | MixClaw 处理方式 | 说明 |
| :--- | :--- | :--- |
| **Docker / Nix 部署支持** | ✂️ **不需要** | **MixClaw 最大优势**：它就是一个单独的 `mixclaw.exe` 或 `mixclaw` 文件，下载即用，无需任何容器或环境！ |
| **内置 Tailscale/SSH 内网穿透** | ✂️ **砍掉** | 越权了，内网穿透由用户自行用其他工具解决。 |

---

## 🎯 MixClaw 终极形态总结
去除 OpenClaw 庞大的客户端矩阵、多媒体解析、自动化调度和冷门聊天软件对接。
我们只保留最核心的骨架：**"大模型大脑 + Go 沙箱执行器 + 内嵌 Web 控制台"**。
最终产物大小预计 `< 30MB`，启动时间 `< 100ms`，内存占用 `< 30MB`。
