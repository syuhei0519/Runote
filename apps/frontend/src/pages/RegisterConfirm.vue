<template>
  <div
    class="max-w-md mx-auto mt-20 p-6 border border-gray-200 rounded shadow-md bg-white dark:bg-gray-800"
  >
    <h2 class="text-xl font-bold mb-6 text-center">登録内容の確認</h2>

    <ul class="mb-6 space-y-3 text-gray-800 dark:text-gray-100">
      <li><strong>ユーザー名：</strong>{{ store.name }}</li>
      <li><strong>パスワード：</strong>（非表示）</li>
    </ul>

    <div class="flex justify-between">
      <BaseButton @click="goBack">戻る</BaseButton>
      <BaseButton @click="submit" :disabled="isSubmitting">
        <span v-if="isSubmitting">送信中...</span>
        <span v-else>登録する</span>
      </BaseButton>
    </div>

    <Toast :visible="showToast" message="登録が完了しました！" />
  </div>
</template>

<script setup lang="ts">
import { useRegisterUserStore } from "@/stores/registerUser";
import { useRouter } from "vue-router";
import { ref } from "vue";
import BaseButton from "@/components/base/BaseButton.vue";
import Toast from "@/components/base/Toast.vue";

const store = useRegisterUserStore();
const router = useRouter();

const isSubmitting = ref(false);
const showToast = ref(false);

const goBack = () => router.push("/register");

const submit = async () => {
  isSubmitting.value = true;

  try {
    console.log("送信データ:", store.name, store.password);
    // await apiClient.post('/auth/register', store.$state)

    showToast.value = true;
    setTimeout(() => {
      showToast.value = false;
      router.push("/login");
    }, 2000);
  } catch (err) {
    alert("登録失敗");
  } finally {
    isSubmitting.value = false;
  }
};
</script>
