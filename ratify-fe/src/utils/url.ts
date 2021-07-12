function validate(
  urlString: string,
  allowInsecure?: boolean,
  allowLocalhost?: boolean
): boolean {
  const url = new URL(urlString);
  return (
    ((url.protocol === "http:" && allowInsecure) ||
      url.protocol === "https:") &&
    (url.hostname !== "localhost" || !!allowLocalhost) &&
    url.origin !== null
  );
}

export const validateURL = (
  allowInsecure?: boolean,
  allowLocalhost?: boolean
) => (urlString: string) => validate(urlString, allowInsecure, allowLocalhost);
