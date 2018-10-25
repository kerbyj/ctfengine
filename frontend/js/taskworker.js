function drawTask() {
    console.log(this.needToDraw);
    $.post({
        type: 'get',
        url: 'http://localhost/api/tasks/getTaskById/'+this.needToDraw,
        success: function(data) {
            // console.log(data);
            document.getElementById("taskName").innerText = data.name;
            document.getElementById("taskDescription").innerText = data.description;
        }
    });
}

function getAlwaysAliveTasks() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/tasks/getAlwaysAliveTasks',
        success: function(data) {
            console.log(data);
            //return;
            let container = document.createElement("div");
            container.classList.add("tableraw");

            $.each(data, function (key, value) {
                let tmpColumn = document.createElement("div");
                tmpColumn.classList.add("column");

                let tmpCategoryName = document.createElement("div");
                tmpCategoryName.classList.add("cell");
                tmpCategoryName.classList.add("headerCategory");
                tmpCategoryName.innerText = key;
                tmpColumn.appendChild(tmpCategoryName);

                value.forEach(function (element, index) {
                    let tmpTask = document.createElement("div");
                    tmpTask.classList.add("cell");
                    tmpTask.classList.add("task");

                    tmpTask.needToDraw = element.id;
                    tmpTask.addEventListener("mouseover", drawTask);
                    tmpTask.addEventListener("click", showModal);

                    tmpTask.innerText = element.value;
                    tmpColumn.appendChild(tmpTask);
                    //console.log(tmpColumn);
                });
                container.appendChild(tmpColumn);
            });
            document.getElementById("tasksContainer").appendChild(container);
        }
    });
}
window.onload = getAlwaysAliveTasks;

// Get the modal
var modal = document.getElementById('taskView');

// Get the button that opens the modal
var btn = document.getElementById("myBtn");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on the button, open the modal
function showModal() {
    modal.style.display = "block";
};

// When the user clicks on <span> (x), close the modal
span.onclick = function() {
    modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}