# Agent-Code

Agent-code a lightweight, coding assistant built in Go that runs in terminal inspired by Claude Code. 
The CLI tool should be designed to help interact with AI models, run terminal commands, safely modify filesystem, and optionally extend its functionality.

---

## Features

- **Command-line Interface (CLI)** – Using commands like `agent-code create`, `agent-code open`, `agent-code read`.
- **TUI Mode** – An interactive terminal UI using Bubble Tea and Lipgross for UI styling on top of Cobra.
- **Safe File & Command Access** – Read, write, and run commands within scope.
---

### Installation

```bash
git clone https://github.com/nathanmbicho/agent-code-assignment

cd agent-code-assignment

make build

./agent-code -h
```

### Scope

So far the project only runs the file and access commands.

#### 
- **CMD** - this holds the cobra TUI commands 
- **pkg**
  - **components** - ui components, (list and input components)
  - **ui** - lipgross ui stylings
- **main** - app execution takes place
- **makefile** - for running code in dev mode



### Next Plan?
Integrate LLM AI agent's to use in ask mode or agent mode
