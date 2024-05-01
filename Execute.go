package main

import (
	"API/Backend/Comandos"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// 					exec -path=/home/daniel/Escritorio/ArchivosPrueba/Calificacion_Proyecto2.script

var logued = false

func Executer(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		texto := fileScanner.Text()
		texto = strings.TrimSpace(texto)
		tk := Comando(texto)
		if texto != "" {
			if Comandos.Comparar(tk, "pause") {
				fmt.Println("************************************** FUNCIÓN PAUSE **************************************")
				var pause string
				Comandos.Mensaje("PAUSE", "Presione \"enter\" para continuar...")
				fmt.Scanln(&pause)
				continue
			} else if string(texto[0]) == "#" {
				fmt.Println("************************************** COMENTARIO **************************************")
				Comandos.Mensaje("COMENTARIO", texto)
				continue
			}
			texto = strings.TrimLeft(texto, tk)
			tokens := SepararTokens(texto)
			funciones(tk, tokens)
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
}

func Comando(text string) string {
	var tkn string
	terminar := false
	for i := 0; i < len(text); i++ {
		if terminar {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !terminar {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				terminar = true
			}
		}
	}
	return tkn
}

func SepararTokens(texto string) []string {
	var tokens []string
	if texto == "" {
		return tokens
	}
	texto += " "
	var token string
	estado := 0
	for i := 0; i < len(texto); i++ {
		c := string(texto[i])
		if estado == 0 && c == "-" {
			estado = 1
		} else if estado == 0 && c == "#" {
			continue
		} else if estado != 0 {
			if estado == 1 {
				if c == "=" {
					estado = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if estado == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					estado = 3
					continue
				} else {
					estado = 4
				}
			} else if estado == 3 {
				if c == "\"" {
					estado = 4
					continue
				}
			} else if estado == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if estado == 4 && c == " " {
				estado = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func funciones(token string, tks []string) string {
	if token != "" {
		if Comandos.Comparar(token, "MKDISK") {
			fmt.Println("*************************************** FUNCIÓN MKDISK **************************************")
			return Comandos.ValidarDatosMKDISK(tks)
		} else if Comandos.Comparar(token, "RMDISK") {
			fmt.Println("*************************************** FUNCIÓN RMDISK **************************************")
			return Comandos.RMDISK(tks)
		} else if Comandos.Comparar(token, "FDISK") {
			fmt.Println("*************************************** FUNCIÓN FDISK  **************************************")
			return Comandos.ValidarDatosFDISK(tks)
		} else if Comandos.Comparar(token, "MOUNT") {
			fmt.Println("*************************************** FUNCIÓN MOUNT  **************************************")
			mensaje := Comandos.ValidarDatosMOUNT(tks)
			Comandos.ListaMount()
			return mensaje
		} else if Comandos.Comparar(token, "MKFS") {
			fmt.Println("*************************************** FUNCIÓN MKFS  **************************************")
			return Comandos.ValidarDatosMKFS(tks)
		} else if Comandos.Comparar(token, "REP") {
			fmt.Println("*************************************** FUNCIÓN REP  **************************************")
			return Comandos.ValidarDatosREP(tks)
		} else if Comandos.Comparar(token, "LOGIN") {
			mensaje := ""
			fmt.Println("*************************************** FUNCIÓN LOGIN  **************************************")
			if logued {
				return Comandos.Error("LOGIN", "Ya hay un usuario en línea.")
			} else {
				logued, mensaje = Comandos.ValidarDatosLOGIN(tks)
				if mensaje != "" {
					return mensaje
				}
			}
		} else if Comandos.Comparar(token, "LOGOUT") {
			fmt.Println("*************************************** FUNCIÓN LOGOUT  **************************************")
			if !logued {
				return Comandos.Error("LOGOUT", "Aún no se ha iniciado sesión.")
			} else {
				logued = Comandos.CerrarSesion()
			}
		} else if Comandos.Comparar(token, "MKGRP") {
			fmt.Println("*************************************** FUNCIÓN MKGRP  **************************************")
			if !logued {
				return Comandos.Error("MKGRP", "Aún no se ha iniciado sesión.")

			} else {
				Comandos.ValidarDatosGrupos(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMGRP") {
			fmt.Println("*************************************** FUNCIÓN RMGRP  **************************************")
			if !logued {
				return Comandos.Error("RMGRP", "Aún no se ha iniciado sesión.")

			} else {
				Comandos.ValidarDatosGrupos(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKUSER") {
			fmt.Println("*************************************** FUNCIÓN MKUSER  **************************************")
			if !logued {
				return Comandos.Error("MKUSER", "Aún no se ha iniciado sesión.")
			} else {
				return Comandos.ValidarDatosUsers(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMUSR") {
			fmt.Println("*************************************** FUNCIÓN RMUSER  **************************************")
			if !logued {
				return Comandos.Error("RMUSER", "Aún no se ha iniciado sesión.")
			} else {
				return Comandos.ValidarDatosUsers(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKDIR") {
			fmt.Println("*************************************** FUNCIÓN MKDIR  **************************************")
			if !logued {
				mensaje := Comandos.Error("MKDIR", "Aún no se ha iniciado sesión.")
				return mensaje
			} else {
				var p string
				particion, err := Comandos.GetMount("MKDIR", Comandos.Logged.Id, &p)
				if err != "" {
					return err
				}
				Comandos.ValidarDatosMKDIR(tks, particion, p)
			}
		} else if Comandos.Comparar(token, "MKFILE") {
			fmt.Println("*************************************** FUNCIÓN MKFILE  **************************************")
			if !logued {
				return Comandos.Error("MKFILE", "Aún no se ha iniciado sesión.")
			} else {
				var p string
				particion, mensaje := Comandos.GetMount("MKDIR", Comandos.Logged.Id, &p)
				if mensaje != "" {
					return mensaje
				}
				Comandos.ValidarDatosMKFILE(tks, particion, p)
			}
		} else if Comandos.Comparar(token, "pause") {
			fmt.Println("************************************** FUNCIÓN PAUSE **************************************")
			var pause string
			Comandos.Mensaje("PAUSE", "Presione \"enter\" para continuar...")
			fmt.Scanln(&pause)

		} else if string(token[0]) == "#" {
			fmt.Println("************************************** COMENTARIO **************************************")
			return Comandos.Mensaje("COMENTARIO", token)
		} else {
			return Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+token+"\"")
		}
	}
	return ""
}

func Exec(texto string) string {

	tk := Comando(texto)
	if texto != "" {
		if Comandos.Comparar(tk, "pause") {
			fmt.Println("************************************** FUNCIÓN PAUSE **************************************")
			var pause string
			Comandos.Mensaje("PAUSE", "Presione \"enter\" para continuar...")
			fmt.Scanln(&pause)
		} else if string(texto[0]) == "#" {
			fmt.Println("************************************** COMENTARIO **************************************")
			Comandos.Mensaje("COMENTARIO", texto)
		}
		texto = strings.TrimLeft(texto, tk)
		tokens := SepararTokens(texto)
		retorno := funciones(tk, tokens)
		return retorno
	}
	return ""

}
