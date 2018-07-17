function danielonload() {
    document.getElementById("date").innerHTML = "Date: " + getdatetimer();

    showlinechart();
}


function showchart() {
    var ctx = document.getElementById("myChart").getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ["Red", "Blue", "Yellow", "Green", "Purple", "Orange"],
            datasets: [{
                label: '# of Votes',
                data: [12, 19, 3, 5, 2, 3],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                    'rgba(153, 102, 255, 0.2)',
                    'rgba(255, 159, 64, 0.2)'
                ],
                borderColor: [
                    'rgba(255,99,132,1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero: true
                    }
                }]
            }
        }
    });
}

function showlinechart() {
    showanychart('line','Y','Y');
}

function showlinechartprice() {
    showanychart('line','Y','N');
}
function showlinechartvolume() {
    showanychart('line','N','Y');
}


function showlinechartdate(){

    var fromDate = document.getElementById('fromDate');
    var toDate = document.getElementById('toDate');
    var currency = document.getElementById('currency');

    var url = "/btcmarketshistorylistdate?currency="+currency.value+"&fromDate="+fromDate.value+"&toDate="+toDate.value;
    window.location.href = url;
}

function showlinechart2(){
    var currency = document.getElementById('currency');
    var url = "/btcmarketshistorylist?currency="+currency.value+"&rows=200";
    window.location.href = url;
}

function showlinechartcoin(coin){
    var url = "/btcmarketshistorylist?currency="+coin+"&rows=200";
    window.location.href = url;
}


function showbarchart() {
    showanychart('bar','Y','Y');
}

