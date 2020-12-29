export enum STATUS {
  PRE_LOADING = "PRE_LOADING",
  LOADING = "LOADING",
  IDLE = "IDLE",
  COMPLETE = "COMPLETE",
  ERROR = "ERROR",
  BAD_REQUEST = "BAD_REQUEST"
}

export const StatusMixin = {
  data: () => ({
    STATUS
  })
};
