<template>
  <div class="session">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">Sessions</h1>
      </v-col>
    </v-row>
    <v-expand-transition>
      <div v-show="formLoadStatus === STATUS.ERROR">
        <v-alert
          type="error"
          text
          dense
          transition="scroll-y-transition"
          class="mt-0"
        >
          Failed revoking session!
        </v-alert>
      </div>
    </v-expand-transition>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-divider inset />
          <div v-for="session in sessions" :key="session.session_id">
            <v-list-item>
              <v-list-item-content>
                <v-row justify="end" align="center" no-gutters>
                  <v-col cols="12" md="">
                    <v-row align="center">
                      <v-col cols="12" sm="auto">
                        <v-icon
                          x-large
                          v-text="
                            session.mobile ? 'mdi-cellphone' : 'mdi-laptop'
                          "
                        />
                      </v-col>
                      <v-col cols="">
                        <v-list-item-title class="text-h5">
                          <span
                            class="d-inline-block text-truncate"
                            style="max-width: 320px;"
                          >
                            {{ `${session.browser} at ${session.os}` }}
                          </span>
                        </v-list-item-title>
                        <v-list-item-subtitle>
                          <div
                            class="d-inline-block text-truncate"
                            style="max-width: 320px;"
                          >
                            From {{ session.ip }}
                          </div>
                        </v-list-item-subtitle>
                        <v-list-item-subtitle>
                          <div
                            class="d-inline-block text-truncate"
                            style="max-width: 320px;"
                          >
                            Issued at
                            {{
                              Intl.DateTimeFormat("default", {
                                dateStyle: "full",
                                timeStyle: "medium"
                              }).format(session.date)
                            }}
                          </div>
                        </v-list-item-subtitle>
                      </v-col>
                    </v-row>
                  </v-col>
                  <v-col cols="auto" md="auto">
                    <div
                      v-if="session.current"
                      class="text-button success--text px-4"
                    >
                      Current Session
                    </div>
                    <v-btn
                      v-else
                      rounded
                      text
                      outlined
                      color="error"
                      :disabled="sessionLoading.includes(session.session_id)"
                      :loading="sessionLoading.includes(session.session_id)"
                      @click="revoke(session)"
                    >
                      Revoke
                    </v-btn>
                  </v-col>
                </v-row>
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
          Failed retrieving sessions!
        </v-alert>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";
import { SessionInfo } from "@/utils/session";

export default Vue.extend({
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    formLoadStatus: STATUS.IDLE,
    sessions: new Array<SessionInfo>(),
    sessionLoading: new Array<string>()
  }),

  created() {
    api.session
      .list()
      .then(response => {
        /* eslint-disable @typescript-eslint/camelcase */
        const activeSessions: SessionInfo[] = response.data.data;
        activeSessions.sort((a, b) => b.issued_at - a.issued_at);
        for (let i = 0; i < activeSessions.length; i++) {
          const date = new Date(activeSessions[i].issued_at * 1000);
          this.sessions.push({
            ...activeSessions[i],
            date: date
          });
        }
        /* eslint-enable @typescript-eslint/camelcase */
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.ERROR;
      });
  },

  methods: {
    revoke(session: SessionInfo) {
      this.formLoadStatus = STATUS.IDLE;
      this.sessionLoading.push(session.session_id);
      api.session
        .revoke(session.session_id)
        .then(() => {
          this.sessions = this.sessions.filter(i => i !== session);
        })
        .catch(() => {
          this.formLoadStatus = STATUS.ERROR;
        })
        .finally(() => {
          this.sessionLoading.splice(
            this.sessionLoading.indexOf(session.session_id),
            1
          );
        });
    }
  }
});
</script>
