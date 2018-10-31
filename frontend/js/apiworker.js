function getinfo() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/user/info',
        success: function(data) {
            console.log(data);

            document.getElementById("name").innerText = data.name;
            delete data.name;
            document.getElementById("commandstatus").innerText = data.command;
            delete data.command;

            var statsUserContainer = document.getElementById("userStats");

            $.each(data, function (key, value) {
                let tmpStat = document.createElement("div");
                tmpStat.classList.add("statsElements");
                tmpStat.innerHTML = `${key} ${value}`;

                statsUserContainer.appendChild(tmpStat);
                statsUserContainer.innerHTML += `<hr width="30%" align="left" style="margin-left:15px; color: #00dcff">`;

            })
        }
    });

    $.post({
        type: 'get',
        url: 'http://localhost/api/board/getstats',
        success: function(data) {
            console.log(data);

            $.each(data, function (key, value) {
                let tmpContainer = document.getElementById(key);

                $.each(value, function (statKey, statValue) {
                    tmpContainer.innerHTML+=`<span class ="stat" style="margin-right: 30px;">${statKey} <b>${statValue}</b></span>`
                })
            })
        }
    });
}
window.onload = getinfo;