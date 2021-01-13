<template>
  <div class="activity">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">Activities</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-timeline align-top dense>
            <v-timeline-item class="pb-10" hide-dot>
              <span class="text-h5">Today</span>
            </v-timeline-item>
            <div v-for="(activity, index) in activities" :key="index">
              <v-timeline-item v-if="activity.separator" class="pb-10" hide-dot>
                <span class="text-h5">
                  {{
                    Intl.DateTimeFormat("default", {
                      dateStyle: "full"
                    }).format(activity.date)
                  }}
                </span>
              </v-timeline-item>
              <v-timeline-item
                v-else
                :color="activity.color"
                :icon="activity.icon"
                :class="activity.end ? 'pb-0' : ' pb-10'"
                fill-dot
              >
                <v-row class="pt-1" dense>
                  <v-col cols="" sm="3">
                    <div class="text-body-1" style="line-height: 32px">
                      {{ activity.date.toLocaleTimeString() }}
                    </div>
                  </v-col>
                  <v-col>
                    <div class="text-h6">{{ activity.title }}</div>
                    <div class="text-subtitle-1 text--secondary">
                      {{ activity.subtitle }}
                    </div>
                  </v-col>
                </v-row>
              </v-timeline-item>
            </div>
          </v-timeline>
        </v-col>
      </v-fade-transition>
    </v-row>
    <v-fade-transition>
      <v-overlay
        opacity="0"
        absolute
        style="height: calc(100vh - 64px)"
        v-show="pageLoadStatus !== STATUS.COMPLETE"
      >
        <v-progress-circular indeterminate size="64" />
      </v-overlay>
    </v-fade-transition>
  </div>
</template>

<script>
import Vue from "vue";
import { STATUS } from "@/constants/status.ts";
import api from "@/apis/api";
import { authManager } from "@/auth/index.ts";

export default Vue.extend({
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    activities: []
  }),

  created() {
    api.log.userActivity().then(response => {
      /* eslint-disable @typescript-eslint/camelcase */
      const logs = response.data.data;
      for (let i = 0; i < logs.length; i++) {
        const desc = JSON.parse(logs[i].description);
        const date = new Date(logs[i].created_at * 1000);
        if (
          i &&
          new Date(date.toDateString()) <
            new Date(
              this.activities[this.activities.length - 1].date.toDateString()
            )
        ) {
          this.activities.push({
            separator: true,
            date: date
          });
        }
        switch (desc.scope) {
          case "oauth::authorize":
            this.activities.push({
              color: {
                I: "success",
                W: "error"
              }[logs[i].severity],
              icon: "mdi-lock",
              title: {
                I: "Successful Sign In",
                W: "Failed Sign In Attempt"
              }[logs[i].severity],
              subtitle: {
                I: `Signed in from ${desc.detail.ip} via ${desc.detail.browser} at ${desc.detail.os}`,
                W: `Incorrect credentials. Attempted from ${desc.detail.ip} via ${desc.detail.browser} at ${desc.detail.os}.`
              }[logs[i].severity],
              date: date
            });
            break;
          case "user::profile":
            this.activities.push({
              color: { I: "info", W: "error" }[logs[i].severity],
              icon: "mdi-account",
              title: "Profile Updated",
              subtitle: "",
              date: date
            });
            break;
        }
      }
      const date = new Date(authManager.getUser().created_at * 1000);
      if (
        this.activities.length &&
        new Date(date.toDateString()) <
          new Date(
            this.activities[this.activities.length - 1].date.toDateString()
          )
      ) {
        this.activities.push({
          separator: true,
          date: date
        });
      }
      this.activities.push({
        color: "info",
        icon: "mdi-account",
        title: "Account Created",
        end: true,
        date: date
      });
      /* eslint-enable @typescript-eslint/camelcase */
      this.pageLoadStatus = STATUS.COMPLETE;
    });
  }
});
</script>
