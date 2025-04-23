<template>
  <div v-if="visible" class="fixed bottom-4 right-4 bg-primary text-white px-4 py-2 rounded shadow-md z-50">
    {{ message }}
  </div>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
  message: string
}>()
</script>