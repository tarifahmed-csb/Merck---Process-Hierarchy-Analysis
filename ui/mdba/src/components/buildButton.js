import React from 'react'
import buildImg from '../assets/Build_Button.svg';

function buildButton() {
  return (
    <div>
        
        <img src={buildImg} alt='Build Database' onClick={() => DBuild('Building ...', 'test')}/>
    </div>
  )
}

async function DBuild(p1, p2) {
  console.log(p1)
  
  var xhr = new XMLHttpRequest();
  var testURL = "http://localhost:1010/build"
  var urls  = ["http://localhost:1010/build"/*, "http://localhost:1011/build", "http://localhost:1012/build"*/]
  // Loops through the declared URLs ^ and makes a request and awaits a response from all of them
  for (var url in urls){
    xhr.open("POST", testURL, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            var json = JSON.parse(xhr.responseText);
            console.log("Status: "+json.status + ", Time: " + json.time + ", Error: " +json.error);
            alert("Status: "+json.status + "\nTime: " + json.time + "\nError: " +json.error);
        }
    };
    var data = JSON.stringify({"type": "build", "name": "test"});
    xhr.send(data);
  }

  // ----------------------------------------------------------------------- KILL ME ----------------------------------------------------------
  //1011 for Dynamo
  /*
  fetch('http://localhost:1011/build', {
    method: 'POST',
    headers:{
      'Content-type': 'application/json'
    },
    body: JSON.stringify({
      Type: 'build',
      Name: 'test',
    }),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Suc:', data);
    var Ddata = data;
  })
  .catch((error) => {console.error('ERR', error)});

  //1012 for PostgreSQL
  fetch('http://localhost:1012/build', {
    method: 'POST',
    headers:{
      'Content-type': 'application/json'
    },
    body: JSON.stringify({
      key1: 'build',
      key2: p2
    }),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Suc:', data);
    var Ndata = data;
  })
  .catch((error) => {console.error('ERR', error)});
  */
}

export default buildButton