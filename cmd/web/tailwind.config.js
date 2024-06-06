/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["view/**/*.templ"],
  theme: {
    extend: {
      gridTemplateRows: {
        layout: "auto minmax(500px, 1fr) auto",
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        dim: {
          ...require("daisyui/src/theming/themes")["dim"],
          // primary: "#a991f7",
          primary: "#37cdbe",
          accent: "#fbd231",
        },
      },
    ],
  },
};
