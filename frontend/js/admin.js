function updateContestLists(){
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT + '/api/tasks/getContestList',
        success: function (data) {
            console.log(data);
            let contestLists = document.getElementsByClassName("tasksList");
            //console.log(contestLists)
            /*
            contestLists.forEach(function (list) {
                list.innerHTML = "<option selected>Select contest</option>";
            });
            */
            for(i = 0; i < contestLists.length; i++){
                //console.log(contestLists[i])
                contestLists[i].innerHTML = "<option value='default' selected>Select contest</option>";
            }

            data.forEach(function (element) {
                //console.log(element.id, element.name);
                for(i = 0; i < contestLists.length; i++){
                    let tmpOption = document.createElement("option");
                    tmpOption.value = element.id;
                    tmpOption.innerText = element.name;

                    contestLists[i].appendChild(tmpOption);
                }
            })
        }
    });
}

function createContest() {
    let contestName = document.getElementById("contest_name").value;
    let contestType = document.getElementById("contest_type").value;
    let contestVisibility = document.getElementById("contest_visibility").checked;
    let contestPermit = document.getElementById("contest_permit").checked;

    console.log(contestName, contestPermit, contestVisibility, contestType);

    contestVisibility === true ? contestVisibility = 1 : contestVisibility = 0;
    contestPermit === true ? contestPermit = 1 : contestPermit = 0;

    console.log(contestName, contestPermit, contestVisibility, contestType);

    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/admin/createContest',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
        data: {
            contest_name: contestName,
            contest_type: contestType,
            visibility: contestVisibility,
            permit: contestPermit
        },
        success: function (data) {
            console.log(data);

            updateContestLists();
        }
    });
}

function createTask() {
    let taskName = document.getElementById("task_name").value;
    let taskFlag = document.getElementById("task_flag").value;
    let taskPrice = document.getElementById("task_price").value;
    let taskCategory = document.getElementById("task_category").value;
    let description = document.getElementById("task_description").value;

    let taskContestId = document.getElementById("createTaskContestList").value;

    console.log(taskContestId, description);

    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/admin/createTask',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
        data: {
            task_name: taskName,
            task_flag: taskFlag,
            task_price: taskPrice,
            task_category: taskCategory,
            task_description: description,
            task_contest: taskContestId
        },
        success: function (data) {
            console.log(data);

            updateContestLists();
        }
    });
}


document.getElementById("newtaskButton").onclick = createTask;
document.getElementById("newcontestButton").onclick = createContest;

document.onload = updateContestLists();

