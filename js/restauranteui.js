// Global variables, storing prices in memory before using cached data or cookie.
//
var prato

// --------------------------------------
// When loading order to add
// --------------------------------------
function onloadorderadd() {
    getdatetime();
}

// --------------------------------------
//       calcular troco
// --------------------------------------
function calcularTroco() {

    var valorpago = document.getElementById("paidAmount");
    var valordaconta = document.getElementById("orderTotal");
    var troco = document.getElementById("troco");

    var valorpagoval = parseFloat(valorpago.value);
    var valordacontaval = parseFloat(valordaconta.value);

    var trocoval = valorpagoval - valordacontaval;

    var totalstrdec = parseFloat(trocoval).toFixed(2);

    troco.value =totalstrdec;
}

function printOrder() {
    window.print();
}


// --------------------------------------
//       Add New Item to Order
// --------------------------------------
function addNewItem() {

    // Check if order has been placed
    var orderID = document.getElementById("orderID");

    if (orderID.value != "") {
        message.value = "Order already placed!"
        return
    }

    // loadPrices()

    var table = document.getElementById("myTable");
    var valueplate = document.getElementById("pratoname");
    var plateqtd = document.getElementById("pratoqtd");

    // Get price
    var pricefromhtml = document.getElementById(valueplate.value);

    var lastRow = table.rows[table.rows.length];
    var lastRowNumber = table.rows.length;

    var row = table.insertRow(lastRow);

    var cell0 = row.insertCell(0);
    var cell1 = row.insertCell(1);
    var cell2 = row.insertCell(2);
    var cell3 = row.insertCell(3);
    var cell4 = row.insertCell(4);

    x = '<input type=checkbox name=row' + lastRowNumber + ' id=checkitem' + lastRowNumber + ' value=' + valueplate + '>';

    var preco = parseFloat(pricefromhtml.value);
    var quantidade = parseFloat(plateqtd.value);
    var total = preco * quantidade;

    var sel = document.getElementById("pratoname").selectedIndex;

    // Checkbox first column 0
    // Dish name column 1
    // Quantidade column 2
    // Preco column 3
    // Total calculado  column 4

    cell0.innerHTML = x;
    cell1.innerHTML = valueplate.options[sel].text;
    // cell1.innerHTML = valueplate.value;
    // cell1.innerHTML = valueplate.nodeValue;
    cell2.innerHTML = plateqtd.value;
    cell3.innerHTML = pricefromhtml.value;

    var totalstrdec = parseFloat(total).toFixed(2)
    cell4.innerHTML = totalstrdec;

    // cell4.innerHTML = total;
    // cell3.innerHTML = plateqtd.value * pricefromhtml.value;
    // cell3.innerHTML = plateqtd.value * prato[sel].price;
}

