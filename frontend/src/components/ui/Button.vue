<script setup>
defineProps({
  variant: {
    type: String,
    default: "default",
    validator: (value) => ["default", "primary", "text"].includes(value),
  },
  disabled: { type: Boolean, default: false },
});

const VARIANT_CLASS = {
  primary:
    "border-[var(--brand-border)] bg-[var(--brand)] font-medium text-white hover:bg-[var(--brand-hover)]",
  default:
    "border-[var(--border-strong)] bg-[var(--bg-elevated)] text-[var(--text-primary)] hover:border-[var(--brand-border)] hover:bg-[var(--bg-hover)]",
};
</script>

<template>
  <button
    v-if="variant === 'text'"
    type="button"
    :disabled="disabled"
    class="shrink-0 cursor-pointer whitespace-nowrap text-[12.5px] text-[var(--text-secondary)] transition-colors hover:text-[var(--brand)] disabled:cursor-not-allowed disabled:opacity-60"
  >
    <slot />
  </button>
  <button
    v-else
    type="button"
    :disabled="disabled"
    class="flex shrink-0 cursor-pointer items-center justify-center gap-1.5 whitespace-nowrap rounded-[8px] border px-3 py-[6px] text-[12.5px] transition-colors disabled:cursor-not-allowed disabled:opacity-60"
    :class="VARIANT_CLASS[variant]"
  >
    <slot />
  </button>
</template>