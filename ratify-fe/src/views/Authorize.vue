<template>
  <v-container fill-height="fill-height">
    <v-layout align-center="align-center" justify-center="justify-center">
      <v-flex class="login-form text-center">
        <div v-if="pageLoadStatus === this.STATUS.PRE_LOADING">
          <v-progress-circular
            indeterminate
            color="primary"
          ></v-progress-circular>
        </div>
        <div v-else-if="pageLoadStatus === this.STATUS.BAD_REQUEST">
          <h1>400 Bad Request</h1>
        </div>
        <div v-else>
          <h1>{{ application.name }}</h1>
          <v-form @submit="submit">
            <v-text-field
              v-model="username"
              :error-messages="usernameErrors"
              label="Username"
              required
              @input="$v.username.$touch()"
              @blur="$v.username.$touch()"
            ></v-text-field>
            <v-text-field
              v-model="password"
              :error-messages="passwordErrors"
              :type="'password'"
              label="Password"
              required
              @input="$v.password.$touch()"
              @blur="$v.password.$touch()"
            ></v-text-field>
            <v-btn
              type="submit"
              block="block"
              :disabled="formLoadStatus === STATUS.LOADING"
            >
              <div v-if="formLoadStatus === STATUS.LOADING">
                <v-progress-circular
                  indeterminate
                  color="primary"
                  :width="3"
                  :size="20"
                ></v-progress-circular>
              </div>
              <div v-else>
                sign in
              </div>
            </v-btn>
          </v-form>
        </div>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { required, alphaNum } from "vuelidate/lib/validators";
import api from "@/apis/api";
import oauth from "@/apis/oauth";
import { STATUS } from "@/constants/status";

import "@/styles/Login.sass";

export default Vue.extend({
  created: function() {
    this.authRequest = {
      clientId: this.$route.query.client_id,
      responseType: this.$route.query.response_type,
      redirectUri: this.$route.query.redirect_uri,
      state: this.$route.query.state,
      codeChallenge: this.$route.query.code_challenge,
      codeChallengeMethod: this.$route.query.code_challenge_method
    };
    if (!this.authRequest.clientId || !this.authRequest.responseType) {
      this.pageLoadStatus = this.PAGE_STATUS.BAD_REQUEST;
      return;
    }
    api.application
      .getOne(this.authRequest.clientId)
      .then(response => {
        this.application = response.data.data;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.BAD_REQUEST;
      });
  },

  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    formLoadStatus: STATUS.IDLE,
    authRequest: {
      clientId: "",
      responseType: "",
      redirectUri: "",
      state: "",
      codeChallenge: "",
      codeChallengeMethod: ""
    },
    application: {
      name: "",
      metadata: ""
    },
    username: "",
    password: ""
  }),

  validations: {
    username: { required, alphaNum },
    password: { required }
  },

  computed: {
    usernameErrors() {
      const errors: string[] = [];
      if (!this.$v.username.$dirty) return errors;
      !this.$v.username.alphaNum && errors.push("Invalid username");
      !this.$v.username.required && errors.push("Username required");
      return errors;
    },
    passwordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.$dirty) return errors;
      !this.$v.password.required && errors.push("Password required");
      return errors;
    }
  },

  methods: {
    submit(e) {
      e.preventDefault();
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.formLoadStatus = STATUS.LOADING;
        oauth
          .authorize({
            /* eslint-disable @typescript-eslint/camelcase */
            client_id: this.authRequest.clientId,
            response_type: this.authRequest.responseType,
            redirect_uri: this.authRequest.redirectUri,
            state: this.authRequest.state,
            code_challenge: this.authRequest.codeChallenge,
            code_challenge_method: this.authRequest.codeChallengeMethod,
            username: this.username,
            password: this.password
            /* eslint-enable @typescript-eslint/camelcase */
          })
          .then(response => {
            console.log(response.data.data);
            window.location.href = response.data.data;
            this.formLoadStatus = STATUS.IDLE;
          })
          .catch(error => {
            console.error(error.response.data.error);
            this.formLoadStatus = STATUS.IDLE;
          });
      }
    }
  }
});
</script>
