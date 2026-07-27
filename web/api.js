export default function createAPI(scheme, host) {
  const baseURL = `${scheme}://${host}`;

  return {
    async showWorkflow(workflowID, options = {}) {
      const {
        responseType = "text",
        timeout = 20_000,
        signal = AbortSignal.timeout(timeout),
        ...requestOptions
      } = options;
      const url = new URL(
        `/workflow/${encodeURIComponent(workflowID)}`,
        baseURL,
      );
      const response = await fetch(url, {
        ...requestOptions,
        method: "GET",
        signal,
      });

      if (!response.ok) {
        throw new Error(
          `GET ${url.pathname} failed: ${response.status} ${response.statusText}`,
        );
      }

      return responseType === "blob" ? response.blob() : response.text();
    },
  };
}
