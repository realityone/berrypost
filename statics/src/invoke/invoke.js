import './invoke.css';
import '../vendor.js';
import CodeMirror from 'codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript.js';
import 'codemirror/mode/shell/shell.js';
import path from 'path-browserify';

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
    const methods = document.getElementsByClassName("service-method");
    for (const m of methods) {
        m.onclick = function () {
            const methodNameInput = document.getElementById("method-name");
            methodNameInput.value = m.dataset.methodName;
            window.requestBodyEditor.setValue(m.dataset.inputSchema);
        }
    }
};

var clickFirstMethod = function () {
    const methods = document.getElementsByClassName("service-method");
    for (const m of methods) {
        m.click();
        return
    }
};

var invokeURL = function (methodName) {
    return path.join("/invoke", methodName)
};

var setupClickSend = function () {
    const sendBtn = document.getElementById("send-button");
    sendBtn.onclick = function () {
        const methodNameInput = document.getElementById("method-name");
        fetch(invokeURL(methodNameInput.value), {
            method: "POST",
            body: window.requestBodyEditor.getValue(),
        }).then((response) => response.json())
            .then((data) => {
                const prettyJSON = JSON.stringify(data, null, 2);
                window.responseBodyEditor.setValue(prettyJSON);
            }).catch((error) => {
                const errorMessage = error.toString();
                window.responseBodyEditor.setValue(errorMessage);
            });
    };
};

window.addEventListener('load', setupCodeMirror);
window.addEventListener('load', fillMethod);
window.addEventListener('load', clickFirstMethod);
window.addEventListener('load', setupClickSend);