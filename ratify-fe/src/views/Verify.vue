<template>
  <v-container fluid class="fill-height gradient-bg">
    <v-col>
      <v-row align="center" justify="center">
        <div v-if="pageLoadStatus !== STATUS.PRE_LOADING" class="login-form">
          <v-scroll-y-transition appear>
            <div>
              <v-card
                elevation="12"
                :loading="formLoadStatus === STATUS.LOADING"
              >
                <div class="pa-6">
                  <div class="text-h2 mt-12 mb-12 text-center">
                    Verify Email
                  </div>
                  <v-expand-transition>
                    <div
                      v-show="
                        !verificationRequest.token &&
                          formLoadStatus === STATUS.COMPLETE
                      "
                    >
                      <v-alert
                        type="success"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Verification email sent!
                      </v-alert>
                      <v-btn
                        plain
                        :ripple="false"
                        class="pa-0"
                        :to="{ name: 'login' }"
                      >
                        <v-icon v-text="'mdi-arrow-left'" class="mr-1" />
                        Back to Login
                      </v-btn>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div
                      v-show="
                        verificationRequest.token &&
                          pageLoadStatus === STATUS.COMPLETE
                      "
                      class="align-content-end"
                    >
                      <v-alert
                        type="success"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Email successfully verified!
                      </v-alert>
                      <div class="text-end">
                        <v-btn
                          plain
                          :ripple="false"
                          class="pa-0"
                          :to="{ name: 'login' }"
                        >
                          Continue to Login
                          <v-icon v-text="'mdi-arrow-right'" class="ml-1" />
                        </v-btn>
                      </div>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div
                      v-show="
                        verificationRequest.token &&
                          pageLoadStatus === STATUS.BAD_REQUEST
                      "
                    >
                      <v-alert
                        type="error"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Email verification link invalid! <br />
                        Click
                        <router-link
                          :to="{ name: 'verify', query: {} }"
                          class="text-link"
                          v-text="'here'"
                        />
                        to resend email.
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div v-show="formLoadStatus === STATUS.ERROR">
                      <v-alert
                        type="error"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Failed sending verification email!
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div
                      v-show="
                        !verificationRequest.token &&
                          formLoadStatus !== STATUS.COMPLETE
                      "
                    >
                      <div>
                        Please enter your email below to send an email
                        verification link.
                      </div>
                      <div>
                        <v-text-field
                          v-model="email"
                          :error-messages="emailErrors"
                          label="Email"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          @input="
                            () => {
                              $v.email.$touch();
                              this.apiResponseCode = '';
                            }
                          "
                          @blur="$v.email.$touch()"
                          :prepend-icon="'mdi-email'"
                        ></v-text-field>
                      </div>
                      <v-btn
                        block
                        rounded
                        color="primaryDim"
                        :disabled="
                          formLoadStatus === STATUS.LOADING ||
                            formLoadStatus === STATUS.COMPLETE
                        "
                        @click="resend"
                      >
                        <div v-if="formLoadStatus === STATUS.LOADING">
                          verifying
                        </div>
                        <div v-else>
                          verify
                        </div>
                      </v-btn>
                    </div>
                  </v-expand-transition>
                </div>
              </v-card>
            </div>
          </v-scroll-y-transition>
        </div>
      </v-row>
    </v-col>
    <v-fade-transition>
      <v-overlay
        v-show="pageLoadStatus === STATUS.PRE_LOADING"
        opacity="0"
        absolute
      >
        <v-progress-circular indeterminate size="64"></v-progress-circular>
      </v-overlay>
    </v-fade-transition>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { email, required } from "vuelidate/lib/validators";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";

import "@/styles/Authorize.sass";

export default Vue.extend({
  data: function() {
    return {
      pageLoadStatus: STATUS.PRE_LOADING,
      formLoadStatus: STATUS.IDLE,
      email: "",
      apiResponseCode: ""
    };
  },

  computed: {
    verificationRequest() {
      return {
        token: this.$route.query.token?.toString()
      };
    },
    emailErrors() {
      const errors: string[] = [];
      if (!this.$v.email.$dirty) return errors;
      !this.$v.email.required && errors.push("Email required");
      !this.$v.email.email && errors.push("Invalid email");
      return errors;
    }
  },

  validations: {
    email: {
      required,
      email
    }
  },

  created: function() {
    if (this.verificationRequest.token) {
      api.user
        .verify(this.verificationRequest.token)
        .then(() => {
          this.pageLoadStatus = STATUS.COMPLETE;
        })
        .catch(() => {
          this.pageLoadStatus = STATUS.BAD_REQUEST;
        });
    } else {
      this.pageLoadStatus = STATUS.COMPLETE;
    }
  },

  methods: {
    resend(e: Event) {
      e.preventDefault();
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.formLoadStatus = STATUS.LOADING;
        api.user
          .resend(this.email)
          .then(() => {
            this.formLoadStatus = STATUS.COMPLETE;
          })
          .catch(error => {
            this.apiResponseCode = error.response.data.code;
            this.formLoadStatus = !this.apiResponseCode
              ? STATUS.ERROR
              : STATUS.IDLE;
          });
      }
    }
  }
});
</script>
