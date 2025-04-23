<template>
  <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
    <div class="bg-white dark:bg-gray-800 p-6 rounded shadow-lg max-w-md w-full">
      <slot name="header">
        <h3 class="text-lg font-semibold mb-4">モーダルタイトル</h3>
      </slot>
      <slot />
      <div class="text-right mt-6">
        <button class="btn-secondary mr-2" @click="$emit('close')">閉じる</button>
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ visible: boolean }>()
defineEmits(['close'])
</script>