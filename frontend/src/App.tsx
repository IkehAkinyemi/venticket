import "./App.css";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { EventListContainer } from "./components/Event/event_list_container";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <div className="container">
          <Route
            path='/'
            element={
              <EventListContainer eventListURL="http://localhost:8181" />
            }
          ></Route>
        </div>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
