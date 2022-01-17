import '../vendor.js';
import {Modal} from 'bootstrap';

let copyBlueprintModal = function(){
    const modelButton = document.getElementById("copy-blueprint-modal");
    modelButton.onclick = function() {
        document.getElementById("copy-service-name").value = document.getElementById("serviceMenu").innerText;
        const myModal = new Modal(document.getElementById('copyBlueprintModal'))
        myModal.show()
    }
}

let copyBlueprintReq = function(){
    const reqButton = document.getElementById("copy-blueprint");
    reqButton.onclick = function() {
        let form = document.getElementById("form-copy-blueprint")
        if (!form.checkValidity()) {
            form.classList.add('was-validated')
            return
        }
        const blueprintName = document.getElementById("copy-blueprint-name").value;
        const FileName = document.getElementById("serviceMenu").innerText;
        fetch("/management/api/blueprint/file/copy", {
            method: "POST",
            body: JSON.stringify({
                'blueprintName' : blueprintName,
                'FileName' : FileName,
            }),
        }).then((response) => {
            if (response.status === 200){
                alert("copy successfully!")
                document.location.reload();
            } else {
                alert("fail to copy")
            }
        })
    }
}

let savetoBlueprintModal = function(){
    const modelButton = document.getElementById("saveto-blueprint-modal");
    modelButton.onclick = function() {
        document.getElementById("saveto-method-name").value = document.getElementById("method-name").value;
        const myModal = new Modal(document.getElementById('savetoBlueprintModal'));
        myModal.show();
    }
}

let savetoBlueprintReq = function(){
    const reqButton = document.getElementById("saveto-blueprint");
    reqButton.onclick = function() {
        let form = document.getElementById("form-saveto-blueprint")
        if (!form.checkValidity()) {
            form.classList.add('was-validated')
            return
        }
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
            if (response.status === 200){
                alert("save successfully!")
                document.location.reload();
            } else {
                alert("fail to save")
            }
        })
    }
}

window.addEventListener('load', copyBlueprintReq);
window.addEventListener('load', copyBlueprintModal);
window.addEventListener('load', savetoBlueprintReq);
window.addEventListener('load', savetoBlueprintModal);
