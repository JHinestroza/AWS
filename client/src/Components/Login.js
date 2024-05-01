import React, { useState } from 'react';
import axios from 'axios'; // Importa Axios
import './Login.css';
import backgroundImage from '../Images/background.jpg';

function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleLogin = async (event) => {
        event.preventDefault(); // Evita el comportamiento predeterminado del formulario
        try {
            const response = await axios.post('http://127.0.0.1:8080/Login', {
                username: username,
                password: password,
            });
            console.log('Usuario autenticado:', response.data); // Aquí manejas la respuesta
        } catch (error) {
            console.error('Error en la autenticación:', error);
        }
    };

    return (
        <div className="login-container" style={{ backgroundImage: `url(${backgroundImage})` }}>
            <form className="login-form" onSubmit={handleLogin}>
                <h2>Login</h2>
                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                />
                <button type="submit">Login</button>
                <div className="options">
                    <label>
                        <input type="checkbox" /> Remember me
                    </label>
                    <a href="#forgot">Forgot password?</a>
                </div>
            </form>
        </div>
    );
}

export default Login;