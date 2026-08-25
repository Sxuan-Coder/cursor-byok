<script setup>
import ModelAdapterCard from "@/components/ModelAdapterCard.vue";
import Sortable from "sortablejs";
import { computed, nextTick, onBeforeUnmount, ref, watch } from "vue";

const props = defineProps({
  group: { type: Object, required: true },
  actionsDisabled: { type: Boolean, default: false },
});

const emit = defineEmits([
  "fetch-models",
  "add-model",
  "edit-provider",
  "delete-provider",
  "test",
  "edit",
  "duplicate",
  "delete",
  "reorder",
]);

const DELETE_ARM_RESET_MS = 3000;

const collapsed = ref(false);
const deleteArmed = ref(false);
const modelGrid = ref(null);
let sortable = null;
let deleteArmTimer = 0;

const adapterOrderKey = computed(() =>
  props.group.adapters.map((adapter) => adapter.id).join("\n"),
);

function typeLabel(type) {
  return type === "anthropic" ? "Anthropic" : "OpenAI";
}

function typeIcon(type) {
  return type === "anthropic" ? "icon-[logos--claude-icon]" : "icon-[bxl--openai]";
}

function maskSecret(value) {
  const text = String(value || "").trim();
  if (!text) {
    return "-";
  }
  if (text.length <= 8) {
    return `${"*".repeat(Math.max(text.length - 2, 0))}${text.slice(-2)}`;
  }
  return `${text.slice(0, 4)}****${text.slice(-4)}`;
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

function armDeleteProvider() {
  if (!deleteArmed.value) {
    deleteArmed.value = true;
    window.clearTimeout(deleteArmTimer);
    deleteArmTimer = window.setTimeout(() => {
      deleteArmed.value = false;
    }, DELETE_ARM_RESET_MS);
    return;
  }
  window.clearTimeout(deleteArmTimer);
  deleteArmed.value = false;
  emit("delete-provider", props.group);
}

function destroySortable() {
  if (!sortable) {
    return;
  }
  sortable.destroy();
  sortable = null;
}

function syncSortable() {
  const element = modelGrid.value;
  if (!element) {
    destroySortable();
    return;
  }
  if (!sortable || sortable.el !== element) {
    destroySortable();
    sortable = Sortable.create(element, {
      animation: 160,
      dataIdAttr: "data-model-id",
      draggable: ".model-sort-item",
      handle: ".model-sort-handle",
      ghostClass: "opacity-40",
      chosenClass: "!border-[#10AD5D]",
      dragClass: "cursor-grabbing",
      onEnd: (event) => {
        handleSortEnd(event);
      },
    });
  }
  sortable.option("disabled", props.actionsDisabled);
  sortable.sort(props.group.adapters.map((adapter) => adapter.id), false);
}

function handleSortEnd(event) {
  const oldIndex = event.oldDraggableIndex ?? event.oldIndex;
  const newIndex = event.newDraggableIndex ?? event.newIndex;
  if (
    !Number.isInteger(oldIndex)
    || !Number.isInteger(newIndex)
    || oldIndex === newIndex
  ) {
    syncSortable();
    return;
  }

  const reordered = props.group.adapters.slice();
  const [movedAdapter] = reordered.splice(oldIndex, 1);
  if (!movedAdapter || newIndex < 0 || newIndex > reordered.length) {
    syncSortable();
    return;
  }
  reordered.splice(newIndex, 0, movedAdapter);
  emit("reorder", { group: props.group, adapters: reordered });
}

watch(
  () => [
    modelGrid.value,
    adapterOrderKey.value,
    props.actionsDisabled,
  ],
  () => {
    void nextTick().then(syncSortable);
  },
  { flush: "post", immediate: true },
);

onBeforeUnmount(() => {
  window.clearTimeout(deleteArmTimer);
  destroySortable();
});
</script>

<template>
  <section class="rounded-[10px] border border-[#2e2e2e] bg-[#1f1f1f]/60">
    <div class="flex flex-wrap items-center gap-x-3 gap-y-2 px-4 py-3">
      <div class="center-row min-w-0 flex-1 gap-2">
        <span class="icon-[mdi--earth] text-[18px] shrink-0 text-[#8f8f8f]"></span>
        <span class="truncate text-[15px] font-medium text-white" :title="group.baseURL">
          {{ formatHost(group.baseURL) }}
        </span>
        <span
          v-for="type in group.types"
          :key="type"
          class="center-row shrink-0 gap-1 rounded-[999px] border border-[#3f3f3f] px-[7px] py-[3px] text-[11px] font-medium text-[#cfcfcf]"
        >
          <span
            :class="[typeIcon(type), 'text-[13px]', type === 'openai' ? '!text-white' : '']"
          ></span>
          <span>{{ typeLabel(type) }}</span>
        </span>
        <span class="shrink-0 text-xs text-[#8f8f8f]">{{ group.adapters.length }} 个模型</span>
        <span
          class="shrink-0 rounded-[6px] border border-[#333] bg-[#232323] px-2 py-[3px] font-mono text-[11px] text-[#8f8f8f]"
          :title="group.apiKey"
        >
          {{ maskSecret(group.apiKey) }}
        </span>
      </div>

      <div class="center-row flex-wrap justify-end gap-2">
        <button
          type="button"
          class="center-row h-8 gap-1.5 whitespace-nowrap rounded-[6px] border border-[#1ca35a] bg-[#123322] px-3 text-sm text-white transition-colors hover:border-[#28c76f] disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="actionsDisabled"
          @click="emit('fetch-models', group)"
        >
          <span class="icon-[mdi--cloud-download-outline] text-[16px]"></span>
          <span>拉取模型</span>
        </button>
        <button
          type="button"
          class="center-row h-8 gap-1.5 whitespace-nowrap rounded-[6px] border border-[#3f3f3f] bg-[#292929] px-3 text-sm text-[#d4d4d4] transition-colors hover:border-[#505050] hover:bg-[#303030] hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="actionsDisabled"
          @click="emit('add-model', group)"
        >
          <span class="icon-[mdi--plus-thick] text-[15px]"></span>
          <span>新增模型</span>
        </button>
        <button
          type="button"
          class="center-row h-8 gap-1.5 whitespace-nowrap rounded-[6px] border border-[#3f3f3f] bg-[#292929] px-3 text-sm text-[#d4d4d4] transition-colors hover:border-[#505050] hover:bg-[#303030] hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="actionsDisabled"
          @click="emit('edit-provider', group)"
        >
          <span class="icon-[mdi--key-outline] text-[15px]"></span>
          <span>编辑凭证</span>
        </button>
        <button
          type="button"
          class="center-row h-8 gap-1.5 whitespace-nowrap rounded-[6px] border px-3 text-sm transition-colors disabled:cursor-not-allowed disabled:opacity-50"
          :class="deleteArmed
            ? 'border-[#b3414a] bg-[#3a1d20] text-[#ff9a9a]'
            : 'border-[#3f3f3f] bg-[#292929] text-[#d4d4d4] hover:border-[#6a3a3f] hover:text-[#ff9a9a]'"
          :disabled="actionsDisabled"
          @click="armDeleteProvider"
        >
          <span class="icon-[mdi--trash-can-outline] text-[15px]"></span>
          <span>{{ deleteArmed ? "确认删除" : "删除" }}</span>
        </button>
        <button
          type="button"
          class="center-row size-8 shrink-0 rounded-[6px] border border-[#3f3f3f] bg-[#292929] text-[#a3a3a3] transition-colors hover:border-[#505050] hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="actionsDisabled"
          :aria-label="collapsed ? '展开分组' : '折叠分组'"
          :title="collapsed ? '展开分组' : '折叠分组'"
          @click="collapsed = !collapsed"
        >
          <span
            class="icon-[mdi--chevron-down] text-[18px] transition-transform duration-150"
            :class="collapsed ? '-rotate-90' : ''"
          ></span>
        </button>
      </div>
    </div>

    <div v-show="!collapsed" class="border-t border-[#2a2a2a] px-4 py-3">
      <div
        ref="modelGrid"
        class="grid gap-3 pb-1 [grid-template-columns:repeat(auto-fill,minmax(250px,1fr))]"
      >
        <ModelAdapterCard
          v-for="adapter in group.adapters"
          :key="adapter.id || `${adapter.baseURL}-${adapter.modelID}`"
          :adapter="adapter"
          :actions-disabled="actionsDisabled"
          @test="emit('test', $event)"
          @edit="emit('edit', $event)"
          @duplicate="emit('duplicate', $event)"
          @delete="emit('delete', $event)"
        />
      </div>
    </div>
  </section>
</template>
