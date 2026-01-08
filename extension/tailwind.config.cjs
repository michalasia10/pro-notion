/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{svelte,ts}"],
  theme: {
    extend: {
      colors: {
        notion: {
          bg: "#fafafa",
          panel: "#ffffff",
          border: "#e6e6e6",
          text: "#1f1f1f",
          muted: "#6b6b6b",
          accent: "#3b82f6"
        }
      }
    }
  },
  plugins: []
}
