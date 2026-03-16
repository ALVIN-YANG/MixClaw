# MixClaw 当前功能说明 (Features)

**MixClaw** 是一个极简、极速的 AI Agent 引擎，融合了 PicoClaw 的 Go 语言轻量级架构和 OpenClaw 风格的丰富交互能力。

## 目前已实现功能 (MVP)

### 1. 核心架构
* **单一可执行文件**: 整个项目后端（API网关）和前端（Web UI）均通过 `go:embed` 编译为一个独立的二进制文件，无需额外配置 Web 服务器或 Node.js 环境。
* **零依赖**: 直接运行编译后的 Go 程序即可使用。

### 2. 内嵌 Web UI
提供了对标 OpenClaw 风格的现代深色极客风 Web 界面，访问 `http://localhost:18790` 即可：
* **对话面板 (Chat)**:
  * 仿终端打字机风格的消息流。
  * 用户和 AI Agent 交互的主界面。
* **大模型配置面板 (Settings)**:
  * 支持自定义 LLM Provider。目前内置预设了：
    * **火山引擎 (Volcengine)** - 完美支持火山引擎特有的 Endpoint ID 和高性价比模型接入。
    * **DeepSeek**
    * **OpenAI**
    * **Ollama** (本地运行)
  * 支持自定义配置 API Base URL、Model Name 和 API Key，完全 **Bring Your Own Key (BYOK)**。
  * 工作区安全沙箱限制配置（限制 Agent 的文件系统访问权限）。
* **Agent 状态面板 (Agents)**:
  * 实时展示主代理（Main Agent）的运行状态及当前绑定的模型。
  * 为未来的子代理 (Subagent) 提供管理入口。

### 3. 数据持久化
* **本地化配置存储**: 用户的模型配置和 API 密钥会安全地以 JSON 格式持久化保存在用户主目录下的 `~/.mixclaw/config.json` 文件中，启动自动加载，无需每次重配。

## 待开发功能 (即将加入)
1. **Agent 模型接入层 (LLM Provider)**: 将 `main.go` 中的对话逻辑从 Mock 替换为真实的大语言模型 API 轮询与流式返回（支持 Tool Calling / Function Calling）。
2. **工具箱系统 (Tools)**:
   * 工作区文件读写工具 (`read_file`, `write_file`)。
   * 安全沙箱终端执行工具 (`exec`)。
3. **多 Agent 协作 (Swarm/Spawn)**:
   * 支持动态创建特定角色的子代理并分配任务。
