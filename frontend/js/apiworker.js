function getinfo() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/users/info',
        success: function(data) {
            //localStorage.setItem("jwtToken", data.token);
            console.log(data);

            document.getElementById("username").innerText = data.name;
            delete data.name;
            document.getElementById("command").innerText = data.command;
            delete data.command;

            var statsUserContainer = document.getElementById("userStats");

            $.each(data, function (key, value) {
                let tmpStat = document.createElement("div");
                tmpStat.classList.add("statsElements");
                tmpStat.innerHTML = `${key}<span class="valueInfo">${value}</span>`;
                statsUserContainer.appendChild(tmpStat);
            })
        }
    });
}
window.onload = getinfo;