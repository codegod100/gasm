/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./*.html",
    "./*_templates.go",
    "./*_templ.go",
    "./main.go"
  ],
  theme: {
    extend: {
      colors: {
        // Catppuccin Mocha palette
        'rosewater': '#f5e0dc',
        'flamingo': '#f2cdcd',
        'pink': '#f5c2e7',
        'mauve': '#cba6f7',
        'red': '#f38ba8',
        'maroon': '#eba0ac',
        'peach': '#fab387',
        'yellow': '#f9e2af',
        'green': '#a6e3a1',
        'teal': '#94e2d5',
        'sky': '#89dceb',
        'sapphire': '#74c7ec',
        'blue': '#89b4fa',
        'lavender': '#b4befe',
        'text': '#cdd6f4',
        'subtext1': '#bac2de',
        'subtext0': '#a6adc8',
        'overlay2': '#9399b2',
        'overlay1': '#7f849c',
        'overlay0': '#6c7086',
        'surface2': '#585b70',
        'surface1': '#45475a',
        'surface0': '#313244',
        'base': '#1e1e2e',
        'mantle': '#181825',
        'crust': '#11111b',
        
        // Legacy colors for compatibility
        'chat-green': '#a6e3a1',
        'chat-green-hover': '#94e2d5',
        'chat-red': '#f38ba8',
        'chat-red-hover': '#eba0ac',
        'chat-blue': '#89b4fa',
        'chat-blue-hover': '#74c7ec'
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}