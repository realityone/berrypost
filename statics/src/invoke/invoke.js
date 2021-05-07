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
    window.requestBodyEditor.on('keyup', () => {
        GeneratePreviewCmdLine();
    });
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

var invokePath = function (methodName) {
    return path.join("/invoke", methodName)
};

var updateStatusBar = function (response) {

};

var setupClickSend = function () {
    const sendBtn = document.getElementById("send-button");
    sendBtn.addEventListener('click', () => {
        const methodNameInput = document.getElementById("method-name");
        if (methodNameInput.value === "") {
            return
        }
        fetch(invokePath(methodNameInput.value), {
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
    })
};

var serviceMenuLiveSearch = function () {
    const searchInput = document.getElementById("serviceMenuSearch");
    searchInput.addEventListener('keyup', () => {
        const filter = searchInput.value;
        const serviceMenuContent = document.getElementById("serviceMenuContent");
        const links = serviceMenuContent.getElementsByTagName("a");
        for (const a of links) {
            const txtValue = a.textContent || a.innerText;
            if (txtValue.indexOf(filter) > -1) {
                a.style.display = "";
                continue;
            }
            a.style.display = "none";
        }
    });
};

var GeneratePreviewCmdLine = function () {
    const methodNameInput = document.getElementById("method-name");
    if (methodNameInput.value === "") {
        return;
    }
    const baseURL = `${location.protocol}//${location.host}`;
    const path = invokePath(methodNameInput.value);
    const body = window.requestBodyEditor.getValue();
    const targetInput = document.getElementById("target-addr");

    const curlCmdLine = `## ${path}
curl -X "POST" "${baseURL}${path}" \\
    -H 'X-Berrypost-Target: ${targetInput.value}' \\
    -H 'Content-Type: text/plain; charset=utf-8' \\
    -d $'${body}'`;
    window.previewEditor.setValue(curlCmdLine);
};

window.addEventListener('load', setupCodeMirror);
window.addEventListener('load', fillMethod);
window.addEventListener('load', clickFirstMethod);
window.addEventListener('load', setupClickSend);
window.addEventListener('load', serviceMenuLiveSearch);
