<template>
  <button @click="handleClick" class="text-blue-500 hover:text-blue-700 underline">
    ← 戻る
  </button>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { defineProps } from "vue";

const props = defineProps<{
  to?: string; // オプションで戻る先を指定できる
}>();

const router = useRouter();

const handleClick = () => {
  if (props.to) {
    router.push(props.to); // 戻る先が指定されていればそこへ
  } else {
    router.back(); // 指定がなければ1つ前の履歴へ
  }
};
</script>
