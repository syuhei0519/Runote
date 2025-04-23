<template>
  <div
    class="max-w-md mx-auto mt-20 p-6 border border-gray-200 rounded shadow-md bg-white dark:bg-gray-800"
  >
    <h2 class="text-xl font-semibold mb-4 text-center">ログイン</h2>
    <form @submit.prevent="onLogin">
      <FormLabel for="username">メールアドレス</FormLabel>
      <BaseInput
        id="username"
        v-model="username"
        type="username"
        placeholder="メールアドレス"
      />

      <FormLabel for="password">パスワード</FormLabel>
      <BaseInput
        id="password"
        v-model="password"
        type="password"
        placeholder="パスワード"
      />

      <div class="mt-4 mb-6">
        <BaseCheckbox v-model="rememberMe"> ログイン状態を保持する </BaseCheckbox>
      </div>

      <BaseButton type="submit">ログイン</BaseButton>
      <BackButton />
    </form>

    <!-- ログイン後にトースト表示（仮） -->
    <Toast :visible="toastVisible" message="ログインしました！" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import BaseInput from "@/components/base/BaseInput.vue";
import BaseButton from "@/components/base/BaseButton.vue";
import BaseCheckbox from "@/components/base/BaseCheckbox.vue";
import FormLabel from "@/components/base/FormLabel.vue";
import Toast from "@/components/base/Toast.vue";
import BackButton from "@/components/base/BackButton.vue";

const username = ref("");
const password = ref("");
const rememberMe = ref(false);
const toastVisible = ref(false);

const onLogin = async () => {
  // TODO: API 連携
  console.log("ログイン情報:", {
    username: username.value,
    password: password.value,
    rememberMe: rememberMe.value,
  });

  // トースト表示（仮）
  toastVisible.value = true;
  setTimeout(() => (toastVisible.value = false), 3000);
};
</script>
