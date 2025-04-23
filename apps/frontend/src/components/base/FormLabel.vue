<template>
  <label :for="for" class="block mb-1 font-medium text-gray-700 dark:text-gray-200">
    <slot />
  </label>
</template>

<script setup lang="ts">
defineProps<{ for: string }>()
</script>