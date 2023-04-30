package main

import (
	"fmt"
	"log"
	"net/http"

	"pack/server"

	"github.com/gorilla/mux"
)

/* func analyzer(s string) {
	if s[0] != '#' {
		s = Tools.DeleteComments(s)
		s = strings.Trim(s, " ")
	}
	cmds := Tools.Split(s)

	if strings.ToLower(cmds[0]) == "mkdisk" {
		fmt.Println("»» " + s)
		cmd.MKDISK(cmds)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "rmdisk" {
		fmt.Println("»» " + s)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "fdisk" {
		fmt.Println("»» " + s)
		cmd.FDISK(cmds)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "mount" {
		fmt.Println("»» " + s)
		cmd.MOUNT(cmds)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "unmount" {
		fmt.Println("»» " + s)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "mkfs" {
		fmt.Println("»» " + s)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "rep" {
		fmt.Println("»» " + s)
		cmd.REP(cmds)
		fmt.Println("-------------------------------------------------------------")
	} else if strings.ToLower(cmds[0]) == "pause" {
		fmt.Println("Press any key to continue...")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune() //Espero un caracter
	} else if s[0] == '#' {
		fmt.Println(s)
	} else if len(s) == 0 {
		fmt.Println("")
	} else {
		fmt.Println("ERROR: el comando \"" + cmds[0] + "\" no es valido.")
	}
} */

/* func readFile(path string) {
	if Tools.Exists(path) {
		//Abrir el archivo
		myfile, err := os.Open(path)
		if err != nil {
			log.Fatalln(err)
		}

		// Crear un scanner para leer el archivo línea por línea
		scanner := bufio.NewScanner(myfile)

		// Iterar sobre cada línea del archivo
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) != 0 {
				analyzer(line)
			}
		}

		// Verificar si ocurrió algún error mientras se leía el archivo
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
		myfile.Close()
	} else {
		fmt.Println("El archivo no existe")
	}
} */

/* func main() {
	//Preparar el scanner para leer la entrada
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("************* ⍟ Consola ⍟ *************")

	for {

		fmt.Print("»» ")
		//Leer la entrada
		reader.Scan()
		//Obtener la entrada
		input := reader.Text()
		// input := `execute >path="/home/jhonx/Descargas/Archivos_de_Entrada_Proyecto_1/Proyecto_1/Parte 1/2-crear-particiones.eea"`
		//Separar la entrada en comandos
		cmds := Tools.Split(strings.Trim(input, " "))
		//Verificar si el comando no es "exit" para leer el archivo
		if input != "exit" {
			if strings.ToLower(cmds[0]) == "execute" {
				readFile(strings.Trim(strings.Split(cmds[1], "=")[1], "\""))
			} else if input != "exit" {
				analyzer(input)
			}
		} else {
			break
		}
	}
} */

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(server.CorsMiddleware)
	router.HandleFunc("/", server.Index)
	router.HandleFunc("/text", server.Inputs)
	router.HandleFunc("/graph", server.Graph)
	fmt.Println("*****************************************************************")
	fmt.Println("*\n*\tServidor corriendo en http://localhost:3000/ \t\t*")
	fmt.Println("*\n*****************************************************************")
	log.Fatal(http.ListenAndServe(":3000", router))
}
