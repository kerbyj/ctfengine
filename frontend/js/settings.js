function changePassword() {

}

function LeaveCommand() {
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT + '/api/users/LeaveCommand',
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
    });
}

function CreateCommand() {
    let commandName = document.getElementById("new_command_name").value;
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/CreateCommand',
        data: {commandname: commandName},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
    });
}

function RenameCommand() {
    let commandName = document.getElementById("rename_input").value;
    console.log(commandName);
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/RenameCommand',
        data: {commandname: commandName},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
    });
}

function DeleteCommand() {
    console.log(this.command_id)
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/DeleteCommand',
        data: {commandid: this.command_id},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
    });
}

function drawMembers() {
    let commandManipulateContainer = document.getElementById("commandManipulateContainer");

    let commandInfo = document.createElement("div");
    commandInfo.innerHTML = `Command name: <b>${this.command_name}</b><p>`;

    commandManipulateContainer.prepend(commandInfo);

    let mainUserTable = document.getElementById("userTable");
    let captainId = this.captain_id;
    let yourId = this.your_id;

    mainUserTable.innerHTML = `<div class="rowUserTable"><div class="cellUserTable">Username</div>
                                <div class="cellUserTable">Status</div>
                                ${captainId == yourId ? `<div class="cellUserTable">Function</div>` : ""}
                                <div> `;


    $.each(this.members, function (key, element) {
        console.log(key, element);
        let userRow = document.createElement("div");
        userRow.classList.add("rowUserTable");

        let rowInnerData = `<div class="cellUserTable">${element}</div>`;

        rowInnerData += `<div class="cellUserTable">${captainId == key ? "Captain" : yourId == key ? "You" : ""}</div>`;

        if (captainId === yourId) {
            rowInnerData += `<div class="cellUserTable">Drop</div>`
        }

        userRow.innerHTML = rowInnerData;
        mainUserTable.appendChild(userRow);
    });

}

function drawManipulateForms() {
    console.log("MANIPULATE");
    let mainCommandContainer = document.getElementById("commandManipulateContainer");

    if (this.command_id === 0) {
        /*
            CREATE COMMAND
         */

        let nameInput = document.createElement("input");
        nameInput.placeholder = "command_name";
        nameInput.id = "new_command_name";
        nameInput.classList.add("settingsInput");
        nameInput.type = "text";
        mainCommandContainer.appendChild(nameInput);

        let CommandButton = document.createElement("div");
        CommandButton.classList.add("buttonSettings");
        CommandButton.addEventListener("click", CreateCommand);

        CommandButton.innerText = "./create_command";
        mainCommandContainer.appendChild(CommandButton);

        /*
            JOIN COMMAND
         */
        let tokenInput = document.createElement("input");
        tokenInput.placeholder = "invite_token";
        tokenInput.classList.add("settingsInput");
        tokenInput.style.marginTop = "25px";
        tokenInput.type = "text";
        mainCommandContainer.appendChild(tokenInput);

        let JoinCommandButton = document.createElement("div");
        JoinCommandButton.classList.add("buttonSettings");

        JoinCommandButton.innerText = "./join_command";
        mainCommandContainer.appendChild(JoinCommandButton);

    } else if ((this.command_id !== 0) && (this.captain_id == this.your_id)) {
        let RenameInput = document.createElement("input");
        RenameInput.placeholder = "new_command_name";
        RenameInput.classList.add("settingsInput");
        RenameInput.style.marginTop = "10px";
        RenameInput.type = "text";
        RenameInput.id = "rename_input";
        RenameInput.needToRename = this.command_id;


        mainCommandContainer.appendChild(RenameInput);

        let ApplyRename = document.createElement("div");
        ApplyRename.classList.add("buttonSettings");
        ApplyRename.addEventListener("click", RenameCommand);
        ApplyRename.innerText = "./apply_rename";
        mainCommandContainer.appendChild(ApplyRename);

        let DeleteButton = document.createElement("div");
        DeleteButton.classList.add("buttonSettings");
        DeleteButton.innerText = "./delete_command";
        DeleteButton.command_id = this.command_id;

        DeleteButton.style.border = "1px #F71D16 solid";
        DeleteButton.addEventListener("click", DeleteCommand);
        mainCommandContainer.appendChild(DeleteButton);

    } else {
        let LeaveButton = document.createElement("div");
        LeaveButton.classList.add("buttonSettings");
        LeaveButton.addEventListener("click", LeaveCommand);

        LeaveButton.innerText = "./leave_command";
        mainCommandContainer.appendChild(LeaveButton);
    }
}

function getCommandInfo() {
    $.post({
        type: 'get',
        url: GLOBAL_ENDPOINT + '/api/users/getCommandStatusForSettings',
        success: function (data) {
            console.log(data);
            //data.drawCaptain = drawCaptainInterface;
            data.drawMembers = drawMembers;

            if (data.command_id !== 0)
                data.drawMembers();

            data.drawManipulateForms = drawManipulateForms;
            data.drawManipulateForms();

            if (data.command_id === 0) {
                console.log("Draw new command interface");
                return;
            }

            if (data.your_id === data.captain_id) {
                console.log("Draw captain interface");
                data.drawCaptain;
            } else {
                console.log("Draw member interface")
            }
        }
    });
}

window.onload = getCommandInfo();

