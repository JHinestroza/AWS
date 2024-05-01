package Comandos

import (
	"API/Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

func Comparar(a string, b string) bool { //creamos el case insensitive
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func Error(op string, mensaje string) string {
	fmt.Println("\tERROR: " + op + "\n\tTIPO: " + mensaje)
	return "\tERROR: " + op + "\n\tTIPO: " + mensaje
}

func Mensaje(op string, mensaje string) string {
	message := "\tCOMANDO: " + op + "\n\tTIPO: " + mensaje
	fmt.Println(message)
	return message
}

func Confirmar(mensaje string) bool {
	fmt.Println(mensaje + " (s/n)")
	var respuesta string
	respuesta = "s"
	if Comparar(respuesta, "s") {
		return true
	}
	return false
}

func ArchivoExiste(ruta string) bool { //verifica que si exista
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerDisco(path string) *Structs.MBR {
	m := Structs.MBR{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &m)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return nil
	}
	var mDir *Structs.MBR = &m
	return mDir
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
