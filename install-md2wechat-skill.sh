#!/usr/bin/env bash
set -euo pipefail

# ==================== 配置 ====================
SKILL_NAME="md2wechat"
REPO_URL="https://github.com/wychl/md2wechat"
DEFAULT_SKILLS_DIR="${HOME}/.openclaw/workspace/skills"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ==================== 辅助函数 ====================
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检测命令是否存在
command_exists() { command -v "$1" &> /dev/null; }

# ==================== 1. 检查并安装 md2wechat ====================
log_info "检查 md2wechat 二进制..."

if ! command_exists md2wechat; then
    log_warn "未找到 md2wechat，尝试自动安装..."
    if command_exists go; then
        log_info "使用 go install 安装最新版本..."
        go install "${REPO_URL}/cmd/md2wechat@latest"
        # 刷新 PATH（如果 GOBIN 不在 PATH 中）
        export PATH="$PATH:$(go env GOPATH)/bin"
        if command_exists md2wechat; then
            log_success "md2wechat 安装成功"
        else
            log_error "自动安装失败，请手动安装后重试"
            exit 1
        fi
    else
        log_error "未检测到 Go 环境，请先安装 Go 或手动下载二进制"
        echo "下载地址: ${REPO_URL}/releases"
        exit 1
    fi
else
    log_success "md2wechat 已安装: $(command -v md2wechat)"
fi

# 显示版本信息（可选）
MD2WECHAT_VERSION=$(md2wechat --version 2>/dev/null || echo "unknown")
log_info "版本: $MD2WECHAT_VERSION"

# ==================== 2. 确定技能安装目录 ====================
SKILLS_DIR="${OPENCLAW_SKILLS_DIR:-$DEFAULT_SKILLS_DIR}"
SKILL_PATH="${SKILLS_DIR}/${SKILL_NAME}"

log_info "技能安装目录: $SKILL_PATH"

# 备份已存在的技能
if [[ -d "$SKILL_PATH" ]]; then
    BACKUP_PATH="${SKILL_PATH}.backup.$(date +%Y%m%d%H%M%S)"
    log_warn "技能已存在，备份至: $BACKUP_PATH"
    mv "$SKILL_PATH" "$BACKUP_PATH"
fi

mkdir -p "$SKILL_PATH"

# ==================== 3. 生成 SKILL.md ====================
log_info "生成 SKILL.md 配置文件..."

MD2WECHAT_BIN=$(command -v md2wechat)

cat > "${SKILL_PATH}/SKILL.md" <<EOF
---
name: ${SKILL_NAME}
description: 将 Markdown 文档转换为微信公众号兼容的 HTML 格式。
version: 1.0.0
---

# ${SKILL_NAME} Skill

## 功能描述
此技能调用本地 \`md2wechat\` 命令行工具，将用户提供的 Markdown 内容或文件转换为适配微信公众号的 HTML 代码。

## 使用场景
- 当用户希望将 Markdown 格式的文章、文档或笔记转换为用于微信公众号发布的 HTML 格式时。
- 当用户明确要求 "转为微信文章"、"生成公众号 HTML" 或使用类似短语时。

## 执行指令
**重要：** 必须使用 \`-skill\` 参数，以确保输出为结构化的 JSON 格式。

使用 \`bash\` 工具执行以下命令：

\`\`\`bash
echo '{{.MarkdownContent}}' | ${MD2WECHAT_BIN} -skill -theme "{{.Theme | default \"default\"}}" -color "{{.Color | default \"#0F4C81\"}}" {{if .ReadingTime}}-reading-time{{end}}
\`\`\`

### 参数说明
- \`MarkdownContent\` (string, 必填): 需要转换的 Markdown 文本内容。
- \`Theme\` (string, 可选): 主题名称，支持 \`default\`、\`simple\`、\`grace\`、\`fresh\`、\`warm\`，默认 \`default\`。
- \`Color\` (string, 可选): 主色调（十六进制），默认 \`#0F4C81\`。
- \`ReadingTime\` (boolean, 可选): 是否显示阅读时间，默认 \`false\`。

### 命令示例
\`\`\`bash
echo '# Hello World' | ${MD2WECHAT_BIN} -skill -theme grace -color "#92617E" -reading-time
\`\`\`

## 输出处理
1. 执行命令后，解析 JSON 格式的输出。
2. 若 \`success\` 为 \`true\`，提取 \`html\` 字段的值，清理多余的代码块标记，仅返回纯净的 HTML 代码。
3. 若 \`success\` 为 \`false\`，向用户报告 \`error\` 字段中的错误信息。
EOF

log_success "SKILL.md 已生成"

# ==================== 4. 验证技能语法（可选） ====================
if command_exists openclaw; then
    log_info "尝试验证技能配置..."
    if openclaw skills validate "${SKILL_PATH}/SKILL.md" 2>/dev/null; then
        log_success "技能配置有效"
    else
        log_warn "无法验证技能（openclaw 命令不支持 validate），请手动检查"
    fi
fi

# ==================== 5. 输出后续指引 ====================
log_success "安装完成！"
echo "===================================================="
echo -e "${YELLOW}下一步操作：${NC}"
echo "1. 重启 OpenClaw 服务以加载新技能："
echo "   openclaw gateway restart"
echo "   # 或 systemctl --user restart openclaw"
echo ""
echo "2. 验证技能是否已加载："
echo "   openclaw skills list | grep ${SKILL_NAME}"
echo ""
echo "3. 测试技能（在 OpenClaw 聊天中）："
echo "   /${SKILL_NAME} 帮我将 '# 你好' 转为微信 HTML"
echo ""
echo "4. 如需卸载，删除技能目录："
echo "   rm -rf ${SKILL_PATH}"
echo "===================================================="

# ==================== 6. 可选：自动重启 OpenClaw ====================
read -p "是否立即重启 OpenClaw 服务? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if command_exists openclaw && openclaw gateway restart; then
        log_success "OpenClaw 已重启"
    else
        log_warn "无法自动重启，请手动重启"
    fi
fi