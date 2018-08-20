// This module exports functions that give access to the amflow API hosted at localhost.
// It uses the axios javascript library for making the actual HTTP requests.
define(['axios'] , function (axios) {
  function merge(obj1, obj2) {
    var obj3 = {};
    for (var attrname in obj1) { obj3[attrname] = obj1[attrname]; }
    for (var attrname in obj2) { obj3[attrname] = obj2[attrname]; }
    return obj3;
  }

  return function (scheme, host, timeout) {
    scheme = scheme || 'http';
    host = host || 'localhost';
    timeout = timeout || 20000;

    // Client is the object returned by this module.
    var client = axios;

    // URL prefix for all API requests.
    var urlPrefix = scheme + '://' + host;

  // Add link
  // path is the request path, the format is "/workflow/:workflowID/links"
  // config is an optional object to be merged into the config built by the function prior to making the request.
  // The content of the config object is described here: https://github.com/axios/axios#request-config
  // This function returns a promise which raises an error if the HTTP response is a 4xx or 5xx.
  client.addLinkWorkflow = function (path, config) {
    var cfg = {
      timeout: timeout,
      url: urlPrefix + path,
      method: 'patch',
      responseType: 'json'
    };
    if (config) {
      cfg = merge(cfg, config);
    }
    return client(cfg);
  }

  // Delete link
  // path is the request path, the format is "/workflow/:workflowID/links/:linkID"
  // config is an optional object to be merged into the config built by the function prior to making the request.
  // The content of the config object is described here: https://github.com/axios/axios#request-config
  // This function returns a promise which raises an error if the HTTP response is a 4xx or 5xx.
  client.deleteLinkWorkflow = function (path, config) {
    var cfg = {
      timeout: timeout,
      url: urlPrefix + path,
      method: 'delete',
      responseType: 'json'
    };
    if (config) {
      cfg = merge(cfg, config);
    }
    return client(cfg);
  }

  // Move link
  // path is the request path, the format is "/workflow/:workflowID/links/:linkID"
  // config is an optional object to be merged into the config built by the function prior to making the request.
  // The content of the config object is described here: https://github.com/axios/axios#request-config
  // This function returns a promise which raises an error if the HTTP response is a 4xx or 5xx.
  client.moveLinkWorkflow = function (path, config) {
    var cfg = {
      timeout: timeout,
      url: urlPrefix + path,
      method: 'patch',
      responseType: 'json'
    };
    if (config) {
      cfg = merge(cfg, config);
    }
    return client(cfg);
  }

  // Read workflow
  // path is the request path, the format is "/workflow/:workflowID"
  // config is an optional object to be merged into the config built by the function prior to making the request.
  // The content of the config object is described here: https://github.com/axios/axios#request-config
  // This function returns a promise which raises an error if the HTTP response is a 4xx or 5xx.
  client.showWorkflow = function (path, config) {
    var cfg = {
      timeout: timeout,
      url: urlPrefix + path,
      method: 'get',
      responseType: 'json'
    };
    if (config) {
      cfg = merge(cfg, config);
    }
    return client(cfg);
  }
  return client;
  };
});
