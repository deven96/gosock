$(function(){
    var socket = null;
    var msgBox = $("#chatbox #message");
    var messages = $("#messages");
    $("#chatbox").submit(function(){
    if (!msgBox.val()) return false;
    if (!socket) {
    alert("Error: There is no socket connection.");
    return false;
    }
    socket.send(JSON.stringify({"Message": msgBox.val()}));
    msgBox.val("");
    return false;
    });
    if (!window["WebSocket"]) {
    alert("Error: Your browser does not support web sockets.")
    } else {
    // dynamically assign host location from golang
    socket = new WebSocket("ws://{{ .Host }}/room");
    socket.onclose = function() {
    alert("Connection has been closed.");
    }
    var i = 0
    socket.onmessage = function(e) {
        var socket_json = JSON.parse(e.data)
        // display_message
        var display_message = " : " +  socket_json.Message
        // create sender text
        var sender = socket_json.Code
        if (socket_json.Code == "#ffffff"){
            sender = "ADMIN"
        }
        var full_message = sender + display_message
        // create list item to house message
        var li_item = $("<li>")
                        .text(full_message)
                        .css("background-color", socket_json.Code);
        
        if (i % 2 == 0) {
            li_item.addClass("float-left")
        }
        else {
            li_item.addClass("float-right")
        }
        messages.append(li_item);
        i = i + 1
        // automatically scroll to last message
        $('#messages').animate({scrollTop: $('#messages').prop("scrollHeight")}, 500);
        return false
    }
    }
    });