<script setup>
import { useCursorAccount } from "@/composables/useCursorAccount";
import { showModal } from "@/composables/useModal";
import { useMessage } from "@/composables/useMessage";
import { autoUpdate, computePosition, flip, offset, shift } from "@floating-ui/dom";
import { nextTick, ref, watch, watchPostEffect } from "vue";

// 侧边栏与顶栏各需要一个账号入口，浮层内容与状态完全一致，这里统一收口。
const props = defineProps({
  placement: { type: String, default: "top-start" },
  variant: { type: String, default: "row" },
  // 触发按钮的悬停提示，用于在顶栏这种窄空间里补充说明登录收益。
  title: { type: String, default: "" },
});

const message = useMessage();

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

const open = ref(false);
const triggerRef = ref(null);
const panelRef = ref(null);
const panelStyle = ref({});

function updatePanelPosition() {
  if (!triggerRef.value || !panelRef.value) {
    return;
  }

  computePosition(triggerRef.value, panelRef.value, {
    placement: props.placement,
    middleware: [offset(8), flip({ padding: 8 }), shift({ padding: 8 })],
  }).then(({ x, y }) => {
    panelStyle.value = { left: `${x}px`, top: `${y}px` };
  });
}

function toggle() {
  open.value = !open.value;
}

function close() {
  open.value = false;
}

function handleDocumentPointerDown(event) {
  const target = event.target;
  if (triggerRef.value?.contains(target) || panelRef.value?.contains(target)) {
    return;
  }
  close();
}

function handleDocumentKeydown(event) {
  if (event.key === "Escape") {
    close();
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

watch(open, (isOpen) => {
  if (isOpen) {
    document.addEventListener("pointerdown", handleDocumentPointerDown, true);
    document.addEventListener("keydown", handleDocumentKeydown, true);
    nextTick(() => updatePanelPosition());
    return;
  }

  document.removeEventListener("pointerdown", handleDocumentPointerDown, true);
  document.removeEventListener("keydown", handleDocumentKeydown, true);
});

watchPostEffect((cleanup) => {
  if (!open.value || !triggerRef.value || !panelRef.value) {
    return;
  }

  const stop = autoUpdate(triggerRef.value, panelRef.value, updatePanelPosition);
  cleanup(() => {
    stop();
  });
});
</script>

<template>
  <button
    ref="triggerRef"
    type="button"
    class="no-drag-region transition-colors"
    :title="props.title || undefined"
    :class="
      variant === 'pill'
        ? [
            'flex shrink-0 items-center gap-1.5 rounded-[8px] border px-2.5 py-[5px] text-[11.5px]',
            open
              ? 'border-[var(--brand-border)] bg-[var(--brand-soft-strong)] text-[var(--brand)]'
              : 'border-[var(--brand-border)] bg-[var(--brand-soft)] text-[var(--brand)] hover:bg-[var(--brand-soft-strong)]',
          ]
        : [
            'flex w-full items-center gap-2.5 rounded-[9px] border px-1.5 py-1.5 text-left',
            open
              ? 'border-[var(--brand-border)] bg-[var(--brand-soft)]'
              : 'border-transparent hover:border-[var(--border)] hover:bg-[var(--bg-hover)]',
          ]
    "
    :aria-expanded="open"
    @click="toggle"
  >
    <template v-if="variant === 'pill'">
      <span class="icon-[mdi--shield-account-outline] shrink-0 text-[14px]"></span>
      <span class="min-w-0 max-w-[132px] truncate">登录解锁插件</span>
    </template>

    <template v-else>
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
        :class="{ 'rotate-180': open }"
      ></span>
    </template>
  </button>

  <Teleport to="body">
    <div
      v-if="open"
      ref="panelRef"
      class="fixed z-[9999] flex w-[286px] flex-col gap-3 rounded-[var(--radius-md)] border border-[var(--border-strong)] bg-[var(--bg-elevated)] p-3.5 shadow-[var(--shadow-pop)]"
      :style="panelStyle"
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

      <div
        class="flex flex-col gap-2 rounded-[var(--radius-sm)] border border-[var(--border-subtle)] bg-[var(--bg-inset)] p-2.5"
      >
        <div
          v-if="identifier"
          class="truncate text-[12px] text-[var(--text-secondary)]"
          :title="identifier"
        >
          {{ maskedIdentifier }}
        </div>
        <p class="text-[11.5px] leading-[17px] text-[var(--text-secondary)]">
          登录后可获得 Cursor 插件 Skill 以及 MCP 商店的官方权限；该账号独立于 Cursor
          客户端当前账号，不影响 Cursor 本身的登录状态。
        </p>
        <p v-if="waiting" class="text-[11.5px] leading-[17px] text-[var(--warn)]">
          请在浏览器完成登录，完成后返回本应用即可自动生效。
        </p>
        <p
          v-if="status.error"
          class="break-all text-[11.5px] leading-[17px] text-[var(--danger)]"
        >
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
</template>