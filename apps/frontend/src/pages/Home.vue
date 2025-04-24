<template>
  <DefaultGuestLayout>
    <section
      class="relative min-h-screen bg-cover bg-center flex items-center justify-center px-6"
      style="background-image: url('/images/runner-7876675_1280.jpg')"
    >
      <!-- オーバーレイ -->
      <div class="absolute inset-0 bg-black/60"></div>

      <!-- コンテンツ -->
      <div class="relative z-10 text-center text-white max-w-xl">
        <h1 class="text-5xl font-extrabold mb-6 drop-shadow-lg">
          🏃‍♂️ Runote へようこそ！
        </h1>
        <p class="text-lg md:text-xl mb-8 text-gray-200">
          ランニングと感情を記録する、<br />新しい日記体験を始めましょう。
        </p>

        <div class="flex justify-center gap-4 flex-wrap">
          <BaseButton
            @click="goToLogin"
            class="bg-blue-500 hover:bg-blue-600 text-white shadow-lg"
          >
            ログイン
          </BaseButton>
          <BaseButton
            @click="goToRegister"
            class="bg-white text-blue-600 hover:bg-gray-100 shadow-lg"
          >
            新規登録
          </BaseButton>
        </div>
      </div>
    </section>
    <section class="bg-white py-20 px-6">
      <div class="max-w-6xl mx-auto text-center">
        <h2 class="text-3xl font-bold mb-4">Runote の特徴</h2>
        <p class="text-gray-600 mb-12">
          あなたのランニングを、よりパーソナルで豊かなものに。
        </p>

        <div class="grid md:grid-cols-3 gap-8 text-left">
          <!-- 特徴① -->
          <div
            class="bg-gray-50 rounded-2xl p-6 shadow hover:shadow-lg transition"
            data-aos="fade-up"
            data-aos-delay="100"
          >
            <div class="mb-4">
              <img src="/icons/pencil.svg" alt="ラン記録" class="h-10 w-10" />
            </div>
            <h3 class="text-xl font-semibold mb-2">ラン記録</h3>
            <p class="text-gray-600 text-sm">
              距離・時間・ペースをワンタップで記録。カレンダー表示で振り返りも簡単。
            </p>
          </div>

          <!-- 特徴② -->
          <div
            class="bg-gray-50 rounded-2xl p-6 shadow hover:shadow-lg transition"
            data-aos="fade-up"
            data-aos-delay="200"
          >
            <div class="mb-4">
              <img src="/icons/pencil.svg" alt="感情ログ" class="h-10 w-10" />
            </div>
            <h3 class="text-xl font-semibold mb-2">感情ログ</h3>
            <p class="text-gray-600 text-sm">
              走った後の気持ちも記録。嬉しさも悔しさも、未来の自分へのメッセージに。
            </p>
          </div>

          <!-- 特徴③ -->
          <div
            class="bg-gray-50 rounded-2xl p-6 shadow hover:shadow-lg transition"
            data-aos="fade-up"
            data-aos-delay="300"
          >
            <div class="mb-4">
              <img src="/icons/pencil.svg" alt="振り返り" class="h-10 w-10" />
            </div>
            <h3 class="text-xl font-semibold mb-2">分析・振り返り</h3>
            <p class="text-gray-600 text-sm">
              週間・月間での推移をグラフ化。自分の成長を“見える化”します。
            </p>
          </div>
        </div>
      </div>
    </section>
  </DefaultGuestLayout>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import BaseButton from "@/components/base/BaseButton.vue";
import DefaultGuestLayout from "@/layouts/DefaultGuestLayout.vue";

const router = useRouter();
const goToLogin = () => router.push("/login");
const goToRegister = () => router.push("/register");
</script>
