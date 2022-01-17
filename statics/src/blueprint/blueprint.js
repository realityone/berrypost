import '../vendor.js';
import { Modal } from 'bootstrap';

let initSel = function() {
    const fileSel = document.getElementById("file-name");
    fileSel.onchange = function() {
        fetch("/management/api/service/methods", {
            method: "POST",
            body: JSON.stringify({
                'fileName' : fileSel.value,
            }),
        }).then((json) => {
            return json.json()
        }).then((data) => {
            initMethod(data);
        })
    }
}

let initMethod = function(methods){
    if (methods === null) {
        return
    }
    const methodSel = document.getElementById("new-method-name");
    const methodOB = $('#new-method-name');
    methodOB.empty();
    for(let j = 0; j < methods.length; j++) {
        methodSel.appendChild(new Option(methods[j], methods[j]));
    }
    methodOB.selectpicker('refresh');
    methodOB.selectpicker('render');
}

let newMethodModal = function(){
    const modelButton = document.getElementById("method-modal");
    modelButton.onclick = function() {
        document.getElementById("new-blueprint-name").value = document.getElementById("serviceMenu").innerText;
        const myModal = new Modal(document.getElementById('newMethodModal'))
        myModal.show()
    }
}

let newMethodReq = function(){
    const reqButton = document.getElementById("add-method");
    reqButton.onclick = function() {
        let form = document.getElementById("form-new-method")
        if (!form.checkValidity()) {
            form.classList.add('was-validated')
            return
        }
        const blueprintName = document.getElementById("serviceMenu").innerText;
        const filename = document.getElementById("file-name").value;
        const method = $("#new-method-name").val();
        fetch("/management/api/blueprint/append-list", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
                'filename' : filename,
                'methodName' : method,
            }),
        }).then((response) => {
            alert("new methods successfully!")
            document.location.reload();
        })
    }
}

let newBlueprintModal = function(){
    const modelButton = document.getElementById("blueprint-modal");
    modelButton.onclick = function() {
        const myModal = new Modal(document.getElementById('newBlueprintModal'))
        myModal.show()
    }
}

let newBlueprintReq = function(){
    const reqButton = document.getElementById("add-blueprint");
    reqButton.onclick = function() {
        let form = document.getElementById("form-new-blueprint")
        if (!form.checkValidity()) {
            form.classList.add('was-validated')
            return
        }
        const blueprintName = document.getElementById("blueprint-name").value;
        fetch("/management/api/blueprint/new", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
            }),
        }).then((response) => {
            if (response.status === 200){
                alert("new blueprint successfully!")
                document.location.replace("/management/blueprint/"+blueprintName);
            } else {
                alert("fail to new blueprint")
            }
        })
    }
}

let deleteMethodModal = function(){
    const modelButton = document.getElementById("delete-method-modal");
    modelButton.onclick = function() {
        document.getElementById("deleteMethodModal-blueprint").value = document.getElementById("serviceMenu").innerText;
        document.getElementById("deleteMethodModal-method").value = document.getElementById("method-name").value;
        const myModal = new Modal(document.getElementById('deleteMethodModal'))
        myModal.show()
    }
}

let deleteMethodReq = function(){
    const reqButton = document.getElementById("delete-method");
    reqButton.onclick = function() {
        const blueprintName = document.getElementById("serviceMenu").innerText;
        const methodNameInput = document.getElementById("method-name");
        fetch("/management/api/blueprint/reduce", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
                'fileName': methodNameInput.dataset.serviceFileName,
                'methodName':methodNameInput.value,
            }),
        }).then((response) => {
            alert("delete successfully!")
            document.location.reload();
        })
    }
}

let deleteBlueprintModal = function(){
    const modelButton = document.getElementById("delete-blueprint-modal");
    modelButton.onclick = function() {
        document.getElementById("deleteBlueprintModal-blueprint").value = document.getElementById("serviceMenu").innerText;
        const myModal = new Modal(document.getElementById('deleteBlueprintModal'))
        myModal.show()
    }
}

let deleteBlueprintReq = function(){
    const reqButton = document.getElementById("delete-blueprint");
    reqButton.onclick = function() {
        const blueprintName = document.getElementById("serviceMenu").innerText;
        fetch("/management/api/blueprint/delete", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
            }),
        }).then((response) => {
            alert("delete successfully!")
            document.location.replace("/management/blueprint");
        })
    }
}

let shareBlueprint = function(){
    const reqButton = document.getElementById("share-blueprint");
    reqButton.onclick = function() {
        const blueprintName = document.getElementById("serviceMenu").innerText;
        fetch("/management/api/blueprint/share", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
            }),
        }).then((json) => {
            return json.json()
        }).then((data) => {
            let protocol = window.location.protocol;
            let host = window.location.host;
            const url = document.getElementById("url");
            url.value = protocol + "//" + host + data;
            const reqButton = document.getElementById("copy-url");
            reqButton.innerText = "Copy"
            const myModal = new Modal(document.getElementById('shareBlueprintModal'));
            myModal.show();
        })
    }
}

let copyUrl = function(){
    const reqButton = document.getElementById("copy-url");
    reqButton.onclick = function() {
        $('#url').select();
        document.execCommand('Copy');
        reqButton.innerText = "Copied"
    }
}

window.addEventListener('load', initSel);
window.addEventListener('load', newMethodModal);
window.addEventListener('load', newMethodReq);
window.addEventListener('load', newBlueprintModal);
window.addEventListener('load', newBlueprintReq);
window.addEventListener('load', deleteMethodModal);
window.addEventListener('load', deleteMethodReq);
window.addEventListener('load', deleteBlueprintModal);
window.addEventListener('load', deleteBlueprintReq);
window.addEventListener('load', shareBlueprint);
window.addEventListener('load', copyUrl);

