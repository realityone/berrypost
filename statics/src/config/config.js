import '../vendor.js';
import '../invoke/invoke.css';
import { Modal } from 'bootstrap'

let UpdateDB = function(){
    const updateButton = document.getElementById("update-button");
    updateButton.onclick = function() {
        const targetAddrInput = document.getElementById("target-addr").value;
        const serviceInput = document.getElementById("serviceMenu").innerText;
        fetch("/management/api/address/update", {
            method: "POST",
            body: JSON.stringify({
                'targetAddrInput' : targetAddrInput,
                'serviceInput': serviceInput,
            }),
        }).then((response) => {
            if (response.status === 200){
                const myModal = new Modal(document.getElementById('successModal'))
                myModal.show()
            } else {
                const myModal = new Modal(document.getElementById('failModal'))
                myModal.show()
            }
        })
    }
}


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

window.addEventListener('load', serviceMenuLiveSearch);
window.addEventListener('load', UpdateDB);

