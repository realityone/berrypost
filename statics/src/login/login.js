import './login.css';
import '../vendor.js';

let signInReq = function(){
    const reqButton = document.getElementById("sign-in");
    reqButton.onclick = function() {
        const userid = document.getElementById("userid").value;
        const password = document.getElementById("password").value;
        fetch("/management/api/signIn", {
            method: "POST",
            body: JSON.stringify({
                'userid' : userid,
                'password' : password,
            }),
            keepalive: true,
        }).then((json) => {
            return json.json()
        }).then((data) => {
            if (data === true) {
                alert("sign in successfully!")
                document.location.replace("/management/invoke");
            } else {
                alert("wrong userid/password!")
            }
        })
    }
}

let signUpReq = function(){
    const reqButton = document.getElementById("sign-up");
    reqButton.onclick = function() {
        const userid = document.getElementById("userid").value;
        const password = document.getElementById("password").value;
        fetch("/management/api/signUp", {
            method: "POST",
            body: JSON.stringify({
                'userid' : userid,
                'password' : password,
            }),
            keepalive: true,
        }).then((json) => {
            return json.json()
        }).then((data) => {
            if (data === true) {
                alert("sign up successfully!")
                document.location.replace("/management/login");
            } else {
                alert("userid already exists, please retry")
            }
        })
    }
}

window.addEventListener('load', signInReq);
window.addEventListener('load', signUpReq);
