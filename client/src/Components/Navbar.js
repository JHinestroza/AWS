import React from 'react';
import './Navbar.css'; // Asegúrate de que la ruta al archivo CSS sea correcta
import { Icon } from 'semantic-ui-react';

const Navbar = () => {
  return (
    <nav className="navbar">
      <ul className="navbar-nav">
        <li className="nav-item">
          <a href="/Screen1" className="nav-link">
          <Icon name='code' />
            <span className="link-text">Comando</span>
          </a>
        </li>
        <li className="nav-item">
          <a href="/Screen2" className="nav-link">
            <Icon name='folder open' />
            <span className="link-text">Pantalla1</span>
          </a>
        </li>
        <li className="nav-item">
          <a href="/Screen3" className="nav-link">
            <Icon name='chart bar' />
            <span className="link-text">Pantalla2</span>
          </a>
        </li>
        <li className="nav-item">
          <a href="/" className="nav-link">
            <Icon name='lock' />
            <span className="link-text">Login</span>
          </a>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;

