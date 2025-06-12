/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      fontFamily: {
        'quicksand': ['Quicksand', 'sans-serif'],
        'perfect-dos-vga': ['PerfectDosVga', 'sans-serif'],
      },
    },
  },
  plugins: [],
} 