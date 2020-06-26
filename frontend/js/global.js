var GLOBAL_ENDPOINT = "http://"+location.host;


function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

let CurrentCsrfToken = getCookie("_csrf");