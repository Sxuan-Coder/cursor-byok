import { computed, ref, watch } from "vue";
import { appState } from "@/state/appState";

// 后端 modelchannel.BuildChannelID 使用 sha256(baseURL\nmodelID\napiKey\nname[\nendpoint])
// 前 16 位十六进制作为渠道 ID；此处复刻同一算法，用于识别“使用中”的模型。
const CHANNEL_ID_HEX_LENGTH = 16;

function trim(value) {
  return String(value ?? "").trim();
}

async function computeChannelID(adapter) {
  const parts = [
    trim(adapter.baseURL),
    trim(adapter.modelID),
    trim(adapter.apiKey),
    trim(adapter.displayName),
  ];
  const endpoint = trim(adapter.openAIEndpoint);
  if (endpoint) {
    parts.push(endpoint);
  }
  const subtle = globalThis.crypto?.subtle;
  if (!subtle || typeof TextEncoder === "undefined") {
    return "";
  }
  const digest = await subtle.digest(
    "SHA-256",
    new TextEncoder().encode(parts.join("\n")),
  );
  return Array.from(new Uint8Array(digest))
    .map((byte) => byte.toString(16).padStart(2, "0"))
    .join("")
    .slice(0, CHANNEL_ID_HEX_LENGTH);
}

const agentAdapterId = ref("");
let computeToken = 0;

async function resolveAgentAdapterId() {
  const token = (computeToken += 1);
  const hash = trim(appState.lastAgentModelHash);
  if (!hash) {
    agentAdapterId.value = "";
    return;
  }
  for (const adapter of appState.modelAdapters) {
    const channelID = await computeChannelID(adapter);
    if (token !== computeToken) {
      return;
    }
    if (channelID && channelID === hash) {
      agentAdapterId.value = adapter.id;
      return;
    }
  }
  if (token === computeToken) {
    agentAdapterId.value = "";
  }
}

watch(
  () => [appState.lastAgentModelHash, appState.modelAdapters],
  () => {
    void resolveAgentAdapterId();
  },
  { deep: true, immediate: true },
);

// useActiveModelIds 返回当前被指派为 Agent 或 Commit 用途的模型适配器 ID 集合。
export function useActiveModelIds() {
  return computed(() => {
    const ids = new Set();
    const commitId = trim(appState.commitModelHash);
    if (commitId) {
      ids.add(commitId);
    }
    if (agentAdapterId.value) {
      ids.add(agentAdapterId.value);
    }
    return ids;
  });
}