import axios from 'axios';

var client = axios;

let showWorkflow = function () {
  return client({
    timeout: 20000,
    url: 'http://127.0.0.1:2323/workflow/default/',
    method: 'get',
    responseType: 'document',
  });
}

export {
  showWorkflow
};