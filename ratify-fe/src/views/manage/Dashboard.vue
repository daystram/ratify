<template>
  <div class="dashboard">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">
          Hi, {{ user.given_name }} {{ user.family_name }}
        </h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        MOCK DASHBOARD
        <Home />
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Home } from "../index";
import { authManager } from "@/auth";

export default Vue.extend({
  components: { Home },
  computed: {
    profile: () => authManager.getUser()
  }
});
</script>
