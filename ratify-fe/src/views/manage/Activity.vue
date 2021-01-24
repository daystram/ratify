<template>
  <div class="activity">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">Activities</h1>
      </v-col>
    </v-row>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-timeline align-top dense>
            <div v-for="(activity, index) in activities" :key="index">
              <v-timeline-item v-if="activity.separator" class="pb-10" hide-dot>
                <span class="text-h5">
                  {{
                    activity.today
                      ? "Today"
                      : Intl.DateTimeFormat("default", {
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
          Failed retrieving activity log!
        </v-alert>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { STATUS } from "@/constants/status.ts";
import api from "@/apis/api";
import { authManager } from "@/auth/index.ts";
import { addDateSeparator, LogInfo, LogSeverityMap } from "@/utils/log.ts";

export default Vue.extend({
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    activities: new Array<{
      color: string;
      icon: string;
      title: string;
      subtitle: string;
      date: Date;
      separator?: boolean;
      today?: boolean;
      end?: boolean;
    }>()
  }),

  created() {
    api.log
      .userActivity()
      .then(response => {
        /* eslint-disable @typescript-eslint/camelcase */
        const logs: LogInfo[] = response.data.data;
        for (let i = 0; i < logs.length; i++) {
          const desc = JSON.parse(logs[i].description);
          const date = new Date(logs[i].created_at * 1000);
          addDateSeparator(date, this.activities);
          switch (desc.scope) {
            case "oauth::authorize":
              this.activities.push({
                color: ({
                  I: "success",
                  W: "error"
                } as LogSeverityMap)[logs[i].severity],
                icon: "mdi-lock",
                title: ({
                  I: "Successful Sign In",
                  W: "Failed Sign In Attempt"
                } as LogSeverityMap)[logs[i].severity],
                subtitle: ({
                  I: `Signed in from ${desc?.detail?.ip} via ${desc?.detail?.browser} at ${desc?.detail?.os}`,
                  W: `Incorrect credentials, attempted from ${desc?.detail?.ip} via ${desc?.detail?.browser} at ${desc?.detail?.os}`
                } as LogSeverityMap)[logs[i].severity],
                date: date
              });
              break;
            case "user::profile":
              this.activities.push({
                color: ({ I: "info", W: "error" } as LogSeverityMap)[
                  logs[i].severity
                ],
                icon: "mdi-account",
                title: "Profile Updated",
                subtitle: "",
                date: date
              });
              break;
            case "user::password":
              this.activities.push({
                color: ({ I: "info", W: "error" } as LogSeverityMap)[
                  logs[i].severity
                ],
                icon: "mdi-key",
                title: ({
                  I: "Password Updated",
                  W: "Failed Password Update Attempt"
                } as LogSeverityMap)[logs[i].severity],
                subtitle: ({
                  I: ``,
                  W: `Incorrect old password, attempted from ${desc?.detail?.ip} via ${desc?.detail?.browser} at ${desc?.detail?.os}`
                } as LogSeverityMap)[logs[i].severity],
                date: date
              });
              break;
            case "user::session":
              this.activities.push({
                color: ({ I: "warning", W: "error" } as LogSeverityMap)[
                  logs[i].severity
                ],
                icon: "mdi-dock-window",
                title: ({
                  I: "Session Revoked",
                  W: "Failed Revoking Session"
                } as LogSeverityMap)[logs[i].severity],
                subtitle: "",
                date: date
              });
              break;
            case "user::mfa":
              this.activities.push({
                color: desc.detail ? "primary" : "warning",
                icon: "mdi-two-factor-authentication",
                title: desc.detail ? "TOTP MFA Enabled" : "TOTP MFA Disabled",
                subtitle: "",
                date: date
              });
              break;
          }
        }
        const date = new Date((authManager.getUser()?.created_at || 0) * 1000);
        addDateSeparator(date, this.activities);
        this.activities.push({
          color: "success",
          icon: "mdi-account",
          title: "Account Created",
          subtitle: "",
          end: true,
          date: date
        });
        /* eslint-enable @typescript-eslint/camelcase */
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.ERROR;
      });
  }
});
</script>
