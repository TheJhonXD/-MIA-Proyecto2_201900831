package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	cmd "pack/packages/Commands"
	Filesystem "pack/packages/FileSystem"
	"pack/packages/Graphviz"
	"pack/packages/Structs"
	"pack/packages/Tools"
	"strings"
)

//Funciones para el manejo de las peticiones

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	response := Structs.Response{Message: "Hello world!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

/* func Console(text string ){
	fmt.Println("************* ⍟ Consola ⍟ *************")
	input := strings.Split(text, "\n")
	for _, i := range input {
		fmt.Println(i)
	}
} */

func analyzer(s string) string {
	messages := ""
	if s[0] != '#' {
		s = Tools.DeleteComments(s)
		s = strings.Trim(s, " ")
	}
	cmds := Tools.Split(s)

	if strings.ToLower(cmds[0]) == "mkdisk" {
		messages += "»» " + s + "\n"
		messages += cmd.MKDISK(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "rmdisk" {
		messages += "»» " + s + "\n"
		messages += cmd.RMDISK(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "fdisk" {
		messages += "»» " + s + "\n"
		messages += cmd.FDISK(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "mount" {
		messages += "»» " + s + "\n"
		messages += cmd.MOUNT(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "unmount" {
		messages += "»» " + s + "\n"
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "mkfs" {
		messages += "»» " + s + "\n"
		messages += cmd.MKFS(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "login" {
		messages += "»» " + s + "\n"
		// messages += cmd.MKFS(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "logout" {
		messages += "»» " + s + "\n"
		// messages += cmd.MKFS(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "mkgrp" {
		messages += "»» " + s + "\n"
		messages += cmd.MKGRP(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "rmgrp" {
		messages += "»» " + s + "\n"
		messages += cmd.RMGRP(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "mkusr" {
		messages += "»» " + s + "\n"
		messages += cmd.MKUSR(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "rmusr" {
		messages += "»» " + s + "\n"
		messages += cmd.RMUSR(cmds)
		messages += "-------------------------------------------------------------" + "\n"
	} else if strings.ToLower(cmds[0]) == "rep" {
		fmt.Println("»» " + s)
		messages += cmd.REP(cmds)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "pause" {
		fmt.Println("Press any key to continue...")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune() //Espero un caracter
	} else if s[0] == '#' {
		messages += s + "\n"
	} else if len(s) == 0 {
		messages += "" + "\n"
		fmt.Println("")
	} else {
		messages += "ERROR: el comando \"" + cmds[0] + "\" no es valido." + "\n"
	}
	return messages
}

func Console(text string) string {
	messages := ""
	input := strings.Split(text, "\n")
	for _, line := range input {
		if len(line) != 0 {
			messages += analyzer(line)
		}
	}
	return messages
}

func Inputs(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	response := Structs.Response{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	response.Message = Console(response.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ruta para crear un nuevo disco
func Graph(w http.ResponseWriter, r *http.Request) {
	response := Structs.ResponseGraph{}
	res := Structs.Response{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	// fmt.Println("ID:", response.Id)
	if response.Name == "graph-disk" {
		res.Message = Graphviz.GetDiskGraph(response.Id)
	} else if response.Name == "graph-sb" {
		res.Message = Graphviz.CreateSBReport(response.Id)
	} else if response.Name == "graph-file" {
		res.Message = Filesystem.UsersReport(response.Id)
	} else {
		fmt.Println("No se encontro el id")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	response := Structs.Login{}
	res := Structs.Response{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	_, message := Filesystem.Login(response.Usr, response.Pwd, response.Id)
	res.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	// response := Structs.Login{}
	res := Structs.Response{}
	// err := json.NewDecoder(r.Body).Decode(&response)
	/* if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	} */
	_, message := Filesystem.Logout()
	res.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