function showanychart(chartype, lineprice, linevolume) {

    var theTbl = document.getElementById('tablecotacao');
    var arr = [];
    var arrvolume = [];
    var labellist = [];

    // find max price and max volume

    var maxpriceBTC = 0;
    var maxvolumeBTC = 0;
    var minpriceBTC = 0;
    var minvolumeBTC = 0;

    var maxpriceLTC = 0;
    var maxvolumeLTC = 0;
    var minpriceLTC = 0;
    var minvolumeLTC = 0;

    var maxpriceXRP = 0;
    var maxvolumeXRP = 0;
    var minpriceXRP = 0;
    var minvolumeXRP = 0;

    var maxpriceETH = 0;
    var maxvolumeETH = 0;
    var minpriceETH = 0;
    var minvolumeETH = 0;

    var maxpriceETC = 0;
    var maxvolumeETC = 0;
    var minpriceETC = 0;
    var minvolumeETC = 0;

    var maxpriceBCH = 0;
    var maxvolumeBCH = 0;
    var minpriceBCH = 0;
    var minvolumeBCH = 0;

    var st = theTbl.rows.length - 1;
    for (var i = st; i >= 1; i--) {

        var valuecol0 = theTbl.rows[i].cells[0].innerHTML.substr(66, 3);
        // Price Coin
        var price = Number(theTbl.rows[i].cells[2].innerHTML);
        // Volume Coin
        var volume = Number(theTbl.rows[i].cells[4].innerHTML);

        if (valuecol0 == "BTC") 
        {
            if (price > maxpriceBTC) {
                maxpriceBTC = price
            }
            if (volume > maxvolumeBTC) {
                maxvolumeBTC = volume
            }
        }
        if (valuecol0 == "LTC") 
        {
            if (price > maxpriceLTC) {
                maxpriceLTC = price
            }
            if (volume > maxvolumeLTC) {
                maxvolumeLTC = volume
            }
        }
        if (valuecol0 == "XRP") 
        {
            if (price > maxpriceXRP) {
                maxpriceXRP = price
            }
            if (volume > maxvolumeXRP) {
                maxvolumeXRP = volume
            }
        }
        if (valuecol0 == "ETH") 
        {
            if (price > maxpriceETH) {
                maxpriceETH = price
            }
            if (volume > maxvolumeETH) {
                maxvolumeETH = volume
            }
        }
        if (valuecol0 == "ETC") 
        {
            if (price > maxpriceETC) {
                maxpriceETC = price
            }
            if (volume > maxvolumeETC) {
                maxvolumeETC = volume
            }
        }
        if (valuecol0 == "BCH") 
        {
            if (price > maxpriceBCH) {
                maxpriceBCH = price
            }
            if (volume > maxvolumeBCH) {
                maxvolumeBCH = volume
            }
        }
    }

    // for (var i = 1; i < theTbl.rows.length; i++) {
    // var st = theTbl.rows.length - 1;
    for (var i = st; i >= 1; i--) {
        // for (var i = 1; i < theTbl.rows.length; i++) {

        // This is the X label
        var Xlabel = theTbl.rows[i].cells[7].innerHTML;
        // var valuecol0 = theTbl.rows[i].cells[0].innerHTML.substr(62, 3);
        var valuecol0 = theTbl.rows[i].cells[0].innerHTML.substr(66, 3);

        if (valuecol0 == "AUD") {
            arr.push(Number(theTbl.rows[i].cells[3].innerHTML));
            arrvolume.push(Number(theTbl.rows[i].cells[5].innerHTML)/1);
        }

            if (valuecol0 == "BTC") {
                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeBTC)*maxpriceBTC);

                arr.push(price);
                arrvolume.push(relativevolume);
            }

            if (valuecol0 == "LTC") {
                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeLTC)*maxpriceLTC);

                arr.push(price);
                arrvolume.push(relativevolume);
            }

            if (valuecol0 == "ETH") {
                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeETH)*maxpriceETH);

                arr.push(price);
                arrvolume.push(relativevolume);
            }

            if (valuecol0 == "XRP") {

                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeXRP)*maxpriceXRP);

                arr.push(price);
                arrvolume.push(relativevolume);
            }

            if (valuecol0 == "ETC") {
                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeETC)*maxpriceETC);

                arr.push(price);
                arrvolume.push(relativevolume);
            }

            if (valuecol0 == "BCH") {
                // Price Coin
                var price = Number(theTbl.rows[i].cells[2].innerHTML);
                // Volume Coin
                var volume = Number(theTbl.rows[i].cells[4].innerHTML);
                // Relative Volume
                var relativevolume = ((volume/maxvolumeBCH)*maxpriceBCH);

                arr.push(price);
                arrvolume.push(relativevolume);
            }
            

            if (valuecol0 == "ALL") {
                arr.push(Number(theTbl.rows[i].cells[2].innerHTML));
                arrvolume.push(Number(theTbl.rows[i].cells[4].innerHTML));
            }

        labellist.push(Xlabel.substr(30, 8));
    }

    if (linevolume == 'N') {
        showchartlineprice(arr, labellist, chartype, 'Price');
    }
    else
    {
        if (lineprice == 'N') {
            showchartlineprice(arrvolume, labellist, chartype, 'Volume');
        }
        else
        {
            showchartline(arr, arrvolume, labellist, chartype);
        }
    }

}

function showchartline(datalist, volumelist, labellist, chartype) {

    var ctx = document.getElementById("myChart").getContext('2d');

    var myChart = new Chart(ctx, {
        type: chartype,
        data: 
        {
            labels: labellist,
            datasets: 
            [
                {
                    label: 'Price',
                    data: datalist,
                    borderColor: "#3cba9f",
                    fill: true
                },
                {
                    label: 'Volume',
                    data: volumelist,
                    borderColor: "#c45850",
                    fill: true
                },
                
            ]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {}
                }]
            }
        }
    });

}



function showchartlineprice(pricelist, labellist, chartype, chartlabel) {

    var ctx = document.getElementById("myChart").getContext('2d');

    var myChart = new Chart(ctx, {
        type: chartype,
        data: 
        {
            labels: labellist,
            datasets: 
            [
                {
                    label: chartlabel,
                    data: pricelist,
                    borderColor: "#3cba9f",
                    fill: false
                },
            ]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {  }
                }]
            }
        }
    });

}
