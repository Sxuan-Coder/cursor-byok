<script setup>
import Button from "@/components/ui/Button.vue";
import Input from "@/components/ui/Input.vue";
import Tooltip from "@/components/ui/Tooltip.vue";
import { updateProviderCredentials } from "@/state/appState";
import { computed, reactive, ref } from "vue";

const props = defineProps({
  mode: { type: String, default: "create" },
  provider: { type: Object, default: null },
});

const emit = defineEmits(["close", "created", "saved"]);

const isCreate = computed(() => props.mode !== "edit");

const typeTabs = [
  { label: "OpenAI", value: "openai", icon: "icon-[bxl--openai]" },
  { label: "Anthropic", value: "anthropic", icon: "icon-[logos--claude-icon]" },
];

const draft = reactive({
  type: "openai",
  baseURL: isCreate.value ? "" : String(props.provider?.baseURL || ""),
  apiKey: isCreate.value ? "" : String(props.provider?.apiKey || ""),
});

const validationError = ref("");
const saveError = ref("");
const submitting = ref(false);

const interfacePlaceholder = computed(() =>
  draft.type === "anthropic" ? "例如：https://api.anthropic.com" : "例如：https://api.openai.com/v1",
);

function validateDraft() {
  const baseURL = draft.baseURL.trim();
  const apiKey = draft.apiKey.trim();
  if (!baseURL) {
    return "接口地址不能为空";
  }
  try {
    const parsed = new URL(baseURL);
    if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
      return "接口地址仅支持 http 或 https";
    }
    if (!parsed.hostname) {
      return "接口地址缺少主机名";
    }
  } catch (_error) {
    return "接口地址不是合法 URL";
  }
  if (!apiKey) {
    return "访问密钥不能为空";
  }
  return "";
}

function handleTypeChange(type) {
  draft.type = type;
}

async function handleSubmit() {
  if (submitting.value) {
    return;
  }
  validationError.value = "";
  saveError.value = "";

  if (isCreate.value) {
    const error = validateDraft();
    if (error) {
      validationError.value = error;
      return;
    }
    emit("created", {
      type: draft.type,
      baseURL: draft.baseURL.trim(),
      apiKey: draft.apiKey.trim(),
    });
    return;
  }

  const error = validateDraft();
  if (error) {
    validationError.value = error;
    return;
  }
  submitting.value = true;
  try {
    const result = await updateProviderCredentials(
      props.provider?.baseURL,
      props.provider?.apiKey,
      {
        baseURL: draft.baseURL.trim(),
        apiKey: draft.apiKey.trim(),
      },
    );
    if (!result.ok) {
      saveError.value = result.error || "更新提供商凭证失败";
      return;
    }
    emit("saved", result.updated);
    emit("close");
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-4 px-4 py-4 text-[#e5e5e5]">
    <p class="text-sm leading-relaxed text-[#a3a3a3]">
      {{ isCreate
        ? "提供商按「接口地址 + 访问密钥」识别。填写后会打开模型拉取列表，至少添加一个模型才会保存。"
        : "更新后会批量替换该提供商下全部模型的接口地址与访问密钥。" }}
    </p>

    <div v-if="isCreate" class="center-row gap-2">
      <button
        v-for="tab in typeTabs"
        :key="tab.value"
        type="button"
        class="center-row gap-2 rounded-[8px] border px-3 py-2 text-sm transition-colors duration-150"
        :class="draft.type === tab.value
          ? 'border-[#1ca35a] bg-[#123322] text-white'
          : 'border-[#343434] bg-[#252525] text-[#a3a3a3] hover:border-[#4a4a4a] hover:text-[#e5e5e5]'"
        @click="handleTypeChange(tab.value)"
      >
        <span :class="[tab.icon, 'text-[16px]']"></span>
        <span>{{ tab.label }}</span>
      </button>
    </div>

    <label class="flex flex-col gap-1">
      <span class="center-row justify-start gap-1.5 text-sm text-[#d4d4d4]">
        <Tooltip content="模型服务的 API 根地址，通常为兼容 OpenAI 或 Anthropic 的接口入口。" />
        <span>接口地址</span>
      </span>
      <input
        v-model="draft.baseURL"
        type="text"
        :placeholder="interfacePlaceholder"
        class="h-9 rounded-[6px] border border-[#3f3f3f] bg-[#232323] px-3 text-sm text-[#e5e5e5] outline-none focus:border-[#10AD5D]"
      />
    </label>

    <label class="flex flex-col gap-1">
      <span class="center-row justify-start gap-1.5 text-sm text-[#d4d4d4]">
        <Tooltip content="调用该模型服务需要使用的访问密钥。" />
        <span>访问密钥</span>
      </span>
      <Input
        v-model="draft.apiKey"
        type="password"
        allow-visibility-toggle
        placeholder="例如：sk-xxxxxx"
        autocomplete="off"
      />
    </label>

    <p v-if="validationError || saveError" class="text-sm text-[#ff9a9a]">
      {{ validationError || saveError }}
    </p>

    <div class="flex items-center justify-end gap-2">
      <Button variant="default" :disabled="submitting" @click="emit('close')">取消</Button>
      <Button variant="primary" :disabled="submitting" @click="handleSubmit">
        {{ submitting ? "保存中..." : (isCreate ? "下一步：选择模型" : "保存凭证") }}
      </Button>
    </div>
  </div>
</template>
