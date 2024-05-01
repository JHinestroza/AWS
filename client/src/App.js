import React from "react";
import {BrowserRouter as Router, Routes, Route} from 'react-router-dom'
import Screen1 from "./Screens/Screen1"
import Screen2 from "./Screens/Screen2"
import Screen3 from "./Screens/Screen3"
import Login from "./Components/Login";
import 'semantic-ui-css/semantic.min.css';

function App() {
  return (
     <Router>
      <Routes>
        <Route  path="/" element={<Login></Login>}></Route>
        <Route  path="/Screen1" element={<Screen1></Screen1>}></Route>
        <Route  path="/Screen2" element={<Screen2></Screen2>}></Route>
        <Route  path="/Screen3" element={<Screen3></Screen3>}></Route>


      </Routes>
     </Router>
  );
}

export default App;
