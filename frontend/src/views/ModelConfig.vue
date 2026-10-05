<script setup>
import Button from "@/components/ui/Button.vue";
import ContentModal from "@/components/ui/ContentModal.vue";
import ModelEditor from "@/components/ModelEditor.vue";
import ModelPickerModal from "@/components/ModelPickerModal.vue";
import ProviderEditorModal from "@/components/ProviderEditorModal.vue";
import ProviderGroupSection from "@/components/ProviderGroupSection.vue";
import { useMessage } from "@/composables/useMessage";
import { useConfigTransfer } from "@/composables/useConfigTransfer";
import {
  appState,
  createEmptyModelAdapter,
  deleteModelAdapterAt,
  deleteProviderAdapters,
  duplicateModelAdapterAt,
  reloadUserConfig,
  runModelAdapterTest,
  saveModelAdapterOrder,
  startModelAdapterTest,
  toUserError,
} from "@/state/appState";
import { buildProviderGroupKey, buildProviderGroups } from "@/state/providerGroups";
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

const BATCH_TEST_CONCURRENCY = 10;
const message = useMessage();
const route = useRoute();
const router = useRouter();

const typeTabs = [
  { label: "全部", value: "all", icon: "icon-[mdi--view-grid-outline]" },
  { label: "OpenAI", value: "openai", icon: "icon-[bxl--openai]" },
  { label: "Anthropic", value: "anthropic", icon: "icon-[logos--claude-icon]" },
];

const activeType = ref("all");
const batchTesting = ref(false);
const batchStopping = ref(false);
const batchTotal = ref(0);
const batchCompleted = ref(0);
const editorOpen = ref(false);
const editorIndex = ref(-1);
const editorAdapter = ref(null);
const editorSession = ref(0);
const sortSaving = ref(false);
const pickerOpen = ref(false);
const pickerSession = ref(0);
const pickerProvider = ref(null);
const pickerAdapters = ref([]);
const providerEditorOpen = ref(false);
const providerEditorSession = ref(0);
const providerEditorMode = ref("create");
const providerEditorProvider = ref(null);
const batchActiveCalls = new Set();
let batchStopRequested = false;

const filteredAdapters = computed(() => (
  activeType.value === "all"
    ? appState.modelAdapters
    : appState.modelAdapters.filter((adapter) => adapter.type === activeType.value)
));
const providerGroups = computed(() => buildProviderGroups(filteredAdapters.value));
const groupActionsDisabled = computed(() => (
  sortSaving.value || appState.configSaving || batchTesting.value
));
const batchButtonText = computed(() => {
  if (batchStopping.value) {
    return "停止中...";
  }
  if (!batchTesting.value) {
    return "测试全部";
  }
  return `停止测试 ${batchCompleted.value}/${batchTotal.value}`;
});
const editorTitle = computed(() => (editorIndex.value >= 0 ? "编辑模型配置" : "新增模型配置"));
const pickerTitle = computed(() => `拉取模型 · ${formatHost(pickerProvider.value?.baseURL)}`);
const providerEditorTitle = computed(() => (
  providerEditorMode.value === "edit" ? "编辑提供商凭证" : "新增提供商"
));
const emptyStateText = computed(() => (
  activeType.value === "all"
    ? "当前还没有配置任何模型。"
    : `当前还没有配置任何 ${typeLabel(activeType.value)} 模型。`
));

watch(
  () => appState.modelAdapters,
  (adapters) => {
    if (
      activeType.value === "all"
      || adapters.some((adapter) => adapter.type === activeType.value)
    ) {
      return;
    }
    activeType.value = "all";
  },
  { deep: true, immediate: true },
);

function showActionError(title, error) {
  const detail = String(error || "服务错误").trim() || "服务错误";
  message(`${title}：${detail}`);
}

const {
  configTransferBusy,
  handleExportConfig,
  handleImportConfig,
} = useConfigTransfer({ message, showActionError });

function typeLabel(type) {
  return type === "anthropic" ? "Anthropic" : "OpenAI";
}

