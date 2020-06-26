function changePassword() {
    let oldPassword = document.getElementById("oldPass").value;
    let newPassword1 = document.getElementById("newPass1").value;
    let newPassword2 = document.getElementById("newPass2").value;

    if(newPassword1 !== newPassword2){
        console.log("Пароли не совпадают");
        return;
    }

    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/ChangePassword',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
        data: {oldPassword: oldPassword, newPassword:newPassword1},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
    });
}

function changeUsername() {
    let newName = document.getElementById("newusernameInput").value;

    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/ChangeUsername',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
        data: {newName: newName},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");
        }
    });
}

function LeaveCommand() {
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/LeaveCommand',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
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
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
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
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
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
    //console.log(this.command_id);
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/DeleteCommand',
        data: {commandid: this.command_id},
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
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
    commandInfo.innerHTML = `Command name: <b>${this.command_name}</b><br>Command invite: <b>${this.command_invite}</b><p>`;

    commandManipulateContainer.prepend(commandInfo);

    let mainUserTable = document.getElementById("userTable");
    let captainId = this.captain_id;
    let yourId = this.your_id;

    mainUserTable.innerHTML = `<div class="rowUserTable"><div class="cellUserTable">Username</div>
                                <div class="cellUserTable">Status</div>
                                <div> `;


    $.each(this.members, function (key, element) {
        console.log(key, element);
        let userRow = document.createElement("div");
        userRow.classList.add("rowUserTable");

        let usernameContainer = document.createElement("div");
        usernameContainer.classList.add("cellUserTable");
        usernameContainer.innerText = element;
        userRow.appendChild(usernameContainer);

        let statusContainer = document.createElement("div");
        statusContainer.classList.add("cellUserTable");
        statusContainer.innerText = captainId == key ? "Captain" : yourId == key ? "You" : "";
        userRow.appendChild(statusContainer);

        /*
        if (captainId === yourId) {
            let dropButtonContainer = document.createElement("div");
            dropButtonContainer.classList.add("cellUserTable");
            dropButtonContainer.innerText = "Drop";
            dropButtonContainer.addEventListener("click");
            userRow.appendChild(dropButtonContainer);
        }
        */

        mainUserTable.appendChild(userRow);
    });

}

function JoinCommand(){
    let invite = document.getElementById("inviteInput").value;
    $.post({
        type: 'post',
        url: GLOBAL_ENDPOINT + '/api/users/JoinCommandViaInvite',
        headers: { 'X-CSRF-Token': CurrentCsrfToken },
        data: {invite: invite},
        success: function (data) {
            console.log(data);
            if (data.status === "success")
                location.reload();
            else
                alert("some error");

        }
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
        tokenInput.id = "inviteInput";
        tokenInput.classList.add("settingsInput");
        tokenInput.style.marginTop = "25px";
        tokenInput.type = "text";
        mainCommandContainer.appendChild(tokenInput);

        let JoinCommandButton = document.createElement("div");
        JoinCommandButton.classList.add("buttonSettings");
        JoinCommandButton.addEventListener("click", JoinCommand);
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
document.getElementById("newpassButton").addEventListener("click", changePassword);
document.getElementById("newusername").addEventListener("click", changeUsername);

document.getElementById("logoutButton").addEventListener("click", function () {
    location.replace("/logout");
});