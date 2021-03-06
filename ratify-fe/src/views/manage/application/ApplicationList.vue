<template>
  <div class="application-list">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm=""><h1 class="text-h2">Applications</h1></v-col>
      <v-col cols="auto">
        <v-dialog
          v-model="create.creating"
          width="546"
          persistent
          overlay-opacity="0"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              color="primary darken-2"
              rounded
              v-bind="attrs"
              v-on="on"
              @click="resetCreate"
            >
              New Application
            </v-btn>
          </template>
          <v-card :loading="create.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  New Application
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="
                      create.formLoadStatus === STATUS.IDLE ||
                        create.formLoadStatus === STATUS.ERROR
                    "
                    text
                    rounded
                    color="error"
                    @click="
                      () => {
                        cancelCreate();
                      }
                    "
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    v-if="create.formLoadStatus !== STATUS.COMPLETE"
                    text
                    rounded
                    class="ml-4"
                    color="success"
                    :disabled="create.formLoadStatus === STATUS.LOADING"
                    @click="confirmCreate"
                  >
                    <div v-if="create.formLoadStatus !== STATUS.LOADING">
                      Create
                    </div>
                    <div v-else>
                      Creating
                    </div>
                  </v-btn>
                  <v-btn
                    v-if="create.formLoadStatus === STATUS.COMPLETE"
                    text
                    rounded
                    color="success"
                    @click="cancelCreate"
                  >
                    Confirm
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-expand-transition>
                <div v-show="create.formLoadStatus === STATUS.ERROR">
                  <v-alert
                    type="error"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Failed creating application!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-show="create.formLoadStatus === STATUS.COMPLETE">
                  <v-alert
                    type="success"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Application successfully created!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-if="create.formLoadStatus !== STATUS.COMPLETE">
                  <v-row>
                    <v-col cols="12">
                      <v-text-field
                        v-model.trim="create.name"
                        :error-messages="nameErrors"
                        :counter="20"
                        label="Name"
                        required
                        :disabled="create.formLoadStatus === STATUS.LOADING"
                        @input="$v.create.name.$touch()"
                        @blur="$v.create.name.$touch()"
                        :prepend-icon="'mdi-application'"
                      />
                      <v-text-field
                        v-model.trim="create.description"
                        :error-messages="descriptionErrors"
                        :counter="50"
                        label="Description"
                        required
                        :disabled="create.formLoadStatus === STATUS.LOADING"
                        @input="$v.create.description.$touch()"
                        @blur="$v.create.description.$touch()"
                        :prepend-icon="'mdi-text'"
                      />
                      <v-text-field
                        v-model.trim="create.loginURL"
                        :error-messages="loginURLErrors"
                        label="Login URL"
                        required
                        hint="Ratify may require to redirect users back to your application's login page"
                        :disabled="create.formLoadStatus === STATUS.LOADING"
                        @input="$v.create.loginURL.$touch()"
                        @blur="$v.create.loginURL.$touch()"
                        :prepend-icon="'mdi-login-variant'"
                      />
                      <v-text-field
                        v-model.trim="create.callbackURL"
                        :error-messages="callbackURLErrors"
                        label="Callback URL"
                        required
                        hint="Use semicolon to separate multiple allowed callback URLs"
                        :disabled="create.formLoadStatus === STATUS.LOADING"
                        @input="$v.create.callbackURL.$touch()"
                        @blur="$v.create.callbackURL.$touch()"
                        :prepend-icon="'mdi-undo-variant'"
                      />
                      <v-text-field
                        v-model.trim="create.logoutURL"
                        :error-messages="logoutURLErrors"
                        label="Logout URL"
                        required
                        hint="Your application's logout URL to trigger global logout"
                        :disabled="create.formLoadStatus === STATUS.LOADING"
                        @input="$v.create.logoutURL.$touch()"
                        @blur="$v.create.logoutURL.$touch()"
                        :prepend-icon="'mdi-logout-variant'"
                      />
                    </v-col>
                  </v-row>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-if="create.formLoadStatus === STATUS.COMPLETE">
                  <v-row>
                    <v-col cols="12">
                      <div>
                        <b>{{ this.create.name }}</b> has been created. Safely
                        store the following <b>client_secret</b>, it
                        <b>cannot</b> be seen again once this prompt is closed.
                        Exposing this secret will leave your application
                        vulnerable.
                      </div>
                      <div class="mt-2">
                        <div class="mb-1 text-overline text--secondary">
                          Client Secret
                        </div>
                        <div>
                          <code>{{ this.create.clientSecret }}</code>
                        </div>
                      </div>
                    </v-col>
                  </v-row>
                </div>
              </v-expand-transition>
            </div>
          </v-card>
        </v-dialog>
      </v-col>
    </v-row>
    <v-fade-transition>
      <v-row v-show="pageLoadStatus === STATUS.COMPLETE">
        <v-col cols="12">
          <v-divider inset />
          <div v-for="application in applications" :key="application.client_id">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-content>
                  <v-row justify="end" align="center" no-gutters>
                    <v-col cols="12" md="">
                      <v-list-item-title class="text-h5">
                        <span
                          class="d-inline-block text-truncate"
                          style="max-width: 320px;"
                        >
                          {{ application.name }}
                        </span>
                      </v-list-item-title>
                      <v-list-item-subtitle>
                        <span
                          class="d-inline-block text-truncate"
                          style="max-width: 320px;"
                        >
                          {{ application.description }}
                        </span>
                      </v-list-item-subtitle>
                    </v-col>
                    <v-col cols="12" md="">
                      <span class="text-overline text-no-wrap mr-1">
                        Client ID
                      </span>
                      <code>{{ application.client_id }}</code>
                    </v-col>
                    <v-col cols="auto" md="auto">
                      <v-btn
                        rounded
                        text
                        outlined
                        color="secondary lighten-1"
                        :to="{
                          name: 'manage:application-detail',
                          params: { clientId: application.client_id }
                        }"
                      >
                        Manage
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-list-item-content>
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
          Failed retrieving application list!
        </v-alert>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";
