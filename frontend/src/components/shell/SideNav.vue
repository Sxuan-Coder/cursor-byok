<script setup>
import AppLogo from "@/components/brand/AppLogo.vue";
import { useCursorAccount } from "@/composables/useCursorAccount";
import { showModal } from "@/composables/useModal";
import { useMessage } from "@/composables/useMessage";
import { appState, checkForAppUpdates } from "@/state/appState";
import { autoUpdate, computePosition, flip, offset, shift } from "@floating-ui/dom";
import { computed, nextTick, ref, watch, watchPostEffect } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const message = useMessage();

const NAV_ITEMS = [
  { key: "overview", label: "概览", icon: "icon-[mdi--view-dashboard-outline]", path: "/" },
  { key: "records", label: "记录", icon: "icon-[mdi--text-box-outline]", path: "/records" },
  { key: "models", label: "模型", icon: "icon-[mdi--view-grid-outline]", path: "/model-config" },
  { key: "usage", label: "用量", icon: "icon-[mdi--chart-donut]", path: "/usage" },
  { key: "settings", label: "设置", icon: "icon-[mdi--cog-outline]", path: "/settings" },
];

const {
  status,
  signedIn,
  waiting,
  busy,
  maskedIdentifier,
  identifier,
  stateText,
  displayName,
  initial: avatarInitial,
  login,
  disconnect,
} = useCursorAccount();

const versionLabel = computed(() => `v${appState.appVersion || "..."}`);

const accountOpen = ref(false);
const accountTriggerRef = ref(null);
const accountPanelRef = ref(null);
const accountPanelStyle = ref({});

function isActive(path) {
  return route.path === path;
}

function updateAccountPanelPosition() {
  if (!accountTriggerRef.value || !accountPanelRef.value) {
    return;
  }

  computePosition(accountTriggerRef.value, accountPanelRef.value, {
    placement: "top-start",
    middleware: [offset(8), flip({ padding: 8 }), shift({ padding: 8 })],
  }).then(({ x, y }) => {
    accountPanelStyle.value = { left: `${x}px`, top: `${y}px` };
  });
}

function toggleAccountPanel() {
  accountOpen.value = !accountOpen.value;
}

function closeAccountPanel() {
  accountOpen.value = false;
}

function handleDocumentPointerDown(event) {
  const target = event.target;
  if (accountTriggerRef.value?.contains(target) || accountPanelRef.value?.contains(target)) {
    return;
  }
  closeAccountPanel();
}

function handleDocumentKeydown(event) {
  if (event.key === "Escape") {
    closeAccountPanel();
  }
}

async function handleLogin() {
  const result = await login();
  if (!result.ok) {
    message(`登录失败：${result.error}`);
  }
}

async function handleDisconnect() {
  const confirmed = await showModal({
    title: "退出登录",
    content: "只会退出 cursor-byok 中的 Cursor 账号，不会退出 Cursor 客户端。是否继续？",
    confirmText: "退出登录",
    cancelText: "取消",
    showCancel: true,
  });
  if (!confirmed) {
    return;
  }

  const result = await disconnect();
  if (!result.ok) {
    message(`退出登录失败：${result.error}`);
  }
}

watch(accountOpen, (open) => {
  if (open) {
    document.addEventListener("pointerdown", handleDocumentPointerDown, true);
    document.addEventListener("keydown", handleDocumentKeydown, true);
    nextTick(() => updateAccountPanelPosition());
    return;
  }

  document.removeEventListener("pointerdown", handleDocumentPointerDown, true);
  document.removeEventListener("keydown", handleDocumentKeydown, true);
});

watchPostEffect((cleanup) => {
  if (!accountOpen.value || !accountTriggerRef.value || !accountPanelRef.value) {
    return;
  }

  const stop = autoUpdate(accountTriggerRef.value, accountPanelRef.value, updateAccountPanelPosition);
  cleanup(() => {
    stop();
  });
});
</script>

