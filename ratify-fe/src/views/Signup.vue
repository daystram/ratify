<template>
  <v-container fluid class="fill-height gradient-bg">
    <v-col>
      <v-row align="center" justify="center">
        <div class="login-form">
          <v-scroll-y-transition appear>
            <div>
              <v-card
                elevation="12"
                :loading="formLoadStatus === STATUS.LOADING"
              >
                <div class="pa-6">
                  <h1 class="text-h2 mt-12 mb-12 text-center">
                    Sign Up
                  </h1>

                  <v-expand-transition>
                    <div v-show="formLoadStatus === STATUS.COMPLETE">
                      <v-alert
                        type="success"
                        text
                        dense
                        transition="scroll-y-transition"
                      >
                        Account successfully created!
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
                        Failed creating account!
                      </v-alert>
                    </div>
                  </v-expand-transition>
                  <v-form
                    v-if="formLoadStatus !== STATUS.COMPLETE"
                    @submit="signup"
                  >
                    <v-row dense>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="firstname"
                          :error-messages="firstnameErrors"
                          :counter="20"
                          label="First name"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          @input="$v.firstname.$touch()"
                          @blur="$v.firstname.$touch()"
                          :prepend-icon="'mdi-account'"
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="lastname"
                          :error-messages="lastnameErrors"
                          :counter="12"
                          label="Last name"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          @input="$v.lastname.$touch()"
                          @blur="$v.lastname.$touch()"
                          :prepend-icon="
                            $vuetify.breakpoint.smAndUp ? '' : 'mdi-blank'
                          "
                        />
                      </v-col>
                    </v-row>
                    <v-text-field
                      v-model="username"
                      :error-messages="usernameErrors"
                      :counter="12"
                      label="Username"
                      hint="Username is permanent"
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
                    />
                    <v-text-field
                      v-model="email"
                      :error-messages="emailErrors"
                      :type="'email'"
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
                    />
                    <v-row dense>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="password"
                          :error-messages="passwordErrors"
                          :type="'password'"
                          label="Password"
                          hint="At least 8 characters"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          :prepend-icon="'mdi-lock'"
                          @input="$v.password.$touch()"
                          @blur="$v.password.$touch()"
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-text-field
                          v-model="confirmPassword"
                          :error-messages="confirmPasswordErrors"
                          :type="'password'"
                          label="Confirm password"
                          required
                          :disabled="
                            formLoadStatus === STATUS.LOADING ||
                              formLoadStatus === STATUS.COMPLETE
                          "
                          @input="$v.confirmPassword.$touch()"
                          @blur="$v.confirmPassword.$touch()"
                        />
                      </v-col>
                    </v-row>
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
                        creating account
                      </div>
                      <div v-else>
                        create account
                      </div>
                    </v-btn>
                  </v-form>
                </div>
              </v-card>
              <p class="text-center mt-4 text-subtitle-1 text--secondary">
                Already have an account?
                <router-link to="/login" class="text-link">Sign In</router-link>
              </p>
            </div>
          </v-scroll-y-transition>
        </div>
      </v-row>
    </v-col>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import {
  required,
  alphaNum,
  minLength,
  sameAs,
  email,
  maxLength
} from "vuelidate/lib/validators";
import api from "@/apis/api";
import { STATUS } from "@/constants/status";

import "@/styles/Authorize.sass";

export default Vue.extend({
  data: function() {
    return {
      formLoadStatus: STATUS.IDLE,
      firstname: "",
      lastname: "",
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      showPassword: false,
      apiResponseCode: ""
    };
  },

  computed: {
    firstnameErrors() {
      const errors: string[] = [];
      if (!this.$v.firstname.$dirty) return errors;
      !this.$v.firstname.required && errors.push("Name required");
      !this.$v.firstname.maxLength && errors.push("Name too long");
      return errors;
    },
    lastnameErrors() {
      const errors: string[] = [];
      if (!this.$v.lastname.$dirty) return errors;
      !this.$v.lastname.required && errors.push("Name required");
      !this.$v.lastname.maxLength && errors.push("Name too long");
      return errors;
    },
    usernameErrors() {
      const errors: string[] = [];
      if (!this.$v.username.$dirty) return errors;
      !this.$v.username.required && errors.push("Username required");
      !this.$v.username.alphaNum && errors.push("Invalid username");
      !this.$v.username.maxLength && errors.push("Username too long");
      !errors.length && // prevent calling API when previous errors are not resolved
        !this.$v.username.isUnique &&
        !this.$v.username.$pending &&
        errors.push("Username already used");
      return errors;
    },
    emailErrors() {
      const errors: string[] = [];
      if (!this.$v.email.$dirty) return errors;
      !this.$v.email.required && errors.push("Email required");
      !this.$v.email.email && errors.push("Invalid email");
      !errors.length &&
        !this.$v.email.isUnique &&
        !this.$v.email.$pending &&
        errors.push("Email already used");
      return errors;
    },
    passwordErrors() {
      const errors: string[] = [];
      if (!this.$v.password.$dirty) return errors;
      !this.$v.password.required && errors.push("Password required");
      !this.$v.password.minLength && errors.push("Password too short");
      !this.$v.password.maxLength && errors.push("Password too long");
      return errors;
    },
    confirmPasswordErrors() {
      const errors: string[] = [];
      if (!this.$v.confirmPassword.$dirty) return errors;
      !this.$v.confirmPassword.required && errors.push("Re-enter password");
      !this.$v.confirmPassword.sameAsPassword &&
        errors.push("Passwords do not match");
      return errors;
    }
  },

  validations: {
    firstname: { required, maxLength: maxLength(20) },
    lastname: { required, maxLength: maxLength(12) },
    username: {
      required,
      alphaNum,
      maxLength: maxLength(12),
      isUnique(value) {
        if (value === "") return true;
        if (this.apiResponseCode === "username_exists") return false;
        return api.form
          .checkUnique({
            field: "user:username",
            value: value
          })
          .then(response => {
            return response.data.data;
          })
          .catch(() => true);
      }
    },
    email: {
      required,
      email,
      maxLength: maxLength(50),
      isUnique(value) {
        if (value === "") return true;
        if (this.apiResponseCode === "email_exists") return false;
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
    },
    password: { required, minLength: minLength(8), maxLength: maxLength(100) },
    confirmPassword: {
      required,
      sameAsPassword: sameAs("password")
    }
  },

  methods: {
    signup(e: Event) {
      e.preventDefault();
      this.$v.$touch();
      if (!this.$v.$invalid) {
        this.formLoadStatus = STATUS.LOADING;
        setTimeout(
          () =>
            api.user
              .signup({
                /* eslint-disable @typescript-eslint/camelcase */
                given_name: this.firstname,
                family_name: this.lastname,
                preferred_username: this.username,
                email: this.email,
                password: this.password
                /* eslint-enable @typescript-eslint/camelcase */
              })
              .then(() => {
                this.formLoadStatus = STATUS.COMPLETE;
              })
              .catch(error => {
                this.apiResponseCode = error.response.data.code;
                this.formLoadStatus = !this.apiResponseCode
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
