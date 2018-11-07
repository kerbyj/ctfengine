var scoreboardContainer = document.getElementById("scoreboardTable");

let latestRequest = 0;
setInterval(function () {
    if(latestRequest === 0)
        return;

    getTopForAllTime(latestRequest);
}, 30000);

function getTopForAllTime(contestId) {
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT+'/api/users/getTopForContest/'+contestId,
        success: function(data) {
            console.log(data);
            latestRequest = contestId;
            let scoreboard = document.getElementById("scoreboard");
            //scoreboard.innerHTML = "";
            scoreboard.innerHTML = `<div class="oneRow">
                        <div class="oneCell rank"></div>
                        <div class="oneCell username">NAME</div>
                        <div class="oneCell score">SCORE</div>
                        <div class="oneCell solvedTasks">SOLVED TASKS</div>
                    </div>`;
            $.each(data, function (key, item) {
                console.log(key, item);
                let topParticipantContainer = document.createElement("div");
                topParticipantContainer.classList.add("oneRow");

                topParticipantContainer.innerHTML = `<div class="oneCell rank">${item.place}</div><div class="oneCell username">${item.name}</div><div class="oneCell score">${item.points}</div><div class="oneCell solvedTasks">${item.solved}/${item.all_tasks_count}</div>`;
                scoreboard.appendChild(topParticipantContainer);
            })
        }
    });
}

function selectContest(){
    //console.log(this.targetContest)
    getTopForAllTime(this.targetContest);
}

function drawContests(){
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT+'/api/tasks/getContestList',
        success: function(data) {
            console.log(data);
            let contests = document.getElementById("contests");
            $.each(data, function (key, item) {
                console.log(key, item);
                let contestContainer = document.createElement("div");
                contestContainer.classList.add("oneRow");
                contestContainer.classList.add("selected");

                contestContainer.targetContest = item.id;
                contestContainer.addEventListener("click", selectContest);

                contestContainer.innerHTML = `<div class="oneCell rank">${item.type}</div><div class="oneCell username">${item.name}</div><div class="oneCell score">${item.tasks_count}</div>`;
                contests.appendChild(contestContainer);
            })
        }
    });
}

window.onload = drawContests();