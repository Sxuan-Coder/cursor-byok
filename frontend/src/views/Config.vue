<script setup>
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import CommitMessageCard from "@/components/CommitMessageCard.vue";
import LocaleSelect from "@/components/LocaleSelect.vue";
import { useMessage } from "@/composables/useMessage";
import {
  appState,
  checkForAppUpdates,
  openConfigWindow,
  openLocalLogsDirectory,
  openModelConfigWindow,
  persistUserConfig,
  reloadUserConfig,
  toUserError,
  updateViewState,
} from "@/state/appState";
import { Browser } from "@wailsio/runtime";
import { computed, onMounted } from "vue";

const AUTHOR_REPOSITORY_URL = "https://github.com/Sxuan-Coder/cursor-byok";
const AUTHOR_LABEL = "@leookun & @Sxuan-Coder";

const message = useMessage();

const versionLabel = computed(() => `v${appState.appVersion || "..."}`);

function showActionError(title, error) {
  const detail = String(error || "服务错误").trim() || "服务错误";
  message(`${title}：${detail}`);
}

async function handleSaveConfig() {
  const result = await persistUserConfig();
  if (!result.ok) {
    showActionError("保存失败", result.error);
    return;
  }
  message("本地配置已保存");
}

async function handleOpenModelConfig() {
  try {
    await openModelConfigWindow();
  } catch (error) {
    showActionError("打开失败", toUserError(error));
  }
}

async function handleOpenSettingsDirectory() {
  try {
    await openConfigWindow();
  } catch (error) {
    showActionError("打开失败", toUserError(error));
  }
}

async function handleOpenLogsDirectory() {
  try {
    await openLocalLogsDirectory();
  } catch (error) {
    showActionError("打开失败", toUserError(error));
  }
}

async function handleCheckForUpdates() {
  if (updateViewState.footerBusy || updateViewState.footerDownloading) {
    return;
  }
  try {
    await checkForAppUpdates();
  } catch (error) {
    showActionError("检查更新失败", toUserError(error));
  }
}

async function handleOpenAuthorHome() {
  try {
    await Browser.OpenURL(AUTHOR_REPOSITORY_URL);
  } catch (error) {
    showActionError("打开作者地址失败", error);
  }
}

onMounted(async () => {
  await reloadUserConfig().catch(() => {});
});
</script>

<template>
  <div class="flex h-full min-h-0 flex-col gap-4 overflow-y-auto px-4 pb-4 pt-4">
    <Card>
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="min-w-0">
          <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">本地配置</h2>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
            运行日志位于 <code class="text-[var(--text-muted)]">~/.cursor-local-assistant-v2/logs/</code>
          </div>
        </div>
        <Button variant="primary" :disabled="appState.configSaving" @click="handleSaveConfig">
          {{ appState.configSaving ? "保存中..." : "保存配置" }}
        </Button>
      </div>
    </Card>

    <Card>
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="min-w-0">
          <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">界面语言</h2>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
            切换当前界面显示语言，设置会立即生效并保存在本机
          </div>
        </div>
        <LocaleSelect wrapper-class="w-[220px] max-w-full" />
      </div>
    </Card>

    <Card>
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="min-w-0">
          <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">模型配置</h2>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
            已配置 {{ appState.modelAdapters.length }} 个模型适配器
          </div>
        </div>
        <Button variant="primary" @click="handleOpenModelConfig">打开模型配置</Button>
      </div>
    </Card>

    <CommitMessageCard />

    <Card>
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="min-w-0">
          <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">目录与日志</h2>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
            打开本地设置目录或运行日志目录，便于排查问题
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="default" @click="handleOpenLogsDirectory">打开日志目录</Button>
          <Button variant="default" @click="handleOpenSettingsDirectory">打开设置目录</Button>
        </div>
      </div>
    </Card>

    <Card>
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="min-w-0">
          <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">关于</h2>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
            当前版本 <span class="font-num">{{ versionLabel }}</span> · BYOK 永久免费
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="default" @click="handleOpenAuthorHome">
            <span class="icon-[mdi--github] text-[15px]"></span>
            <span>{{ AUTHOR_LABEL }}</span>
          </Button>
          <Button
            variant="primary"
            :disabled="updateViewState.footerBusy || updateViewState.footerDownloading"
            @click="handleCheckForUpdates"
          >
            {{ updateViewState.footerDownloading ? updateViewState.footerProgressText : "检查更新" }}
          </Button>
        </div>
      </div>
    </Card>
  </div>
</template>