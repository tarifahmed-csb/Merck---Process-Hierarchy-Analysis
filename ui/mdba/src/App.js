import './App.css';
import BuildButton from "./components/buildButton.js";
import QueryDrop from "./components/queryDrop.js";
import LogData from "./components/log.js"

function App() {
  return (
    <div className="App" wi="true">
      <h1>
        Merck Heirarchy Data Access
      </h1>

      <div id='bigBuild'>
        <input id='buildNum' name='buildNum' className='input' placeholder='How many' />
        <br/>
        <button onClick={() => bigBuild()}>
          <span className="transition"></span>
          <span className="gradient"></span>
          <span className="label">Mass Build</span>
        </button>
      </div> 
      <div>
        <BuildButton />
      </div>
      <br />
      <div>
        <QueryDrop />
      </div>
      <div>
        <LogData />
      </div>
    </div>
  );
}

async function bigBuild() {
  var num = 50;
  num = document.getElementById('buildNum').value;
  console.log('---------- Running %d build cycles ---------- ', num);
  var xhr = new XMLHttpRequest();
  var urls  = ["http://localhost:1010/build", "http://localhost:1011/build", "http://localhost:1012/build"]
  // Loops through the declared URLs ^ and makes a request and awaits a response from all of them
  for (var j = 0; j < num; j++){
    var procName = 'process'+(j+1);
    for (var i = 0; i < urls.length; i++){
      console.log(urls[i])
      xhr.open("POST", urls[i], true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.onreadystatechange = function () {
          if (xhr.readyState === 4 && xhr.status === 200) {
              var json = JSON.parse(xhr.responseText);
              console.log("Status: "+json.status + ", Time: " + json.time + ", Error: " +json.error);
              alert("Status: "+json.status + "\nTime: " + json.time + "\nError: " +json.error);
          }
      };
      var data = JSON.stringify({"type": "build", "name": procName});
      xhr.send(data);
    }
  }


}

export default App;