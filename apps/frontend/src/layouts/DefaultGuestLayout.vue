<template>
  <div class="flex flex-col min-h-screen">
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from "@/components/layout/AppHeaderGuest.vue";
import AppFooter from "@/components/layout/AppFooterGuest.vue";
</script>
