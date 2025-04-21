import type { Config } from 'tailwindcss'

const config: Config = {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: '#3B82F6',  // blue-500
        secondary: '#10B981', // emerald-500
        runoteGray: '#6B7280', // gray-500
      },
    },
  },
  plugins: [],
}

export default config