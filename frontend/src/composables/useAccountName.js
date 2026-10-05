import { useCursorAccount } from "@/composables/useCursorAccount";

// useAccountName 复用侧边栏与横幅的账号展示名，来源为已登录 Cursor 账号邮箱前缀。
export function useAccountName() {
  const { displayName, initial, refresh } = useCursorAccount();

  return { displayName, initial, refresh };
}