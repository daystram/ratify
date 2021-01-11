<template>
  <div class="profile">
    <v-row class="mb-8" align="center">
      <v-col cols="12" sm="">
        <h1 class="text-h2">Your Profile</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-fade-transition>
        <v-col v-show="pageLoadStatus === STATUS.COMPLETE" cols="12">
          <v-card :loading="profile.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Profile
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="
                      profile.editing &&
                        profile.formLoadStatus !== STATUS.LOADING
                    "
                    text
                    rounded
                    color="error"
                    @click="cancelProfile"
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    :disabled="
                      !profileUpdated ||
                        profile.formLoadStatus === STATUS.LOADING
                    "
                    :color="profile.editing ? 'success' : 'secondary lighten-1'"
                    @click="saveProfile"
                  >
                    <div v-if="!profile.editing">Edit</div>
                    <div
                      v-else-if="
                        profile.editing &&
                          profile.formLoadStatus !== STATUS.LOADING
                      "
                    >
                      Save
                    </div>
                    <div
                      v-else-if="
                        profile.editing &&
                          profile.formLoadStatus === STATUS.LOADING
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
                <div v-show="profile.successAlert">
                  <v-alert
                    type="info"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Profile updated!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-show="profile.formLoadStatus === STATUS.ERROR">
                  <v-alert
                    type="error"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Failed updating user!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-row align="center" justify="center">
                <v-col cols="auto">
                  <v-avatar
                    color="primaryDim"
                    size="128"
                    style="user-select: none; font-size: x-large"
                  >
                    {{
                      (profile.givenName &&
                        profile.givenName[0].toUpperCase()) +
                        (profile.familyName &&
                          profile.familyName[0].toUpperCase())
                    }}
                  </v-avatar>
                </v-col>
                <v-col cols="12" sm="">
                  <div v-if="!profile.editing">
                    <h2 class="text-h3">
                      {{ profile.givenName }} {{ profile.familyName }}
                    </h2>
                  </div>
                  <div v-else>
                    <v-row dense>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model.trim="profile.givenName"
                          :error-messages="givenNameErrors"
                          :counter="20"
                          label="First name"
                          required
                          :disabled="profile.formLoadStatus === STATUS.LOADING"
                          @input="$v.profile.givenName.$touch()"
                          @blur="$v.profile.givenName.$touch()"
                          :prepend-icon="'mdi-account'"
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model.trim="profile.familyName"
                          :error-messages="familyNameErrors"
                          :counter="12"
                          label="Last name"
                          required
                          :disabled="profile.formLoadStatus === STATUS.LOADING"
                          @input="$v.profile.familyName.$touch()"
                          @blur="$v.profile.familyName.$touch()"
                          :prepend-icon="
                            $vuetify.breakpoint.smAndUp ? '' : 'mdi-blank'
                          "
                        />
                      </v-col>
                    </v-row>
                  </div>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12" sm="6">
                  <div v-if="!profile.editing">
                    <div class="mb-1 text-overline text--secondary">Email</div>
                    <div>
                      {{ profile.email }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="profile.email"
                      :error-messages="emailErrors"
                      :type="'email'"
                      label="Email"
                      required
                      :disabled="profile.formLoadStatus === STATUS.LOADING"
                      @input="
                        () => {
                          $v.profile.email.$touch();
                          this.profile.apiResponseCode = '';
                        }
                      "
                      @blur="$v.profile.email.$touch()"
                      :prepend-icon="'mdi-email'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!profile.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Username
                    </div>
                    <div>
                      {{ profile.username }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model.trim="profile.username"
                      label="Username"
                      required
                      :disabled="true"
                      :prepend-icon="'mdi-identifier'"
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
          <v-card :loading="password.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Update Password
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    :disabled="password.formLoadStatus === STATUS.LOADING"
                    color="success"
                    @click="savePassword"
                  >
                    <div v-if="password.formLoadStatus !== STATUS.LOADING">
                      Update
                    </div>
                    <div v-else-if="password.formLoadStatus === STATUS.LOADING">
                      Updating
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-expand-transition>
                <div v-show="password.successAlert">
                  <v-alert
                    type="info"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Password updated!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-show="password.formLoadStatus === STATUS.ERROR">
                  <v-alert
                    type="error"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Failed changing password!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-expand-transition>
                <div v-show="!$v.password.oldPassword.correct">
                  <v-alert
                    type="error"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    Incorrect old password!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-row dense>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model="password.oldPassword"
                    :error-messages="oldPasswordErrors"
                    :type="'password'"
                    label="Old password"
                    required
                    :disabled="password.formLoadStatus === STATUS.LOADING"
                    :prepend-icon="'mdi-lock'"
                    @input="
                      () => {
                        $v.password.oldPassword.$touch();
                        this.password.apiResponseCode = '';
                      }
                    "
                    @blur="$v.password.oldPassword.$touch()"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model="password.newPassword"
                    :error-messages="newPasswordErrors"
                    :type="'password'"
                    label="New password"
                    hint="At least 8 characters"
                    required
                    :disabled="password.formLoadStatus === STATUS.LOADING"
                    :prepend-icon="
                      $vuetify.breakpoint.smAndUp ? '' : 'mdi-blank'
                    "
                    @input="$v.password.newPassword.$touch()"
                    @blur="$v.password.newPassword.$touch()"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model="password.confirmNewPassword"
                    :error-messages="confirmNewPasswordErrors"
                    :type="'password'"
                    label="Confirm new password"
                    required
                    :disabled="password.formLoadStatus === STATUS.LOADING"
                    :prepend-icon="
                      $vuetify.breakpoint.smAndUp ? '' : 'mdi-blank'
                    "
                    @input="$v.password.confirmNewPassword.$touch()"
                    @blur="$v.password.confirmNewPassword.$touch()"
                  />
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
          <v-card :loading="mfa.formLoadStatus === STATUS.LOADING">
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Multi-Factor Authentication
                </v-col>
              </v-row>
            </v-card-title>
            <v-divider inset />
            <div class="v-card__body">
              <v-expand-transition>
                <div v-show="mfa.successAlert">
                  <v-alert
                    type="info"
                    text
                    dense
                    transition="scroll-y-transition"
                  >
                    MFA {{ mfa.enabled ? "enabled" : "disabled" }}!
                  </v-alert>
                </div>
              </v-expand-transition>
              <v-row justify="end" align="center">
                <v-col cols="">
                  <div>
                    TOTP
                  </div>
                  <div class="text--secondary">
                    Time-based one-time password
                  </div>
                </v-col>
                <v-col cols="auto">
                  <v-dialog
                    v-model="mfa.prompt"
                    width="500"
                    :persistent="mfa.formLoadStatus === STATUS.LOADING"
                    @input="v => v || cancelMFA()"
                  >
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn
                        rounded
                        text
                        outlined
                        :color="mfa.enabled ? 'error' : 'success'"
                        v-bind="attrs"
                        v-on="on"
                        @click="enableMFA"
                      >
                        {{ mfa.enabled ? "Disable" : "Enable" }}
                      </v-btn>
                    </template>
                    <v-card
                      v-if="mfa.modeEnable"
                      :loading="mfa.formLoadStatus === STATUS.LOADING"
                    >
                      <v-card-title>
                        <v-row no-gutters align="center">
                          <v-col cols="auto">
                            Disable TOTP
                          </v-col>
                          <v-spacer />
                          <v-col cols="auto">
                            <v-btn text icon color="grey" @click="cancelMFA">
                              <v-icon v-text="'mdi-close'" />
                            </v-btn>
                          </v-col>
                        </v-row>
                      </v-card-title>
                      <v-divider inset />
                      <div class="v-card__body">
                        <v-expand-transition>
                          <div v-show="mfa.formLoadStatus === STATUS.ERROR">
                            <v-alert
                              type="error"
                              text
                              dense
                              transition="scroll-y-transition"
                            >
                              Failed disabling MFA!
                            </v-alert>
                          </div>
                        </v-expand-transition>
                        <v-row align="center">
                          <v-col>
                            <div class="mb-4">
                              Are you sure you want to disable TOTP?
                            </div>
                            <v-btn
                              rounded
                              block
                              outlined
                              color="error"
                              @click="disableMFA"
                            >
                              Disable
                            </v-btn>
                          </v-col>
                        </v-row>
                      </div>
                    </v-card>
                    <v-card
                      v-else
                      :loading="mfa.formLoadStatus === STATUS.LOADING"
                    >
                      <v-card-title>
                        <v-row no-gutters align="center">
                          <v-col cols="auto">
                            Enable TOTP
                          </v-col>
                          <v-spacer />
                          <v-col cols="auto">
                            <v-btn
                              text
                              icon
                              color="grey"
                              @click="cancelMFA"
                              :disabled="mfa.formLoadStatus === STATUS.LOADING"
                            >
                              <v-icon v-text="'mdi-close'" />
                            </v-btn>
                          </v-col>
                        </v-row>
                      </v-card-title>
                      <v-divider inset />
                      <div class="v-card__body">
                        <v-expand-transition>
                          <div v-show="mfa.formLoadStatus === STATUS.ERROR">
                            <v-alert
                              type="error"
                              text
                              dense
                              transition="scroll-y-transition"
                            >
                              Failed enabling MFA!
                            </v-alert>
                          </div>
                        </v-expand-transition>
                        <v-expand-transition>
                          <div v-show="!$v.mfa.otp.correct">
                            <v-alert
                              type="error"
                              text
                              dense
                              transition="scroll-y-transition"
                            >
                              Incorrect OTP code!
                            </v-alert>
                          </div>
                        </v-expand-transition>
                        <v-row align="center">
                          <v-col>
                            <div class="mb-4">
                              Scan the following QR code on your authenticator
                              app. You can use apps like Google Authenticator or
                              Authy.
                            </div>
                            <div class="text-center">
                              <v-expand-transition>
                                <div v-if="mfa.uri">
                                  <qrcode-vue
                                    :value="mfa.uri"
                                    :size="280"
                                    :margin="2"
                                    renderAs="svg"
                                    foreground="#121212"
                                    level="L"
                                  />
                                </div>
                              </v-expand-transition>
                              <v-expand-transition>
                                <div v-if="!mfa.uri">
                                  <v-progress-circular
                                    indeterminate
                                    size="48"
                                    class="ma-4"
                                  />
                                </div>
                              </v-expand-transition>
                            </div>
                            <div class="mt-4">
                              Enter the generated code below.
                            </div>
                            <div>
                              <v-text-field
                                v-model="mfa.otp"
                                class="pt-0"
                                placeholder="OTP"
                                :disabled="
                                  mfa.formLoadStatus === STATUS.LOADING
                                "
                                :prepend-icon="'mdi-two-factor-authentication'"
                                :error-messages="otpErrors"
                                @input="
                                  () => {
                                    $v.mfa.otp.$touch();
                                    this.mfa.apiResponseCode = '';
                                  }
                                "
                                @blur="$v.mfa.otp.$touch()"
                              />
                            </div>
                            <v-btn
                              rounded
                              block
                              outlined
                              color="success"
                              @click="confirmMFA"
                            >
                              Enable
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
import QrcodeVue from "qrcode.vue";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";
import {
  and,
  email,
  maxLength,
  minLength,
  numeric,
  required,
  sameAs
} from "vuelidate/lib/validators";

export default Vue.extend({
  components: {
    QrcodeVue
  },
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    profile: {
      givenName: "",
      familyName: "",
      username: "",
      email: "",
      editing: false,
      before: {
        givenName: "",
        familyName: "",
        email: ""
      },
      formLoadStatus: STATUS.IDLE,
      apiResponseCode: "",
      successAlert: false
    },
    password: {
      oldPassword: "",
      newPassword: "",
      confirmNewPassword: "",
      formLoadStatus: STATUS.IDLE,
      apiResponseCode: "",
      successAlert: false
    },
    mfa: {
      prompt: false,
      modeEnable: false,
      enabled: false,
      uri: "",
      otp: "",
      formLoadStatus: STATUS.IDLE,
      apiResponseCode: "",
      successAlert: false
    }
  }),

  computed: {
    profileUpdated: {
      cache: false,
      get: function() {
        return (
          this.profile.givenName !== this.profile.before.givenName ||
          this.profile.familyName !== this.profile.before.familyName ||
          this.profile.email !== this.profile.before.email
        );
      }
    },
    givenNameErrors() {
      const errors: string[] = [];
      // ?. operator fixed annoying TS strict null checks on nested Vuelidate validators
      if (!this.$v.profile.givenName?.$dirty) return errors;
      !this.$v.profile.givenName.required && errors.push("Name required");
      !this.$v.profile.givenName.maxLength && errors.push("Name too long");
      return errors;
    },
    familyNameErrors() {
      const errors: string[] = [];
      if (!this.$v.profile.familyName?.$dirty) return errors;
      !this.$v.profile.familyName.required && errors.push("Name required");
      !this.$v.profile.familyName.maxLength && errors.push("Name too long");
      return errors;
    },
    emailErrors() {
      const errors: string[] = [];
      if (!this.$v.profile.email?.$dirty) return errors;
      !this.$v.profile.email.required && errors.push("Email required");
      !this.$v.profile.email.email && errors.push("Invalid email");
      !errors.length &&
        !this.$v.profile.email.isUnique &&
        !this.$v.profile.email.$pending &&
        errors.push("Email already used");
      return errors;
    },
    oldPasswordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.oldPassword?.$dirty) return errors;
      !this.$v.password.oldPassword.required &&
        errors.push("Password required");
      !this.$v.password.oldPassword.correct && errors.push("");
      return errors;
    },
    newPasswordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.newPassword?.$dirty) return errors;
      !this.$v.password.newPassword?.required &&
        errors.push("Password required");
      !this.$v.password.newPassword?.minLength &&
        errors.push("Password too short");
      !this.$v.password.newPassword?.maxLength &&
        errors.push("Password too long");
      return errors;
    },
    confirmNewPasswordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.confirmNewPassword?.$dirty) return errors;
      !this.$v.password.confirmNewPassword?.required &&
        errors.push("Re-enter password");
      !this.$v.password.confirmNewPassword?.sameAsPassword &&
        errors.push("Passwords do not match");
      return errors;
    },
    otpErrors() {
      const errors: string[] = [];
      if (!this.$v.mfa.otp?.$dirty) return errors;
      !this.$v.mfa.otp?.required && errors.push("OTP Required");
      !this.$v.mfa.otp?.length && errors.push("Invalid OTP");
      !this.$v.mfa.otp?.numeric && errors.push("Invalid OTP");
      return errors;
    }
  },

  validations: {
    profile: {
      givenName: { required, maxLength: maxLength(20) },
      familyName: { required, maxLength: maxLength(12) },
      email: {
        required,
        email,
        maxLength: maxLength(50),
        isUnique(value) {
          if (value === this.profile.before.email) return true;
          if (this.profile.apiResponseCode === "email_exists") return false;
          return api.form
            .checkUnique({
              field: "user:email",
              value: value
            })
            .then(response => {
              return response.data.data;
            })
            .catch(() => true);
        }
      }
    },
    password: {
      oldPassword: {
        required,
        correct() {
          return (
            this.$data.password.apiResponseCode !== "incorrect_credentials"
          );
        }
      },
      newPassword: {
        required,
        minLength: minLength(8),
        maxLength: maxLength(100)
      },
      confirmNewPassword: {
        required,
        sameAsPassword: sameAs("newPassword")
      }
    },
    mfa: {
      otp: {
        required,
        numeric,
        length: and(minLength(6), maxLength(6)),
        correct() {
          return this.$data.mfa.apiResponseCode !== "incorrect_credentials";
        }
      }
    }
  },

  created() {
    api.user
      .detail()
      .then(response => {
        this.profile.givenName = response.data.data.given_name;
        this.profile.familyName = response.data.data.family_name;
        this.profile.username = response.data.data.preferred_username;
        this.profile.email = response.data.data.email;
        this.mfa.enabled = response.data.data.mfa_enabled;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(() => {
        this.pageLoadStatus = STATUS.ERROR;
      });
  },

  methods: {
    cancelProfile() {
      this.profile.editing = false;
      this.profile.givenName = this.profile.before.givenName;
      this.profile.familyName = this.profile.before.familyName;
      this.profile.email = this.profile.before.email;
      this.profile.email = this.profile.before.email;
      this.profile.formLoadStatus = STATUS.IDLE;
      this.profile.apiResponseCode = "";
      this.profile.before = {
        givenName: "",
        familyName: "",
        email: ""
      };
      this.$v.profile.$reset();
    },
    saveProfile() {
      if (!this.profile.editing) {
        this.profile.editing = true;
        this.profile.before.givenName = this.profile.givenName;
        this.profile.before.familyName = this.profile.familyName;
        this.profile.before.email = this.profile.email;
        this.profile.successAlert = false;
        return;
      }
      this.$v.profile.$touch();
      if (!this.$v.profile.$invalid) {
        this.profile.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.user
              .update({
                /* eslint-disable @typescript-eslint/camelcase */
                given_name: this.profile.givenName.trim(),
                family_name: this.profile.familyName.trim(),
                email: this.profile.email.trim()
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(() => {
                this.profile.editing = false;
                this.profile.formLoadStatus = STATUS.COMPLETE;
                this.profile.successAlert = true;
                setTimeout(() => {
                  this.profile.successAlert = false;
                }, 5000);
              })
              .catch(error => {
                this.profile.editing = true;
                this.profile.apiResponseCode = error.response.data.code;
                this.profile.formLoadStatus = !this.profile.apiResponseCode
                  ? STATUS.ERROR
                  : STATUS.IDLE;
              }),
          2000
        );
      }
    },
    savePassword() {
      this.password.successAlert = false;
      this.$v.password.$touch();
      if (!this.$v.password.$invalid) {
        this.password.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.user
              .updatePassword({
                /* eslint-disable @typescript-eslint/camelcase */
                old_password: this.password.oldPassword,
                new_password: this.password.newPassword
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(() => {
                this.password.formLoadStatus = STATUS.COMPLETE;
                this.password.oldPassword = "";
                this.password.newPassword = "";
                this.password.confirmNewPassword = "";
                this.$v.password.$reset();
                this.password.successAlert = true;
                setTimeout(() => {
                  this.password.successAlert = false;
                }, 5000);
              })
              .catch(error => {
                this.password.apiResponseCode = error.response.data.code;
                this.password.formLoadStatus = !this.password.apiResponseCode
                  ? STATUS.ERROR
                  : STATUS.IDLE;
              }),
          2000
        );
      }
    },
    enableMFA() {
      this.mfa.successAlert = false;
      this.mfa.modeEnable = this.mfa.enabled;
      if (!this.mfa.enabled) {
        api.mfa
          .enable()
          .then(response => {
            this.mfa.uri = response.data.data;
          })
          .catch(() => {
            this.mfa.formLoadStatus = STATUS.ERROR;
          });
      }
    },
    confirmMFA() {
      this.$v.mfa.$touch();
      if (!this.$v.mfa.$invalid) {
        this.mfa.formLoadStatus = STATUS.LOADING;
        api.mfa
          .confirm(this.mfa.otp)
          .then(() => {
            this.mfa.formLoadStatus = STATUS.COMPLETE;
            this.mfa.enabled = true;
            this.mfa.successAlert = true;
            setTimeout(() => {
              this.mfa.successAlert = false;
            }, 5000);
            this.cancelMFA();
          })
          .catch(error => {
            this.mfa.apiResponseCode = error.response.data.code;
            this.mfa.formLoadStatus = !this.mfa.apiResponseCode
              ? STATUS.ERROR
              : STATUS.IDLE;
          });
      }
    },
    disableMFA() {
      this.mfa.formLoadStatus = STATUS.LOADING;
      api.mfa
        .disable()
        .then(() => {
          this.mfa.formLoadStatus = STATUS.COMPLETE;
          this.mfa.enabled = false;
          this.mfa.successAlert = true;
          setTimeout(() => {
            this.mfa.successAlert = false;
          }, 5000);
          this.cancelMFA();
        })
        .catch(error => {
          this.mfa.apiResponseCode = error.response.data.code;
          this.mfa.formLoadStatus = !this.mfa.apiResponseCode
            ? STATUS.ERROR
            : STATUS.IDLE;
        });
    },
    cancelMFA() {
      this.mfa.prompt = false;
      this.mfa.otp = "";
      this.mfa.uri = "";
      this.mfa.formLoadStatus = STATUS.IDLE;
      this.$v.mfa.$reset();
    }
  }
});
</script>
