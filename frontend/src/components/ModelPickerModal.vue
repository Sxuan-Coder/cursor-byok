<script setup>
import Button from "@/components/ui/Button.vue";
import Select from "@/components/ui/Select.vue";
import Tooltip from "@/components/ui/Tooltip.vue";
import {
  addModelAdapters,
  createModelAdapterFromTemplate,
  fetchAvailableModelIDs,
} from "@/state/appState";
import { computed, onMounted, ref } from "vue";

const props = defineProps({
  baseURL: { type: String, required: true },
  apiKey: { type: String, required: true },
  initialType: { type: String, default: "openai" },
  providerAdapters: { type: Array, default: () => [] },
});

const emit = defineEmits(["close", "added", "manual"]);

const typeTabs = [
  { label: "OpenAI", value: "openai", icon: "icon-[bxl--openai]" },
  { label: "Anthropic", value: "anthropic", icon: "icon-[logos--claude-icon]" },
];

const activeType = ref(props.initialType === "anthropic" ? "anthropic" : "openai");
const models = ref([]);
const loading = ref(false);
const loadError = ref("");
const searchText = ref("");
const selectedModelIDs = ref(new Set());
const templateID = ref("");
const adding = ref(false);
const addActionError = ref("");
const requestSeq = ref(0);

const existingModelIDs = computed(() => new Set(
  props.providerAdapters
    .filter((adapter) => adapter.type === activeType.value)
    .map((adapter) => adapter.modelID),
));
const filteredModels = computed(() => {
  const keyword = searchText.value.trim().toLowerCase();
  const list = models.value;
  if (!keyword) {
    return list;
  }
  return list.filter((modelID) => modelID.toLowerCase().includes(keyword));
});
const selectableModels = computed(() => filteredModels.value.filter((modelID) => !existingModelIDs.value.has(modelID)));
const selectableCount = computed(() => selectableModels.value.length);
const allSelectableSelected = computed(() => (
  selectableCount.value > 0
  && selectableModels.value.every((modelID) => selectedModelIDs.value.has(modelID))
));
const templateOptions = computed(() => [
  { label: "不继承（使用默认配置）", value: "", icon: "icon-[mdi--minus-circle-outline]" },
  ...props.providerAdapters
    .filter((adapter) => adapter.type === activeType.value)
    .map((adapter) => ({
      label: `${adapter.displayName}（${adapter.modelID}）`,
      value: adapter.id,
      icon: "icon-[mdi--content-copy]",
    })),
]);
const selectedTemplate = computed(() =>
  props.providerAdapters.find((adapter) => adapter.id === templateID.value) ?? null,
);
const canSubmit = computed(() => !adding.value && selectedModelIDs.value.size > 0);

function resetSelection() {
  selectedModelIDs.value = new Set();
  templateID.value = "";
  addActionError.value = "";
}

async function refreshModelList() {
  const baseURL = String(props.baseURL ?? "").trim();
  const apiKey = String(props.apiKey ?? "").trim();
  if (!baseURL || !apiKey) {
    loadError.value = "缺少接口地址或访问密钥，无法拉取模型列表";
    models.value = [];
    return;
  }

  const representative = props.providerAdapters.find((adapter) => adapter.type === activeType.value)
    ?? props.providerAdapters[0];
  const requestType = activeType.value;
  const seq = requestSeq.value + 1;
  requestSeq.value = seq;
  loading.value = true;
  loadError.value = "";
  models.value = [];
  try {
    const list = await fetchAvailableModelIDs({
      type: requestType,
      baseURL,
      apiKey,
      customHeadersEnabled: representative?.customHeadersEnabled ?? false,
      customHeadersJSON: representative?.customHeadersJSON ?? "",
    });
    if (seq !== requestSeq.value || requestType !== activeType.value) {
      return;
    }
    if (list.length === 0) {
      loadError.value = "接口没有返回任何模型";
    }
    models.value = list;
  } catch (error) {
    if (seq !== requestSeq.value || requestType !== activeType.value) {
      return;
    }
    models.value = [];
    loadError.value = String(error?.message || error || "拉取模型列表失败").trim() || "拉取模型列表失败";
  } finally {
    if (seq === requestSeq.value) {
      loading.value = false;
    }
  }
}

