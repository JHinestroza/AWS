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

func ValidarDatosGrupos(context []string, action string) string {
	name := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "name") {
			name = tk[1]
		}
	}
	if name == "" {
		return Error(action+"GRP", "No se encontró el parámetro name en el comando.")
	}
	if Comparar(action, "MK") {
		mkgrp(name)
	} else if Comparar(action, "RM") {
		rmgrp(name)
	} else {
		return Error(action+"GRP", "No se reconoce este comando.")
	}
	return ""
}

func mkgrp(n string) string {
	if !Comparar(Logged.User, "root") {

		return Error("MKGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
	}

	var path string
	partition, mensaje := GetMount("MKGRP", Logged.Id, &path)
	if mensaje != "" {
		return mensaje
	}
	if string(partition.Part_status) == "0" {
		return Error("MKGRP", "No se encontró la partición montada con el id: "+Logged.Id)
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco.")
		return Error("MKGRP", "No se ha encontrado el disco.")

	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo")
		return Error("MKGRP", "Error al leer el archivo")

	}
	inode := Structs.NewInodo()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodo{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodo{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo")
		return Error("MKGRP", "Error al leer el archivo")

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
			Error("MKGRP", "Error al leer el archivo")
			return Error("MKGRP", "Error al leer el archivo")

		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'G' || linea[2] == 'g' {
			c++
			in := strings.Split(linea, ",")
			if in[2] == n {
				if linea[0] != '0' {
					return Error("MKGRP", "EL nombre "+n+", ya está en uso.")
				}
			}
		}
	}
	txt += strconv.Itoa(c+1) + ",G," + n + "\n"

	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		return Error("MKGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		return Error("MKGRP", "No se ha encontrado el disco.")
	}

	for i := 0; i < len(cadenasS); i++ {

		var fbAux Structs.BloqueArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))+int64(unsafe.Sizeof(Structs.BloqueArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodo{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	file.Close()
	return Mensaje("MKGRP", "Grupo "+n+", creado correctamente!")
}

func rmgrp(n string) string {
	if !Comparar(Logged.User, "root") {
		return Error("RMGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
	}

	var path string
	partition, mensaje := GetMount("RMGRP", Logged.Id, &path)
	if mensaje != "" {
		return mensaje
	}
	if string(partition.Part_status) == "0" {
		return Error("RMGRP", "No se encontró la partición montada con el id: "+Logged.Id)
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		return Error("RMGRP", "No se ha encontrado el disco.")
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		return Error("RMGRP", "Error al leer el archivo")
	}
	inode := Structs.NewInodo()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodo{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodo{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		return Error("RMGRP", "Error al leer el archivo")
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
			return Error("RMGRP", "Error al leer el archivo")
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == n {
				existe = true
				aux += strconv.Itoa(0) + ",G," + in[2] + "\n"
				continue
			}
		}
		aux += linea + "\n"
	}
	if !existe {
		return Error("RMGRP", "No se encontró el grupo \""+n+"\".")
	}
	txt = aux

	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		return Error("RMGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		return Error("RMGRP", "No se ha encontrado el disco.")
	}
	for i := 0; i < len(cadenasS); i++ {

		var fbAux Structs.BloqueArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))+int64(unsafe.Sizeof(Structs.BloqueArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloqueCarpetas{}))+int64(unsafe.Sizeof(Structs.BloqueArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodo{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	file.Close()
	return Mensaje("RMGRP", "Grupo "+n+", eliminado correctamente!")
}
