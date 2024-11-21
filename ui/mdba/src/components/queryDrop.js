import React from 'react'
import './queryDrop.css'

function queryDrop() {
    return(
        <div>
            
            
            <br/>
            <br/>
            <div id='querySel'>
                <label form='queries'>Select a query type</label>
                <br/>
                <select name='queries' id='queries' defaultValue='default'>
                    <option value='measures'>Find all Measurements</option>
                    <option value='processes'>Find all Processes</option>
                    <option value='rawMat'>Find all Raw Materials</option>
                </select>
            </div>
            <div id='procEntry'>
                <label form='process'>Enter a process to query</label>
                <br/>
                <input type='text' name='process' id='process' defaultValue='***'/>
            </div>
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
    var procName = document.getElementById('process').value
    var e = document.getElementById(type);
    var query_type = e.value;
    var urls  = ["http://localhost:1010/query", "http://localhost:1011/query", "http://localhost:1012/query"]
    // Loops through the declared URLs ^ and makes a request and awaits a response from all of them
    console.log(query_type);
    var xhr = new XMLHttpRequest();
    if (query_type == 'measures') { // --------------------------------- Measurement Query -----------------------------------------------
        for (var i = 0; i < urls.length; i++){
            xhr.open("POST", urls[i], true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.onreadystatechange = function () {
                if (xhr.readyState === 4 && xhr.status === 200) {
                    var json = JSON.parse(xhr.responseText);
                    console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                    alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
                }
            };
            var data = JSON.stringify({"type": "measures", "name": procName});
            xhr.send(data);
        }
    } else if (query_type == 'processes') { // --------------------------------- Process Query --------------------------------------
        for (var i = 0; i < urls.length; i++){
            xhr.open("POST", urls[i], true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.onreadystatechange = function () {
                if (xhr.readyState === 4 && xhr.status === 200) {
                    var json = JSON.parse(xhr.responseText);
                    console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                    alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
                }
            };
            var data = JSON.stringify({"type": "process", "name": procName});
            xhr.send(data);
    }
    } else if (query_type == 'rawMat') { // --------------------------------- Raw Material Query ---------------------------------------
        for (var i = 0; i < urls.length; i++){
            xhr.open("POST", urls[i], true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.onreadystatechange = function () {
                if (xhr.readyState === 4 && xhr.status === 200) {
                    var json = JSON.parse(xhr.responseText);
                    console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                    alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
                }
            };
            var data = JSON.stringify({"type": "rawMat", "name": procName});
            xhr.send(data);
        }
    } else {
        console.log('INVALID TYPE')
    }
    console.log('working?');
}

export default queryDrop