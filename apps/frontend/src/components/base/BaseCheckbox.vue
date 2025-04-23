<template>
  <label class="inline-flex items-center space-x-2">
    <input
      type="checkbox"
      :checked="modelValue"
      @change="$emit('update:modelValue', $event.target.checked)"
      class="rounded border-gray-300 text-primary focus:ring-primary"
    />
    <span><slot /></span>
  </label>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: boolean;
}>();

defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
</script>
