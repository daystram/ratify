import Vue from "vue";
import Vuetify from "vuetify/lib";

import "@mdi/font/css/materialdesignicons.css";

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    dark: true,
    themes: {
      dark: {
        primary: "#00c3c3",
        primaryDim: "#008686",
        secondary: "#c1842f"
      }
    },
    options: { customProperties: true }
  },
  icons: {
    iconfont: "mdi"
  }
});
