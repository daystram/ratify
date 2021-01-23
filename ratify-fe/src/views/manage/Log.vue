<template>
  <div class="log">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">Logs</h1>
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
import { STATUS } from "@/constants/status";
import api from "@/apis/api";
import { addDateSeparator, LogInfo, LogSeverityMap } from "@/utils/log";

export default Vue.extend({
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    logs: new Array<{
      color: string;
      icon: string;
      title: string;
      subtitle: string;
      date: Date;
      separator?: boolean;
      today?: boolean;
    }>()
  }),

  created() {
    api.log
      .adminActivity()
      .then(response => {
        /* eslint-disable @typescript-eslint/camelcase */
        const logs: LogInfo[] = response.data.data;
        for (let i = 0; i < logs.length; i++) {
          const desc = JSON.parse(logs[i].description);
          const date = new Date(logs[i].created_at * 1000);
          addDateSeparator(date, this.logs);
          switch (desc.scope) {
            case "application::detail":
              this.logs.push({
                color: "info",
                icon: "mdi-application",
                title: "Application Detail Updated",
                subtitle: `${logs[i].application_name} updated by ${logs[i].preferred_username}`,
                date: date
              });
              break;
            case "application::create":
              this.logs.push({
                color: ({ I: "success", W: "error" } as LogSeverityMap)[
                  logs[i].severity
                ],
                icon: "mdi-application",
                title: ({
                  I: "Added New Application",
                  W: "Removed Application"
                } as LogSeverityMap)[logs[i].severity],
                subtitle: ({
                  I: `${logs[i]?.application_name} created by ${logs[i].preferred_username}`,
                  W: `${desc?.detail?.name} deleted by ${logs[i].preferred_username}`
                } as LogSeverityMap)[logs[i].severity],
                date: date
              });
              break;
            case "application::secret":
              this.logs.push({
                color: "warning",
                icon: "mdi-key",
                title: "Application Secret Key Revoked",
                subtitle: `Secret key for ${logs[i].application_name} revoked by ${logs[i].preferred_username}`,
                date: date
              });
              break;
          }
        }
        /* eslint-enable @typescript-eslint/camelcase */
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.ERROR;
      });
  }
});
</script>
