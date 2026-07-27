import './style.css'
import panzoom from 'panzoom';
import { saveAs } from 'file-saver';
import createAPI from './api.js'
import getLocation from './location.js'

const { scheme, host } = getLocation(window);
const api = createAPI(scheme, host);

const parse = (blob) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(blob, "text/xml");
  return doc.documentElement;
}

const loadDocument = (options) => {
  return api.showWorkflow("default", options);
}

let panzoomInstance;

const loadGraph = () => {
  loadDocument().then((workflowDocument) => {
    try {
      const viewport = document.querySelector("#viewport");
      viewport.innerHTML = "";
      viewport.appendChild(parse(workflowDocument));
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
  loadDocument({ responseType: "blob" }).then((workflowDocument) => {
    saveAs(workflowDocument, "workflow.svg", {type: "image/svg+xml"});
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
    <header class="toolbar">
      <div class="brand" aria-label="amflow">
        <span class="brand-mark" aria-hidden="true">a</span>
        <span class="brand-name">amflow</span>
      </div>
      <label class="search-field">
        <span class="search-icon" aria-hidden="true"></span>
        <span class="visually-hidden">Search workflow</span>
        <input
          class="search"
          type="search"
          placeholder="Search workflow"
          autocomplete="off"
          spellcheck="false"
        >
        <kbd class="search-shortcut" aria-hidden="true">⇧ /</kbd>
      </label>
      <div class="toolbar-actions">
        <button class="reload" type="button">Reload</button>
        <button class="save" type="button">Save SVG</button>
      </div>
    </header>
    <main class="content">
      <div id="viewport"></div>
    </main>
  </div>
`

// Focus the search box when shift and `/` keys are pressed simultaneously.
window.addEventListener('keydown', (event) => {
  if (event.shiftKey && event.code === 'Slash') {
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
