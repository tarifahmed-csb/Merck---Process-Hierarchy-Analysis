import React from 'react'
import './queryDrop.css'

function queryDrop() {
    return(
        <div>
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
                <label form='queryProc'>Enter a process to query</label>
                <br/>
                <input type='text' name='queryProc' id='queryProc' defaultValue='***'/>
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
    var procName = document.getElementById('queryProc').value
    var e = document.getElementById(type);
    var query_type = e.value;
    const urls  = ["http://localhost:1010/query", "http://localhost:1011/query", "http://localhost:1012/query"]

    for(var i = 0; i < 3; i++){
        var err = genQ(query_type, procName, urls[i]);
        if(err !== 0){
            console.log('ERROR querying %s\nERR CODE: %d', urls[i], err);
        } else {
            console.log('QUERIED %s', urls[i]);
        }
    }
    // genQ(query_type, procName, urls[1]); // query 1010
    // genQ(query_type, procName, urls[2]); // query 1011
    // genQ(query_type, procName, urls[3]); // query 1012
    function genQ(type, procName, url) {
        try{
            var xhr = new XMLHttpRequest();
            xhr.open("POST", url, true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.onreadystatechange = function () {
                if (xhr.readyState === 4 && xhr.status === 200) {
                    var json = JSON.parse(xhr.responseText);
                    console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
                    alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
                }
            };
            var data = JSON.stringify({"type": type, "name": procName});
            console.log("URL: "+url+"\nType: "+type+"\nName: "+procName);
            xhr.send(data);
            return 0;
        } catch (e) {
            return e;
        }
    }


//     console.log('-------------------- END NEW STUFF -------------------------------------------------');
//     var xhr;
//     if (query_type === 'measures') { // --------------------------------- Measurement Query ----------------------------------------------- 
//             //---------------------------------------------------------- to :1010 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1010/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             var data = JSON.stringify({"type": "measures", "name": procName});
//             console.log("URL: http://localhost:1010/query\nType: measures\nName: "+procName)
//             xhr.send(data);

//             //--------------------------------------------------------- to :1011 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1011/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "measures", "name": procName});
//             console.log("URL: http://localhost:1011/query\nType: measures\nName: "+procName)
//             xhr.send(data);

//             //--------------------------------------------------------- to :1012 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1012/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "measures", "name": procName});
//             console.log("URL: http://localhost:1012/query\nType: measures\nName: "+procName)
//             xhr.send(data);
//     } else if (query_type === 'processes') { // --------------------------------- Process Query --------------------------------------
//         //--------------------------------------------------------- to :1010 --------------------------------------------------------
//         xhr = new XMLHttpRequest();    
//         xhr.open("POST", "http://localhost:1010/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "process", "name": procName});
//             console.log("URL: http://localhost:1010/query\nType: process\nName: "+procName)
//             xhr.send(data);

//             //------------------------------------------------------ to :1011 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1011/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "process", "name": procName});
//             console.log("URL: http://localhost:1011/query\nType: process\nName: "+procName)
//             xhr.send(data);

//             //------------------------------------------------------- to :1012 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1012/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "process", "name": procName});
//             console.log("URL: http://localhost:1012/query\nType: process\nName: "+procName)
//             xhr.send(data);
//     } else if (query_type === 'rawMat') { // --------------------------------- Raw Material Query ---------------------------------------
//         //------------------------------------------------------------ to :1010 --------------------------------------------------------
//         xhr = new XMLHttpRequest();    
//         xhr.open("POST", "http://localhost:1010/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "rawMat", "name": procName});
//             console.log("URL: http://localhost:1010/query\nType: rawMat\nName: "+procName)
//             xhr.send(data);

//             //---------------------------------------------------------- to :1011 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1011/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "rawMat", "name": procName});
//             console.log("URL: http://localhost:1011/query\nType: rawMat\nName: "+procName)
//             xhr.send(data);

//             //----------------------------------------------------------- to :1012 --------------------------------------------------------
//             xhr = new XMLHttpRequest();
//             xhr.open("POST", "http://localhost:1012/query", true);
//             xhr.setRequestHeader("Content-Type", "application/json");
//             xhr.onreadystatechange = function () {
//                 if (xhr.readyState === 4 && xhr.status === 200) {
//                     var json = JSON.parse(xhr.responseText);
//                     console.log("Status: "+json.status + ", Time: " + json.time + ", Results: \n"+ json.results +"\nError: " +json.error);
//                     alert("Status: "+json.status + "\nTime: " + json.time + ", Results: \n"+ json.results + "\nError: " +json.error);
//                 }
//             };
//             data = JSON.stringify({"type": "rawMat", "name": procName});
//             console.log("URL: http://localhost:1012/query\nType: rawMat\nName: "+procName)
//             xhr.send(data);
//     } else {
//         console.log('INVALID TYPE')
//     }
//     console.log('working?');
}



export default queryDrop