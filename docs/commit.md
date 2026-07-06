# Commit Message 规范

---

## 一、好的 Commit Message 应该遵循的规范

### 1. 结构（推荐 Conventional Commits）

```plaintext
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

- **type**：`feat`（新功能）、`fix`（Bug 修复）、`docs`（文档）、`style`（格式）、`refactor`（重构）、`perf`（性能）、`test`（测试）、`chore`（构建/工具）
- **scope**（可选）：影响的模块（如 `convert`、`theme`）
- **subject**：简短描述（不超过 50 个字符，祈使语气，不加句号）
- **body**（可选）：详细说明“为什么”和“如何”
- **footer**（可选）：关闭 Issue 或 Breaking Change 说明

### 2. 原则

- **描述意图，而非动作**：不说“修改了代码”，而说“修复了列表项在公众号中换行的问题”
- **解释“为什么”**：在 body 中说明决策背景
- **保持原子性**：一个 commit 只做一件事

---

## 二、如何让 AI（比如我）帮你生成 Commit Message

### 方法 1：直接提供 `git diff` 输出

运行：

```bash
git diff --cached   # 如果已暂存
# 或
git diff            # 如果未暂存
```

将输出粘贴给 AI，并说：“请根据这些改动生成 Conventional Commit message。”

### 方法 2：用 AI 工具集成到 Git 钩子

可以编写脚本调用 OpenAI API 自动生成 commit message，但个人建议保持人工审核。

---

## 三、更多技巧

### 1. 使用 `git commit -m` 带多行（shell 支持）

```bash
git commit -m "fix(convert): 自定义列表项渲染器防止公众号换行" -m "- 仅对 <ul> 下的 <li> 包裹 span" -m "- 保留 <strong>"
```

### 2. 利用 `commitizen` 等工具辅助交互式撰写

```bash
npm install -g commitizen
git cz   # 交互式生成符合规范的 message
```

### 3. 将 AI 生成的 message 写入文件再提交

```bash
# 将 AI 输出的内容保存到 .git/COMMIT_EDITMSG
git commit -F .git/COMMIT_EDITMSG
```

---

## 五、快速参考：常用 Type 缩写

| Type     | 使用场景                         |
|----------|----------------------------------|
| `feat`   | 新功能                           |
| `fix`    | Bug 修复                         |
| `docs`   | 文档更新                         |
| `style`  | 代码格式（不影响逻辑）           |
| `refactor`| 重构（既非新功能也非修复）       |
| `perf`   | 性能优化                         |
| `test`   | 测试相关                         |
| `chore`  | 构建、工具、依赖更新             |

---

## 六、一句话总结

**让 AI 写 Commit Message 的方法**：  
👉 给 AI 提供 `git diff` 或改动说明，AI 就会按照规范帮你生成清晰、可维护的 commit message。  
你可以直接复制我上面提供的 message，执行 `git commit -m "..."`。

如果还有其他需求（如生成多行 body），告诉我，我再帮你调整。 😊
