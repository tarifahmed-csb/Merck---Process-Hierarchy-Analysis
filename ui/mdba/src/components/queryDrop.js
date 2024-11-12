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
    if (query_type == 'measures') {
        
    } else if (query_type == 'processes') {
        
    } else if (query_type == 'rawMat') {
        
    } else {
        console.log('INVALID TYPE')
    }
    console.log('working?');
}

export default queryDrop