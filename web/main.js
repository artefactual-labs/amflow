import './style.css'
import panzoom from 'panzoom';
import createClient from './client.es6.js'
import getLocation from './location.js'

const { scheme, host } = getLocation(window);
const client = createClient(scheme, host);

const parse = function(blob) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(blob, 'text/xml');
  return doc.documentElement;
}

const loadGraph = function() {
  client.showWorkflow("/workflow/default").then((resp) => {
    try {
      const viewport = document.querySelector('#viewport');
      viewport.innerHTML = '';
      viewport.appendChild(parse(resp.data));
      panzoom(viewport, { minZoom: 0.01, zoomSpeed: 0.3, autocenter: true }).zoomAbs(0, 0, 0.1);
    } catch (error) {
      console.error(error);
      alert('An error occurred while loading the graph.');
    }
  }).catch((error) => {
    console.error(error);
    alert('An error occurred while retrieving the data.');
  })
};

document.querySelector('#app').innerHTML = `
  <div class="wrapper">
    <div class="box sidebar">
      <button class="reload">Reload</button>
    </div>
    <div class="box content">
      <div id="viewport"></div>
    </div>
  </div>
`

document.querySelector('.reload').addEventListener('click', loadGraph);

document.body.addEventListener('click', function(event) {
  const el = event.target.closest('.node');
  if (el) {
    el.classList.toggle('active');
  }
});

loadGraph();
