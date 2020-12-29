import axios, { AxiosInstance, AxiosResponse } from "axios";
import qs from "qs";

const oauthClient: AxiosInstance = axios.create({
  baseURL: "oauth/"
});

export default {
  authorize: function(authRequest: unknown): Promise<AxiosResponse> {
    return oauthClient.post(`authorize`, authRequest);
  },
  token: function(tokenRequest: unknown): Promise<AxiosResponse> {
    return oauthClient.post(`token`, qs.stringify(tokenRequest), {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    });
  },
  revoke: function(revokeRequest: unknown): Promise<AxiosResponse> {
    return oauthClient.post(`revoke`, revokeRequest);
  }
};
