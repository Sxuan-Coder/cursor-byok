<script setup>
import { useConfigTransfer } from "@/composables/useConfigTransfer";
import { useMessage } from "@/composables/useMessage";
import { launchCursor } from "@/services/clientApi";
import { toUserError } from "@/state/appState";
import { useRouter } from "vue-router";

const router = useRouter();
const message = useMessage();

function showActionError(title, error) {
  const detail = String(error || "服务错误").trim() || "服务错误";
  message(`${title}：${detail}`);
}

const { configTransferBusy, handleExportConfig, handleImportConfig } = useConfigTransfer({
  message,
  showActionError,
});

const ACTIONS = [
  {
    key: "launchCursor",
    title: "打开 Cursor",
    description: "启动 Cursor 客户端",
    icon: "icon-[simple-icons--cursor]",
  },
  {
    key: "add",
    title: "新增模型",
    description: "新增模型配置",
    icon: "icon-[mdi--plus-box-outline]",
  },
  {
    key: "import",
    title: "导入配置",
    description: "从文件导入配置",
    icon: "icon-[mdi--file-import-outline]",
  },
  {
    key: "export",
    title: "导出配置",
    description: "导出当前配置",
    icon: "icon-[mdi--file-export-outline]",
  },
  {
    key: "test",
    title: "测试全部",
    description: "验证所有模型可用性",
    icon: "icon-[mdi--shield-check-outline]",
  },
  {
    key: "commitPrompt",
    title: "Commit 提示词",
    description: "设置提交信息提示词",
    icon: "icon-[mdi--comment-text-outline]",
  },
];

async function handleLaunchCursor() {
  try {
    await launchCursor();
  } catch (error) {
    showActionError("打开 Cursor 失败", toUserError(error));
  }
}

function runAction(key) {
  if (configTransferBusy.value) {
    return;
  }
  switch (key) {
    case "launchCursor":
      void handleLaunchCursor();
      break;
    case "add":
      router.push({ path: "/model-config", query: { new: "1" } });
      break;
    case "import":
      void handleImportConfig();
      break;
    case "export":
      void handleExportConfig();
      break;
    case "test":
      router.push({ path: "/model-config", query: { test: "all" } });
      break;
    case "commitPrompt":
      router.push({ path: "/settings", query: { focus: "commit" } });
      break;
    default:
      break;
  }
}
</script>

<template>
  <section class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]">
    <header
      class="flex h-[42px] items-center gap-2 border-b border-[var(--border-subtle)] px-3.5"
    >
      <span class="icon-[mdi--lightning-bolt-outline] text-[15px] text-[var(--brand)]"></span>
      <h2 class="text-[13px] font-medium text-[var(--text-primary)]">快捷操作</h2>
    </header>

    <div class="grid grid-cols-3 gap-2 p-2.5 max-[1000px]:grid-cols-2">
      <button
        v-for="action in ACTIONS"
        :key="action.key"
        type="button"
        class="flex min-w-0 items-center gap-2.5 rounded-[10px] border border-[var(--border)] bg-[var(--bg-card-soft)] px-2.5 py-2.5 text-left transition-colors hover:border-[var(--brand-border)] hover:bg-[var(--brand-soft)] disabled:opacity-60"
        :disabled="configTransferBusy"
        @click="runAction(action.key)"
      >
        <span
          class="flex h-[28px] w-[28px] shrink-0 items-center justify-center rounded-[8px] bg-[var(--brand-soft-strong)]"
        >
          <span :class="[action.icon, 'text-[16px] text-[var(--brand)]']"></span>
        </span>
        <span class="flex min-w-0 flex-col gap-[2px]">
          <span class="truncate text-[12.5px] font-medium text-[var(--text-primary)]">
            {{ action.title }}
          </span>
          <span class="truncate text-[10.5px] text-[var(--text-faint)]">
            {{ action.description }}
          </span>
        </span>
      </button>
    </div>
  </section>
</template>