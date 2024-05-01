import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Icon, List } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css';

function FileSelector() {
    const [files, setFiles] = useState([]);

    useEffect(() => {
        axios.get('http://127.0.0.1:8080/SelectFile')
            .then(response => {
                setFiles(response.data);
            })
            .catch(error => {
                console.error('Error fetching files:', error);
            });
    }, []);

    return (
        <div style={{  paddingLeft: '200px' }}>
            <h1>Seleccione un archivo .dsk</h1>
            <List>
                {files.map((file, index) => (
                    <List.Item key={index}>
                        <Icon name="disk"  size='huge'/> {/* Cambiado a 'disk' para un ícono más específico */}
                        {file}
                    </List.Item>
                ))}
            </List>
        </div>
    );
}
export default FileSelector;