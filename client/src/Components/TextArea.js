import React, { useState } from 'react';
import './Codigo.css';
import axios from 'axios';


function Codigo() {
    // Estado para almacenar los comandos ingresados
    const [commands, setCommands] = useState([]);
    // Estado para el valor actual del input
    const [currentCommand, setCurrentCommand] = useState('');

    // Función para manejar el envío de comandos
    const handleCommandSubmit = (event) => {
      event.preventDefault();
      console.log(currentCommand);
      if (currentCommand !== "") {
          const data = {
              code: currentCommand,  // Cambiado de 'comandos' a 'code' para alinearse con tu backend
          };
    
          axios.post('http://127.0.0.1:8080/interpreter', data)
              .then(response => {
                  console.log('Backend response:', response.data);
                  // Actualiza el estado de los comandos con la respuesta del servidor si es necesario
                  setCommands(prevCommands => [...prevCommands, response.data.message]);  // Asumiendo que quieres guardar la respuesta
              })
              .catch(error => {
                  console.error('Error al enviar el comando:', error);
              });
    
          // Limpiar el campo de entrada después de enviar
          setCurrentCommand('');
      }
  };

     
    return (
        <div className="command-interface">
            <div className="command-history">
                <ul>
                    {commands.map((cmd, index) => (
                        <li key={index}>{cmd}</li>
                    ))}
                </ul>
            </div>
            <form onSubmit={handleCommandSubmit} className="command-input">
                <input
                    type="text"
                    value={currentCommand}
                    onChange={(e) => setCurrentCommand(e.target.value)}
                    placeholder="Ingrese su comando"
                />
                <button type="submit">Enviar</button>
            </form>
        </div>
    );
}

export default Codigo;