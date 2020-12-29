import axios, { AxiosInstance, AxiosResponse } from "axios";

const apiClient: AxiosInstance = axios.create({
  baseURL: "api/v1/"
});

export default {
  application: {
    getOne: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.get(`application/${clientId}`);
    }
  },
  form: {
    checkUnique: function(uniqueRequest: unknown): Promise<AxiosResponse> {
      return apiClient.post("form/unique", uniqueRequest);
    }
  },
  user: {
    signup: function(userSignup: unknown): Promise<AxiosResponse> {
      return apiClient.post("user", userSignup);
    }
  }
};
