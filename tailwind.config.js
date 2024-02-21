/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/*.tmpl'],
  theme: {
    extend: {
      colors: {
        'light-bg': '#fff7f0',
        'lighter-bg': '#fafaf9',

        'primary-btn': '#ec4899',
        'dark-primary-btn': '#be185d',
        'secondary-btn': '#44403c',

        'primary-txt': '#ec4899',
        'dark-txt': '#374151',
        'darker-txt': '#111827',

        'primary-br': '#be185d',
        'dark-br': '#e5e7eb',
      },
    },
  },
  plugins: [],
}
