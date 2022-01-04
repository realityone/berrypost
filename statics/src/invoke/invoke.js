import './invoke.css';
import '../vendor.js';
import CodeMirror from 'codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript.js';
import 'codemirror/mode/shell/shell.js';
import path from 'path-browserify';
import {Modal} from "bootstrap";

var setupCodeMirror = function() {
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

var fillMethod = function() {
    const methods = document.getElementsByClassName("service-method");
    for (const m of methods) {
        m.onclick = function() {
            const methodNameInput = document.getElementById("method-name");
            methodNameInput.value = m.dataset.grpcMethodName;
            const targetAddressInput = document.getElementById("target-addr");
            targetAddressInput.value = m.dataset.preferTarget;
            methodNameInput.dataset.serviceMethod = m.dataset.serviceMethod;
            methodNameInput.dataset.serviceFileName = m.dataset.serviceFileName;
            window.requestBodyEditor.setValue(m.dataset.inputSchema);
        }
    }
};

var clickFirstMethod = function() {
    const methods = document.getElementsByClassName("service-method");
    for (const m of methods) {
        m.click();
        return
    }
};

var invokePath = function(methodName) {
    return path.join("/invoke", methodName)
};

var startRequestSentAction = function(serviceMethod, grpcMethodName) {
    cleanRequestActionBadge();
    const actionText = document.getElementById("request-action-text");
    actionText.innerText = `Processing ${serviceMethod}`;

    const actionDescription = document.getElementById("request-action-description");
    actionDescription.innerText = `Invoke ${grpcMethodName}`;
};

var cleanRequestActionBadge = function() {
    const actionBadge = document.getElementById("request-action-badge");
    if (actionBadge) {
        actionBadge.parentElement.removeChild(actionBadge);
    }
};

var onReceiveResponse = function(response, serviceMethod) {
    cleanRequestActionBadge();

    const actionText = document.getElementById("request-action-text");
    actionText.innerText = `Sent ${serviceMethod}`;

    var badgeText = `${response.status}`;
    if ((response.status >= 200) && (response.status < 300)) {
        badgeText = `${response.status} ${response.statusText}`
    }
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

var setupClickSend = function() {
    const sendBtn = document.getElementById("send-button");
    sendBtn.addEventListener('click', () => {
        const methodNameInput = document.getElementById("method-name");
        if (methodNameInput.value === "") {
            return;
        }
        const targetInput = document.getElementById("target-addr");
        const metadataTable = document.getElementById("metadata-table");

        const headers = {
            'Content-Type': 'application/json',
            'X-Berrypost-Target': targetInput.value,
        };
        for (const row of Array.from(metadataTable.rows).slice(1)) {
            const inputs = row.getElementsByTagName("input");
            const enabledCheckbox = inputs[0];
            const nameInput = inputs[1];
            const valueInput = inputs[2];
            if (!enabledCheckbox.checked) {
                continue;
            }
            headers[`X-Berrypost-Md-${nameInput.value}`] = valueInput.value;
        }

        startRequestSentAction(
            methodNameInput.dataset.serviceMethod,
            methodNameInput.value,
        );
        fetch(invokePath(methodNameInput.value), {
            method: "POST",
            body: window.requestBodyEditor.getValue(),
            headers: headers,
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

var serviceMenuLiveSearch = function() {
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

var metadataHeaderArgs = function(metadataTable) {
    var args = [""];
    for (const row of Array.from(metadataTable.rows).slice(1)) {
        const inputs = row.getElementsByTagName("input");
        const enabledCheckbox = inputs[0];
        const nameInput = inputs[1];
        const valueInput = inputs[2];
        if (!enabledCheckbox.checked) {
            continue;
        }
        args.push(`    -H 'X-Berrypost-Md-${nameInput.value}: ${valueInput.value}' \\`)
    }
    return args.join("\n");
}

var GeneratePreviewCmdLine = function() {
    const methodNameInput = document.getElementById("method-name");
    if (methodNameInput.value === "") {
        return;
    }
    const baseURL = `${location.protocol}//${location.host}`;
    const path = invokePath(methodNameInput.value);
    const body = window.requestBodyEditor.getValue();
    const targetInput = document.getElementById("target-addr");
    const metadataTable = document.getElementById("metadata-table");
    const metadataHeaderString = metadataHeaderArgs(metadataTable);

    const curlCmdLine = `## ${methodNameInput.value}
curl -X "POST" "${baseURL}${path}" \\
    -H 'X-Berrypost-Target: ${targetInput.value}' \\
    -H 'Content-Type: application/json' \\${metadataHeaderString}
    -d $'${body}'`;
    window.previewEditor.setValue(curlCmdLine);
};

var setupPreviewTrigger = function() {
    window.requestBodyEditor.on('change', (instance, changeObj) => GeneratePreviewCmdLine());
    const targetInput = document.getElementById("target-addr");
    targetInput.addEventListener('keyup', () => GeneratePreviewCmdLine());
};

var setupRequestEditor = function() {
    const navList = document.getElementById("request-editor-nav");
    const links = navList.getElementsByTagName("a");
    for (const a of links) {
        a.addEventListener('click', function(e) {
            e.preventDefault();
        });
    }

    const metadataTable = document.getElementById("metadata-table");
    metadataTable.addEventListener('input', (e) => {
        if (e.target.name !== "name") {
            return
        }
        const row = e.target.parentElement.parentElement;
        if (row.rowIndex === (metadataTable.rows.length - 1)) {
            const newRow = metadataTable.insertRow();
            newRow.innerHTML = row.innerHTML;
            const enableCheckbox = newRow.cells[0].getElementsByTagName("input")[0];
            enableCheckbox.setAttribute("hidden", "");
            enableCheckbox.removeAttribute("checked");
        }
        const enableCheckbox = row.cells[0].getElementsByTagName("input")[0];
        enableCheckbox.setAttribute("checked", "");
        enableCheckbox.removeAttribute("hidden");
    });
    metadataTable.addEventListener('input', () => GeneratePreviewCmdLine());
    metadataTable.addEventListener('click', (e) => {
        if (e.target.parentElement.classList.contains("delete-metadata")) {
            const row = e.target.parentElement.parentElement.parentElement;
            if (row.rowIndex !== 1) {
                metadataTable.deleteRow(row.rowIndex);
            }
        }
    });
};

let copyBlueprintModal = function(){
    const modelButton = document.getElementById("copy-blueprint-modal");
    modelButton.onclick = function() {
        // todo check user token
        const myModal = new Modal(document.getElementById('copyBlueprintModal'))
        myModal.show()
    }
}


let copyBlueprintReq = function(){
    const reqButton = document.getElementById("copy-blueprint");
    reqButton.onclick = function() {
        const blueprintName = document.getElementById("copy-blueprint-name").value;
        const FileName = document.getElementById("serviceMenu").innerText;
        fetch("/management/api/blueprint/copy", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
                'FileName' : FileName,
            }),
        }).then((response) => {
            alert("copy successfully!")
            document.location.reload();
        })
    }
}


let savetoBlueprintModal = function(){
    const modelButton = document.getElementById("saveto-blueprint-modal");
    modelButton.onclick = function() {
        // todo check user token
        const myModal = new Modal(document.getElementById('savetoBlueprintModal'))
        myModal.show()
    }
}

let savetoBlueprintReq = function(){
    const reqButton = document.getElementById("saveto-blueprint");
    reqButton.onclick = function() {
        const blueprintName = document.getElementById("saveto-blueprint-name").value;
        const filename = document.getElementById("serviceMenu").innerText;
        const methodNameInput = document.getElementById("method-name").value;
        fetch("/management/api/blueprint/append", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
                'fileName':filename,
                'methodName':methodNameInput,
            }),
        }).then((response) => {
            alert("save successfully!")
            document.location.reload();
        })
    }
}

window.addEventListener('load', setupCodeMirror);
window.addEventListener('load', setupPreviewTrigger);
window.addEventListener('load', fillMethod);
window.addEventListener('load', setupClickSend);
window.addEventListener('load', serviceMenuLiveSearch);
window.addEventListener('load', clickFirstMethod);
window.addEventListener('load', setupRequestEditor);
window.addEventListener('load', copyBlueprintReq);
window.addEventListener('load', copyBlueprintModal);
window.addEventListener('load', savetoBlueprintReq);
window.addEventListener('load', savetoBlueprintModal);
