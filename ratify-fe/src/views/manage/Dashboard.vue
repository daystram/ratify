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
        <v-card>
          <v-card-title>
            <v-row no-gutters align="center">
              <v-col cols="auto">
                Your Activty
              </v-col>
            </v-row>
          </v-card-title>
          <v-divider inset />
          <div class="v-card__body">
            <v-expand-transition>
              <v-row v-if="activity.formLoadStatus === STATUS.COMPLETE">
                <v-col cols="12" sm="6" md="3" order-md="1">
                  <div class="mb-1 text-overline text--secondary text-center">
                    Total Sign Ins
                  </div>
                  <div class="text-h4 text-center">
                    {{ activity.signInCount }}
                  </div>
                </v-col>
                <v-col cols="12" sm="12" md="6" order-sm="1" order-md="2">
                  <div class="mb-1 text-overline text--secondary text-center">
                    Last Signed In
                  </div>
                  <div class="text-h4 text-center">
                    {{
                      Intl.DateTimeFormat("default", {
                        dateStyle: "medium",
                        timeStyle: "short"
                      }).format(activity.lastSignIn)
                    }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6" md="3" order-md="3">
                  <div class="mb-1 text-overline text--secondary text-center">
                    Active Sessions
                  </div>
                  <div class="text-h4 text-center">
                    {{ activity.sessionCount }}
                  </div>
                </v-col>
              </v-row>
            </v-expand-transition>
            <v-expand-transition>
              <div v-if="activity.formLoadStatus === STATUS.ERROR">
                <v-alert
                  type="error"
                  text
                  dense
                  transition="scroll-y-transition"
                  class="mb-0"
                >
                  Failed retrieving your activity!
                </v-alert>
              </div>
            </v-expand-transition>
          </div>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row no-gutters align="center">
              <v-col cols="auto">
                Updates
              </v-col>
            </v-row>
          </v-card-title>
          <v-divider inset />
          <div class="v-card__body">
            <v-expand-transition>
              <v-row v-if="activity.formLoadStatus === STATUS.COMPLETE">
                <v-col v-if="!this.updates.length">
                  <div class="text--disabled font-italic text-center">
                    You're all set!
                  </div>
                </v-col>
                <div
                  v-for="(update, i) in updates"
                  :key="i"
                  style="width: 100%"
                >
                  <v-col>
                    <v-alert
                      :type="update.severity"
                      text
                      prominent
                      dense
                      :icon="update.icon"
                      transition="scroll-y-transition"
                      class="ma-0"
                    >
                      <h3 class="text-h6">
                        {{ update.title }}
                      </h3>
                      <div>
                        {{ update.detail }}
                      </div>
                    </v-alert>
                  </v-col>
                </div>
              </v-row>
            </v-expand-transition>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { authManager } from "@/auth";
import { STATUS } from "@/constants/status";
import api from "@/apis/api";

export default Vue.extend({
  data() {
    return {
      activity: {
        formLoadStatus: STATUS.LOADING,
        signInCount: 0,
        lastSignIn: new Date(),
        sessionCount: 0
      },
      updates: new Array<{
        severity: string;
        title: string;
        detail: string;
        icon: string;
      }>()
    };
  },

  computed: {
    user: () => authManager.getUser()
  },

  created() {
    api.dashboard
      .fetch()
      .then(response => {
        this.activity.formLoadStatus = STATUS.COMPLETE;
        this.activity.signInCount = response.data.data.signin_count;
        this.activity.lastSignIn = new Date(
          response.data.data.last_signin * 1000
        );
        this.activity.sessionCount = response.data.data.session_count;
        if (response.data.data.recent_failure) {
          this.updates.push({
            severity: "error",
            title: "Failed Sign In Attempt",
            detail:
              "There has been a recent failed sign in attempt to your account. Go to Activities page to view more.",
            icon: "mdi-account-alert"
          });
        }
        if (!response.data.data.mfa_enabled) {
          this.updates.push({
            severity: "warning",
            title: "MFA Disabled",
            detail:
              "Multi-factor authentication is not enabled. Enable in your Profile page to ehance your account security.",
            icon: "mdi-two-factor-authentication"
          });
        }
      })
      .catch(() => {
        this.activity.formLoadStatus = STATUS.ERROR;
      });
  }
});
</script>
