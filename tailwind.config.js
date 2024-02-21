/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/*.tmpl'],
  theme: {
    extend: {
      colors: {
        'light-bg': '#fff7f0', //custom
        'lighter-bg': '#f9fafb', // gray-50

        'primary-btn': '#ec4899', // pink-500
        'dark-primary-btn': '#be185d', // pink-600
        'secondary-btn': '#6b7280', // gray-500
        'dark-secondary-btn': '#4b5563', // gray-600

        'primary-txt': '#ec4899', // pink-500
        'dark-txt': '#374151', // gray-700
        'darker-txt': '#111827', // gray-900

        'primary-br': '#db2777', //pink-600
        'dark-br': '#e5e7eb', // gray-200
      },
    },
  },
  plugins: [],
}
