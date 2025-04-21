// frontend/.eslintrc.cjs
module.exports = {
    root: true,
    env: {
      browser: true,
      es2021: true,
      node: true
    },
    extends: [
      'plugin:vue/vue3-recommended',
      '@vue/eslint-config-typescript',
      'prettier'
    ],
    parserOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module'
    },
    rules: {
      // シングルファイルコンポーネント名の強制を無効化
      'vue/multi-word-component-names': 'off',
  
      // v-html を許可（XSS注意だが自己責任で使う場合）
      'vue/no-v-html': 'off',
  
      // 未使用の変数で、引数が _ で始まる場合は警告なし
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }]
    }
}