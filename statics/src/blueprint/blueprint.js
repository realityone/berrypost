import '../vendor.js';
import '../invoke/invoke.css';
import { Modal } from 'bootstrap';

let newMethodModal = function(){
    const modelButton = document.getElementById("method-modal");
    modelButton.onclick = function() {
        const myModal = new Modal(document.getElementById('newMethodModal'))
        myModal.show()
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
        const blueprintName = document.getElementById("blueprint-name").value;
        fetch("/management/api/blueprint/new", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
            }),
        }).then((response) => {
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
        const blueprintName = document.getElementById("blueprint-list").value;
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

let deleteMethodModal = function(){
    const modelButton = document.getElementById("delete-method-modal");
    modelButton.onclick = function() {
        // todo check user token
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
        // todo check user token
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

window.addEventListener('load', newMethodModal);
window.addEventListener('load', newBlueprintModal);
window.addEventListener('load', newBlueprintReq);
window.addEventListener('load', savetoBlueprintReq);
window.addEventListener('load', savetoBlueprintModal);
window.addEventListener('load', deleteMethodModal);
window.addEventListener('load', deleteMethodReq);
window.addEventListener('load', deleteBlueprintModal);
window.addEventListener('load', deleteBlueprintReq);

