import React from 'react'
import data from "../assets/logDat.txt"
import './log.css'

function log(){
    

    return(
        <div className='log' wi="true">
            <h1>
                <button onClick={loadData}>
                    <span className="transition"></span>
                    <span className="gradient"></span>
                    <span className="label">Logged Requests</span>
                </button>
            </h1>
            <table id='logTable' align='center'>
            </table>
        </div>
    );
}

async function loadData() {
    //adapted from https://stackoverflow.com/questions/55643149/how-to-store-fetch-api-json-response-in-a-javascript-object
    fetch(data).then(r => r.text()).then(txt => displayData(txt));
}

async function displayData(data) {
    console.log(data);
    data = data.split(/\r\n|\r|\n/);
    console.log(data);

    //adapted from https://stackoverflow.com/questions/57613694/separating-an-array-of-strings-to-create-a-table-row
    var tabHTML = "";
    var index = "1";
    if(data[1] != null){
        for (let ln in data){
            tabHTML += '<tr>';
            tabHTML += '<td>'+index+"| "+data[ln]+'</td>';
            tabHTML += '</tr>';
            index++;
        }
        document.getElementById("logTable").innerHTML = tabHTML;
    }
}

export default log;