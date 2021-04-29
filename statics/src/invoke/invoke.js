import './invoke.css';
import '../vendor.js';
import CodeMirror from 'codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript.js';
import 'codemirror/mode/shell/shell.js';

window.addEventListener('load', function () {
    var requestBodyEditor = CodeMirror.fromTextArea(requestBody, {
        lineNumbers: true,
        mode: { name: "javascript", json: true },
    });
    requestBodyEditor.setSize('100%', '300px');
    var previewEditor = CodeMirror.fromTextArea(preview, {
        lineNumbers: true,
        mode: "shell",
    });
    previewEditor.setSize('100%', '100%');
    var responseBodyEditor = CodeMirror.fromTextArea(responseBody, {
        lineNumbers: true,
        mode: { name: "javascript", json: true },
    });
    responseBodyEditor.setSize('100%', '100%');
});