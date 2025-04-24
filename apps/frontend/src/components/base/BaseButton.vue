<template>
  <button
    :type="type"
    class="btn px-6 py-3 rounded-full font-semibold transition-all duration-300"
    :class="{
      'opacity-50 cursor-not-allowed': disabled,
    }"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
}>();
</script>
