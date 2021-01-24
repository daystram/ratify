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
              </v-row>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-fade-transition>
    <!-- <v-fade-transition>
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
                    Delete application
                  </div>
                  <div class="text--secondary">
                    You cannot un-delete an application. Take extreme caution.
                  </div>
                </v-col>
                <v-col cols="auto">
                  <v-dialog
                    v-model="deleting.prompt"
                    width="500"
                    @input="v => v || cancelDelete()"
                  >
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn
                        rounded
                        outlined
                        text
                        color="error"
                        v-bind="attrs"
                        v-on="on"
                        :disabled="detail.locked"
                      >
                        Delete
                      </v-btn>
                    </template>
                    <v-card class="danger-border">
                      <v-card-title>
                        <v-row no-gutters align="center">
                          <v-col cols="auto" class="error--text">
                            Delete Application
                          </v-col>
                          <v-spacer />
                          <v-col cols="auto">
                            <v-btn text icon color="grey" @click="cancelDelete">
                              <v-icon v-text="'mdi-close'" />
                            </v-btn>
                          </v-col>
                        </v-row>
                      </v-card-title>
                      <v-divider inset />
                      <div class="v-card__body">
                        <v-alert type="warning" text dense>
                          You are about to delete this application!
                        </v-alert>
                        <v-row align="center">
                          <v-col>
                            <div>
                              Are you sure you want to permanently delete
                              <b>{{ detail.name }}</b
                              >? This is action is <b>irreversible</b> and all
                              of this application's clients will not be able to
                              user Ratify authentication service.
                            </div>
                            <div class="mt-4">
                              Type <b>{{ detail.name }}</b> to confirm.
                            </div>
                            <div>
                              <v-text-field
                                v-model="deleting.confirmName"
                                class="py-2"
                                :prepend-icon="'mdi-application'"
                              />
                            </div>
                            <v-btn
                              rounded
                              block
                              outlined
                              color="error"
                              :disabled="deleting.confirmName !== detail.name"
                              @click="confirmDelete"
                            >
                              Delete
                            </v-btn>
                          </v-col>
                        </v-row>
                      </div>
                    </v-card>
                  </v-dialog>
                </v-col>
              </v-row>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-fade-transition> -->
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

export default Vue.extend({
  data: () => ({
    metric: {
      lastSignIn: 0,
      signInCount: 0
    },
    detail: {},
    pageLoadStatus: STATUS.PRE_LOADING
  }),

  created() {
    api.user
      .detail(this.$route.params.subject)
      .then(response => {
        this.metric.lastSignIn = response.data.data.last_signin;
        this.metric.signInCount = response.data.data.signin_count;
        this.detail = response.data.data;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(error => {
        if (error.response.status === 404) {
          this.$router.push({ name: "manage:user" });
          return;
        }
        this.pageLoadStatus = STATUS.ERROR;
      });
  }
});
</script>
