<template>
  <v-container fill-height fluid class="gradient-bg">
    <v-col>
      <v-row>
        <v-col>
          <h1 class="text-sm-h1 text-h3 text-center">
            <Logo />
          </h1>
          <p
            class="mt-12 mb-12 text-subtitle-1 ma-auto text-center app-subtitle"
          >
            Ratify is a Central Authentication Service (CAS) implementing OAuth
            2.0 and OpenID Connect (OID) protocols.
          </p>
        </v-col>
      </v-row>
      <v-row align="center" justify="center">
        <div v-if="!isAuthenticated">
          <v-btn
            :to="{ name: 'login' }"
            elevation="6"
            rounded
            x-large
            class="ma-2"
            v-text="'Sign In'"
          />
          <v-btn
            :to="{ name: 'signup' }"
            elevation="6"
            rounded
            x-large
            color="primaryDim"
            class="ma-2"
            v-text="'Sign Up'"
          />
        </div>
        <div v-else>
          <v-btn
            :to="{ name: 'manage:dashboard' }"
            elevation="6"
            rounded
            x-large
            color="primaryDim"
            class="ma-2"
            v-text="'Manage'"
          />
        </div>
      </v-row>
    </v-col>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Logo from "@/components/Logo.vue";
import { isAuthenticated } from "@/auth";

export default Vue.extend({
  components: { Logo },

  data: function() {
    return {
      isAuthenticated: false
    };
  },

  updated: function() {
    this.isAuthenticated = isAuthenticated();
  },

  mounted: function() {
    this.isAuthenticated = isAuthenticated();
  }
});
</script>
