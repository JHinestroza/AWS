package Comandos

import (
	"API/Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type UsuarioActivo struct {
	User     string
	Password string
	Id       string
	Uid      int
	Gid      int
}

var Logged UsuarioActivo

func ValidarDatosLOGIN(context []string) (bool, string) {
	id := ""
	user := ""
	pass := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "id") {
			id = tk[1]
		} else if Comparar(tk[0], "user") {
			user = tk[1]
		} else if Comparar(tk[0], "pass") {
			pass = tk[1]
		}
	}
	if id == "" || user == "" || pass == "" {
		return false, Error("LOGIN", "Se necesitan parámetros obligatorios para el comando LOGIN.")
	}
	bander, mensaje := sesionActiva(user, pass, id)
	return bander, mensaje
}

func sesionActiva(u string, p string, id string) (bool, string) {
	var path string
	partition, mensaje := GetMount("LOGIN", id, &path)
	if mensaje != "" {
		return false, mensaje
	}
	if string(partition.Part_status) == "0" {

		return false, Error("LOGIN", "No se encontró la partición montada con el id: "+id)
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {

		return false, Error("LOGIN", "No se ha encontrado el disco.")
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		return false, Error("LOGIN", "Error al leer el archivo")
	}
	inode := Structs.NewInodo()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodo{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodo{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		return false, Error("LOGIN", "Error al leer el archivo")
	}

	var fb Structs.BloqueArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))+int64(unsafe.Sizeof(Structs.BloqueArchivos{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(Structs.BloqueArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			return false, Error("LOGIN", "Error al leer el archivo")
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'U' || linea[2] == 'u' {
			in := strings.Split(linea, ",")
			if Comparar(in[3], u) && Comparar(in[4], p) && in[0] != "0" {
				idGrupo := "0"
				existe := false
				for j := 0; j < len(vctr)-1; j++ {
					line := vctr[j]
					if (line[2] == 'G' || line[2] == 'g') && line[0] != '0' {
						inG := strings.Split(line, ",")
						if inG[2] == in[2] {
							idGrupo = inG[0]
							existe = true
							break
						}
					}
				}
				if !existe {
					return false, Error("Login", "No se encontró el grupo \""+in[2]+"\".")
				}

				Logged.Id = id
				Logged.User = u
				Logged.Password = p
				Logged.Uid, _ = strconv.Atoi(in[0])
				Logged.Gid, _ = strconv.Atoi(idGrupo)
				return true, Mensaje("LOGIN", "logueado correctamente")
			}
		}
	}

	return false, Error("LOGIN", "No se encontró el usuario "+u)
}

func CerrarSesion() bool {
	Mensaje("LOGOUT", "El usuario "+Logged.User+", ha finalizado sesion")
	Logged = UsuarioActivo{}
	return false
}
