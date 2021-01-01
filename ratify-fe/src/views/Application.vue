<template>
  <div class="application">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm=""><h1 class="text-h2">Applications</h1></v-col>
      <v-col cols="auto">
        <v-btn color="primary darken-2" rounded>
          New Application
        </v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-divider :inset="false" />
          <div v-for="application in applications" :key="application.client_id">
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
      </v-fade-transition>
    </v-row>
    <v-fade-transition>
      <v-overlay
        v-show="pageLoadStatus !== STATUS.COMPLETE"
        opacity="0"
        absolute
      >
        <v-progress-circular indeterminate size="64"></v-progress-circular>
      </v-overlay>
    </v-fade-transition>
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
      pageLoadStatus: STATUS.PRE_LOADING
    };
  },

  created() {
    api.application
      .getAll()
      .then(response => {
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
