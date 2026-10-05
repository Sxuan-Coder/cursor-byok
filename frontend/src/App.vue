<template>
  <MainLayout />
  <MessageProvider />
  <Modal

    :visible="modalState.visible"
    :title="modalState.title"
    :content="modalState.content"
    :confirm-text="modalState.confirmText"
    :cancel-text="modalState.cancelText"
    :show-cancel="modalState.showCancel"
    :confirm-disabled="modalState.confirmDisabled"
    @confirm="resolveModal(true)"
    @cancel="resolveModal(false)"
  />
  <Modal
    :visible="appState.updatePromptVisible"
    :title="updateViewState.promptTitle"
    :content="updateViewState.promptContent"
    :confirm-text="updateViewState.promptConfirmText"
    :cancel-text="updateViewState.promptCancelText"
    :show-cancel="updateViewState.promptShowCancel"
    :confirm-disabled="appState.updatePromptBusy"
    @confirm="confirmUpdatePrompt"
    @cancel="dismissUpdatePrompt"
  />
  <InputModal
    :visible="inputModalState.visible"
    :title="inputModalState.title"
    :content="inputModalState.content"
    :placeholder="inputModalState.placeholder"
    :model-value="inputModalState.value"
    @update:model-value="inputModalState.value = $event"
    @confirm="resolveInputModal(true)"
    @cancel="resolveInputModal(false)"
  />
</template>
<script setup>
import MainLayout from "@/layouts/MainLayout.vue";
import Modal from "@/components/ui/Modal.vue";
import MessageProvider from "@/components/ui/MessageProvider.vue";
import { modalState, resolveModal } from "@/composables/useModal";

import InputModal from "@/components/ui/InputModal.vue";
import { inputModalState, resolveInputModal } from "@/composables/useInputModal";
import { appState, confirmUpdatePrompt, dismissUpdatePrompt, updateViewState } from "@/state/appState";

// 更新弹窗曾用「当前是否首页」判断主窗口，导致切到其它 Tab 时弹窗不渲染。
// 现在所有页面都在同一主窗口内切换路由，弹窗改为全局无条件渲染。
// import { computed } from "vue";
// import { useRoute } from "vue-router";
//
// const route = useRoute();
// const isMainWindow = computed(() => route.path === "/");
</script>
