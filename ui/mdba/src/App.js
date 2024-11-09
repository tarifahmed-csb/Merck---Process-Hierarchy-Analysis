import './App.css';
import Button from "./components/buildButton.js";

function App() {
  return (
    <div className="App" wi>
      <h1>
        Merck Heirarchy Data Access
      </h1>

      <div>
        <Button />
      </div>

      <div>
        <label form='queries'>Select a query type</label>
        <br/>
        <select name='queries' id='queries' defaultValue='default'>
          <option value='measures'>Find all Measurements</option>
          <option value='processes'>Find all Processes</option>
          <option value='rawMat'>Find all Raw Materials</option>
        </select>
        <br/>
        <button type='submit' form='queries' formMethod='post' onSubmit="query('queries')">
          <span class="transition"></span>
          <span class="gradient"></span>
          <span class="label">Run Query</span>
        </button>
      </div>
    </div>
  );
}

async function query(s1){
var e = document.getElementById(s1);
var query_type = e.value;
console.log(query_type);
}

export default App;