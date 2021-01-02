<template>
  <div class="application-detail">
    <v-row align="center">
      <v-col cols="12" sm="">
        <v-btn
          plain
          :ripple="false"
          class="pa-0"
          @click="() => $router.push({ name: 'manage:application' })"
        >
          <v-icon v-text="'mdi-arrow-left'" class="mr-1" />
          Back
        </v-btn>
      </v-col>
    </v-row>
    <v-row class="mb-8" align="center">
      <v-col cols="12">
        <h1 class="text-h2">{{ application.name }}</h1>
        <div class="text-subtitle-1 text--secondary">
          {{ application.description }}
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-card>
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Details
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="
                      application.editing &&
                        application.formLoadStatus !== STATUS.LOADING
                    "
                    text
                    rounded
                    color="error"
                    @click="cancelApplication"
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    :disabled="!applicationUpdated"
                    :color="
                      application.editing ? 'success' : 'secondary lighten-1'
                    "
                    @click="saveApplication"
                  >
                    <div v-if="!application.editing">Edit</div>
                    <div
                      v-else-if="
                        application.editing &&
                          application.formLoadStatus !== STATUS.LOADING
                      "
                    >
                      Save
                    </div>
                    <div
                      v-else-if="
                        application.editing &&
                          application.formLoadStatus === STATUS.LOADING
                      "
                    >
                      <v-progress-circular
                        indeterminate
                        color="success"
                        size="16"
                        class="mr-2"
                      />
                      <span>
                        Saving
                      </span>
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-expand-transition>
                <div v-show="application.formLoadStatus === STATUS.ERROR">
                  <v-alert
                    type="error"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Failed updating application!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-row>
                <v-col cols="12" sm="6">
                  <div>
                    <div class="mb-1 text-overline text--secondary">
                      Client ID
                    </div>
                    <div>
                      <code>{{ application.clientId }}</code>
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div>
                    <div class="mb-1 text-overline text--secondary">
                      Client Secret
                    </div>
                    <div style="height: 24px">
                      <code>
                        ••••••••••••••••
                      </code>
                      <v-btn
                        :ripple="false"
                        plain
                        color="error"
                        class="my-n4"
                        :disabled="application.editing"
                      >
                        Revoke
                      </v-btn>
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!application.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Name
                    </div>
                    <div>
                      {{ application.name }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="application.name"
                      :error-messages="nameErrors"
                      :counter="20"
                      label="Name"
                      required
                      :disabled="application.formLoadStatus === STATUS.LOADING"
                      @input="$v.application.name.$touch()"
                      @blur="$v.application.name.$touch()"
                      :prepend-icon="'mdi-application'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!application.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Description
                    </div>
                    <div>
                      {{ application.description }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="application.description"
                      :error-messages="descriptionErrors"
                      :counter="50"
                      label="Description"
                      required
                      :disabled="application.formLoadStatus === STATUS.LOADING"
                      @input="$v.application.description.$touch()"
                      @blur="$v.application.description.$touch()"
                      :prepend-icon="'mdi-text'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!application.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Login URL
                    </div>
                    <div>
                      {{ application.loginURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="application.loginURL"
                      :error-messages="loginURLErrors"
                      label="Login URL"
                      required
                      hint="Ratify may require to redirect users back to your application's login page"
                      :disabled="application.formLoadStatus === STATUS.LOADING"
                      @input="$v.application.loginURL.$touch()"
                      @blur="$v.application.loginURL.$touch()"
                      :prepend-icon="'mdi-login-variant'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!application.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Callback URL
                    </div>
                    <div>
                      {{ application.callbackURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="application.callbackURL"
                      :error-messages="callbackURLErrors"
                      label="Callback URL"
                      required
                      hint="Use semicolon to separate multiple allowed callback URLs"
                      :disabled="application.formLoadStatus === STATUS.LOADING"
                      @input="$v.application.callbackURL.$touch()"
                      @blur="$v.application.callbackURL.$touch()"
                      :prepend-icon="'mdi-undo-variant'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!application.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Login URL
                    </div>
                    <div>
                      {{ application.logoutURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="application.logoutURL"
                      :error-messages="logoutURLErrors"
                      label="Logout URL"
                      required
                      hint="Your application's logout URL to trigger global logout"
                      :disabled="application.formLoadStatus === STATUS.LOADING"
                      @input="$v.application.logoutURL.$touch()"
                      @blur="$v.application.logoutURL.$touch()"
                      :prepend-icon="'mdi-logout-variant'"
                    />
                  </div>
                </v-col>
              </v-row>
            </div>
          </v-card>
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
import { STATUS } from "@/constants/status";
import api from "@/apis/api";
import { maxLength, required, url } from "vuelidate/lib/validators";

export default Vue.extend({
  data() {
    return {
      application: {
        name: "",
        clientId: "",
        description: "",
        loginURL: "",
        callbackURL: "",
        logoutURL: "",
        editing: false,
        before: {
          name: "",
          description: "",
          loginURL: "",
          callbackURL: "",
          logoutURL: ""
        },
        formLoadStatus: STATUS.IDLE,
        apiResponseCode: ""
      },
      pageLoadStatus: STATUS.PRE_LOADING
    };
  },

  validations: {
    application: {
      name: { required, maxLength: maxLength(20) },
      description: { required, maxLength: maxLength(50) },
      loginURL: { required, url },
      callbackURL: { required, url },
      logoutURL: { required, url }
    }
  },

  computed: {
    applicationUpdated: {
      cache: false,
      get: function() {
        return (
          this.application.name !== this.application.before.name ||
          this.application.description !==
            this.application.before.description ||
          this.application.loginURL !== this.application.before.loginURL ||
          this.application.callbackURL !==
            this.application.before.callbackURL ||
          this.application.logoutURL !== this.application.before.logoutURL
        );
      }
    },
    nameErrors() {
      const errors: string[] = [];
      if (!this.$v.application.name?.$dirty) return errors;
      !this.$v.application.name.required && errors.push("Name required");
      !this.$v.application.name.maxLength && errors.push("Name too long");
      return errors;
    },
    descriptionErrors() {
      const errors: string[] = [];
      if (!this.$v.application.description?.$dirty) return errors;
      !this.$v.application.description.required &&
        errors.push("Description required");
      !this.$v.application.description.maxLength &&
        errors.push("Description too long");
      return errors;
    },
    loginURLErrors() {
      const errors: string[] = [];
      if (!this.$v.application.loginURL?.$dirty) return errors;
      !this.$v.application.loginURL.required &&
        errors.push("Login URL required");
      !this.$v.application.loginURL.url && errors.push("Invalid URL");
      return errors;
    },
    callbackURLErrors() {
      const errors: string[] = [];
      if (!this.$v.application.callbackURL?.$dirty) return errors;
      !this.$v.application.callbackURL.required &&
        errors.push("Callback URL required");
      !this.$v.application.callbackURL.url && errors.push("Invalid URL");
      return errors;
    },
    logoutURLErrors() {
      const errors: string[] = [];
      if (!this.$v.application.logoutURL?.$dirty) return errors;
      !this.$v.application.logoutURL.required &&
        errors.push("Logout URL required");
      !this.$v.application.logoutURL.url && errors.push("Invalid URL");
      return errors;
    }
  },

  created() {
    api.application.detail(this.$route.params.clientId, true).then(response => {
      this.application.name = response.data.data.name;
      this.application.clientId = response.data.data.client_id;
      this.application.description = response.data.data.description;
      this.application.loginURL = response.data.data.login_url;
      this.application.callbackURL = response.data.data.callback_url;
      this.application.logoutURL = response.data.data.logout_url;
      this.pageLoadStatus = STATUS.COMPLETE;
    });
  },

  methods: {
    cancelApplication() {
      this.application.editing = false;
      this.application.name = this.application.before.name;
      this.application.description = this.application.before.description;
      this.application.loginURL = this.application.before.loginURL;
      this.application.callbackURL = this.application.before.callbackURL;
      this.application.logoutURL = this.application.before.logoutURL;
      this.application.formLoadStatus = STATUS.IDLE;
      this.application.apiResponseCode = "";
      this.application.before = {
        name: "",
        description: "",
        loginURL: "",
        callbackURL: "",
        logoutURL: ""
      };
      this.$v.$reset();
    },
    saveApplication() {
      if (!this.application.editing) {
        this.application.editing = true;
        this.application.before.name = this.application.name;
        this.application.before.description = this.application.description;
        this.application.before.loginURL = this.application.loginURL;
        this.application.before.callbackURL = this.application.callbackURL;
        this.application.before.logoutURL = this.application.logoutURL;
        return;
      }
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.application.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.application
              .update(this.application.clientId, {
                /* eslint-disable @typescript-eslint/camelcase */
                name: this.application.name,
                description: this.application.description,
                login_url: this.application.loginURL,
                callback_url: this.application.callbackURL,
                logout_url: this.application.logoutURL
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(response => {
                console.log(response.data);
                this.application.editing = false;
                this.application.formLoadStatus = STATUS.SUCCESS;
              })
              .catch(error => {
                console.error(error.response.data);
                this.application.editing = true;
                this.application.apiResponseCode = error.response.data.code;
                this.application.formLoadStatus = !this.application
                  .apiResponseCode
                  ? STATUS.ERROR
                  : STATUS.IDLE;
              }),
          2000
        );
      }
    }
  }
});
</script>
