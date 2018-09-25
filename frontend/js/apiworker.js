function getinfo() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/users/info',
        success: function(data) {
            //localStorage.setItem("jwtToken", data.token);
            console.log(data);

            document.getElementById("username").innerText = data.name;
            document.getElementById("command").innerText = data.command;
        }
    });
}


window.onload = getinfo;