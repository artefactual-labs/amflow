import axios from 'axios';

var client = axios;

let showWorkflow = function () {
  return client({
    timeout: 20000,
    url: '/workflow/default/',
    method: 'get',
    responseType: 'document',
  });
}

export {
  showWorkflow
};
