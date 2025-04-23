<template>
  <form @submit.prevent="confirm">
    <FormLabel for="name">ユーザー名</FormLabel>
    <BaseInput id="name" v-model="store.name" type="text" placeholder="ユーザー名" />

    <FormLabel for="password">パスワード</FormLabel>
    <BaseInput
      id="password"
      v-model="store.password"
      type="password"
      placeholder="パスワード"
    />

    <FormLabel for="confirm">パスワード確認</FormLabel>
    <BaseInput
      id="confirm"
      v-model="confirmPassword"
      type="password"
      placeholder="もう一度パスワードを入力"
    />

    <p v-if="confirmPassword && !isPasswordMatch" class="text-sm text-red-500 mt-1">
      パスワードが一致しません
    </p>
    <p v-if="confirmPassword && isPasswordMatch" class="text-sm text-green-600 mt-1">
      ✅ パスワードが一致しました
    </p>

    <BaseButton type="submit" class="mt-6 w-full" :disabled="!isPasswordMatch">
      次へ
    </BaseButton>
    <BackButton />
  </form>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRegisterUserStore } from "@/stores/registerUser";
import { useRouter } from "vue-router";
import BaseInput from "@/components/base/BaseInput.vue";
import BaseButton from "@/components/base/BaseButton.vue";
import FormLabel from "@/components/base/FormLabel.vue";
import BackButton from "@/components/base/BackButton.vue";

const store = useRegisterUserStore();
const router = useRouter();

const confirmPassword = ref("");
const isPasswordMatch = computed(
  () =>
    store.password && confirmPassword.value && store.password === confirmPassword.value
);

const confirm = () => {
  if (!isPasswordMatch.value) return;
  router.push("/register/confirm");
};
</script>
