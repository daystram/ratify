<template>
  <div class="application-detail">
    <v-row align="center">
      <v-col cols="12" sm="">
        <v-btn
          plain
          :ripple="false"
          class="pa-0"
          :to="{ name: 'manage:application' }"
        >
          <v-icon v-text="'mdi-arrow-left'" class="mr-1" />
          Back
        </v-btn>
      </v-col>
    </v-row>
    <v-row class="mb-8" align="center">
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <h1
            :class="
              'text-h2 text-truncate ' + (detail.name ? '' : 'text--disabled')
            "
          >
            {{ detail.name || "Application Name" }}
          </h1>
          <div
            :class="
              'text-subtitle-1 text-truncate ' +
                (detail.description ? 'text--secondary' : 'text--disabled')
            "
          >
            {{ detail.description || "Application description" }}
          </div>
        </v-col>
      </v-fade-transition>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-card :loading="detail.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Details
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="
                      detail.editing && detail.formLoadStatus !== STATUS.LOADING
                    "
                    text
                    rounded
                    color="error"
                    @click="cancelDetail"
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    :disabled="
                      !detailUpdated || detail.formLoadStatus === STATUS.LOADING
                    "
                    :color="detail.editing ? 'success' : 'secondary lighten-1'"
                    @click="saveDetail"
                  >
                    <div v-if="!detail.editing">Edit</div>
                    <div
                      v-else-if="
                        detail.editing &&
                          detail.formLoadStatus !== STATUS.LOADING
                      "
                    >
                      Save
                    </div>
                    <div
                      v-else-if="
                        detail.editing &&
                          detail.formLoadStatus === STATUS.LOADING
                      "
                    >
                      Saving
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-expand-transition>
                <div v-show="detail.successAlert">
                  <v-alert
                    type="info"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Application updated!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-show="detail.formLoadStatus === STATUS.ERROR">
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
                      <code>{{ detail.clientId }}</code>
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div>
                    <div class="mb-1 text-overline text--secondary">
                      Client Secret
                    </div>
                    <div style="height: 24px">
                      <code>••••••••••••••••</code>
                      <v-dialog
                        v-model="revoke.prompt"
                        width="545"
                        :persistent="revoke.formLoadStatus === STATUS.COMPLETE"
                        @input="v => v || cancelRevoke()"
                      >
                        <template v-slot:activator="{ on, attrs }">
                          <v-btn
                            :ripple="false"
                            plain
                            color="error"
                            class="my-n4"
                            :disabled="detail.editing"
                            v-bind="attrs"
                            v-on="on"
                            @click="() => (revoke.formLoadStatus = STATUS.IDLE)"
                          >
                            Revoke
                          </v-btn>
                        </template>
                        <v-card class="danger-border">
                          <v-card-title>
                            <v-row no-gutters align="center">
                              <v-col cols="auto" class="error--text">
                                Revoke Client Secret
                              </v-col>
                              <v-spacer />
                              <v-col cols="auto">
                                <v-btn
                                  v-if="
                                    revoke.formLoadStatus !== STATUS.COMPLETE
                                  "
                                  text
                                  icon
                                  color="grey"
                                  @click="cancelRevoke"
                                >
                                  <v-icon v-text="'mdi-close'" />
                                </v-btn>
                                <v-btn
                                  v-if="
                                    revoke.formLoadStatus === STATUS.COMPLETE
                                  "
                                  text
                                  rounded
                                  color="success"
                                  @click="
                                    () => {
                                      cancelRevoke();
                                      revoke.clientSecret = '';
                                    }
                                  "
                                >
                                  Confirm
                                </v-btn>
                              </v-col>
                            </v-row>
                          </v-card-title>
                          <v-divider inset />
                          <div class="v-card__body">
                            <v-expand-transition>
                              <div
                                v-if="revoke.formLoadStatus === STATUS.ERROR"
                              >
                                <v-alert type="error" text dense>
                                  Failed revoking application client secret!
                                </v-alert>
                              </div>
                            </v-expand-transition>
                            <v-expand-transition>
                              <div
                                v-if="revoke.formLoadStatus === STATUS.COMPLETE"
                              >
                                <v-alert type="info" text dense>
                                  New client secret issued!
                                </v-alert>
                              </div>
                            </v-expand-transition>
                            <v-expand-transition>
                              <div
                                v-if="revoke.formLoadStatus !== STATUS.COMPLETE"
                              >
                                <v-row>
                                  <v-col>
                                    <div>
                                      Are you sure you want to revoke the client
                                      secret for
                                      <b>{{ detail.name }}</b
                                      >? This is action will render the
                                      previously issued client secret unusable.
                                    </div>
                                    <div class="mt-4">
                                      Type <b>{{ detail.name }}</b> to confirm.
                                    </div>
                                    <v-text-field
                                      v-model="revoke.confirmName"
                                      class="py-2"
                                      :prepend-icon="'mdi-application'"
                                    />
                                    <v-btn
                                      rounded
                                      block
                                      outlined
                                      color="error"
                                      :disabled="
                                        revoke.confirmName !== detail.name
                                      "
                                      @click="confirmRevoke"
                                    >
                                      Revoke
                                    </v-btn>
                                  </v-col>
                                </v-row>
                              </div>
                            </v-expand-transition>
                            <v-expand-transition>
                              <div
                                v-if="revoke.formLoadStatus === STATUS.COMPLETE"
                              >
                                <v-row>
                                  <v-col cols="12">
                                    <div>
                                      Safely store the following
                                      <b>client_secret</b>, it <b>cannot</b> be
                                      seen again once this prompt is closed.
                                      Exposing this secret will leave your
                                      application vulnerable.
                                    </div>
                                    <div class="mt-2">
                                      <div
                                        class="mb-1 text-overline text--secondary"
                                      >
                                        Client Secret
                                      </div>
                                      <div>
                                        <code>
                                          {{ revoke.clientSecret }}
                                        </code>
                                      </div>
                                    </div>
                                  </v-col>
                                </v-row>
                              </div>
                            </v-expand-transition>
                          </div>
                        </v-card>
                      </v-dialog>
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!detail.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Name
                    </div>
                    <div>
                      {{ detail.name }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="detail.name"
                      :error-messages="nameErrors"
                      :counter="20"
                      label="Name"
                      required
                      :disabled="detail.formLoadStatus === STATUS.LOADING"
                      @input="$v.detail.name.$touch()"
                      @blur="$v.detail.name.$touch()"
                      :prepend-icon="'mdi-application'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!detail.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Description
                    </div>
                    <div>
                      {{ detail.description }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="detail.description"
                      :error-messages="descriptionErrors"
                      :counter="50"
                      label="Description"
                      required
                      :disabled="detail.formLoadStatus === STATUS.LOADING"
                      @input="$v.detail.description.$touch()"
                      @blur="$v.detail.description.$touch()"
                      :prepend-icon="'mdi-text'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!detail.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Login URL
                    </div>
                    <div>
                      {{ detail.loginURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="detail.loginURL"
                      :error-messages="loginURLErrors"
                      label="Login URL"
                      required
                      hint="Ratify may require to redirect users back to your application's login page"
                      :disabled="
                        detail.formLoadStatus === STATUS.LOADING ||
                          detail.locked
                      "
                      @input="$v.detail.loginURL.$touch()"
                      @blur="$v.detail.loginURL.$touch()"
                      :prepend-icon="'mdi-login-variant'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!detail.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Callback URL
                    </div>
                    <div>
                      {{ detail.callbackURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="detail.callbackURL"
                      :error-messages="callbackURLErrors"
                      label="Callback URL"
                      required
                      hint="Use semicolon to separate multiple allowed callback URLs"
                      :disabled="
                        detail.formLoadStatus === STATUS.LOADING ||
                          detail.locked
                      "
                      @input="$v.detail.callbackURL.$touch()"
                      @blur="$v.detail.callbackURL.$touch()"
                      :prepend-icon="'mdi-undo-variant'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!detail.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Logout URL
                    </div>
                    <div>
                      {{ detail.logoutURL }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="detail.logoutURL"
                      :error-messages="logoutURLErrors"
                      label="Logout URL"
                      required
                      hint="Your application's logout URL to trigger global logout"
                      :disabled="
                        detail.formLoadStatus === STATUS.LOADING ||
                          detail.locked
                      "
                      @input="$v.detail.logoutURL.$touch()"
                      @blur="$v.detail.logoutURL.$touch()"
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
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
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
      </v-fade-transition>
    </v-row>
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
          class="mt-3"
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
import { maxLength, required, url } from "vuelidate/lib/validators";

export default Vue.extend({
  data: () => ({
    detail: {
      name: "",
      clientId: "",
      description: "",
      loginURL: "",
      callbackURL: "",
      logoutURL: "",
      locked: false,
      editing: false,
      before: {
        name: "",
        description: "",
        loginURL: "",
        callbackURL: "",
        logoutURL: ""
      },
      formLoadStatus: STATUS.IDLE,
      apiResponseCode: "",
      successAlert: false
    },
    revoke: {
      prompt: false,
      confirmName: "",
      clientSecret: "",
      formLoadStatus: STATUS.IDLE
    },
    deleting: {
      prompt: false,
      confirmName: "",
      formLoadStatus: STATUS.IDLE
    },
    pageLoadStatus: STATUS.PRE_LOADING
  }),

  validations: {
    detail: {
      name: { required, maxLength: maxLength(20) },
      description: { required, maxLength: maxLength(50) },
      loginURL: { required, url },
      callbackURL: { required, url },
      logoutURL: { required, url }
    }
  },

  computed: {
    detailUpdated: {
      cache: false,
      get: function() {
        return (
          this.detail.name !== this.detail.before.name ||
          this.detail.description !== this.detail.before.description ||
          this.detail.loginURL !== this.detail.before.loginURL ||
          this.detail.callbackURL !== this.detail.before.callbackURL ||
          this.detail.logoutURL !== this.detail.before.logoutURL
        );
      }
    },
    nameErrors() {
      const errors: string[] = [];
      if (!this.$v.detail.name?.$dirty) return errors;
      !this.$v.detail.name.required && errors.push("Name required");
      !this.$v.detail.name.maxLength && errors.push("Name too long");
      return errors;
    },
    descriptionErrors() {
      const errors: string[] = [];
      if (!this.$v.detail.description?.$dirty) return errors;
      !this.$v.detail.description.required &&
        errors.push("Description required");
      !this.$v.detail.description.maxLength &&
        errors.push("Description too long");
      return errors;
    },
    loginURLErrors() {
      const errors: string[] = [];
      if (!this.$v.detail.loginURL?.$dirty) return errors;
      !this.$v.detail.loginURL.required && errors.push("Login URL required");
      !this.$v.detail.loginURL.url && errors.push("Invalid URL");
      return errors;
    },
    callbackURLErrors() {
      const errors: string[] = [];
      if (!this.$v.detail.callbackURL?.$dirty) return errors;
      !this.$v.detail.callbackURL.required &&
        errors.push("Callback URL required");
      !this.$v.detail.callbackURL.url && errors.push("Invalid URL");
      return errors;
    },
    logoutURLErrors() {
      const errors: string[] = [];
      if (!this.$v.detail.logoutURL?.$dirty) return errors;
      !this.$v.detail.logoutURL.required && errors.push("Logout URL required");
      !this.$v.detail.logoutURL.url && errors.push("Invalid URL");
      return errors;
    }
  },

  created() {
    api.application
      .detail(this.$route.params.clientId, true)
      .then(response => {
        this.detail.name = response.data.data.name;
        this.detail.clientId = response.data.data.client_id;
        this.detail.description = response.data.data.description;
        this.detail.loginURL = response.data.data.login_url;
        this.detail.callbackURL = response.data.data.callback_url;
        this.detail.logoutURL = response.data.data.logout_url;
        this.detail.locked = response.data.data.locked;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(error => {
        if (error.response.status === 404) {
          this.$router.push({ name: "manage:application" });
          return;
        }
        this.pageLoadStatus = STATUS.ERROR;
      });
  },

  methods: {
    cancelDetail() {
      this.detail.editing = false;
      this.detail.name = this.detail.before.name;
      this.detail.description = this.detail.before.description;
      this.detail.loginURL = this.detail.before.loginURL;
      this.detail.callbackURL = this.detail.before.callbackURL;
      this.detail.logoutURL = this.detail.before.logoutURL;
      this.detail.formLoadStatus = STATUS.IDLE;
      this.detail.apiResponseCode = "";
      this.detail.before = {
        name: "",
        description: "",
        loginURL: "",
        callbackURL: "",
        logoutURL: ""
      };
      this.$v.detail.$reset();
    },
    saveDetail() {
      if (!this.detail.editing) {
        this.detail.editing = true;
        this.detail.before.name = this.detail.name;
        this.detail.before.description = this.detail.description;
        this.detail.before.loginURL = this.detail.loginURL;
        this.detail.before.callbackURL = this.detail.callbackURL;
        this.detail.before.logoutURL = this.detail.logoutURL;
        this.detail.successAlert = false;
        return;
      }
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.detail.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.application
              .update(this.detail.clientId, {
                /* eslint-disable @typescript-eslint/camelcase */
                name: this.detail.name.trim(),
                description: this.detail.description.trim(),
                login_url: this.detail.loginURL.trim(),
                callback_url: this.detail.callbackURL.trim(),
                logout_url: this.detail.logoutURL.trim()
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(() => {
                this.detail.editing = false;
                this.detail.formLoadStatus = STATUS.COMPLETE;
                this.detail.successAlert = true;
                setTimeout(() => {
                  this.detail.successAlert = false;
                }, 5000);
              })
              .catch(error => {
                this.detail.editing = true;
                this.detail.apiResponseCode = error.response.data.code;
                this.detail.formLoadStatus = !this.detail.apiResponseCode
                  ? STATUS.ERROR
                  : STATUS.IDLE;
              }),
          2000
        );
      }
    },
    confirmDelete() {
      this.deleting.formLoadStatus = STATUS.LOADING;
      api.application
        .delete(this.detail.clientId)
        .then(() => {
          this.$router.push({ name: "manage:application" });
        })
        .catch(() => {
          this.deleting.formLoadStatus = STATUS.ERROR;
        });
    },
    cancelDelete() {
      this.deleting.prompt = false;
      this.deleting.confirmName = "";
      this.deleting.formLoadStatus = STATUS.IDLE;
    },
    confirmRevoke() {
      this.revoke.formLoadStatus = STATUS.LOADING;
      api.application
        .revoke(this.detail.clientId)
        .then(response => {
          this.revoke.clientSecret = response.data.data.client_secret;
          this.revoke.formLoadStatus = STATUS.COMPLETE;
        })
        .catch(() => {
          this.revoke.formLoadStatus = STATUS.ERROR;
        });
    },
    cancelRevoke() {
      this.revoke.prompt = false;
      this.revoke.confirmName = "";
      this.revoke.formLoadStatus = STATUS.IDLE;
    }
  }
});
</script>
