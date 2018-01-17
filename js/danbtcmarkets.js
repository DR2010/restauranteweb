// --------------------------------------
//          Save order JSON
// --------------------------------------
function saveBTCPreOrder() {
    debugger;
    var minXRP = document.getElementById("minXRP");
    var maxXRP = document.getElementById("maxXRP");
    var minETH = document.getElementById("minETH");
    var maxETH = document.getElementById("maxETH");
    var message = document.getElementById("message");

    var coinsmax = new Array();
    coinsmax[0] = { Currency: "XRP", Min: minXRP.value, Max: maxXRP.value };
    coinsmax[1] = { Currency: "ETH", Min: minETH.value, Max: maxETH.value };

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/btcpreorderadd";


    var paramsjson = JSON.stringify({
        coinsmax
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            var json_data = http.responseText; q
            var contact = JSON.parse(json_data);
        }
    }

    http.send(paramsjson);

    message.value = "Ok."
}