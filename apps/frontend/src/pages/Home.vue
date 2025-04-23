<template>
  <div class="min-h-screen flex flex-col items-center justify-center text-center p-8">
    <h1 class="text-3xl font-bold mb-4">🏃‍♂️ Runote へようこそ！</h1>
    <p class="text-gray-600 dark:text-gray-300 mb-6">
      ランニングと感情を記録する、新しい日記体験を始めましょう。
    </p>
    <div class="space-x-4">
      <BaseButton @click="goToLogin">ログイン</BaseButton>
      <BaseButton @click="goToRegister">新規登録</BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import BaseButton from "@/components/base/BaseButton.vue";

const router = useRouter();
const goToLogin = () => router.push("/login");
const goToRegister = () => router.push("/register");
</script>
