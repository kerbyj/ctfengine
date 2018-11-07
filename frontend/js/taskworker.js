function drawTask() {
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT+'/api/tasks/getTaskById/'+this.needToDraw,
        success: function(data) {
            console.log(data);
            document.getElementById("taskName").innerText = `${data.name}`;
            document.getElementById("taskName").taskid = data.id;
            document.getElementById("taskDescriptionText").innerText = data.description;
            document.getElementById("taskCategory").innerHTML = `<span class="statElementKey">Category</span><span class="statElementValue">${data.category}</span>`;
            document.getElementById("taskContest").innerHTML = `<span class="statElementKey">Contest</span><span class="statElementValue">${data.contest}</span>`;
            document.getElementById("taskValue").innerHTML = `<span class="statElementKey">Value</span><span class="statElementValue">${data.value}</span>`;

            let attachmentsContainer = document.getElementById("attachments");
            if(data.attachments.length !== 0){
                data.attachments.forEach(function (attachment) {
                   let attachElement = document.createElement("span");
                   attachElement.classList.add("attachment");
                   attachElement.innerHTML = `<span class="attachment"><i class="material-icons md-14">cloud_upload</i> <a target="_blank" href="${GLOBAL_ENDPOINT}/files/${attachment.name}">${attachment.name}</a></span><br>`

                    attachmentsContainer.appendChild(attachElement)
                });
            }
        }
    });
}

function getAlwaysAliveTasks() {
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT+'/api/tasks/getAlwaysAliveTasks',
        success: function(data) {
            console.log(data);
            $.each(data, function (key, value) {
                let globalContainer = document.getElementById("tasksContainer");
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
                //console.log(tmpTasks)
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
    //document.getElementById("tasksContainer").style.visibility = "hidden";
};

span.onclick = function() {
    modal.style.display = "none";
};

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
        //document.getElementById("tasksContainer").style.visibility = "visible";
    }
};


function checkFlag(){
    let flagValue = document.getElementById("flagValue").value;
    let taskid = document.getElementById("taskName").taskid;
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT+'/api/tasks/checkFlag',
        data: {flag: flagValue, taskid: taskid},
        success: function(data) {
            console.log(data);
            if(data.result === true)
                document.getElementById(taskid).classList.add("solved");
        }
    });
}
document.getElementById("sendFlag").onclick = checkFlag;

