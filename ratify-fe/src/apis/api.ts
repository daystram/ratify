import axios, { AxiosInstance, AxiosResponse } from "axios";
import { authManager, refreshAuth } from "@/auth";
import { ACCESS_TOKEN } from "@/auth/AuthManager";
import router from "@/router";

const apiClient: AxiosInstance = axios.create({
  baseURL: "/api/v1/"
});

apiClient.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response.status === 401) {
      refreshAuth(router.currentRoute.fullPath);
    }
    return Promise.reject(error);
  }
);

const withAuth = () => ({
  headers: {
    Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
  }
});

export default {
  application: {
    detail: function(
      clientId: string,
      complete?: boolean
    ): Promise<AxiosResponse> {
      return apiClient.get(
        `application/${clientId}`,
        complete ? withAuth() : {}
      );
    },
    update: function(
      clientId: string,
      application: object
    ): Promise<AxiosResponse> {
      return apiClient.put(`application/${clientId}`, application, withAuth());
    },
    revoke: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.put(`application/${clientId}/revoke`, {}, withAuth());
    },
    delete: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.delete(`application/${clientId}`, withAuth());
    },
    list: function(): Promise<AxiosResponse> {
      return apiClient.get(`application/`, withAuth());
    },
    register: function(application: object): Promise<AxiosResponse> {
      return apiClient.post(`application`, application, withAuth());
    }
  },
  form: {
    checkUnique: function(uniqueRequest: object): Promise<AxiosResponse> {
      return apiClient.post(`form/unique`, uniqueRequest);
    }
  },
  user: {
    detail: function(subject?: string): Promise<AxiosResponse> {
      return apiClient.get(`user/${subject || ""}`, withAuth());
    },
    update: function(user: object): Promise<AxiosResponse> {
      return apiClient.put(`user`, user, withAuth());
    },
    signup: function(userSignup: object): Promise<AxiosResponse> {
      return apiClient.post("user", userSignup);
    },
    verify: function(token: string): Promise<AxiosResponse> {
      return apiClient.post("user/verify", { token });
    },
    resend: function(email: string): Promise<AxiosResponse> {
      return apiClient.post("user/resend", { email });
    }
  }
};
