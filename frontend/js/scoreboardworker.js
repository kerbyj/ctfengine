var scoreboardContainer = document.getElementById("scoreboardTable");

function getTopForAllTime() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/users/topForAllTime',
        success: function(data) {
            console.log(data);

            $.each(data, function (key, item) {
                let userRow = document.createElement("tr");
                userRow.innerHTML = `<td>${item.username}</td><td>${item.command}</td><td>${item.points}</td>`;

                scoreboardContainer.appendChild(userRow);
            })
        }
    });
}
window.onload = getTopForAllTime();