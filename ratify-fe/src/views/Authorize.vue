<template>
  <v-container fluid class="fill-height gradient-bg">
    <v-col>
      <v-row align="center" justify="center">
        <div v-if="pageLoadStatus === STATUS.BAD_REQUEST">
          <v-fade-transition>
            <div>
              <h1 class="text-h2 mb-8 text-center">400 Bad Request</h1>
              <p class="text-button text--secondary text-center">
                {{ redirectCountdown }}
              </p>
            </div>
          </v-fade-transition>
        </div>
        <div v-else-if="pageLoadStatus === STATUS.COMPLETE" class="login-form">
          <v-scroll-y-transition appear>
            <div>
              <v-card
                elevation="12"
                :loading="formLoadStatus === STATUS.LOADING"
              >
                <div class="pa-6">
                  <div class="text-h2 mt-12 mb-12 text-center">
                    <div
                      class="text-subtitle-1 text-center text--disabled font-weight-thin mb-2"
                    >
                      Sign in to
                    </div>
                    {{ application.name }}
                  </div>
                  <v-expand-transition>
                    <div v-show="formLoadStatus === STATUS.COMPLETE">
                      <v-alert
                        type="success"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Successfully signed in!
                      </v-alert>
                      <div
                        class="text-button text--disabled text-center mt-4 mb-1"
                      >
                        {{ redirectCountdown }}
                      </div>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div
                      v-show="requireOTP && formLoadStatus !== STATUS.COMPLETE"
                    >
                      <v-alert
                        type="info"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Please enter the code your authenticator app generated.
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div v-show="!$v.login.username.correct">
                      <v-alert type="error" text dense>
                        Incorrect
                        {{ this.requireOTP ? "code" : "credentials" }}!
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div
                      v-show="
                        apiResponseCode === 'email_unverified' &&
                          formLoadStatus !== STATUS.COMPLETE
                      "
                    >
                      <v-alert type="warning" text dense>
                        Email has not been verified! <br />
                        Click
                        <router-link
                          to="/verify"
                          class="text-link"
                          v-text="'here'"
                        />
                        to re-verify.
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
                        >Failed signing in!
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div v-if="formLoadStatus !== STATUS.COMPLETE">
                      <v-form @submit="loginUser">
                        <v-expand-transition>
                          <div v-if="!requireOTP">
                            <v-text-field
                              v-model="login.username"
                              :error-messages="usernameErrors"
                              label="Username or email"
                              required
                              :disabled="
                                formLoadStatus === STATUS.LOADING ||
                                  formLoadStatus === STATUS.COMPLETE
                              "
                              @input="
                                () => {
                                  $v.login.username.$touch();
                                  this.apiResponseCode = '';
                                }
                              "
                              @blur="$v.login.username.$touch()"
                              :prepend-icon="'mdi-identifier'"
                            />
                            <v-text-field
                              v-model="login.password"
                              :error-messages="passwordErrors"
                              :type="'password'"
                              label="Password"
                              required
                              :disabled="
                                formLoadStatus === STATUS.LOADING ||
                                  formLoadStatus === STATUS.COMPLETE
                              "
                              @input="
                                () => {
                                  $v.login.password.$touch();
                                  this.apiResponseCode = '';
                                }
                              "
                              @blur="$v.login.password.$touch()"
                              :prepend-icon="'mdi-lock'"
                            />
                          </div>
                        </v-expand-transition>
                        <v-expand-transition>
                          <div v-if="requireOTP">
                            <v-text-field
                              v-model="mfa.otp"
                              :error-messages="otpErrors"
                              label="Code"
                              autofocus
                              required
                              :disabled="
                                formLoadStatus === STATUS.LOADING ||
                                  formLoadStatus === STATUS.COMPLETE
                              "
                              @input="
                                () => {
                                  $v.mfa.otp.$touch();
                                  this.apiResponseCode = '';
                                }
                              "
                              @blur="$v.mfa.otp.$touch()"
                              :prepend-icon="'mdi-two-factor-authentication'"
                            />
                          </div>
                        </v-expand-transition>
                        <v-btn
                          type="submit"
                          block
                          rounded
                          color="primaryDim"
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          class="mt-2"
                        >
                          <div v-if="formLoadStatus === STATUS.LOADING">
                            signing in
                          </div>
                          <div v-else>
                            sign in
                          </div>
                        </v-btn>
                      </v-form>
                    </div>
                  </v-expand-transition>
                </div>
              </v-card>
              <p class="text-center mt-4 text-subtitle-1 text--secondary">
                Don't have an account?
                <router-link to="/signup" class="text-link">
                  Sign Up
                </router-link>
              </p>
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
import {
  and,
  maxLength,
  minLength,
  numeric,
  required
} from "vuelidate/lib/validators";
import api from "@/apis/api";
import oauth from "@/apis/oauth";
import { STATUS } from "@/constants/status";

import "@/styles/Authorize.sass";

