<script setup>
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import CommitMessageCard from "@/components/CommitMessageCard.vue";
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
import { computed, nextTick, onMounted } from "vue";
import { useRoute } from "vue-router";

const AUTHOR_REPOSITORY_URL = "https://github.com/Sxuan-Coder/cursor-byok";
const AUTHOR_PULL_REQUEST_URL = `${AUTHOR_REPOSITORY_URL}/pulls`;
const AUTHOR_LABEL = "@leookun & @Sxuan-Coder";

// 「模型配置」入口打开的是独立子窗口，模型管理走侧栏 Tab 即可，此卡片暂时隐藏。
const SHOW_MODEL_CONFIG_CARD = false;

const route = useRoute();
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

async function handleOpenGitHubStar() {
  try {
    await Browser.OpenURL(AUTHOR_REPOSITORY_URL);
  } catch (error) {
    showActionError("打开失败", toUserError(error));
  }
}

async function handleOpenPullRequests() {
  try {
    await Browser.OpenURL(AUTHOR_PULL_REQUEST_URL);
  } catch (error) {
    showActionError("打开失败", toUserError(error));
  }
}

// 首页快捷操作会带 ?focus=commit 跳转过来，定位滚动到 Commit 配置卡片。
async function focusCommitCard() {
  if (route.query.focus !== "commit") {
    return;
  }
  await nextTick();
  document.getElementById("commit-prompt-card")?.scrollIntoView({
    behavior: "smooth",
    block: "center",
  });
}

onMounted(async () => {
  void focusCommitCard();
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

    <Card v-if="SHOW_MODEL_CONFIG_CARD">
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

    <!-- 多根组件无法透传 id，用外层节点承载锚点 -->
    <div id="commit-prompt-card">
      <CommitMessageCard />
    </div>

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

    <Card>
      <div class="min-w-0">
        <h2 class="text-[13.5px] font-medium text-[var(--text-primary)]">支持本项目</h2>
        <div class="mt-1 text-[12px] text-[var(--text-secondary)]">
          如果这个工具帮到了你，欢迎点亮 Star，或提交 PR 一起共建
        </div>
      </div>

      <div class="mt-3 grid grid-cols-2 gap-2 max-[760px]:grid-cols-1">
        <div
          class="flex items-center justify-between gap-3 rounded-[10px] border border-[var(--border)] bg-[var(--bg-card-soft)] p-3"
        >
          <div class="flex min-w-0 items-center gap-2.5">
            <span
              class="flex h-[34px] w-[34px] shrink-0 items-center justify-center rounded-[8px] bg-[var(--brand-soft-strong)]"
            >
              <span class="icon-[mdi--star-outline] text-[17px] text-[var(--brand)]"></span>
            </span>
            <div class="flex min-w-0 flex-col gap-[2px]">
              <span class="truncate text-[12.5px] font-medium text-[var(--text-primary)]">
                为仓库点 Star
              </span>
              <span class="truncate text-[10.5px] text-[var(--text-faint)]">
                支持作者持续维护本项目
              </span>
            </div>
          </div>
          <Button variant="primary" @click="handleOpenGitHubStar">前往 Star</Button>
        </div>

        <div
          class="flex items-center justify-between gap-3 rounded-[10px] border border-[var(--border)] bg-[var(--bg-card-soft)] p-3"
        >
          <div class="flex min-w-0 items-center gap-2.5">
            <span
              class="flex h-[34px] w-[34px] shrink-0 items-center justify-center rounded-[8px] bg-[var(--brand-soft-strong)]"
            >
              <span class="icon-[mdi--source-pull] text-[17px] text-[var(--brand)]"></span>
            </span>
            <div class="flex min-w-0 flex-col gap-[2px]">
              <span class="truncate text-[12.5px] font-medium text-[var(--text-primary)]">
                提交 PR
              </span>
              <span class="truncate text-[10.5px] text-[var(--text-faint)]">
                欢迎通过 Pull Request 参与共建
              </span>
            </div>
          </div>
          <Button variant="default" @click="handleOpenPullRequests">提交 PR</Button>
        </div>
      </div>
    </Card>
  </div>
</template>