function handleTypeChange(type) {
  if (activeType.value === type) {
    return;
  }
  activeType.value = type;
  resetSelection();
  searchText.value = "";
  models.value = [];
  void refreshModelList();
}

function toggleModel(modelID) {
  const next = new Set(selectedModelIDs.value);
  if (next.has(modelID)) {
    next.delete(modelID);
  } else {
    next.add(modelID);
  }
  selectedModelIDs.value = next;
}

function toggleSelectAll() {
  if (allSelectableSelected.value) {
    selectedModelIDs.value = new Set();
    return;
  }
  selectedModelIDs.value = new Set(selectableModels.value);
}

async function handleConfirm() {
  if (!canSubmit.value) {
    return;
  }
  const modelIDs = Array.from(selectedModelIDs.value);
  adding.value = true;
  addActionError.value = "";
  try {
    const drafts = modelIDs.map((modelID) => createModelAdapterFromTemplate(selectedTemplate.value, {
      type: activeType.value,
      baseURL: props.baseURL.trim(),
      apiKey: props.apiKey.trim(),
      modelID,
      displayName: modelID,
    }));
    const result = await addModelAdapters(drafts);
    if (!result.ok) {
      addActionError.value = result.error || "新增模型失败";
      return;
    }
    emit("added", result.added);
    emit("close");
  } finally {
    adding.value = false;
  }
}

function handleManual() {
  emit("manual", {
    type: activeType.value,
    baseURL: props.baseURL.trim(),
    apiKey: props.apiKey.trim(),
  });
  emit("close");
}

onMounted(() => {
  void refreshModelList();
});
</script>

