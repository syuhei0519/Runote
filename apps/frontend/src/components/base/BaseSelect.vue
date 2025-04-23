<template>
  <select
    :value="modelValue"
    @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    class="input"
  >
    <option v-for="option in options" :key="option.value" :value="option.value">
      {{ option.label }}
    </option>
  </select>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
  options: { value: string; label: string }[];
}>();

defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();
</script>
