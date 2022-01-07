import '../vendor.js';
import { Modal } from 'bootstrap';

let copyBlueprintModal = function(){
    const modelButton = document.getElementById("copy-blueprint-modal");
    modelButton.onclick = function() {
        const myModal = new Modal(document.getElementById('copyBlueprintModal'))
        myModal.show()
    }
}

let copyBlueprintReq = function(){
    const newName = document.getElementById("copy-blueprint-name")
    newName.value = document.getElementById("from-blueprint-name").innerText
    const reqButton = document.getElementById("copy-blueprint");
    reqButton.onclick = function() {
        fetch("/management/api/blueprint/copy", {
            method: "POST",
            body: JSON.stringify({
                'newName' : newName.value,
                'token' : getQueryVariable("token"),
            }),
        }).then((json) => {
            return json.json()
        }).then((data) => {
            if (data === true) {
                alert("copy successfully!")
                document.location.reload();
            } else {
                alert("blueprint name already exists, please retry")
            }
        })
    }
}


function getQueryVariable(variable)
{
    let query = window.location.search.substring(1);
    let vars = query.split("&");
    for (let i = 0; i < vars.length; i++) {
        let pair = vars[i].split("=");
        if(pair[0] === variable){
            return pair[1];
        }
    }
    return false;
}

window.addEventListener('load', copyBlueprintReq);
window.addEventListener('load', copyBlueprintModal);
