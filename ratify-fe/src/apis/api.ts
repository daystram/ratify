import axios, { AxiosInstance, AxiosResponse } from "axios";
import { ACCESS_TOKEN } from "@daystram/ratify-client";
import { authManager, refreshAuth } from "@/auth";
import router from "@/router";

const apiClient: AxiosInstance = axios.create({
  baseURL: `${
    process.env.NODE_ENV === "development"
      ? process.env.VUE_APP_DEV_BASE_API
      : ""
  }/api/v1/`
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
      return apiClient.post(`application/`, application, withAuth());
    }
  },
  form: {
    checkUnique: function(uniqueRequest: object): Promise<AxiosResponse> {
      return apiClient.post(`form/unique`, uniqueRequest);
    }
  },
  user: {
    detail: function(subject: string): Promise<AxiosResponse> {
      return apiClient.get(`user/${subject}`, withAuth());
    },
    list: function(): Promise<AxiosResponse> {
      return apiClient.get(`user/`, withAuth());
    },
    update: function(user: object): Promise<AxiosResponse> {
      return apiClient.put(`user/`, user, withAuth());
    },
    updatePassword: function(passwords: object): Promise<AxiosResponse> {
      return apiClient.put(`user/password`, passwords, withAuth());
    },
    updateSuperuser: function(superuser: object): Promise<AxiosResponse> {
      return apiClient.put(`user/superuser`, superuser, withAuth());
    },
    signup: function(userSignup: object): Promise<AxiosResponse> {
      return apiClient.post(`user/`, userSignup);
    },
    verify: function(token: string): Promise<AxiosResponse> {
      return apiClient.post(`user/verify`, { token });
    },
    resend: function(email: string): Promise<AxiosResponse> {
      return apiClient.post(`user/resend`, { email });
    }
  },
  session: {
    list: function(): Promise<AxiosResponse> {
      return apiClient.get(`session/`, withAuth());
    },
    revoke: function(sessionId: string): Promise<AxiosResponse> {
      // eslint-disable-next-line @typescript-eslint/camelcase
      return apiClient.post(`session/`, { session_id: sessionId }, withAuth());
    }
  },
  mfa: {
    enable: function(): Promise<AxiosResponse> {
      return apiClient.post(`mfa/enable`, {}, withAuth());
    },
    confirm: function(otp: string): Promise<AxiosResponse> {
      return apiClient.post(`mfa/confirm`, { otp }, withAuth());
    },
    disable: function(): Promise<AxiosResponse> {
      return apiClient.post(`mfa/disable`, {}, withAuth());
    }
  },
  log: {
    userActivity: function(): Promise<AxiosResponse> {
      return apiClient.get(`log/user_activity`, withAuth());
    },
    adminActivity: function(): Promise<AxiosResponse> {
      return apiClient.get(`log/admin_activity`, withAuth());
    }
  },
  dashboard: {
    fetch: function(): Promise<AxiosResponse> {
      return apiClient.get(`dashboard/`, withAuth());
    }
  }
};
