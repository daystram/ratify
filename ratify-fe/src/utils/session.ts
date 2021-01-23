export interface SessionInfo {
  session_id: string;
  ip: string;
  browser: string;
  os: string;
  mobile: boolean;
  issued_at: number;
  current: boolean;
  date?: Date;
}
