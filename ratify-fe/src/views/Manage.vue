<template>
  <div class="manage">
    <v-app-bar app>
      <h1><v-icon style="font-size: inherit">mdi-lock</v-icon> Ratify</h1>
      <router-link to="/logout">Logout</router-link>
    </v-app-bar>
    <v-main>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-main>
  </div>
</template>
