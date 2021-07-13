function validate(urlString: string, allowInsecure?: boolean): boolean {
  try {
    const url = new URL(urlString);
    const regex = /^([-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,})$|^(localhost)$/;
    return (
      ((url.protocol === "http:" && allowInsecure) ||
        url.protocol === "https:") &&
      regex.test(url.hostname) &&
      url.origin !== null
    );
  } catch (e) {
    return false;
  }
}

export const validateURL = (allowInsecure?: boolean) => (urlString: string) =>
  validate(urlString, allowInsecure);
