<template>
  <button
    :type="type"
    class="btn"
    :class="{
      'opacity-50 cursor-not-allowed': disabled
    }"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
}>()
</script>