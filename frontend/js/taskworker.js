function drawTask() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/tasks/getTaskById/'+this.needToDraw,
        success: function(data) {
            //console.log(data);
            document.getElementById("taskName").innerText = `${data.name} - ${data.value}`;
            document.getElementById("taskName").taskid = data.id;
            document.getElementById("taskDescriptionText").innerText = data.description;
        }
    });
}

function getAlwaysAliveTasks() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/tasks/getAlwaysAliveTasks',
        success: function(data) {
            $.each(data, function (key, value) {
                let globalContainer = document.getElementById("container");
                let categoryHeader = document.createElement("span");
                categoryHeader.classList.add("tasks");
                categoryHeader.classList.add("header");
                categoryHeader.innerText = key;

                globalContainer.appendChild(categoryHeader);

                let tmpTasks = {};
                $.each(value, function (key, value) {
                    let category = value.category;
                    if(tmpTasks[category] === undefined)
                        tmpTasks[category] = [];

                    tmpTasks[category].push(value);
                });
                let container = document.createElement("div");
                container.classList.add("tableraw");


                let containerForTask = document.createElement("div");
                containerForTask.classList.add("tasks");
                $.each(tmpTasks, function (category, tasks) {
                    let tmpColumn = document.createElement("div");
                    tmpColumn.classList.add("column");
                    let tmpCategoryName = document.createElement("div");
                    tmpCategoryName.classList.add("cell");
                    tmpCategoryName.classList.add("headerCategory");
                    tmpCategoryName.innerText = category;
                    tmpColumn.appendChild(tmpCategoryName);

                    tasks.forEach(function (element, index) {
                        let tmpTask = document.createElement("div");

                        if(element.solved)
                            tmpTask.classList.add("solved");

                        tmpTask.classList.add("cell");
                        tmpTask.classList.add("task");

                        tmpTask.id = element.id;
                        tmpTask.title = element.Name;

                        tmpTask.needToDraw = element.id;

                        tmpTask.innerText = element.value;
                        tmpColumn.appendChild(tmpTask);

                        tmpTask.addEventListener("mouseover", drawTask);
                        tmpTask.addEventListener("click", showModal);
                    });
                    container.appendChild(tmpColumn);
                    containerForTask.appendChild(container);
                });
                globalContainer.appendChild(containerForTask);
            });
        }
    });
}
window.onload = getAlwaysAliveTasks;

var modal = document.getElementById('taskView');
var btn = document.getElementById("myBtn");
var span = document.getElementsByClassName("close")[0];
function showModal() {
    modal.style.display = "block";
};

span.onclick = function() {
    modal.style.display = "none";
};

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
};


function checkFlag(){
    let flagValue = document.getElementById("flagValue").value;
    let taskid = document.getElementById("taskName").taskid;
    $.post({
        type: 'post',
        url: 'http://localhost/api/tasks/checkFlag',
        data: {flag: flagValue, taskid: taskid},
        success: function(data) {
            console.log(data);
            if(data.result === true)
                document.getElementById(taskid).classList.add("solved");
        }
    });
}
document.getElementById("sendFlag").onclick = checkFlag;

