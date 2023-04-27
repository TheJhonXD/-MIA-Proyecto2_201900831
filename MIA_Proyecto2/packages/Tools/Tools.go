package Tools

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

// Divide una cadena en un vector usando como delimitador los espacios
// Recibe una cadena de entrada
func Split(s string) []string {
	re := regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
	words := re.FindAllString(s, -1)
	return words
}

// Verifica si un archivo existe
// Recibe la ruta del archivo
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Elimina los comentarios de una cadena
// Recibe una cadena de entrada
func DeleteComments(s string) string {
	// re := regexp.MustCompile(`#.*`)
	re := regexp.MustCompile(`(?m)#.*$`)
	return re.ReplaceAllString(s, "")
}

// Remueve el nombre de archivo de la ruta y devuelve la ruta
func GetPath(p string) string {
	return filepath.Dir(p)
}

// Devuelve el nombre del archivo
func GetFileName(p string) string {
	return filepath.Base(p)
}

// Devuelve la extensión del archivo
func GetFileExt(p string) string {
	return filepath.Ext(p)
}

// Crea los directorios de la ruta ingresada si no existen
func CreateDir(path string) bool {
	if !Exists(GetPath(path)) {
		// fmt.Println(GetPath(path))
		if err := os.MkdirAll(path, 0777); err != nil {
			fmt.Println("Directorio creado")
		} else {
			fmt.Println("ERROR: No se pudo crear el directorio")
			return false
		}
	}
	return true
}

func GetPercentage(size float64, disk_size float64) int {
	return int(math.Round((size / disk_size) * 100))
}
