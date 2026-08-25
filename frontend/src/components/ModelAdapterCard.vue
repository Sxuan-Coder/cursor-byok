<script setup>
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import ModelAdapterTestCard from "@/components/ModelAdapterTestCard.vue";
import { getModelAdapterTestResultByID } from "@/state/appState";

defineProps({
  adapter: { type: Object, required: true },
  actionsDisabled: { type: Boolean, default: false },
});

defineEmits(["test", "edit", "duplicate", "delete"]);

function typeLabel(type) {
  return type === "anthropic" ? "Anthropic" : "OpenAI";
}

function getAdapterTestResult(adapter) {
  return getModelAdapterTestResultByID(adapter?.id);
}

function isAdapterTesting(adapter) {
  return getAdapterTestResult(adapter)?.status === "running";
}
</script>

<template>
  <Card
    class="model-sort-item group relative pb-2"
    :data-model-id="adapter.id"
  >
    <button
      type="button"
      class="model-sort-handle w-[30px] h-[30px] center-row justify-center absolute left-2 top-2 z-10 shrink-0 touch-none cursor-grab rounded-[6px] border border-transparent bg-transparent text-transparent opacity-0 outline-none transition-[opacity,color,border-color,background-color] focus-visible:border-[#10AD5D] focus-visible:bg-[#333333] focus-visible:text-white focus-visible:opacity-100 active:cursor-grabbing group-hover:border-[#454545] group-hover:bg-[#333333] group-hover:text-white group-hover:opacity-100 disabled:cursor-not-allowed disabled:opacity-30"
      :disabled="actionsDisabled"
      aria-label="拖拽排序"
      title="拖拽排序"
      @click.stop
    >
      <span class="icon-[icon-park-outline--drag] text-[20px]"></span>
    </button>
    <div class="flex h-[150px] flex-col justify-between gap-3">
      <div class="flex flex-col gap-2.5">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0 flex-1">
            <div class="truncate text-base font-medium text-white">{{ adapter.displayName }}</div>
            <div class="mt-1 truncate text-sm text-[#8f8f8f]">{{ adapter.modelID }}</div>
          </div>
          <span
            class="center-row shrink-0 gap-1 rounded-[999px] border border-[#3f3f3f] px-[7px] py-[4px] text-[11px] font-medium text-[#cfcfcf]"
          >
            <span class="icon-[bxl--openai] text-[14px] !text-white" v-if="adapter.type === 'openai'"></span>
            <span class="icon-[logos--claude-icon] text-[14px]" v-else></span>
            <span>{{ typeLabel(adapter.type) }}</span>
          </span>
        </div>

        <ModelAdapterTestCard
          compact
          title="测试"
          empty-text="未测试"
          :result="getAdapterTestResult(adapter)"
        />
      </div>

      <div class="center-row flex-wrap justify-end gap-2 pt-0">
        <Button
          variant="default"
          :disabled="actionsDisabled || isAdapterTesting(adapter)"
          @click="$emit('test', adapter)"
        >
          {{ isAdapterTesting(adapter) ? "测试中..." : "测试" }}
        </Button>
        <Button variant="default" :disabled="actionsDisabled" @click="$emit('edit', adapter)">编辑</Button>
        <Button variant="default" :disabled="actionsDisabled" @click="$emit('duplicate', adapter)">复制</Button>
        <Button variant="text" :disabled="actionsDisabled" @click="$emit('delete', adapter)">删除</Button>
      </div>
    </div>
  </Card>
</template>
