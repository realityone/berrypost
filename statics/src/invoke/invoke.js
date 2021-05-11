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
            methodNameInput.value = m.dataset.grpcMethodName;
            methodNameInput.dataset.serviceMethod = m.dataset.serviceMethod;
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

var startRequestSentAction = function (serviceMethod, grpcMethodName) {
    cleanRequestActionBadge();
    const actionText = document.getElementById("request-action-text");
    actionText.innerText = `Processing ${serviceMethod}`;

    const actionDescription = document.getElementById("request-action-description");
    actionDescription.innerText = `Invoke ${grpcMethodName}`;
};

var cleanRequestActionBadge = function () {
    const actionBadge = document.getElementById("request-action-badge");
    if (actionBadge) {
        actionBadge.parentElement.removeChild(actionBadge);
    }
};

var onReceiveResponse = function (response, serviceMethod) {
    cleanRequestActionBadge();

    const actionText = document.getElementById("request-action-text");
    actionText.innerText = `Sent ${serviceMethod}`;

    const badgeText = `${response.status} ${response.statusText}`;
    const actionSpan = document.createElement("span");
    actionSpan.id = "request-action-badge";
    actionSpan.innerText = badgeText;
    actionSpan.classList.add("badge", "border", "float-end", "mt-1");
    if ((response.status >= 200) && (response.status < 400)) {
        actionSpan.classList.add("border-success", "text-success");
    } else if ((response.status >= 400) && (response.status < 500)) {
        actionSpan.classList.add("border-warning", "text-warning");
    } else {
        actionSpan.classList.add("border-danger", "text-danger");
    }

    actionText.appendChild(actionSpan);
};

var setupClickSend = function () {
    const sendBtn = document.getElementById("send-button");
    sendBtn.addEventListener('click', () => {
        const methodNameInput = document.getElementById("method-name");
        if (methodNameInput.value === "") {
            return;
        }
        const targetInput = document.getElementById("target-addr");

        startRequestSentAction(
            methodNameInput.dataset.serviceMethod,
            methodNameInput.value,
        );
        fetch(invokePath(methodNameInput.value), {
            method: "POST",
            body: window.requestBodyEditor.getValue(),
            headers: {
                'Content-Type': 'application/json',
                'X-Berrypost-Target': targetInput.value,
            },
        }).then((response) => {
            onReceiveResponse(response, methodNameInput.dataset.serviceMethod);
            return response.json();
        }).then((data) => {
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

    const curlCmdLine = `## ${methodNameInput.value}
curl -X "POST" "${baseURL}${path}" \\
    -H 'X-Berrypost-Target: ${targetInput.value}' \\
    -H 'Content-Type: application/json' \\
    -d $'${body}'`;
    window.previewEditor.setValue(curlCmdLine);
};

var setupPreviewTrigger = function () {
    window.requestBodyEditor.on('change', (instance, changeObj) => GeneratePreviewCmdLine());
    const targetInput = document.getElementById("target-addr");
    targetInput.addEventListener('keyup', () => GeneratePreviewCmdLine());
};

window.addEventListener('load', setupCodeMirror);
window.addEventListener('load', setupPreviewTrigger);
window.addEventListener('load', fillMethod);
window.addEventListener('load', setupClickSend);
window.addEventListener('load', serviceMenuLiveSearch);
window.addEventListener('load', clickFirstMethod);
