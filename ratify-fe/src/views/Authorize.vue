<template>
  <v-container fluid class="fill-height gradient-bg">
    <v-col>
      <v-row align="center" justify="center">
        <div v-if="pageLoadStatus === STATUS.BAD_REQUEST">
          <v-fade-transition>
            <div>
              <h1 class="text-h2 mb-8 text-center">400 Bad Request</h1>
              <p class="text-subtitle-1 text--secondary text-center">
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
                        class="text-subtitle-1 text--disabled text-center mt-4 mb-1"
                      >
                        {{ redirectCountdown }}
                      </div>
                    </div>
                  </v-expand-transition>
                  <v-expand-transition>
                    <div v-show="!$v.username.correct">
                      <v-alert type="error" text dense
                        >Incorrect username or password!
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
                    <div v-show="formLoadStatus !== STATUS.COMPLETE">
                      <v-form @submit="login">
                        <v-text-field
                          v-model="username"
                          :error-messages="usernameErrors"
                          label="Username"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          @input="
                            () => {
                              $v.username.$touch();
                              this.apiResponseCode = '';
                            }
                          "
                          @blur="$v.username.$touch()"
                          :prepend-icon="'mdi-identifier'"
                        ></v-text-field>
                        <v-text-field
                          v-model="password"
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
                              $v.password.$touch();
                              this.apiResponseCode = '';
                            }
                          "
                          @blur="$v.password.$touch()"
                          :prepend-icon="'mdi-lock'"
                        ></v-text-field>
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
import { alphaNum, required } from "vuelidate/lib/validators";
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
      username: "",
      password: "",
      redirectCountdown: "",
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
      if (!this.$v.username.$dirty) return errors;
      !this.$v.username.alphaNum && errors.push("Invalid username");
      !this.$v.username.required && errors.push("Username required");
      !this.$v.username.correct && errors.push("");
      return errors;
    },
    passwordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.$dirty) return errors;
      !this.$v.password.required && errors.push("Password required");
      !this.$v.password.correct && errors.push("");
      return errors;
    }
  },

  validations: {
    username: {
      required,
      alphaNum,
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
        this.authorizeUser()
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
    login(e: Event) {
      e.preventDefault();
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.formLoadStatus = STATUS.LOADING;
        this.authorizeUser().catch(error => {
          this.apiResponseCode = error.response.data.code;
          this.formLoadStatus = !this.apiResponseCode
            ? STATUS.ERROR
            : STATUS.IDLE;
        });
      }
    },
    authorizeUser() {
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
          preferred_username: this.username,
          password: this.password
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
