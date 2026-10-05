<script setup>
import { computed } from "vue";

const props = defineProps({
  rate: { type: Number, default: null },
  size: { type: Number, default: 128 },
});

const GEOMETRY = {
  cx: 66,
  cy: 66,
  radius: 50,
  stroke: 11,
};

// 弧线本体按 132×78 绘制，容器再加高到 132×94，把百分比读数放到弧线下方，
// 避免数字贴住圆弧内缘。
const VIEWBOX_WIDTH = 132;
const ARC_HEIGHT = 78;
const CONTAINER_HEIGHT = 94;

const normalizedRate = computed(() => {
  const value = Number(props.rate);
  if (!Number.isFinite(value)) {
    return null;
  }
  return Math.max(0, Math.min(1, value));
});

const percentage = computed(() =>
  normalizedRate.value === null ? 0 : normalizedRate.value * 100,
);

const label = computed(() =>
  normalizedRate.value === null ? "--" : `${percentage.value.toFixed(2)}%`,
);

const containerHeight = computed(
  () => (props.size * CONTAINER_HEIGHT) / VIEWBOX_WIDTH,
);
const arcHeight = computed(() => (props.size * ARC_HEIGHT) / VIEWBOX_WIDTH);

function polarPoint(ratio) {
  const { cx, cy, radius } = GEOMETRY;
  const angle = Math.PI * (1 - ratio);
  return {
    x: cx + radius * Math.cos(angle),
    y: cy - radius * Math.sin(angle),
  };
}

const trackPath = computed(() => {
  const { cx, cy, radius } = GEOMETRY;
  return `M ${cx - radius} ${cy} A ${radius} ${radius} 0 0 1 ${cx + radius} ${cy}`;
});

const valuePath = computed(() => {
  const rate = normalizedRate.value;
  if (rate === null || rate <= 0) {
    return "";
  }
  const { cx, cy, radius } = GEOMETRY;
  const end = polarPoint(rate);
  // 表盘只占 180°，弧线永远不超过半圈，large-arc-flag 必须固定为 0。
  return `M ${cx - radius} ${cy} A ${radius} ${radius} 0 0 1 ${end.x} ${end.y}`;
});

const needle = computed(() => {
  const rate = normalizedRate.value;
  if (rate === null || rate <= 0) {
    return null;
  }
  const { cx, cy, radius, stroke } = GEOMETRY;
  const angle = Math.PI * (1 - rate);
  const inner = radius - stroke / 2 - 1;
  const outer = radius + stroke / 2 + 1;
  return {
    x1: cx + inner * Math.cos(angle),
    y1: cy - inner * Math.sin(angle),
    x2: cx + outer * Math.cos(angle),
    y2: cy - outer * Math.sin(angle),
  };
});
</script>

<template>
  <div
    class="relative shrink-0"
    :style="{ width: `${size}px`, height: `${containerHeight}px` }"
    role="img"
    :aria-label="`缓存命中率 ${label}`"
  >
    <svg
      viewBox="0 0 132 78"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="w-full"
      :style="{ height: `${arcHeight}px` }"
    >
      <path
        :d="trackPath"
        stroke="#2c2c2c"
        :stroke-width="GEOMETRY.stroke"
        stroke-linecap="round"
      />
      <path
        v-if="valuePath"
        :d="valuePath"
        stroke="#10ad5d"
        :stroke-width="GEOMETRY.stroke"
        stroke-linecap="round"
        class="transition-[d] duration-500"
      />
      <line
        v-if="needle"
        :x1="needle.x1"
        :y1="needle.y1"
        :x2="needle.x2"
        :y2="needle.y2"
        stroke="#eaffe8"
        stroke-width="2.4"
        stroke-linecap="round"
      />
    </svg>
    <div
      class="pointer-events-none absolute inset-x-0 bottom-0 flex justify-center"
    >
      <span
        class="font-num text-[18px] leading-none font-medium text-[#f2f2f2]"
      >
        {{ label }}
      </span>
    </div>
  </div>
</template>