import './style.css'
import panzoom from 'panzoom';
import { saveAs } from 'file-saver';
import createClient from './client.es6.js'
import getLocation from './location.js'

const { scheme, host } = getLocation(window);
const client = createClient(scheme, host);

const parse = (blob) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(blob, "text/xml");
  return doc.documentElement;
}

const loadDocument = (config) => {
  return client.showWorkflow("/workflow/default", config);
}

let panzoomInstance;

const loadGraph = () => {
  loadDocument().then((resp) => {
    try {
      const viewport = document.querySelector("#viewport");
      viewport.innerHTML = "";
      viewport.appendChild(parse(resp.data));
      panzoomInstance = panzoom(viewport, { minZoom: 0.01, zoomSpeed: 0.3, autocenter: false });
      panzoomInstance.zoomAbs(0, 0, 0.1);
    } catch (error) {
      console.error(error);
      alert("An error occurred while loading the graph.");
    }
  }).catch((error) => {
    console.error(error);
    alert("An error occurred while retrieving the data.");
  })
};

const saveGraph = () => {
  loadDocument({ responseType: "blob" }).then((resp) => {
    saveAs(resp.data, "workflow.svg", {type: "image/svg+xml"});
  }).catch((error) => {
    console.error(error);
    alert("An error occurred while retrieving the data.");
  })
}

const search = (text) => {
  const svg = document.querySelector('svg');
  const textElements = svg.querySelectorAll('*');
  const matchingNodesSet = new Set();
  const lctext = text.toLowerCase();

  for (let element of textElements) {
    if (element.textContent.toLowerCase().includes(lctext)) {
      let parent = element.parentElement;
      while (parent) {
        if (parent.classList.contains('node')) {
          matchingNodesSet.add(parent);
          break;
        }
        parent = parent.parentElement;
      }
    }
  }

  return Array.from(matchingNodesSet);
}

const zoomElement = (element) => {
  const svg = document.querySelector('svg');
  const rect = element.getBBox();
  const point = svg.createSVGPoint();
  point.x = rect.x + rect.width / 2;
  point.y = rect.y + rect.height / 2;
  const globalPoint = point.matrixTransform(element.getScreenCTM());
  const viewportWidth = window.innerWidth || document.documentElement.clientWidth;
  const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
  const panX = (viewportWidth / 2 - globalPoint.x);
  const panY = (viewportHeight / 2 - globalPoint.y);
  panzoomInstance.moveBy(panX, panY, false);
}

// Write the initial layout of the application.
document.querySelector('#app').innerHTML = `
  <div class="wrapper">
    <div class="box sidebar">
      <button class="reload">Reload</button>
      <button class="save">Save</button>
      <input class="search" type="text" placeholder="Search (shift + /)"/>
    </div>
    <div class="box content">
      <div id="viewport"></div>
    </div>
  </div>
`

// Focus the search box when shift and `/` keys are pressed simultaneously.
window.addEventListener('keydown', (event) => {
  if (event.shiftKey && event.key === '/') {
    event.preventDefault();
    const search = document.querySelector('.search');
    if (search) {
      search.focus();
    }
  }
});

let results;

document.querySelector(".search").addEventListener("keydown", function(event) {
  if (event.key == "Escape") {
    event.preventDefault();
    this.value = "";
  } else if (event.key == "Enter") {
    event.preventDefault();
    if (this.value.trim() != "") {
      results = search(this.value);
      if (results.length > 0) {
        zoomElement(results[0]);
        results.forEach((item) => console.log(this.value, item));
      }
    }
  }
})

// Highlight graph nodes when they are clicked by the user.
document.body.addEventListener('click', function(event) {
  const el = event.target.closest('.node');
  if (el) {
    el.classList.toggle('active');
  }
});

document.querySelector('.reload').addEventListener('click', loadGraph);
document.querySelector('.save').addEventListener('click', saveGraph);

loadGraph();