export default Vue.extend({
  data: function() {
    return {
      pageLoadStatus: STATUS.PRE_LOADING,
      formLoadStatus: STATUS.IDLE,
      application: {
        name: "",
        metadata: ""
      },
      login: {
        username: "",
        password: ""
      },
      mfa: {
        otp: ""
      },
      requireOTP: false,
      redirectCountdown: "Redirecting now",
      redirectCounter: 0,
      apiResponseCode: ""
    };
  },

  computed: {
    authRequest() {
      return {
        clientId: this.$route.query.client_id?.toString(),
        responseType: this.$route.query.response_type?.toString(),
        redirectUri: this.$route.query.redirect_uri?.toString(),
        scope: this.$route.query.scope?.toString(),
        state: this.$route.query.state?.toString(),
        codeChallenge: this.$route.query.code_challenge?.toString(),
        codeChallengeMethod: this.$route.query.code_challenge_method?.toString(),
        immediate: this.$route.query.immediate === "true"
      };
    },
    usernameErrors() {
      const errors: string[] = [];
      if (!this.$v.login.username?.$dirty) return errors;
      !this.$v.login.username?.required &&
        errors.push("Username or email required");
      !this.$v.login.username?.correct && errors.push("");
      return errors;
    },
    passwordErrors() {
      const errors: string[] = [];
      if (!this.$v.login.password?.$dirty) return errors;
      !this.$v.login.password?.required && errors.push("Password required");
      !this.$v.login.password?.correct && errors.push("");
      return errors;
    },
    otpErrors() {
      const errors: string[] = [];
      if (!this.$v.mfa.otp?.$dirty) return errors;
      !this.$v.mfa.otp?.required && errors.push("Code required");
      !this.$v.mfa.otp?.length && errors.push("Invalid code");
      !this.$v.mfa.otp?.numeric && errors.push("Invalid code");
      !this.$v.mfa.otp?.correct && errors.push("");
      return errors;
    }
  },

  validations: {
    login: {
      username: {
        required,
        correct() {
          return this.$data.apiResponseCode !== "incorrect_credentials";
        }
      },
      password: {
        required,
        correct() {
          return this.$data.apiResponseCode !== "incorrect_credentials";
        }
      }
    },
    mfa: {
      otp: {
        required,
        numeric,
        length: and(minLength(6), maxLength(6)),
        correct() {
          return this.$data.apiResponseCode !== "incorrect_credentials";
        }
      }
    }
  },

  created: function() {
    if (
      !this.authRequest.clientId ||
      !this.authRequest.responseType ||
      !this.authRequest.redirectUri
    ) {
      this.pageLoadStatus = STATUS.BAD_REQUEST;
      this.startRedirectCounter(() => this.$router.replace({ name: "home" }));
    }
    api.application
      .detail(this.authRequest.clientId)
      .then(response => {
        this.application = response.data.data;
        // attempt authorization via session_id cookie
        this.authorizeUser(true)
          .then(() => {
            this.pageLoadStatus = STATUS.COMPLETE;
          })
          .catch(() => {
            this.pageLoadStatus = STATUS.COMPLETE;
          });
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.BAD_REQUEST;
        this.startRedirectCounter(() => this.$router.replace({ name: "home" }));
      });
  },

  methods: {
    loginUser(e: Event) {
      e.preventDefault();
      this.$v.login.$touch();
      this.requireOTP && this.$v.mfa.$touch();
      if (
        !(this.$v.login.$invalid || (this.requireOTP && this.$v.mfa.$invalid))
      ) {
        this.formLoadStatus = STATUS.LOADING;
        this.authorizeUser(false).catch(error => {
          this.apiResponseCode = error.response.data.code;
          this.formLoadStatus = !this.apiResponseCode
            ? STATUS.ERROR
            : STATUS.IDLE;
          this.requireOTP =
            this.requireOTP || this.apiResponseCode === "missing_otp";
        });
      }
    },
    authorizeUser(useSession: boolean) {
      return oauth
        .authorize({
          /* eslint-disable @typescript-eslint/camelcase */
          client_id: this.authRequest.clientId,
          response_type: this.authRequest.responseType,
          redirect_uri: this.authRequest.redirectUri,
          scope: this.authRequest.scope,
          state: this.authRequest.state,
          code_challenge: this.authRequest.codeChallenge,
          code_challenge_method: this.authRequest.codeChallengeMethod,
          preferred_username: this.login.username,
          password: this.login.password,
          otp: this.requireOTP ? this.mfa.otp : "",
          use_session: useSession
          /* eslint-enable @typescript-eslint/camelcase */
        })
        .then(response => {
          this.formLoadStatus = STATUS.COMPLETE;
          this.startRedirectCounter(() =>
            window.location.replace(response.data.data)
          );
        });
    },
    startRedirectCounter(callback: () => void) {
      let count = this.authRequest.immediate ? 0 : 3;
      this.redirectCountdown = count
        ? `Redirecting in ${count} second${count === 1 ? "" : "s"}`
        : `Redirecting now`;
      this.redirectCounter = setInterval(() => {
        if (count) {
          this.redirectCountdown = `Redirecting in ${count} second${
            count === 1 ? "" : "s"
          }`;
        } else {
          this.redirectCountdown = `Redirecting now`;
          clearInterval(this.redirectCounter);
          callback();
        }
        count--;
      }, 1000);
    }
  },

  beforeDestroy() {
    clearInterval(this.redirectCounter);
  }
});
</script>
