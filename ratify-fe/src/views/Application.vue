<template>
  <div class="application">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm=""><h1 class="text-h2">Applications</h1></v-col>
      <v-col cols="auto">
        <v-btn color="primary" rounded>
          New Application
        </v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-divider :inset="false" />
        <div v-for="application in applications" :key="application">
          <v-list-item>
            <v-list-item-content>
              <v-list-item-content>
                <v-row justify="center" align="center" no-gutters>
                  <v-col>
                    <v-list-item-title class="text-h5">
                      {{ application.name }}
                    </v-list-item-title>
                    <v-list-item-subtitle>
                      {{ application.description }}
                    </v-list-item-subtitle>
                  </v-col>
                  <v-col>
                    <span class="text-overline text-no-wrap mr-1">
                      Client ID
                    </span>
                    <code>{{ application.client_id }}</code>
                  </v-col>
                  <v-col cols="auto">
                    <v-btn rounded text color="secondary lighten-1"
                      >manage</v-btn
                    >
                  </v-col>
                </v-row>
              </v-list-item-content>
            </v-list-item-content>
          </v-list-item>
          <v-divider :inset="false" />
        </div>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";

export default Vue.extend({
  data: function() {
    return {
      applications: [],
      pageLoadStatus: STATUS.LOADING
    };
  },

  created() {
    api.application
      .getAll()
      .then(response => {
        console.log(response.data.data);
        this.applications = response.data.data;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(error => {
        console.error(error);
        this.pageLoadStatus = STATUS.ERROR;
      });
  }
});
</script>
