function dayDiff(from: Date, to: Date): number {
  return (
    new Date(from.toDateString()).valueOf() -
    new Date(to.toDateString()).valueOf()
  );
}

export function addDateSeparator(
  date: Date,
  activities: { separator?: boolean; today?: boolean; date: Date }[]
) {
  const diff = dayDiff(
    date,
    activities.length ? activities[activities.length - 1].date : new Date()
  );
  if (diff < 0 || (!activities.length && diff === 0)) {
    activities.push({
      separator: true,
      today: diff === 0,
      date: date
    });
  }
}

export interface LogSeverityMap {
  I: string;
  W: string;
  E: string;
  F: string;
}

export interface LogInfo {
  preferred_username: string;
  application_name: string | undefined;
  client_id: string | undefined;
  severity: keyof LogSeverityMap;
  description: string;
  created_at: number;
}