import { validateURL } from "@/utils/url";
import { maxLength, required } from "vuelidate/lib/validators";

export default Vue.extend({
  data: () => ({
    applications: [],
    create: {
      creating: false,
      name: "",
      description: "",
      loginURL: "",
      callbackURL: "",
      logoutURL: "",
      clientSecret: "",
      formLoadStatus: STATUS.IDLE,
      apiResponseCode: ""
    },
    pageLoadStatus: STATUS.PRE_LOADING
  }),

  computed: {
    nameErrors() {
      const errors: string[] = [];
      if (!this.$v.create.name?.$dirty) return errors;
      !this.$v.create.name.required && errors.push("Name required");
      !this.$v.create.name.maxLength && errors.push("Name too long");
      return errors;
    },
    descriptionErrors() {
      const errors: string[] = [];
      if (!this.$v.create.description?.$dirty) return errors;
      !this.$v.create.description.required &&
        errors.push("Description required");
      !this.$v.create.description.maxLength &&
        errors.push("Description too long");
      return errors;
    },
    loginURLErrors() {
      const errors: string[] = [];
      if (!this.$v.create.loginURL?.$dirty) return errors;
      !this.$v.create.loginURL.required && errors.push("Login URL required");
      !this.$v.create.loginURL.url && errors.push("Invalid URL");
      return errors;
    },
    callbackURLErrors() {
      const errors: string[] = [];
      if (!this.$v.create.callbackURL?.$dirty) return errors;
      !this.$v.create.callbackURL.required &&
        errors.push("Callback URL required");
      !this.$v.create.callbackURL.url && errors.push("Invalid URL");
      return errors;
    },
    logoutURLErrors() {
      const errors: string[] = [];
      if (!this.$v.create.logoutURL?.$dirty) return errors;
      !this.$v.create.logoutURL.required && errors.push("Logout URL required");
      !this.$v.create.logoutURL.url && errors.push("Invalid URL");
      return errors;
    }
  },

  validations: {
    create: {
      name: { required, maxLength: maxLength(20) },
      description: { required, maxLength: maxLength(50) },
      loginURL: { required, url: validateURL(true) },
      callbackURL: { required, url: validateURL(true) },
      logoutURL: { required, url: validateURL(true) }
    }
  },

  created() {
    this.loadApplications();
  },

  methods: {
    loadApplications() {
      this.pageLoadStatus = STATUS.PRE_LOADING;
      api.application
        .list()
        .then(response => {
          this.applications = response.data.data;
          this.applications.sort((a: { name: string }, b: { name: string }) =>
            a["name"].localeCompare(b["name"])
          );
          this.pageLoadStatus = STATUS.COMPLETE;
        })
        .catch(() => {
          this.pageLoadStatus = STATUS.ERROR;
        });
    },
    resetCreate() {
      this.create.name = "";
      this.create.description = "";
      this.create.loginURL = "";
      this.create.callbackURL = "";
      this.create.logoutURL = "";
      this.create.clientSecret = "";
      this.create.formLoadStatus = STATUS.IDLE;
      this.$v.create.$reset();
    },
    confirmCreate() {
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.create.formLoadStatus = STATUS.LOADING;
        api.application
          .register({
            /* eslint-disable @typescript-eslint/camelcase */
            name: this.create.name.trim(),
            description: this.create.description.trim(),
            login_url: this.create.logoutURL.trim(),
            callback_url: this.create.callbackURL.trim(),
            logout_url: this.create.logoutURL.trim()
            /* eslint-enable @typescript-eslint/camelcase */
          })
          .then(response => {
            this.create.clientSecret = response.data.data.client_secret;
            this.create.formLoadStatus = STATUS.COMPLETE;
            this.loadApplications();
          })
          .catch(error => {
            this.create.apiResponseCode = error.response.data.code;
            this.create.formLoadStatus = !this.create.apiResponseCode
              ? STATUS.ERROR
              : STATUS.IDLE;
          });
      }
    },
    cancelCreate() {
      this.create.creating = false;
      this.create.clientSecret = "";
    }
  }
});
</script>