// --------------------------------------
//       Clean up fields
// --------------------------------------
function newOrderX() {

    var orderID = document.getElementById("orderID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var eatmode = document.getElementById("EatMode");
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    orderID.value = "";
    orderClientName.value = "";
    orderDate.value = getdate();
    orderTime.value = "";
    eatmode.value = "Eatin";
    status.value = "New Order";
    message.value = "Place new order";

}

function newOrder(){
    var url = "/orderadddisplay";
    window.location.href = url;
}


// --------------------------------------
//   Show Selected Rows - Debug only
// --------------------------------------

function showSelectedRows() {

    var selchbox = []; // array that will store the value of selected checkboxes

    var table = document.getElementById("myTable");

    var lastRowNumber = table.rows.length - 1;

    for (var i = lastRowNumber; i >= 0; i--) {
        var chk = document.getElementById('checkitem' + i);

        var col2 = table.rows[i].cells[2].innerHTML;

        if (chk != null)
            if (chk.checked) alert(col2);

    }
    return selchbox;
}

// --------------------------------------
//   Remove Selected Rows from Order
// --------------------------------------

function removeSelectedRows() {
    // JavaScript & jQuery Course - http://coursesweb.net/javascript/
    var selchbox = []; // array that will store the value of selected checkboxes

    var table = document.getElementById("myTable");
    var lastRowNumber = table.rows.length;

    for (var i = lastRowNumber; i >= 0; i--) {
        var chk = document.getElementById('checkitem' + i);

        if (chk != null)
            if (chk.checked) table.deleteRow(i);

    }

    return selchbox;
}


// --------------------------------------
//          Save order JSON
// --------------------------------------
function saveOrder() {

    var orderID = document.getElementById("orderID");
    var userID = document.getElementById("userID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    // var eatmode = document.getElementById("EatMode");
    var eatmode = "EatIn";
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    if (orderClientName.value == "") {
        message.value = "Order name is mandatory!"
        orderClientName.focus();
        return
    }


    if (orderID.value != "") {
        message.value = "Order already placed!"
        alert("Order already placed.")
        return
    }

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    if (rowLength == 1) {
        message.value = "You haven't added any items yet. You can only place an order with items."
        return
    }

    // Get Confirmation
    var answer = confirm("Confirm order placement?")
    if (answer) {
        //some code
    }
    else {
        return
    }



    var pratosselected = new Array();

    //loops through rows    
    // Skip row = 0 pois e' o header.
    // Porem o novo array comeca com zero

    v = 0;
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var action = "";
        var pratoname = "";
        var quantidade = "";
        var preco = "";
        var total = "";

    // Checkbox first column 0
    // Dish name column 1
    // Quantidade column 2
    // Preco column 3
    // Total calculado  column 4

        for (var j = 0; j < cellLength; j++) {

            var cellVal = oCells.item(j).innerHTML;
            if (j == 0) {
                action = cellVal;
            }
            if (j == 1) {
                pratoname = cellVal;
            }
            if (j == 2) {
                quantidade = cellVal;
            }
            if (j == 3) {
                preco = cellVal;
            }
            if (j == 4) {
                total = cellVal;
            }

        }

        if (action == "") continue;
        if (action == "Action") continue;

        pratosselected[v] = { pratoname: pratoname, quantidade: quantidade, price: preco, total: total };
        v++;
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";

    status.value = "Placed";

    var paramsjson = JSON.stringify({
        ID: orderID.value,
        ClientName: orderClientName.value,
        ClientID: userID.innerText,
        Date: orderDate.value,
        Time: orderTime.value,
        EatMode: eatmode.value,
        Status: status.value,
        Items: pratosselected
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."
            status.value = "Placed"

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(paramsjson);

}


// --------------------------------------------------------------
//   20-Mar-2018 - Save order
//   This version includes the user ID generated by the backend
// --------------------------------------------------------------
function saveOrder2() {

    // This is the user ID logged on or the Generated User ID
    // If it is set to Anonymous a new user ID will be generated
    //
    var userID = document.getElementById("userID");
    var orderID = document.getElementById("orderID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    // var eatmode = document.getElementById("EatMode");
    var eatmode = "EatIn";
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    if (orderClientName.value == "") {
        message.value = "Order name is mandatory!"
        orderClientName.focus();
        return
    }


    if (orderID.value != "") {
        message.value = "Order already placed!"
        alert("Order already placed.")
        return
    }

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    if (rowLength == 1) {
        message.value = "You haven't added any items yet. You can only place an order with items."
        return
    }

    // Get Confirmation
    var answer = confirm("Confirm order placement?")
    if (answer) {
        //some code
    }
    else {
        return
    }



    var pratosselected = new Array();

    //loops through rows    
    // Skip row = 0 pois e' o header.
    // Porem o novo array comeca com zero

    v = 0;
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var action = "";
        var pratoname = "";
        var quantidade = "";
        var preco = "";
        var total = "";

    // Checkbox first column 0
    // Dish name column 1
    // Quantidade column 2
    // Preco column 3
    // Total calculado  column 4

        for (var j = 0; j < cellLength; j++) {

            var cellVal = oCells.item(j).innerHTML;
            if (j == 0) {
                action = cellVal;
            }
            if (j == 1) {
                pratoname = cellVal;
            }
            if (j == 2) {
                quantidade = cellVal;
            }
            if (j == 3) {
                preco = cellVal;
            }
            if (j == 4) {
                total = cellVal;
            }

        }

        if (action == "") continue;
        if (action == "Action") continue;

        pratosselected[v] = { pratoname: pratoname, quantidade: quantidade, price: preco, total: total };
        v++;
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderclientadd";

    status.value = "Placed";

    var paramsjson = JSON.stringify({
        ID: orderID.value,
        ClientName: orderClientName.value,
        ClientID: userID.innerText,
        Date: orderDate.value,
        Time: orderTime.value,
        EatMode: eatmode.value,
        Status: status.value,
        Items: pratosselected
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."
            status.value = "Placed"

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(paramsjson);

}



// --------------------------------------------------------
//          Save order JSON and Open WebSockets Connection
// --------------------------------------------------------
function saveOrderWithWebSockets() {

    var orderID = document.getElementById("orderID");
    var userID = document.getElementById("userID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    // var eatmode = document.getElementById("EatMode");
    var eatmode = "EatIn";
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    if (orderClientName.value == "") {
        message.value = "Order name is mandatory!"
        orderClientName.focus();
        return
    }


    if (orderID.value != "") {
        message.value = "Order already placed!"
        alert("Order already placed.")
        return
    }

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    if (rowLength == 1) {
        message.value = "You haven't added any items yet. You can only place an order with items."
        return
    }

    // Get Confirmation
    var answer = confirm("Confirm order placement?")
    if (answer) {
        //some code
    }
    else {
        return
    }



    var pratosselected = new Array();

    //loops through rows    
    // Skip row = 0 pois e' o header.
    // Porem o novo array comeca com zero

    v = 0;
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var action = "";
        var pratoname = "";
        var quantidade = "";
        var preco = "";
        var total = "";

    // Checkbox first column 0
    // Dish name column 1
    // Quantidade column 2
    // Preco column 3
    // Total calculado  column 4

        for (var j = 0; j < cellLength; j++) {

            var cellVal = oCells.item(j).innerHTML;
            if (j == 0) {
                action = cellVal;
            }
            if (j == 1) {
                pratoname = cellVal;
            }
            if (j == 2) {
                quantidade = cellVal;
            }
            if (j == 3) {
                preco = cellVal;
            }
            if (j == 4) {
                total = cellVal;
            }

        }

        if (action == "") continue;
        if (action == "Action") continue;

        pratosselected[v] = { pratoname: pratoname, quantidade: quantidade, price: preco, total: total };
        v++;
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";

    status.value = "Placed";

    var paramsjson = JSON.stringify({
        ID: orderID.value,
        ClientName: orderClientName.value,
        ClientID: userID.innerText,
        Date: orderDate.value,
        Time: orderTime.value,
        EatMode: eatmode.value,
        Status: status.value,
        Items: pratosselected
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."
            status.value = "Placed"

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(paramsjson);
    
    // Open the websocket connection
    //

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
}

// --------------------------------------
//          Save order JSON
// --------------------------------------
function backToList() {

    window.location.replace("/orderlist");
}

function pad(num, size) {
    return ('000000000' + num).substr(-size);
}

function getdatetime() {
    var date = document.getElementById("orderDate");
    var time = document.getElementById("orderTime");

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd;
    date.value = today;

    var hour = hh + ':' + min;
    time.value = hour;

}

function getdate() {

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd;

    return today

}

function getdatetimer() {

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd + ' ~ ' + hh + ':' + min;

    return today

}

// --------------------------------------
//          Login
// --------------------------------------
function login() {

    var userid = document.getElementById("userid");
    var password = document.getElementById("password");
    var message = document.getElementById("message");

    if (userid.value == "") {
        userid.focus();
        return
    }

    if (password.value == "") {
        password.focus()
        return
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/login";

    var paramsjson = JSON.stringify({
        userid: userid.value,
        password: password.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."

            var json_data = http.responseText;
            var contact = JSON.parse(json_data);
        } else {
            message.value = "Invalid Credentials."
        }
    }

    http.send(paramsjson);

}

