import {basePath, registerPath} from "/assets/js/config.js";

class User {

    constructor(username, email, password) {
        this.username = username;
        this.email = email;
        this.password = password;
    }
}

let notify = document.getElementById("notify")

export default function signUp() {
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    let confirmPassword = document.getElementById("confirm_password").value
    let email = document.getElementById("email").value;

    let isValidUsername = validateField(username === "", "Username is required")
    if (!isValidUsername) {
        return
    }

    let isValidEmail = validateField(email === "", "Email is required")
    if (!isValidEmail) {
        return
    }

    let isValidPassword = validateField(password === "", "Password is required")
    if (!isValidPassword) {
        return
    }

    let isValidConfirmPassword = validateField(password !== confirmPassword, "Password Does not match")
    if (!isValidConfirmPassword) {
        return
    }

    let user = new User(
        username,
        email,
        password
    )

    fetch(basePath + registerPath, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(user)
    }).then(data => {
        console.log(data.status)
        if (data.status === 201) {
            notify.style.display = 'block'
            notify.style.color = "green"
            notify.innerHTML = "Registration Successful"
            setTimeout(() => window.location.replace("/login"), 5000)
        } else {
            notify.style.display = 'block'
            notify.style.color = "red"
            data.json().then(value => notify.innerHTML = value.message)
        }
    })
}

function validateField(condition, message) {
    if (condition) {
        notify.style.display = 'block'
        notify.style.color = "red"
        notify.innerHTML = message
        setTimeout(() => notify.style.display = 'none', 5000)
        return false
    } else {
        return true
    }
}