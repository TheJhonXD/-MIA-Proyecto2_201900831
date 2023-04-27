package Graphviz

import (
	"fmt"
	"os"
	"os/exec"
	"pack/packages/Disks"
	"pack/packages/Structs"
	"pack/packages/Tools"
	"strconv"
)

// Crea las columnas y filas de la tabla para crear el mbr, espacios libres y ocupados
// Recibe la ruta del disco y el mbr
func repDiskTR(path string, m Structs.MBR) string {
	graph := ""
	start := "\t\t\t<tr>\n"
	ss := Disks.BlockSize(path, m)
	var ssExt []Structs.SpaceSize
	graph += start
	graph += "\t\t\t\t<td colspan=\"1\" rowspan=\"12\">MBR</td>\n"

	ep := Disks.GetExtPart(path)

	contInUse := 0

	if ep.Part_s > 0 {
		ssExt = Disks.ExtBlockSize(path)
		for _, sExt := range ssExt {
			if sExt.In_use == 's' {
				contInUse++
			}
		}
	}
	tamExt := len(ssExt)

	for _, s := range ss {
		pct := Tools.GetPercentage(float64(s.Part_s), float64(m.Mbr_tamano))
		if s.In_use == 'n' {
			graph += "\t\t\t\t<td colspan=\"2\" rowspan=\"12\">Libre<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
		} else if s.In_use == 's' {
			if s.Type == 'p' {
				graph += "\t\t\t\t<td colspan=\"2\" rowspan=\"12\">Primaria<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
			} else if s.Type == 'e' {
				graph += "\t\t\t\t<td colspan=\"" + strconv.Itoa(tamExt+contInUse) + "\" rowspan=\"12\">Extendida<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
			}
		}
	}
	graph += "\t\t\t</tr>\n"

	if ep.Part_s > 0 && contInUse > 0 {
		graph += "\t\t\t<tr>\n"
		for _, sExt := range ssExt {
			pct := Tools.GetPercentage(float64(sExt.Part_s), float64(ep.Part_s))
			if sExt.In_use == 'n' {
				graph += "\t\t\t\t<td colspan=\"1\" rowspan=\"11\">Libre<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
			} else if sExt.In_use == 's' {
				graph += "\t\t\t\t<td colspan=\"1\" rowspan=\"11\">EBR</td>\n"
				graph += "\t\t\t\t<td colspan=\"1\" rowspan=\"11\">Logica<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
			}
		}
		graph += "\t\t\t</tr>\n"
	}
	return graph
}

// Crea el maquetado del inicio, intermedio, el final de la grafica y la retorna en un string
// Recibe como parametro la ruta del disco
func diskChart(path string) string {
	m := Structs.GetMBR(path)
	graph := ""
	start := "digraph {\n\tnode [ shape=none fontname=Arial fontsize=12];\n\n\tn4 [ label = <\n"
	st_table := "\t\t<table border=\"1\" color=\"dodgerblue3\">\n"
	end := "\t\t</table>\n\t> ];\n\n\t{rank=same n4};\n}"
	graph += start
	graph += st_table
	graph += repDiskTR(path, m)
	graph += end
	return graph
}

// Crea el archivo dot y la imagen del reporte de disco
// Recibe como parametro la ruta donde se guardará la imagen y el id del disco
func GetDiskGraph(path string, id string) {
	ruta := Tools.GetPath(path) + Tools.GetFileName(path) + ".dot"
	var stringCmd string

	if Disks.IdExists(id) {
		if Tools.CreateDir(Tools.GetPath(path)) {
			text := diskChart(Disks.GetDiskMtd(id).Path)
			myfile, err := os.Create(ruta)
			if err != nil {
				fmt.Println("ERROR: El archivo dot no pudo ser creado")
				return
			}

			myfile.Write([]byte(text))
			myfile.Close()
			stringCmd = "dot -T" + Tools.GetFileExt(path) + " " + ruta + " -o " + Tools.GetPath(ruta) + Tools.GetFileName(path)
			cmd := exec.Command(stringCmd)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("ERROR: ", err)
				return
			}
			fmt.Println(string(output))
			fmt.Println("Grafo creado correctamente")
		}
	} else {
		fmt.Println("ERROR: La particion no está montada")
	}
}
