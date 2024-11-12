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
  // 1010 for Graph
  fetch('http://localhost:1010/build', {
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
    var Gdata = data;
  })
  .catch((error) => {console.error('ERR', error)});

  //1011 for Dynamo
  fetch('http://localhost:1011/build', { 
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
}

export default buildButton

// On Click will need to access the DB backends and run their mains