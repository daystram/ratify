<template>
  <div class="manage fill-height">
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon
        class="hidden-lg-and-up"
        @click.stop="drawer = !drawer"
      />
      <h1>
        <router-link
          class="text-md-h4 text-h5 text-center"
          style="text-decoration: none; color: inherit;"
          :to="{ name: 'home' }"
        >
          <Logo />
        </router-link>
      </h1>
      <v-spacer />
      <v-menu
        right
        nudge-bottom="12px"
        offset-y
        min-width="280px"
        max-width="280px"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-avatar
            color="primaryDim"
            size="32"
            v-on="on"
            v-bind="attrs"
            style="user-select: none"
          >
            {{ user.given_name[0] + user.family_name[0] }}
          </v-avatar>
        </template>
        <v-list rounded>
          <v-list-item-group color="primary">
            <v-list-item two-line disabled>
              <v-list-item-avatar color="primaryDim" size="48">
                <div class="text-center flex-fill text--primary">
                  {{ user.given_name[0] + user.family_name[0] }}
                </div>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="text--primary">
                  {{ user.given_name + " " + user.family_name }}
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ user.preferred_username }}
                </v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
            <v-list-item :to="{ name: 'manage:profile' }" dense>
              <v-list-item-icon>
                <v-icon>mdi-account</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Account</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item :to="{ name: 'logout' }" dense>
              <v-list-item-icon>
                <v-icon color="error">mdi-logout-variant</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title class="error--text">
                  Logout
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item :to="{ name: 'logout:global' }" dense>
              <v-list-item-icon>
                <v-icon color="error">mdi-logout-variant</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title class="error--text">
                  Logout Everywhere
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>
    </v-app-bar>
    <v-navigation-drawer app clipped v-model="drawer">
      <v-list nav dense rounded>
        <v-list-item-group color="primary">
          <v-list-item :to="{ name: 'manage:dashboard' }">
            <v-list-item-icon>
              <v-icon v-text="'mdi-view-dashboard'" />
            </v-list-item-icon>
            <v-list-item-title v-text="'Dashboard'" />
          </v-list-item>
          <v-list-group
            v-if="user.is_superuser"
            :value="true"
            no-action
            prepend-icon="mdi-account-supervisor-circle"
          >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Admin</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item :to="{ name: 'manage:user' }">
              <v-list-item-title v-text="'Users'" />
              <v-list-item-icon>
                <v-icon v-text="'mdi-account-multiple'" />
              </v-list-item-icon>
            </v-list-item>
            <v-list-item :to="{ name: 'manage:application' }">
              <v-list-item-title v-text="'Applications'" />
              <v-list-item-icon>
                <v-icon v-text="'mdi-application'" />
              </v-list-item-icon>
            </v-list-item>
            <v-list-item :to="{ name: 'manage:log' }">
              <v-list-item-title v-text="'Logs'" />
              <v-list-item-icon>
                <v-icon v-text="'mdi-format-list-bulleted'" />
              </v-list-item-icon>
            </v-list-item>
          </v-list-group>
          <v-list-item :to="{ name: 'manage:session' }">
            <v-list-item-icon>
              <v-icon v-text="'mdi-dock-window'" />
            </v-list-item-icon>
            <v-list-item-title v-text="'Sessions'" />
          </v-list-item>
          <v-list-item :to="{ name: 'manage:activity' }">
            <v-list-item-icon>
              <v-icon v-text="'mdi-timeline-alert'" />
            </v-list-item-icon>
            <v-list-item-title v-text="'Activities'" />
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>
    <v-main class="fill-height">
      <v-container
        class="mx-auto pa-4 pa-sm-6 pa-md-8"
        style="max-width: 960px"
        fluid
      >
        <router-view></router-view>
      </v-container>
    </v-main>
  </div>
</template>

<script>
import Vue from "vue";
import Logo from "@/components/Logo.vue";
import { authManager } from "@/auth";

export default Vue.extend({
  components: { Logo },

  data() {
    return {
      drawer: null
    };
  },

  computed: {
    user: () => authManager.getUser()
  }
});
</script>