function formatHost(value) {
  const text = String(value || "").trim();
  if (!text) {
    return "-";
  }
  try {
    const parsed = new URL(text);
    return parsed.host || text;
  } catch {
    return text.replace(/^https?:\/\//, "");
  }
}

function openEditor(index = -1, draft = null) {
  editorIndex.value = index;
  if (index >= 0) {
    editorAdapter.value = appState.modelAdapters[index];
  } else if (draft) {
    editorAdapter.value = { ...createEmptyModelAdapter(), ...draft };
  } else {
    editorAdapter.value = {
      ...createEmptyModelAdapter(),
      type: activeType.value === "anthropic" ? "anthropic" : "openai",
    };
  }
  editorSession.value += 1;
  editorOpen.value = true;
}

function closeEditor() {
  if (appState.configSaving) {
    return;
  }
  editorOpen.value = false;
}

function handleEditorSaved(adapter) {
  if (activeType.value !== "all" && adapter?.type) {
    activeType.value = adapter.type;
  }
}

function openPicker({ baseURL, apiKey, type = "openai", providerAdapters = [] }) {
  pickerProvider.value = { baseURL, apiKey, type };
  pickerAdapters.value = providerAdapters;
  pickerSession.value += 1;
  pickerOpen.value = true;
}

function handlePickerAdded(added) {
  message(`已新增 ${added} 个模型`);
}

function handlePickerManual({ type, baseURL, apiKey }) {
  openEditor(-1, { type, baseURL, apiKey });
}

function openProviderEditor(mode, provider = null) {
  providerEditorMode.value = mode;
  providerEditorProvider.value = provider;
  providerEditorSession.value += 1;
  providerEditorOpen.value = true;
}

function closeProviderEditor() {
  if (appState.configSaving) {
    return;
  }
  providerEditorOpen.value = false;
}

function handleProviderCreated({ type, baseURL, apiKey }) {
  openPicker({ baseURL, apiKey, type, providerAdapters: [] });
}

function handleProviderSaved(updated) {
  message(`已更新 ${updated} 个模型的凭证`);
}

function handleGroupFetchModels(group) {
  const preferredType = activeType.value !== "all" && group.types.includes(activeType.value)
    ? activeType.value
    : group.types[0];
  openPicker({
    baseURL: group.baseURL,
    apiKey: group.apiKey,
    type: preferredType || "openai",
    providerAdapters: group.adapters,
  });
}

function handleGroupAddModel(group) {
  const preferredType = activeType.value !== "all" && group.types.includes(activeType.value)
    ? activeType.value
    : group.types[0];
  openEditor(-1, {
    type: preferredType || "openai",
    baseURL: group.baseURL,
    apiKey: group.apiKey,
  });
}

async function handleDeleteProvider(group) {
  const result = await deleteProviderAdapters(group.baseURL, group.apiKey);
  if (!result.ok) {
    showActionError("删除提供商失败", result.error);
  }
}

async function restoreModelOrder(previousAdapters) {
  try {
    await reloadUserConfig({ modelAdaptersOnly: true });
  } catch (_error) {
    appState.modelAdapters = previousAdapters;
  }
}

async function handleGroupReorder({ group, adapters: reorderedAdapters }) {
  const previousAdapters = appState.modelAdapters.slice();
  let cursor = 0;
  const nextAdapters = previousAdapters
    .map((adapter) => {
      if (buildProviderGroupKey(adapter) !== group.key) {
        return adapter;
      }
      const nextAdapter = reorderedAdapters[cursor];
      cursor += 1;
      return nextAdapter;
    })
    .map((adapter, index) => ({
      ...adapter,
      sort: index + 1,
    }));
  if (cursor !== group.adapters.length) {
    return;
  }

  sortSaving.value = true;
  appState.modelAdapters = nextAdapters;
  try {
    const result = await saveModelAdapterOrder(nextAdapters.map((adapter) => adapter.id));
    if (!result.ok) {
      await restoreModelOrder(previousAdapters);
      showActionError("排序失败", result.error);
    }
  } catch (error) {
    await restoreModelOrder(previousAdapters);
    showActionError("排序失败", toUserError(error));
  } finally {
    sortSaving.value = false;
  }
}

async function handleTestModelAdapter(adapter) {
  try {
    await runModelAdapterTest(adapter);
  } catch (_error) {
    // 失败结果会通过事件同步到界面，这里不再额外弹窗打断用户。
  }
}

async function handleEditModelAdapter(adapter) {
  openEditor(appState.modelAdapters.indexOf(adapter));
}

async function handleDuplicateModelAdapter(adapter) {
  const index = appState.modelAdapters.indexOf(adapter);
  if (index < 0) {
    showActionError("复制失败", "模型配置不存在，无法复制");
    return;
  }
  const result = await duplicateModelAdapterAt(index);
  if (!result.ok) {
    showActionError("复制失败", result.error);
  }
}

async function handleDeleteModelAdapter(adapter) {
  const index = appState.modelAdapters.indexOf(adapter);
  if (index < 0) {
    showActionError("删除失败", "模型配置不存在，无法删除");
    return;
  }
  const result = await deleteModelAdapterAt(index);
  if (!result.ok) {
    showActionError("删除失败", result.error);
  }
}

function isCancelError(error) {
  return String(error?.name || "").trim() === "CancelError";
}

async function stopBatchTesting() {
  if (!batchTesting.value || batchStopping.value) {
    return;
  }
  batchStopRequested = true;
  batchStopping.value = true;
  const activeCalls = Array.from(batchActiveCalls);
  await Promise.allSettled(
    activeCalls.map((call) => (typeof call?.cancel === "function" ? call.cancel("batch-stop") : undefined)),
  );
}

async function handleTestAllModelAdapters() {
  if (batchTesting.value) {
    await stopBatchTesting();
    return;
  }
  const adapters = filteredAdapters.value.slice();
  if (adapters.length === 0) {
    return;
  }
  batchStopRequested = false;
  batchTesting.value = true;
  batchStopping.value = false;
  batchTotal.value = adapters.length;
  batchCompleted.value = 0;
  let nextIndex = 0;
  try {
    const workers = Array.from({ length: Math.min(BATCH_TEST_CONCURRENCY, adapters.length) }, async () => {
      while (!batchStopRequested) {
        const currentIndex = nextIndex;
        nextIndex += 1;
        if (currentIndex >= adapters.length) {
          return;
        }
        const adapter = adapters[currentIndex];
        const call = startModelAdapterTest(adapter);
        batchActiveCalls.add(call);
        try {
          await call;
        } catch (error) {
          if (!isCancelError(error) && !batchStopRequested) {
            // 单个失败结果由卡片自行展示，这里继续后续测试。
          }
        } finally {
          batchActiveCalls.delete(call);
          batchCompleted.value += 1;
        }
      }
    });
    await Promise.allSettled(workers);
  } finally {
    batchActiveCalls.clear();
    batchStopRequested = false;
    batchTesting.value = false;
    batchStopping.value = false;
  }
}

// 概览页的快捷操作通过 query 参数跳转到本页，这里消费一次后立即清理，避免刷新时重复触发。
function consumeQueryActions() {
  const wantsNew = route.query.new === "1";
  const wantsTestAll = route.query.test === "all";
  if (!wantsNew && !wantsTestAll) {
    return;
  }
  if (wantsNew) {
    openProviderEditor("create");
  }
  if (wantsTestAll) {
    void handleTestAllModelAdapters();
  }
  router.replace({ path: route.path });
}

watch(
  () => route.query,
  () => consumeQueryActions(),
);

onMounted(async () => {
  await reloadUserConfig({ modelAdaptersOnly: true }).catch(() => { });
  consumeQueryActions();
});

onBeforeUnmount(() => {
  void stopBatchTesting();
});
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden pt-4 text-[var(--text-primary)]">
    <div class="shrink-0 pb-4">
      <div class="flex items-center justify-between gap-4 px-4">
        <div class="center-row gap-2">
          <button
            v-for="tab in typeTabs"
            :key="tab.value"
            type="button"
            class="center-row gap-2 rounded-[8px] border px-3 py-2 text-sm transition-colors duration-150"
            :class="activeType === tab.value
              ? 'border-[var(--brand-border)] bg-[var(--brand-soft-strong)] text-[var(--brand)]'
              : 'border-[var(--border)] bg-[var(--bg-card-soft)] text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]'"
            @click="activeType = tab.value"
          >
            <span :class="[tab.icon, 'text-[16px]']"></span>
            <span>{{ tab.label }}</span>
          </button>
        </div>
        <div class="center-row gap-2">
          <Button
            variant="default"
            :disabled="sortSaving || appState.configSaving || batchTesting || configTransferBusy || appState.serviceRunning || appState.backendRunning || appState.proxyRunning"
            @click="handleImportConfig"
          >
            导入配置
          </Button>
          <Button
            variant="default"
            :disabled="sortSaving || appState.configSaving || batchTesting || configTransferBusy"
            @click="handleExportConfig"
          >
            导出配置
          </Button>
          <Button
            variant="default"
            :disabled="sortSaving || appState.configSaving || configTransferBusy || (!batchTesting && filteredAdapters.length === 0)"
            @click="handleTestAllModelAdapters"
          >
            {{ batchButtonText }}
          </Button>
          <Button variant="primary" :disabled="sortSaving || appState.configSaving || batchTesting || configTransferBusy" @click="openProviderEditor('create')">
            新增提供商
          </Button>
        </div>
      </div>
    </div>

    <div class="min-h-0 flex-1">
      <div v-if="filteredAdapters.length === 0"
        class="flex h-full min-h-[220px] flex-col items-center justify-center gap-3 rounded-[8px] px-4 text-sm text-[#a3a3a3]">
        <span>{{ emptyStateText }}</span>
        <div class="center-row gap-2">
          <Button variant="primary" :disabled="sortSaving || appState.configSaving || configTransferBusy" @click="openProviderEditor('create')">
            <span class="center-row gap-1.5">
              <span class="icon-[mdi--server-plus] text-[16px]"></span>
              <span>新增提供商</span>
            </span>
          </Button>
          <Button variant="default" :disabled="sortSaving || appState.configSaving || configTransferBusy" @click="openEditor()">手动新增模型</Button>
        </div>
      </div>

      <div v-else class="h-full min-h-0 overflow-y-auto scroll-shadow-bottom p-4 pt-0">
        <div class="flex flex-col gap-4 pb-1">
          <ProviderGroupSection
            v-for="group in providerGroups"
            :key="group.key"
            :group="group"
            :actions-disabled="groupActionsDisabled"
            @fetch-models="handleGroupFetchModels"
            @add-model="handleGroupAddModel"
            @edit-provider="openProviderEditor('edit', { baseURL: $event.baseURL, apiKey: $event.apiKey })"
            @delete-provider="handleDeleteProvider"
            @test="handleTestModelAdapter"
            @edit="handleEditModelAdapter"
            @duplicate="handleDuplicateModelAdapter"
            @delete="handleDeleteModelAdapter"
            @reorder="handleGroupReorder"
          />
        </div>
      </div>
    </div>
  </div>

  <ContentModal
    :open="editorOpen"
    :title="editorTitle"
    size="xl"
    :close-disabled="appState.configSaving"
    @close="closeEditor"
  >
    <ModelEditor
      v-if="editorOpen && editorAdapter"
      :key="editorSession"
      :index="editorIndex"
      :adapter="editorAdapter"
      @saved="handleEditorSaved"
      @close="closeEditor"
    />
  </ContentModal>

  <ContentModal
    :open="pickerOpen"
    :title="pickerTitle"
    size="xl"
    :close-disabled="appState.configSaving"
    @close="pickerOpen = false"
  >
    <ModelPickerModal
      v-if="pickerOpen && pickerProvider"
      :key="pickerSession"
      :baseURL="pickerProvider.baseURL"
      :apiKey="pickerProvider.apiKey"
      :initialType="pickerProvider.type"
      :providerAdapters="pickerAdapters"
      @added="handlePickerAdded"
      @manual="handlePickerManual"
      @close="pickerOpen = false"
    />
  </ContentModal>

  <ContentModal
    :open="providerEditorOpen"
    :title="providerEditorTitle"
    size="md"
    :close-disabled="appState.configSaving"
    @close="closeProviderEditor"
  >
    <ProviderEditorModal
      v-if="providerEditorOpen"
      :key="providerEditorSession"
      :mode="providerEditorMode"
      :provider="providerEditorProvider"
      @created="handleProviderCreated"
      @saved="handleProviderSaved"
      @close="closeProviderEditor"
    />
  </ContentModal>
</template>
