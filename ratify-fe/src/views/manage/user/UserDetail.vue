<template>
  <div class="user-detail">
    <v-row align="center">
      <v-col cols="12" sm="">
        <v-btn plain :ripple="false" class="pa-0" :to="{ name: 'manage:user' }">
          <v-icon v-text="'mdi-arrow-left'" class="mr-1" />
          Back
        </v-btn>
      </v-col>
    </v-row>
    <v-fade-transition>
      <v-row
        v-show="pageLoadStatus === STATUS.COMPLETE"
        class="mb-8"
        align="center"
      >
        <v-col cols="12">
          <h1 class="text-h2 text-truncate pb-2">
            {{ `${detail.given_name} ${detail.family_name}` }}
          </h1>
          <div class="text-subtitle-1 text-truncate text--secondary">
            {{ detail.preferred_username }}
          </div>
        </v-col>
      </v-row>
    </v-fade-transition>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-card>
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Metrics
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-row>
                <v-col cols="12" sm="4">
                  <div class="mb-1 text-overline text--secondary text-center">
                    Total Sign Ins
                  </div>
                  <div class="text-h4 text-center">
                    {{ metric.signInCount }}
                  </div>
                </v-col>
                <v-col cols="12" sm="8">
                  <div class="mb-1 text-overline text--secondary text-center">
                    Last Authorized
                  </div>
                  <div class="text-h4 text-center">
                    {{
                      metric.lastSignIn
                        ? Intl.DateTimeFormat("default", {
                            dateStyle: "medium",
                            timeStyle: "short"
                          }).format(new Date(metric.lastSignIn * 1000))
                        : "Never"
                    }}
                  </div>
                </v-col>
              </v-row>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-fade-transition>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-card :loading="detail.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Details
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-row>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    Given Name
                  </div>
                  <div>
                    {{ detail.given_name }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    Family Name
                  </div>
                  <div>
                    {{ detail.family_name }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    Username
                  </div>
                  <div>
                    {{ detail.preferred_username }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    Email
                  </div>
                  <div>
                    {{ detail.email }}
                    <v-chip
                      v-if="detail.email_verified"
                      color="success"
                      outlined
                      pill
                      small
                      class="ml-4"
                    >
                      Verified
                    </v-chip>
                    <v-chip
                      v-else
                      color="warning"
                      outlined
                      pill
                      small
                      class="ml-4"
                    >
                      Unverified
                    </v-chip>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    MFA
                  </div>
                  <div>
                    <v-chip
                      v-if="detail.mfa_enabled"
                      color="success"
                      outlined
                      pill
                      small
                    >
                      TOTP Enabled
                    </v-chip>
                    <v-chip v-else color="error" outlined pill small>
                      Disabled
                    </v-chip>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    User ID
                  </div>
                  <div>
                    <code>{{ detail.sub }}</code>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="mb-1 text-overline text--secondary">
                    Created At
                  </div>
                  <div>
                    {{
                      detail.created_at &&
                        Intl.DateTimeFormat("default", {
                          dateStyle: "full",
                          timeStyle: "medium"
                        }).format(new Date(detail.created_at * 1000))
                    }}
                  </div>
                </v-col>
              </v-row>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-fade-transition>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-card class="danger-border">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto" class="error--text">
                  Danger Zone
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-row justify="end" align="center">
                <v-col cols="">
                  <div>
                    Ratify Admin
                  </div>
                  <div class="text--secondary">
                    An admin is able to manage applications and users. When
                    changing a user's admin status, they will be logged out of
                    all active sessions.
                  </div>
                </v-col>
                <v-col cols="auto">
                  <v-switch
                    v-model="superuser.superuser"
                    :disabled="
                      superuser.disabled ||
                        superuser.formLoadStatus === STATUS.LOADING
                    "
                    :loading="superuser.formLoadStatus === STATUS.LOADING"
                    inset
                    color="warning"
                    @change="updateSuperuser"
                  />
                </v-col>
              </v-row>
            </div>
          </v-card>
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
          Failed retrieving application detail!
        </v-alert>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { STATUS } from "@/constants/status";
import api from "@/apis/api";
import { authManager } from "@/auth";

export default Vue.extend({
  data: () => ({
    metric: {
      lastSignIn: 0,
      signInCount: 0
    },
    detail: {},
    superuser: {
      superuser: false,
      disabled: false,
      formLoadStatus: STATUS.IDLE
    },
    pageLoadStatus: STATUS.PRE_LOADING
  }),

  created() {
    api.user
      .detail(this.$route.params.subject)
      .then(response => {
        this.metric.lastSignIn = response.data.data.last_signin;
        this.metric.signInCount = response.data.data.signin_count;
        this.detail = response.data.data;
        this.superuser.superuser = response.data.data.superuser;
        this.superuser.disabled =
          authManager.getUser()?.sub === response.data.data.sub;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(error => {
        if (error.response.status === 404) {
          this.$router.push({ name: "manage:user" });
          return;
        }
        this.pageLoadStatus = STATUS.ERROR;
      });
  },

  methods: {
    updateSuperuser() {
      this.superuser.formLoadStatus = STATUS.LOADING;
      setTimeout(() => {
        api.user
          .updateSuperuser({
            sub: (this.detail as { sub: string }).sub,
            superuser: this.superuser.superuser
          })
          .catch(() => (this.superuser.superuser = !this.superuser.superuser))
          .finally(() => (this.superuser.formLoadStatus = STATUS.IDLE));
      }, 2000);
    }
  }
});
</script>
