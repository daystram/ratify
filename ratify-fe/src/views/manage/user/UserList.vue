<template>
  <div class="user-list">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm=""><h1 class="text-h2">Users</h1></v-col>
    </v-row>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-divider inset />
          <div v-for="user in users" :key="user.subject">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-content>
                  <v-row justify="end" align="center" no-gutters>
                    <v-col cols="12" md="">
                      <v-list-item-title class="text-h5">
                        <span
                          class="d-inline-block text-truncate"
                          style="max-width: 320px;"
                        >
                          {{ `${user.given_name} ${user.family_name}` }}
                        </span>
                      </v-list-item-title>
                      <v-list-item-subtitle>
                        <span
                          class="d-inline-block text-truncate"
                          style="max-width: 320px;"
                        >
                          {{ user.preferred_username }}
                        </span>
                      </v-list-item-subtitle>
                    </v-col>
                    <v-col cols="12" md="">
                      <span class="text-overline text-no-wrap mr-1">
                        User ID
                      </span>
                      <code>{{ user.sub }}</code>
                    </v-col>
                    <v-col cols="auto" md="auto">
                      <v-btn
                        rounded
                        text
                        outlined
                        color="secondary lighten-1"
                        :to="{
                          name: 'manage:user-detail',
                          params: { subject: user.sub }
                        }"
                      >
                        Manage
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-list-item-content>
              </v-list-item-content>
            </v-list-item>
            <v-divider inset />
          </div>
        </v-col>
      </v-row>
    </v-fade-transition>
    <v-fade-transition>
      <v-overlay
        v-show="
          pageLoadStatus === STATUS.PRE_LOADING ||
            pageLoadStatus === STATUS.LOADING
        "
        opacity="0"
        absolute
      >
        <v-progress-circular indeterminate size="64" />
      </v-overlay>
    </v-fade-transition>
    <v-expand-transition>
      <div v-show="pageLoadStatus === STATUS.ERROR">
        <v-alert
          type="error"
          text
          dense
          transition="scroll-y-transition"
          class="mt-0"
        >
          Failed retrieving application list!
        </v-alert>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";

export default Vue.extend({
  data: () => ({
    users: [],
    pageLoadStatus: STATUS.PRE_LOADING
  }),

  created() {
    this.loadUsers();
  },

  methods: {
    loadUsers() {
      this.pageLoadStatus = STATUS.PRE_LOADING;
      api.user
        .list()
        .then(response => {
          this.users = response.data.data;
          this.users.sort((a: { given_name: string }, b: { given_name: string }) =>
            a["given_name"].localeCompare(b["given_name"])
          );
          this.pageLoadStatus = STATUS.COMPLETE;
        })
        .catch(() => {
          this.pageLoadStatus = STATUS.ERROR;
        });
    }
  }
});
</script>
