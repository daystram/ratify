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
          <v-card>
            <v-card-title>
              <v-row no-gutters align="center">
                <v-col cols="auto">
                  Profile
                </v-col>
                <v-spacer />
                <v-col cols="auto">
                  <v-btn
                    v-if="
                      user.editing && user.formLoadStatus !== STATUS.LOADING
                    "
                    text
                    rounded
                    color="error"
                    @click="cancelUser"
                  >
                    Cancel
                  </v-btn>
                  <v-btn
                    text
                    rounded
                    class="ml-4"
                    :disabled="!userUpdated"
                    :color="user.editing ? 'success' : 'secondary lighten-1'"
                    @click="saveUser"
                  >
                    <div v-if="!user.editing">Edit</div>
                    <div
                      v-else-if="
                        user.editing && user.formLoadStatus !== STATUS.LOADING
                      "
                    >
                      Save
                    </div>
                    <div
                      v-else-if="
                        user.editing && user.formLoadStatus === STATUS.LOADING
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
                <div v-show="user.formLoadStatus === STATUS.ERROR">
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
                      (user.givenName && user.givenName[0].toUpperCase()) +
                        (user.familyName && user.familyName[0].toUpperCase())
                    }}
                  </v-avatar>
                </v-col>
                <v-col cols="12" sm="">
                  <div v-if="!user.editing">
                    <h2 class="text-h3">
                      {{ user.givenName }} {{ user.familyName }}
                    </h2>
                  </div>
                  <div v-else>
                    <v-row dense>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="user.givenName"
                          :error-messages="givenNameErrors"
                          :counter="20"
                          label="First name"
                          required
                          :disabled="user.formLoadStatus === STATUS.LOADING"
                          @input="$v.user.givenName.$touch()"
                          @blur="$v.user.givenName.$touch()"
                          :prepend-icon="'mdi-account'"
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="user.familyName"
                          :error-messages="familyNameErrors"
                          :counter="12"
                          label="Last name"
                          required
                          :disabled="user.formLoadStatus === STATUS.LOADING"
                          @input="$v.user.familyName.$touch()"
                          @blur="$v.user.familyName.$touch()"
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
                  <div v-if="!user.editing">
                    <div class="mb-1 text-overline text--secondary">Email</div>
                    <div>
                      {{ user.email }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="user.email"
                      :error-messages="emailErrors"
                      :type="'email'"
                      label="Email"
                      required
                      :disabled="user.formLoadStatus === STATUS.LOADING"
                      @input="
                        () => {
                          $v.user.email.$touch();
                          this.apiResponseCode = '';
                        }
                      "
                      @blur="$v.user.email.$touch()"
                      :prepend-icon="'mdi-email'"
                    />
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div v-if="!user.editing">
                    <div class="mb-1 text-overline text--secondary">
                      Username
                    </div>
                    <div>
                      {{ user.username }}
                    </div>
                  </div>
                  <div v-else>
                    <v-text-field
                      v-model="user.username"
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
import { email, maxLength, required } from "vuelidate/lib/validators";

export default Vue.extend({
  data: () => ({
    pageLoadStatus: STATUS.PRE_LOADING,
    user: {
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
      apiResponseCode: ""
    },
    password: {
      oldPassword: "",
      newPassword: "",
      confirmNewPassword: ""
    }
  }),

  computed: {
    userUpdated: {
      cache: false,
      get: function() {
        return (
          this.user.givenName !== this.user.before.givenName ||
          this.user.familyName !== this.user.before.familyName ||
          this.user.email !== this.user.before.email
        );
      }
    },
    givenNameErrors() {
      const errors: string[] = [];
      // ?. operator fixed annoying TS strict null checks on nested Vuelidate validators
      if (!this.$v.user.givenName?.$dirty) return errors;
      !this.$v.user.givenName.required && errors.push("Name required");
      !this.$v.user.givenName.maxLength && errors.push("Name too long");
      return errors;
    },
    familyNameErrors() {
      const errors: string[] = [];
      if (!this.$v.user.familyName?.$dirty) return errors;
      !this.$v.user.familyName.required && errors.push("Name required");
      !this.$v.user.familyName.maxLength && errors.push("Name too long");
      return errors;
    },
    emailErrors() {
      const errors: string[] = [];
      if (!this.$v.user.email?.$dirty) return errors;
      !this.$v.user.email.required && errors.push("Email required");
      !this.$v.user.email.email && errors.push("Invalid email");
      !errors.length &&
        !this.$v.user.email.isUnique &&
        !this.$v.user.email.$pending &&
        errors.push("Email already used");
      return errors;
    }
  },

  validations: {
    user: {
      givenName: { required, maxLength: maxLength(20) },
      familyName: { required, maxLength: maxLength(12) },
      email: {
        required,
        email,
        maxLength: maxLength(50),
        isUnique(value) {
          if (value === this.user.before.email) return true;
          if (this.user.apiResponseCode === "email_exists") return false;
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
    }
  },

  created() {
    api.user
      .detail()
      .then(response => {
        this.user.givenName = response.data.data.given_name;
        this.user.familyName = response.data.data.family_name;
        this.user.username = response.data.data.preferred_username;
        this.user.email = response.data.data.email;
        this.pageLoadStatus = STATUS.COMPLETE;
      })
      .catch(error => {
        console.error(error.data);
        this.pageLoadStatus = STATUS.ERROR;
      });
  },

  methods: {
    cancelUser() {
      this.user.editing = false;
      this.user.givenName = this.user.before.givenName;
      this.user.familyName = this.user.before.familyName;
      this.user.email = this.user.before.email;
      this.user.email = this.user.before.email;
      this.user.formLoadStatus = STATUS.IDLE;
      this.user.apiResponseCode = "";
      this.user.before = {
        givenName: "",
        familyName: "",
        email: ""
      };
      this.$v.$reset();
    },
    saveUser() {
      if (!this.user.editing) {
        this.user.editing = true;
        this.user.before.givenName = this.user.givenName;
        this.user.before.familyName = this.user.familyName;
        this.user.before.email = this.user.email;
        return;
      }
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.user.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.user
              .update({
                /* eslint-disable @typescript-eslint/camelcase */
                given_name: this.user.givenName,
                family_name: this.user.familyName,
                email: this.user.email
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(() => {
                this.user.editing = false;
                this.user.formLoadStatus = STATUS.COMPLETE;
              })
              .catch(error => {
                console.error(error.response.data);
                this.user.editing = true;
                this.user.apiResponseCode = error.response.data.code;
                this.user.formLoadStatus = !this.user.apiResponseCode
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
