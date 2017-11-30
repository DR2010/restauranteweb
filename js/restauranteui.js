var prato

// --------------------------------------
//    Load Prices in Memory or Cache
// --------------------------------------
function loadPrices() {
    prato = new Array();
    prato[0] = { name: "Coxinha", price: "7" };
    prato[1] = { name: "BolodeCenoura", price: "3" };
    prato[2] = { name: "Refrigerante", price: "5" };
    prato[3] = { name: "BolodeAipim", price: "9" };
    prato[4] = { name: "PasteldeQueijo", price: "5" };
    prato[5] = { name: "Brigadeiro", price: "8" };
}

// --------------------------------------
//       Add New Item to Order
// --------------------------------------
function addNewItem() {

    loadPrices()

    var table = document.getElementById("myTable");
    var valueplate = document.getElementById("pratoname");
    var plateqtd = document.getElementById("pratoqtd");

    var lastRow = table.rows[table.rows.length];
    var lastRowNumber = table.rows.length;

    var row = table.insertRow(lastRow);

    var cell0 = row.insertCell(0);
    var cell1 = row.insertCell(1);
    var cell2 = row.insertCell(2);
    var cell3 = row.insertCell(3);

    x = '<input type=checkbox name=row' + lastRowNumber + ' id=checkitem' + lastRowNumber + ' value=' + valueplate + '>';

    var sel = document.getElementById("pratoname").selectedIndex;
    cell0.innerHTML = x;
    cell1.innerHTML = valueplate.nodeValue;
    cell1.innerHTML = valueplate.options[sel].text;
    cell2.innerHTML = plateqtd.value;
    cell3.innerHTML = plateqtd.value * prato[sel].price;
}

// --------------------------------------
//       Clean up fields
// --------------------------------------
function newOrder() {

    var orderID = document.getElementById("orderID");
    var orderClientID = document.getElementById("orderClientID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var foodeatplace = document.getElementById("foodeatplace");
    var status = document.getElementById("status");

    orderID.value = "";
    orderClientID.value = "";
    orderClientName.value = "";
    orderDate.value = "";
    orderTime.value = "";
    foodeatplace.value = "";
    status.value = "New Order.";

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
    var orderClientID = document.getElementById("orderClientID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var foodeatplace = document.getElementById("foodeatplace");
    var status = document.getElementById("status");

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    var pratosselected = new Array();

    //loops through rows    
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var pratoname = "";
        var quantidade = "";
        var preco = "";


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

            pratosselected[i] = { pratoname: pratoname, quantidade: quantidade, preco: preco };
        }
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";

    var paramsjson = JSON.stringify({
        orderID: orderID.value,
        orderClientID: orderClientID.value,
        orderClientName: orderClientName.value,
        orderDate: orderDate.value,
        orderTime: orderTime.value,
        foodeatplace: foodeatplace.value,
        status: status.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            status.value = "Order placed successfully."

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(paramsjson);

}

// --------------------------------------
//          Save order JSON fetch
// --------------------------------------
function saveOrderJSONfetch() {

    var orderID = document.getElementById("orderID");
    var orderClientID = document.getElementById("orderClientID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var foodeatplace = document.getElementById("foodeatplace");
    var status = document.getElementById("status");

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    var pratosselected = new Array();

    //loops through rows    
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var pratoname = "";
        var quantidade = "";
        var preco = "";


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

            pratosselected[i] = { pratoname: pratoname, quantidade: quantidade, preco: preco };
        }
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";

    var paramsjson = JSON.stringify({
        orderID: orderID.value,
        orderClientID: orderClientID.value,
        orderClientName: orderClientName.value,
        orderDate: orderDate.value,
        orderTime: orderTime.value,
        foodeatplace: foodeatplace.value,
        status: status.value
    });


    fetch(url, {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: paramsjson
        })
        .then(function(res) {
            console.log(res)
        })
        .catch(function(res) {
            console.log(res)
        })
}


// --------------------------------------
//          Save order FORM
// --------------------------------------
function saveOrderX() {

    var orderID = document.getElementById("orderID");
    var orderClientID = document.getElementById("orderClientID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var foodeatplace = document.getElementById("foodeatplace");
    var status = document.getElementById("status");

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    var pratosselected = new Array();

    //loops through rows    
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var pratoname = "";
        var quantidade = "";
        var preco = "";


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

            pratosselected[i] = { pratoname: pratoname, quantidade: quantidade, preco: preco };
        }
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";
    var params =
        "orderClientID=" + orderClientID.value +
        "&orderClientName=" + orderClientName.value +
        "&orderDate=" + orderDate.value +
        "&orderTime=" + orderTime.value +
        "&foodeatplace=" + foodeatplace.value;

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    // http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            status.value = "Order placed successfully."

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(params);

}