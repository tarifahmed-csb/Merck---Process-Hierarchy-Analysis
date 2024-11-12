import './App.css';
import BuildButton from "./components/buildButton.js";
import QueryDrop from "./components/queryDrop.js";

function App() {
  return (
    <div className="App" wi="true">
      <h1>
        Merck Heirarchy Data Access
      </h1>
      <div>
        <BuildButton />
      </div>
      <br/>
      <div>
        <QueryDrop />
      </div>
    </div>
  );
}

export default App;