<template>
  <div class="flex h-full min-h-0 flex-col text-[#e5e5e5]">
    <div class="shrink-0 space-y-3 px-4 pt-4">
      <div class="flex items-center justify-between gap-3">
        <div class="center-row gap-2">
          <button
            v-for="tab in typeTabs"
            :key="tab.value"
            type="button"
            class="center-row gap-2 rounded-[8px] border px-3 py-2 text-sm transition-colors duration-150"
            :class="activeType === tab.value
              ? 'border-[#1ca35a] bg-[#123322] text-white'
              : 'border-[#343434] bg-[#252525] text-[#a3a3a3] hover:border-[#4a4a4a] hover:text-[#e5e5e5]'"
            @click="handleTypeChange(tab.value)"
          >
            <span :class="[tab.icon, 'text-[16px]']"></span>
            <span>{{ tab.label }}</span>
          </button>
        </div>
        <Button variant="default" :disabled="loading || adding" @click="refreshModelList">
          <span class="center-row gap-1.5">
            <span class="icon-[mdi--refresh] text-[15px]" :class="{ 'animate-spin': loading }"></span>
            <span>刷新列表</span>
          </span>
        </Button>
      </div>

      <div class="relative">
        <span class="icon-[mdi--magnify] pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-[16px] text-[#8f8f8f]"></span>
        <input
          v-model="searchText"
          type="text"
          placeholder="搜索模型标识"
          class="h-9 w-full rounded-[6px] border border-[#3f3f3f] bg-[#232323] pl-9 pr-3 text-sm text-[#e5e5e5] outline-none focus:border-[#10AD5D]"
        />
      </div>

      <div class="flex items-center gap-3">
        <span class="center-row justify-start gap-1.5 text-sm text-[#d4d4d4]">
          <Tooltip content="批量新增的模型会继承所选模板的高级参数（推理强度、额外参数 JSON、自定义请求头、Token 上限等），模型标识与显示名称使用拉取到的模型名。" />
          <span>继承配置模板</span>
        </span>
        <div class="min-w-0 flex-1">
          <Select
            v-model="templateID"
            :options="templateOptions"
            :disabled="adding"
            aria-label="选择配置模板"
          />
        </div>
      </div>

      <p v-if="providerAdapters.length === 0" class="rounded-[6px] border border-[#3a3223] bg-[#292512] px-3 py-2 text-xs text-[#d8c68a]">
        该提供商还没有任何模型，确认新增后才会保存到配置。
      </p>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto px-4 py-3 scroll-shadow-bottom">
      <div v-if="loading" class="flex h-full min-h-[160px] items-center justify-center text-sm text-[#8f8f8f]">
        正在拉取模型列表...
      </div>
      <div v-else-if="loadError" class="flex h-full min-h-[160px] flex-col items-center justify-center gap-3 text-sm text-[#c9c9c9]">
        <span>{{ loadError }}</span>
        <Button variant="default" @click="refreshModelList">重试</Button>
      </div>
      <div v-else-if="filteredModels.length === 0" class="flex h-full min-h-[160px] items-center justify-center text-sm text-[#8f8f8f]">
        没有匹配的模型
      </div>
      <div v-else class="flex flex-col gap-1">
        <label
          class="center-row sticky top-0 z-[1] cursor-pointer gap-2 rounded-[6px] border border-[#333] bg-[#232323] px-3 py-2 text-sm text-[#d4d4d4]"
          :class="allSelectableSelected ? 'text-white' : ''"
        >
          <input
            type="checkbox"
            class="size-4 accent-[#10AD5D]"
            :checked="allSelectableSelected"
            :disabled="selectableCount === 0 || adding"
            @change="toggleSelectAll"
          />
          <span>{{ allSelectableSelected ? "取消全选" : "全选可选模型" }}</span>
          <span class="ml-auto text-xs text-[#8f8f8f]">可选 {{ selectableCount }} / 共 {{ filteredModels.length }}</span>
        </label>
        <label
          v-for="modelID in filteredModels"
          :key="modelID"
          class="center-row cursor-pointer gap-2 rounded-[6px] border border-transparent px-3 py-2 text-sm transition-colors"
          :class="existingModelIDs.has(modelID)
            ? 'cursor-not-allowed opacity-50'
            : 'hover:border-[#3f3f3f] hover:bg-[#252525]'"
          :title="existingModelIDs.has(modelID) ? '该模型已在此提供商下配置' : modelID"
        >
          <input
            type="checkbox"
            class="size-4 accent-[#10AD5D]"
            :checked="selectedModelIDs.has(modelID)"
            :disabled="existingModelIDs.has(modelID) || adding"
            @change="toggleModel(modelID)"
          />
          <span class="min-w-0 flex-1 truncate font-mono text-[13px] text-[#e5e5e5]">{{ modelID }}</span>
          <span
            v-if="existingModelIDs.has(modelID)"
            class="shrink-0 rounded-[6px] border border-[#2f4a35] bg-[#16281c] px-2 py-[2px] text-[11px] text-[#7fce9d]"
          >
            已添加
          </span>
        </label>
      </div>
    </div>

    <div class="flex shrink-0 flex-wrap items-center justify-between gap-2 px-4 py-3">
      <span v-if="addActionError" class="min-w-0 flex-1 truncate text-sm text-[#ff9a9a]" :title="addActionError">
        {{ addActionError }}
      </span>
      <span v-else class="min-w-0 flex-1 truncate text-xs text-[#8f8f8f]">
        已选 {{ selectedModelIDs.size }} 个模型
      </span>
      <div class="center-row gap-2">
        <Button variant="text" :disabled="adding" @click="handleManual">手动配置单个模型</Button>
        <Button variant="default" :disabled="adding" @click="emit('close')">取消</Button>
        <Button variant="primary" :disabled="!canSubmit" @click="handleConfirm">
          {{ adding ? "新增中..." : "确认新增" }}
        </Button>
      </div>
    </div>
  </div>
</template>
