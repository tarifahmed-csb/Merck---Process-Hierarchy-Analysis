import React from 'react'
import buildImg from '../assets/Build_Button.svg';

function buildButton() {
  return (
    <div>
        <img src={buildImg} alt='Build Database' onClick={() => DBuild('Building ...')}/>
    </div>
  )
}

async function DBuild(p1) {
  console.log(p1)
  fetch('http://localhost:8080', {
    method: 'POST',
    headers:{
      'Content-type': 'application/json'
    },
    body: JSON.stringify({
      key1: 'build'
    }),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Suc:', data);
  })
  .catch((error) => {console.error('ERR', error)});
}

export default buildButton

// On Click will need to access the DB backends and run their mains