export enum STATUS {
  PRE_LOADING = "PRE_LOADING",
  LOADING = "LOADING",
  IDLE = "IDLE",
  COMPLETE = "COMPLETE",
  BAD_REQUEST = "BAD_REQUEST"
}

export const StatusMixin = {
  data: () => ({
    STATUS
  })
};
