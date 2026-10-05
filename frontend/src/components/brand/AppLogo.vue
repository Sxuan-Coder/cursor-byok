<script setup>
defineProps({
  size: { type: [Number, String], default: 24 },
});
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="appLogoTop" x1="16" y1="2.5" x2="16" y2="16">
        <stop offset="0%" stop-color="#3ee58a" />
        <stop offset="100%" stop-color="#16b866" />
      </linearGradient>
      <linearGradient id="appLogoLeft" x1="4" y1="10" x2="16" y2="29">
        <stop offset="0%" stop-color="#14a75c" />
        <stop offset="100%" stop-color="#0a6b3d" />
      </linearGradient>
      <linearGradient id="appLogoRight" x1="28" y1="10" x2="16" y2="29">
        <stop offset="0%" stop-color="#0f8a4c" />
        <stop offset="100%" stop-color="#075c33" />
      </linearGradient>
    </defs>
    <path d="M16 2.6 L27.6 9.3 L16 16 L4.4 9.3 Z" fill="url(#appLogoTop)" />
    <path d="M4.4 9.3 L16 16 L16 29.4 L4.4 22.7 Z" fill="url(#appLogoLeft)" />
    <path d="M27.6 9.3 L16 16 L16 29.4 L27.6 22.7 Z" fill="url(#appLogoRight)" />
    <path
      d="M16 2.6 L27.6 9.3 L16 16 L4.4 9.3 Z M4.4 9.3 L16 16 L16 29.4 L4.4 22.7 Z M27.6 9.3 L16 16 L16 29.4 L27.6 22.7 Z"
      stroke="#2b2b2b"
      stroke-width="0.75"
      stroke-linejoin="round"
    />
  </svg>
</template>