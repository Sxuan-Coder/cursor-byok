<script setup>
// 概览横幅右侧装饰：点阵底纹 + 会话卡片 + 绿色校验徽标（纯矢量，无位图）。
</script>

<template>
  <svg
    viewBox="0 0 380 132"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
    class="h-full w-full"
    preserveAspectRatio="xMidYMid slice"
  >
    <defs>
      <radialGradient id="heroGlow" cx="0.62" cy="0.45" r="0.62">
        <stop offset="0%" stop-color="#10ad5d" stop-opacity="0.2" />
        <stop offset="100%" stop-color="#10ad5d" stop-opacity="0" />
      </radialGradient>
      <linearGradient id="heroCardFace" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#22322a" />
        <stop offset="100%" stop-color="#182520" />
      </linearGradient>
      <linearGradient id="heroBadge" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stop-color="#3ee58a" />
        <stop offset="100%" stop-color="#10ad5d" />
      </linearGradient>
      <pattern id="heroDots" width="13" height="13" patternUnits="userSpaceOnUse">
        <circle cx="1.1" cy="1.1" r="1.05" fill="#3d3d3d" />
      </pattern>
      <mask id="heroDotFade">
        <linearGradient id="heroDotFadeGrad" x1="0" y1="0" x2="1" y2="0">
          <stop offset="0%" stop-color="#000000" />
          <stop offset="38%" stop-color="#ffffff" />
          <stop offset="100%" stop-color="#ffffff" />
        </linearGradient>
        <rect width="380" height="132" fill="url(#heroDotFadeGrad)" />
      </mask>
    </defs>

    <rect width="380" height="132" fill="url(#heroDots)" opacity="0.55" mask="url(#heroDotFade)" />
    <rect width="380" height="132" fill="url(#heroGlow)" />

    <!-- 会话预览卡片 -->
    <g transform="translate(184 22)">
      <rect
        width="132"
        height="88"
        rx="11"
        fill="url(#heroCardFace)"
        stroke="#2f4a3b"
        stroke-width="1"
      />
      <rect x="14" y="17" width="30" height="5.5" rx="2.75" fill="#37795a" />
      <rect x="14" y="33" width="86" height="4" rx="2" fill="#3b4d44" />
      <rect x="14" y="44" width="104" height="4" rx="2" fill="#35453d" />
      <rect x="14" y="55" width="70" height="4" rx="2" fill="#314038" />
      <rect x="14" y="66" width="92" height="4" rx="2" fill="#2c3a33" />
    </g>

    <!-- 校验徽标 -->
    <g transform="translate(300 44)">
      <circle cx="0" cy="0" r="21" fill="#10ad5d" opacity="0.14" />
      <circle cx="0" cy="0" r="15.5" fill="url(#heroBadge)" />
      <path
        d="M-6.2 0.4 L-1.9 4.7 L6.4 -3.6"
        stroke="#08301b"
        stroke-width="2.6"
        stroke-linecap="round"
        stroke-linejoin="round"
        fill="none"
      />
    </g>
  </svg>
</template>