<template>
  <nav class="flex h-full w-[var(--sidebar-w)] shrink-0 flex-col border-r border-[var(--border-subtle)] bg-[var(--bg-sidebar)] px-3 pb-3 pt-3">
    <div class="flex items-center gap-2 px-1 pb-3">
      <span class="flex h-[30px] w-[30px] shrink-0 items-center justify-center rounded-[9px] border border-[var(--brand-border)] bg-[var(--brand-soft)]">
        <AppLogo :size="17" />
      </span>
      <span class="flex min-w-0 flex-col leading-tight">
        <span class="truncate text-[13px] font-semibold text-[var(--text-primary)]">
          Cursor助手
        </span>
        <span class="truncate text-[9.5px] text-[var(--text-faint)]">
          BYOK 多模型代理客户端
        </span>
      </span>
    </div>

    <ul class="flex flex-col gap-1">
      <li v-for="item in NAV_ITEMS" :key="item.key">
        <button
          type="button"
          class="relative flex w-full items-center gap-2.5 rounded-[9px] px-3 py-[9px] text-left text-[13px] transition-colors"
          :class="
            isActive(item.path)
              ? 'bg-[var(--brand-soft-strong)] font-medium text-[var(--brand)]'
              : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]'
          "
          :aria-current="isActive(item.path) ? 'page' : undefined"
          @click="router.push(item.path)"
        >
          <span
            v-if="isActive(item.path)"
            class="absolute left-0 top-1/2 h-[16px] w-[3px] -translate-y-1/2 rounded-full bg-[var(--brand)]"
          ></span>
          <span :class="[item.icon, 'shrink-0 text-[17px]']"></span>
          <span class="truncate">{{ item.label }}</span>
        </button>
      </li>
    </ul>

    <div class="mt-auto flex flex-col gap-2 pt-3">
      <button
        ref="accountTriggerRef"
        type="button"
        class="flex w-full items-center gap-2.5 rounded-[9px] border px-1.5 py-1.5 text-left transition-colors"
        :class="
          accountOpen
            ? 'border-[var(--brand-border)] bg-[var(--brand-soft)]'
            : 'border-transparent hover:border-[var(--border)] hover:bg-[var(--bg-hover)]'
        "
        :aria-expanded="accountOpen"
        @click="toggleAccountPanel"
      >
        <span
          class="flex h-[26px] w-[26px] shrink-0 items-center justify-center rounded-full bg-[var(--brand-soft-strong)] text-[11.5px] font-semibold text-[var(--brand)]"
        >
          {{ avatarInitial }}
        </span>
        <span class="min-w-0 flex-1 truncate text-[12.5px] text-[var(--text-primary)]">
          {{ displayName }}
        </span>
        <span
          class="icon-[mdi--chevron-down] shrink-0 text-[15px] text-[var(--text-faint)] transition-transform"
          :class="{ 'rotate-180': accountOpen }"
        ></span>
      </button>

      <button
        type="button"
        class="flex items-center gap-1.5 rounded-[7px] px-1 py-0.5 text-[11px] text-[var(--text-faint)] transition-colors hover:text-[var(--text-secondary)] disabled:opacity-60"
        :disabled="appState.updateState === 'checking'"
        @click="checkForAppUpdates()"
      >
        <span>{{ versionLabel }}</span>
        <span>检查更新</span>
      </button>
    </div>

    <Teleport to="body">
      <div
        v-if="accountOpen"
        ref="accountPanelRef"
        class="fixed z-[9999] flex w-[286px] flex-col gap-3 rounded-[var(--radius-md)] border border-[var(--border-strong)] bg-[var(--bg-elevated)] p-3.5 shadow-[var(--shadow-pop)]"
        :style="accountPanelStyle"
      >
        <div class="flex items-center justify-between gap-2">
          <h3 class="text-[13px] font-medium text-[var(--text-primary)]">Cursor 控制面账号</h3>
          <span
            class="rounded-full border px-2 py-[2px] text-[10.5px]"
            :class="
              signedIn
                ? 'border-[var(--brand-border)] text-[var(--brand)]'
                : 'border-[var(--border-strong)] text-[var(--text-secondary)]'
            "
          >
            {{ stateText }}
          </span>
        </div>

        <div class="flex flex-col gap-2 rounded-[var(--radius-sm)] border border-[var(--border-subtle)] bg-[var(--bg-inset)] p-2.5">
          <div v-if="identifier" class="truncate text-[12px] text-[var(--text-secondary)]" :title="identifier">
            {{ maskedIdentifier }}
          </div>
          <p class="text-[11.5px] leading-[17px] text-[var(--text-secondary)]">
            登录后可用插件市场、Skills 与 MCP 等控制面能力；该账号独立于 Cursor 客户端当前账号，不影响 Cursor 本身的登录状态。
          </p>
          <p v-if="waiting" class="text-[11.5px] leading-[17px] text-[var(--warn)]">
            请在浏览器完成登录，完成后返回本应用即可自动生效。
          </p>
          <p v-if="status.error" class="break-all text-[11.5px] leading-[17px] text-[var(--danger)]">
            {{ status.error }}
          </p>
        </div>

        <button
          v-if="signedIn"
          type="button"
          class="flex h-[32px] w-full items-center justify-center rounded-[var(--radius-sm)] border border-[var(--border-strong)] text-[12.5px] text-[var(--text-primary)] transition-colors hover:border-[var(--danger)] hover:text-[var(--danger)] disabled:opacity-60"
          :disabled="busy"
          @click="handleDisconnect"
        >
          退出登录
        </button>
        <button
          v-else
          type="button"
          class="flex h-[32px] w-full items-center justify-center gap-1.5 rounded-[var(--radius-sm)] bg-[var(--brand)] text-[12.5px] font-medium text-[#0b1a12] transition-colors hover:bg-[var(--brand-hover)] disabled:opacity-60"
          :disabled="busy || waiting"
          @click="handleLogin"
        >
          <span class="icon-[mdi--login-variant] text-[15px]"></span>
          <span>{{ waiting ? "等待登录..." : "登录 Cursor" }}</span>
        </button>
      </div>
    </Teleport>
  </nav>
</template>