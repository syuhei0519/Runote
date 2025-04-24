<template>
  <header class="bg-white shadow-md py-4 px-6 flex justify-between items-center">
    <h1 class="text-xl font-bold text-blue-600">Runote</h1>
    <nav class="space-x-6">
      <a href="#features" class="text-gray-600 hover:text-blue-600 transition">特徴</a>
      <a href="#about" class="text-gray-600 hover:text-blue-600 transition">概要</a>
      <a href="#contact" class="text-gray-600 hover:text-blue-600 transition"
        >お問い合わせ</a
      >
    </nav>
    <div class="space-x-4">
      <BaseButton @click="goToLogin">ログイン</BaseButton>
      <BaseButton @click="goToRegister" class="bg-blue-500 text-white hover:bg-blue-600"
        >新規登録</BaseButton
      >
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import BaseButton from "@/components/base/BaseButton.vue";

const router = useRouter();
const goToLogin = () => router.push("/login");
const goToRegister = () => router.push("/register");
</script>
