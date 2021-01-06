import axios, { AxiosInstance, AxiosResponse } from "axios";
import qs from "qs";
import { authManager } from "@/auth";
import { ACCESS_TOKEN } from "@/auth/AuthManager";

const oauthClient: AxiosInstance = axios.create({
  baseURL: "/oauth/"
});

const withAuth = () => ({
  headers: {
    Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
  }
});

export default {
  authorize: function(authRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`authorize`, authRequest);
  },
  token: function(tokenRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`token`, qs.stringify(tokenRequest), {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    });
  },
  logout: function(logoutRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`logout`, logoutRequest, withAuth());
  }
};
