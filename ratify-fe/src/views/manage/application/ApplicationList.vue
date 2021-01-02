<template>
  <div class="application-list">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm=""><h1 class="text-h2">Applications</h1></v-col>
      <v-col cols="auto">
        <v-dialog
          v-model="create.creating"
          width="500"
          persistent
          overlay-opacity="0"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn color="primary darken-2" rounded v-bind="attrs" v-on="on">
              New Application
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  New Application
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="create.formLoadStatus !== STATUS.LOADING"
                    text
                    rounded
                    color="error"
                    @click="resetCreate"
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    color="success"
                    @click="createApplication"
                  >
                    <div v-if="create.formLoadStatus !== STATUS.LOADING">
                      Create
                    </div>
                    <div v-else>
                      <v-progress-circular
                        indeterminate
                        color="success"
                        size="16"
                        class="mr-2"
                      />
                      <span>
                        Creating
                      </span>
                    </div>
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
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    v-model="create.name"
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
                    v-model="create.description"
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
                    v-model="create.loginURL"
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
                    v-model="create.callbackURL"
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
                    v-model="create.logoutURL"
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
          </v-card>
        </v-dialog>
      </v-col>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-divider :inset="false" />
          <div v-for="application in applications" :key="application.client_id">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-content>
                  <v-row justify="end" align="center" no-gutters>
                    <v-col cols="12" md="">
                      <v-list-item-title class="text-h5">
                        {{ application.name }}
                      </v-list-item-title>
                      <v-list-item-subtitle>
                        {{ application.description }}
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
                        color="secondary lighten-1"
                        @click="
                          $router.push({
                            name: 'manage:application-detail',
                            params: { clientId: application.client_id }
                          })
                        "
                      >
                        manage
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-list-item-content>
              </v-list-item-content>
            </v-list-item>
            <v-divider :inset="false" />
          </div>
        </v-col>
      </v-fade-transition>
    </v-row>
    <v-fade-transition>
      <v-overlay
        v-show="pageLoadStatus !== STATUS.COMPLETE"
        opacity="0"
        absolute
      >
        <v-progress-circular indeterminate size="64" />
      </v-overlay>
    </v-fade-transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";
import { maxLength, required, url } from "vuelidate/lib/validators";

export default Vue.extend({
  data() {
    return {
      applications: [],
      create: {
        creating: false,
        name: "",
        description: "",
        loginURL: "",
        callbackURL: "",
        logoutURL: "",
        formLoadStatus: STATUS.IDLE,
        apiResponseCode: ""
      },
      pageLoadStatus: STATUS.PRE_LOADING
    };
  },

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
      loginURL: { required, url },
      callbackURL: { required, url },
      logoutURL: { required, url }
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
          this.applications.sort((a, b) => b["created_at"] - a["created_at"]);
          this.pageLoadStatus = STATUS.COMPLETE;
        })
        .catch(error => {
          console.error(error);
          this.pageLoadStatus = STATUS.ERROR;
        });
    },
    resetCreate() {
      this.create.creating = false;
      this.create.name = "";
      this.create.description = "";
      this.create.loginURL = "";
      this.create.callbackURL = "";
      this.create.logoutURL = "";
      this.create.formLoadStatus = STATUS.IDLE;
      this.$v.$reset();
    },
    createApplication() {
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.create.formLoadStatus = STATUS.LOADING;
        api.application
          .register({
            /* eslint-disable @typescript-eslint/camelcase */
            name: this.create.name,
            description: this.create.description,
            login_url: this.create.logoutURL,
            callback_url: this.create.callbackURL,
            logout_url: this.create.logoutURL
            /* eslint-enable @typescript-eslint/camelcase */
          })
          .then(() => {
            this.resetCreate();
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
    }
  }
});
</script>
