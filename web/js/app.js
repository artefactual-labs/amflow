'use strict';

import $ from 'jquery';
import panzoom from 'panzoom';
import * as client from './client';

function load() {
    // let c = client('http', 'localhost:2323');
    client.showWorkflow().then(function(resp) {
        let $area = $('.viewport')
        $area.empty();
        $area.append(resp.data.firstElementChild)
        panzoom($area.get(0), {
            minZoom: 0.01,
            zoomSpeed: 0.3,
            autocenter: true
        }).zoomAbs(0, 0, 0.1);

        $(".node", $(document)).on("click", function(event) {
            $(this).toggleClass("active");
        });
    });
}

function setupReloadButton() {
    document.querySelector('.reload-btn').addEventListener('click', function() {
        load();
    });
}

$(document).ready(function() {
    load();
    setupReloadButton();
});