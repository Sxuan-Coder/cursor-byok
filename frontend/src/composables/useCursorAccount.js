import {
  disconnectCursorAccount,
  getCursorAccountStatus,
  startCursorAccountLogin,
} from "@/services/clientApi";
import { toUserError } from "@/state/appState";
import { computed, onMounted, ref, watch } from "vue";

// Cursor 控制面账号是全局单例状态：侧边栏与设置页共用同一份数据与同一次请求。
const SIGNED_OUT_STATUS = { state: "signed_out", authId: "", email: "", error: "" };
const WAITING_POLL_MS = 1500;

const status = ref({ ...SIGNED_OUT_STATUS });
const busy = ref(false);
let loadPromise = null;
let pollTimer = null;

function normalizeStatus(payload) {
  const value = payload ?? {};
  return {
    state: String(value.state || SIGNED_OUT_STATUS.state),
    authId: String(value.authId || ""),
    email: String(value.email || ""),
    error: String(value.error || ""),
  };
}

const state = computed(() => status.value.state);
const signedIn = computed(() => state.value === "signed_in");
const waiting = computed(() => state.value === "waiting");

const identifier = computed(() => {
  if (!signedIn.value) {
    return "";
  }
  return String(status.value.email || status.value.authId || "").trim();
});

function maskIdentifier(value) {
  const raw = String(value || "").trim();
  if (!raw) {
    return "";
  }

  const atIndex = raw.indexOf("@");
  if (atIndex > 0 && atIndex < raw.length - 1) {
    const localPart = raw.slice(0, atIndex);
    const domain = raw.slice(atIndex + 1);
    const maskedLocalPart =
      localPart.length <= 2
        ? `${localPart[0]}***`
        : `${localPart[0]}***${localPart.at(-1)}`;
    return `${maskedLocalPart}@${domain}`;
  }

  if (raw.length <= 8) {
    return "****";
  }
  return `${raw.slice(0, 4)}****${raw.slice(-4)}`;
}

const maskedIdentifier = computed(() => maskIdentifier(identifier.value));

const stateText = computed(() => {
  if (signedIn.value) {
    return "已经登录";
  }
  if (waiting.value) {
    return "等待浏览器登录";
  }
  return "未连接";
});

const displayName = computed(() => {
  if (!identifier.value) {
    return "本地用户";
  }
  const atIndex = identifier.value.indexOf("@");
  return atIndex > 0 ? identifier.value.slice(0, atIndex) : identifier.value;
});

const initial = computed(() => (displayName.value.slice(0, 1) || "?").toUpperCase());

async function refresh() {
  try {
    status.value = normalizeStatus(await getCursorAccountStatus());
  } catch (_error) {
    status.value = { ...SIGNED_OUT_STATUS };
  }
  return status.value;
}

function ensureLoaded() {
  loadPromise ??= refresh();
  return loadPromise;
}

async function runAction(runner) {
  busy.value = true;
  try {
    status.value = normalizeStatus(await runner());
    await refresh();
    return { ok: true, error: null };
  } catch (error) {
    await refresh();
    return { ok: false, error: toUserError(error) };
  } finally {
    busy.value = false;
  }
}

function login() {
  return runAction(() => startCursorAccountLogin());
}

function disconnect() {
  return runAction(() => disconnectCursorAccount());
}

// 等待浏览器完成登录期间轮询状态，登录成功后停止。
watch(
  waiting,
  (isWaiting) => {
    if (!isWaiting && pollTimer) {
      window.clearInterval(pollTimer);
      pollTimer = null;
      return;
    }
    if (isWaiting && !pollTimer) {
      pollTimer = window.setInterval(() => {
        void refresh();
      }, WAITING_POLL_MS);
    }
  },
  { immediate: true },
);

export function useCursorAccount() {
  onMounted(() => {
    void ensureLoaded();
  });

  return {
    status,
    busy,
    signedIn,
    waiting,
    identifier,
    maskedIdentifier,
    stateText,
    displayName,
    initial,
    refresh,
    login,
    disconnect,
  };
}