<script setup>
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import { useCursorAccount } from "@/composables/useCursorAccount";
import { showModal } from "@/composables/useModal";
import { useMessage } from "@/composables/useMessage";

const message = useMessage();

const {
  status,
  busy,
  signedIn,
  waiting,
  maskedIdentifier,
  stateText,
  login,
  disconnect,
} = useCursorAccount();

async function handleCursorAccountLogin() {
  const result = await login();
  if (!result.ok) {
    message(`登录失败：${result.error}`);
  }
}

async function handleCursorAccountDisconnect() {
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
</script>

<template>
  <Card>
    <div class="flex flex-col gap-3">
      <div class="flex items-center justify-between gap-4">
        <div class="flex min-w-0 flex-wrap items-center gap-2">
          <h2 class="text-base font-medium text-white">Cursor 控制面账号</h2>
          <span
            class="rounded-full border border-[#3a3a3a] bg-[#202020] px-2 py-0.5 text-xs text-[#b8b8b8]"
          >
            {{ stateText }}
          </span>
        </div>
      </div>

      <div class="flex items-end justify-between gap-4">
        <div class="min-w-0">
          <div v-if="maskedIdentifier" class="truncate text-sm text-[#d0d0d0]">
            {{ maskedIdentifier }}
          </div>
          <div class="mt-1 text-sm text-[#a3a3a3]">
            独立用于插件、Skills 和 MCP；不会改变 Cursor 客户端当前账号
          </div>
          <div v-if="waiting" class="mt-1 text-sm text-[#d6a84b]">
            请在浏览器完成登录，完成后返回 Cursor 重新打开插件市场
          </div>
          <div v-if="status.error" class="mt-1 break-all text-sm text-[#e06c75]">
            {{ status.error }}
          </div>
        </div>
        <Button
          v-if="signedIn"
          class="shrink-0"
          :disabled="busy"
          @click="handleCursorAccountDisconnect"
        >
          退出登录
        </Button>
        <Button
          v-else
          class="shrink-0"
          variant="primary"
          :disabled="busy || waiting"
          @click="handleCursorAccountLogin"
        >
          {{ waiting ? "等待登录..." : "登录 Cursor" }}
        </Button>
      </div>
    </div>
  </Card>
</template>