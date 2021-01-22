import axios, { AxiosInstance, AxiosResponse } from "axios";

const oauthClient: AxiosInstance = axios.create({
  baseURL: `${
    process.env.NODE_ENV === "development"
      ? process.env.VUE_APP_DEV_BASE_API
      : ""
  }/oauth/`
});

export default {
  authorize: function(authRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`authorize`, authRequest);
  }
};
