import {basePath, clientId, clientSecret, loginPath} from "/assets/js/config.js";

class LoginRequest {

    constructor(username, password, clientId, clientSecret) {
        this.username = username;
        this.password = password;
        this.client_id = clientId;
        this.client_secret = clientSecret;
    }
}

export default function login() {

    console.log("base = {}",  basePath)

    let username = document.getElementById("username");
    let password = document.getElementById("password");
    let notify = document.getElementById("notify")
    let loginRequest = new LoginRequest(
        username.value,
        password.value,
        clientId,
        clientSecret
    )

    fetch("http://localhost:8080/login", {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(loginRequest)
    }).then(data => {
        if (data.status === 200) {
            notify.style.display = 'block';
            notify.style.color = "green"

            data.json().then(value => {
                notify.style.display = 'block'
                notify.style.color = "green"
                notify.innerText = value.message
                let accessToken = value.payload.access_token

                fetch("/", {
                    method: "GET",
                    mode: "cors",
                    headers: {
                        "Authorization": accessToken
                    }
                }).then(v => {
                    if (v.status === 200) {
                        v.text().then(text => {
                            setTimeout(() => document.write(text), 1000);
                        })
                    }
                })
            })


        } else {
            notify.style.display = 'block';
            notify.style.color = "red"
        }

        data.json().then(value => {
            notify.innerText = value.message
        })
        setTimeout(() => notify.style.display = 'none', 5000)
    })
}