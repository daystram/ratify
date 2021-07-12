function validate(
  urlString: string,
  allowInsecure?: boolean,
  allowLocalhost?: boolean
): boolean {
  try {
    const url = new URL(urlString);
    console.log(url);
    return (
      ((url.protocol === "http:" && allowInsecure) ||
        url.protocol === "https:") &&
      (url.hostname !== "localhost" || !!allowLocalhost) &&
      url.origin !== null
    );
  } catch (e) {
    return false;
  }
}

export const validateURL = (
  allowInsecure?: boolean,
  allowLocalhost?: boolean
) => (urlString: string) => validate(urlString, allowInsecure, allowLocalhost);
