var tasksContainer = document.getElementById("tasksContainer");

function getAlwaysAliveTasks() {
    $.post({
        type: 'get',
        url: 'http://localhost/api/tasks/getAlwaysAliveTasks',
        success: function(data) {
            console.log(data);

            tasksContainer.innerHTML += `<h2>Always alive</h2>`;

            var alwaysAliveTasksContainer = document.createElement("div");
            alwaysAliveTasksContainer.id = "alwaysAliveTasks";
            alwaysAliveTasksContainer.classList.add("block");

            tasksContainer.appendChild(alwaysAliveTasksContainer);

            $.each(data, function (key, value) {
                let newCategory = document.createElement("h3");
                newCategory.innerText = key.toUpperCase();
                alwaysAliveTasksContainer.appendChild(newCategory);

                value.forEach(function (item) {
                    let newTask = document.createElement("div");
                    newTask.classList.add("statsElements");
                    newTask.innerHTML = `${item.Name}<span class="valueInfo">${item.value}</span>`;

                    alwaysAliveTasksContainer.appendChild(newTask);
                })
            });
        }
    });
}
window.onload = getAlwaysAliveTasks;