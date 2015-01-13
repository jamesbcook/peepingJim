// Took from orig peepingtom
// https://bitbucket.org/LaNMaSteR53/peepingtom/src/a6bf47438d25f4e5578bd9253dd03a951cafe488/capture.js?at=master
var page = require('webpage').create();

var url = phantom.args[0];
var filename = phantom.args[1];
page.viewportSize = { width: 1024, height: 768 };
page.clipRect = { top: 0, left: 0, width: 1024, height: 768 };

var callback = function(status) {
    if (status !== 'success') {
        console.log('Unable to load the address: ' + url);
        phantom.exit(1)
    }
    page.render(filename);
    console.log('Successfully rendered: ' + url);
    phantom.exit(0);
};
// set connection timeout below
var timer = window.setTimeout(callback, parseInt(phantom.args[2]), 'timeout');
page.open(url, callback);
