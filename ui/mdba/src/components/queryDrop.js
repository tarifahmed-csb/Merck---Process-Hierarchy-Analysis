import React from 'react'

function queryDrop() {
    return(
        <div>
            <label form='queries'>Select a query type</label>
            <br/>
            <select name='queries' id='queries' defaultValue='default'>
            <option value='measures'>Find all Measurements</option>
            <option value='processes'>Find all Processes</option>
            <option value='rawMat'>Find all Raw Materials</option>
            </select>
            <br/>
            <br/>
            <button type='button' form='queries' formMethod='post' onClick={() => query('queries')}>
            <span className="transition"></span>
            <span className="gradient"></span>
            <span className="label">Run Query</span>
            </button>
        </div>
    )
}

async function query(type) {
    var e = document.getElementById(type);
    var query_type = e.value;
    console.log(query_type);
    if (query_type == 'measures') { // --------------------------------- Measurement Query -----------------------------------------------
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:1010/query";
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
                var json = JSON.parse(xhr.responseText);
                console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
            }
        };
        var data = JSON.stringify({"type": "measures", "name": "test"});
        xhr.send(data);
    } else if (query_type == 'processes') { // --------------------------------- Process Query --------------------------------------
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:1010/query";
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
                var json = JSON.parse(xhr.responseText);
                console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
            }
        };
        var data = JSON.stringify({"type": "process", "name": "***"});
        xhr.send(data);
    } else if (query_type == 'rawMat') { // --------------------------------- Raw Material Query ---------------------------------------
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:1010/query";
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
                var json = JSON.parse(xhr.responseText);
                console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
            }
        };
        var data = JSON.stringify({"type": "rawMat", "name": "test"});
        xhr.send(data);
    } else {
        console.log('INVALID TYPE')
    }
    console.log('working?');
}

export default queryDrop