<template>
  <div class="manage">
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon
        class="hidden-lg-and-up"
        @click.stop="drawer = !drawer"
      />
      <h1
        class="text-md-h4 text-h5 text-center"
        @click="$router.push({ name: 'home' })"
        style="cursor: pointer"
      >
        <Logo />
      </h1>
      <v-spacer />
      <v-menu right nudge-bottom="12px" offset-y>
        <template v-slot:activator="{ on, attrs }">
          <v-avatar
            color="grey darken-1 shrink"
            size="32"
            v-on="on"
            v-bind="attrs"
            style="user-select: none"
          >
            DA
          </v-avatar>
        </template>
        <v-list>
          <v-list-item-group color="primary">
            <v-list-item two-line disabled>
              <v-list-item-avatar>
                <img src="https://randomuser.me/api/portraits/men/83.jpg" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>Danny August Ramaputra</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider />
            <v-list-item :to="{ name: 'manage:profile' }">
              <v-list-item-icon>
                <v-icon>mdi-account</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Account</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item @click="$router.push({ name: 'logout' })">
              <v-list-item-icon>
                <v-icon>mdi-key</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Logout</v-list-item-title>
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
          <v-list-item :to="{ name: 'manage:incident' }">
            <v-list-item-icon>
              <v-icon v-text="'mdi-timeline-alert'" />
            </v-list-item-icon>
            <v-list-item-title v-text="'Incidents'" />
          </v-list-item>
          <v-list-item :to="{ name: 'manage:setting' }">
            <v-list-item-icon>
              <v-icon v-text="'mdi-cog'" />
            </v-list-item-icon>
            <v-list-item-title v-text="'Settings'" />
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>
    <v-main>
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

export default Vue.extend({
  components: { Logo },

  data: () => ({
    drawer: null
  })
});
</script>
