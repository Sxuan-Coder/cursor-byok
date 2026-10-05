<script setup>
import HeroArt from "@/components/brand/HeroArt.vue";
import { useAccountName } from "@/composables/useAccountName";

const { displayName } = useAccountName();
</script>

<template>
  <section
    class="relative h-[68px] shrink-0 overflow-hidden rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]"
  >
    <div class="pointer-events-none absolute inset-y-0 right-0 w-[44%]">
      <HeroArt />
    </div>

    <div class="relative z-10 flex flex-col px-5 pt-[11px]">
      <h1 class="text-[17.5px] font-semibold leading-[22px] text-[var(--text-primary)]">
        你好，<span class="text-[var(--brand)]">{{ displayName }}</span>
      </h1>
      <p class="mt-[2px] text-[12px] leading-[16px] text-[var(--text-secondary)]">
        一站式管理你的 BYOK 账号，享受更自由的 AI 编程体验
      </p>
    </div>

    <div
      class="absolute bottom-[8px] right-5 z-10 flex items-center gap-[5px] text-[11px] text-[var(--text-faint)]"
    >
      <span class="h-[5px] w-[5px] rounded-full bg-[var(--brand)]"></span>
      <span>Cursor · 多模型 · 本地代理</span>
    </div>
  </section>
</template>