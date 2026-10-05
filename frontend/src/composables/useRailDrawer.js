import { ref } from "vue";

// 概览页右侧栏在窄窗口下会收起，此处在顶栏提供一个抽屉开关共用同一状态。
const railDrawerOpen = ref(false);

export function useRailDrawer() {
  function openRailDrawer() {
    railDrawerOpen.value = true;
  }

  function closeRailDrawer() {
    railDrawerOpen.value = false;
  }

  function toggleRailDrawer() {
    railDrawerOpen.value = !railDrawerOpen.value;
  }

  return { railDrawerOpen, openRailDrawer, closeRailDrawer, toggleRailDrawer };
}