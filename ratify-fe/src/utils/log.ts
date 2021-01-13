function dayDiff(from: Date, to: Date): number {
  return (
    new Date(from.toDateString()).valueOf() -
    new Date(to.toDateString()).valueOf()
  );
}

export function addDateSeparator(
  date: Date,
  activities: [{ separator: boolean; today: boolean; date: Date }]
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
