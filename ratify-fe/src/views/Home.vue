<template>
  <v-container fill-height fluid class="gradient-bg">
    <v-col>
      <v-row>
        <v-col>
          <h1 class="text-sm-h1 text-h3 text-center">
            <v-icon style="font-size: inherit">mdi-lock</v-icon>Ratify
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
          <v-btn to="/login" elevation="6" rounded x-large class="ma-2"
            >Sign In</v-btn
          >
          <v-btn
            to="/signup"
            elevation="6"
            rounded
            x-large
            color="primaryDim"
            class="ma-2"
            >Sign Up</v-btn
          >
        </div>
        <div v-else>
          <v-btn to="/manage">Manage</v-btn>
        </div>
      </v-row>
    </v-col>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { isAuthenticated } from "@/auth";

export default Vue.extend({
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
