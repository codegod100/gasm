/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./*.html",
    "./*_templates.go",
    "./*_templ.go"
  ],
  theme: {
    extend: {
      colors: {
        'chat-green': '#4caf50',
        'chat-green-hover': '#45a049',
        'chat-red': '#f44336',
        'chat-red-hover': '#da190b',
        'chat-blue': '#008cba',
        'chat-blue-hover': '#007b9a'
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}