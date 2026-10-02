<script setup>
import { computed, ref } from "vue";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import ContentModal from "@/components/ui/ContentModal.vue";
import Select from "@/components/ui/Select.vue";
import { useMessage } from "@/composables/useMessage";
import { getCommitPromptDefault } from "@/services/clientApi";
import { appState, saveCommitConfig, toUserError } from "@/state/appState";

const message = useMessage();
const promptOpen = ref(false);
const promptDraft = ref("");
const defaultPrompt = ref("");
const saving = ref(false);

const modelOptions = computed(() => [
  { label: "直连（默认）", value: "" },
  ...appState.modelAdapters.map((adapter) => ({
    label:
      adapter.displayName && adapter.displayName !== adapter.modelID
        ? `${adapter.displayName}（${adapter.modelID}）`
        : adapter.displayName || adapter.modelID,
    value: adapter.id,
  })),
]);

async function handleModelChange(value) {
  saving.value = true;
  try {
    const result = await saveCommitConfig(value, appState.commitPrompt);
    if (!result.ok) {
      throw new Error(result.error || "服务错误");
    }
    message(value ? "提交信息已切换为本地模型生成" : "提交信息已恢复直连");
  } catch (error) {
    message(`保存失败：${toUserError(error)}`);
  } finally {
    saving.value = false;
  }
}

async function openPromptEditor() {
  if (!defaultPrompt.value) {
    try {
      defaultPrompt.value = await getCommitPromptDefault();
    } catch (_error) {
      defaultPrompt.value = "";
    }
  }
  promptDraft.value = (appState.commitPrompt || "").trim() || defaultPrompt.value;
  promptOpen.value = true;
}

function resetPrompt() {
  promptDraft.value = defaultPrompt.value;
}

async function savePrompt() {
  const draft = promptDraft.value.trim();
  const normalized =
    defaultPrompt.value && draft === defaultPrompt.value.trim() ? "" : draft;
  saving.value = true;
  try {
    const result = await saveCommitConfig(appState.commitModelHash || "", normalized);
    if (!result.ok) {
      throw new Error(result.error || "服务错误");
    }
    promptOpen.value = false;
    message("提交信息提示词已保存");
  } catch (error) {
    message(`保存失败：${toUserError(error)}`);
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Card>
    <div class="flex items-center justify-between gap-4">
      <div class="min-w-0">
        <h2 class="text-base font-medium text-white">Commit 提交信息</h2>
        <div class="text-sm text-[#a3a3a3]">
          选择生成提交信息的模型；选择“直连（默认）”时保持官方通道转发，保存后立即生效
        </div>
      </div>
      <div class="flex shrink-0 items-center gap-3">
        <div class="w-[240px] max-w-full">
          <Select
            :model-value="appState.commitModelHash || ''"
            :options="modelOptions"
            :disabled="saving || appState.configSaving"
            aria-label="提交信息生成模型"
            @update:model-value="handleModelChange"
          />
        </div>
        <Button
          variant="text"
          :disabled="saving || appState.configSaving"
          @click="openPromptEditor"
        >
          提示词设置
        </Button>
      </div>
    </div>
  </Card>

  <ContentModal :open="promptOpen" title="提交信息提示词" size="xl" @close="promptOpen = false">
    <div class="flex h-full min-h-0 flex-col gap-3 p-4">
      <textarea
        v-model="promptDraft"
        class="min-h-0 w-full flex-1 resize-none rounded-[6px] border border-[#3f3f3f] bg-[#232323] px-3 py-2 font-mono text-xs leading-relaxed text-[#e5e5e5] outline-none focus:border-[#10AD5D]"
        spellcheck="false"
        aria-label="提交信息提示词"
      />
      <div class="flex shrink-0 items-center justify-between gap-3">
        <Button variant="text" @click="resetPrompt">恢复默认</Button>
        <div class="flex items-center gap-2">
          <Button :disabled="saving" @click="promptOpen = false">取消</Button>
          <Button variant="primary" :disabled="saving" @click="savePrompt">
            {{ saving ? "保存中..." : "保存" }}
          </Button>
        </div>
      </div>
    </div>
  </ContentModal>
</template>