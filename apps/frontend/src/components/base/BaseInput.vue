<template>
  <input
    :type="type"
    :placeholder="placeholder"
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    class="input"
  />
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
  type?: string;
  placeholder?: string;
}>();

defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();
</script>
