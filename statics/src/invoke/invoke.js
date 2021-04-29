import './invoke.css';
import '../vendor.js';
import CodeMirror from 'codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript.js';
import 'codemirror/mode/shell/shell.js';

var setupCodeMirror = function () {
    window.requestBodyEditor = CodeMirror.fromTextArea(requestBody, {
        lineNumbers: true,
        mode: { name: "javascript", json: true },
    });
    window.requestBodyEditor.setSize('100%', '300px');
    window.previewEditor = CodeMirror.fromTextArea(preview, {
        lineNumbers: true,
        mode: "shell",
    });
    window.previewEditor.setSize('100%', '100%');
    window.responseBodyEditor = CodeMirror.fromTextArea(responseBody, {
        lineNumbers: true,
        mode: { name: "javascript", json: true },
    });
    window.responseBodyEditor.setSize('100%', '100%');
};

var fillMethod = function () {
    var methods = document.getElementsByClassName("service-method");
    for (const m of methods) {
        m.onclick = function () {
            var methodNameInput = document.getElementById("method-name");
            methodNameInput.value = m.dataset.methodName;
            window.requestBodyEditor.setValue(m.dataset.inputSchema);
        }
    }
};

window.addEventListener('load', setupCodeMirror);
window.addEventListener('load', fillMethod);