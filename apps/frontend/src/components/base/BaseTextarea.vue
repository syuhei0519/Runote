<template>
  <textarea
    :placeholder="placeholder"
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    class="input h-32 resize-none"
  />
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
  placeholder?: string;
}>();

defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();
</script>
