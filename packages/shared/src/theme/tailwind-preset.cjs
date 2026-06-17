/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        brand: {
          red: "#E84142",
          "red-dark": "#C5302F",
        },
        surface: {
          app: "#0A0A0B",
          base: "#0F0F11",
          card: "#141416",
          raised: "#1F1F22",
          border: "#27272A",
        },
        status: {
          draft: "#71717A",
          submitted: "#F59E0B",
          changes: "#FB923C",
          resubmitted: "#FBBF24",
          approved: "#38BDF8",
          banner: "#A78BFA",
          published: "#34D399",
          "payment-initiated": "#2DD4BF",
          "payment-confirmed": "#10B981",
        },
      },
      fontFamily: {
        sans: [
          "Inter",
          "ui-sans-serif",
          "system-ui",
          "-apple-system",
          "Segoe UI",
          "Roboto",
          "sans-serif",
        ],
      },
      boxShadow: {
        "glow-red": "0 0 40px -8px rgba(232, 65, 66, 0.45)",
      },
      borderRadius: {
        xl2: "1rem",
      },
    },
  },
  plugins: [],
};
