package Graphviz

import (
	"fmt"
	"pack/packages/Disks"
	"pack/packages/Structs"
	"pack/packages/Tools"
	"strconv"
	"unsafe"
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
				graph += "\t\t\t\t<td colspan=\"" + strconv.Itoa(tamExt+contInUse) + "\">Extendida<br/><font point-size=\"8\">" + strconv.Itoa(pct) + "% del disco</font></td>\n"
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
func GetDiskGraph(id string) string {
	text := ""
	// messages := ""
	// ruta := Tools.GetPath(path) + Tools.GetFileName(path) + ".dot"
	// fmt.Println(Disks.GetDisksMounted())
	if Disks.IdExists(id) {
		text = diskChart(Disks.GetDiskMtd(id).Path)
		/* flag, message := Tools.CreateDir(Tools.GetPath(path))
		messages += message
		if flag {
			text := diskChart(Disks.GetDiskMtd(id).Path)
			myfile, err := os.Create(ruta)
			if err != nil {
				messages += "ERROR: El archivo dot no pudo ser creado" + "\n"
				return messages
			}

			myfile.Write([]byte(text))
			myfile.Close()

			cmd := exec.Command("dot", "-T"+Tools.GetFileExt(path)[1:], ruta, "-o", Tools.GetPath(ruta)+Tools.GetFileName(path)+Tools.GetFileExt(path))
			err = cmd.Run()
			if err != nil {
				messages += "ERROR: " + string(err.Error()) + "\n"
				return messages
			}
			messages += "Grafo creado correctamente" + "\n"
		} */
	} else {
		fmt.Println("El id no existe")
	}
	return text
}

/**************************************** REPORTE SB ****************************************/
//Retorna el maquetado de una fila de informacion que se le indique
//Recibe como parametro el nombre del parametro y su valor
func createRowSBRpt(param string, valor string) string {
	graph := "\t\t\t<tr>\n"
	graph += "\t\t\t\t<td border=\"0\" width=\"250\">" + param + "</td>\n"
	graph += "\t\t\t\t<td border=\"0\" width=\"250\">" + valor + "</td>\n"
	graph += "\t\t\t</tr>\n"
	return graph
}

func getTableSBInfo(id string) string {
	graph := ""
	md := Disks.GetDiskMtd(id)
	m := Structs.GetMBR(md.Path)
	sb := Structs.RSBV()

	if Disks.IsPrimPart(m, md.Name) || Disks.IsExtPart(m, md.Name) {
		p := Disks.GetPartByName(md.Path, md.Name)
		sb = Structs.GetSuperBlock(md.Path, p.Part_start)
	} else if Disks.IsLogPart(md.Path, md.Name) {
		e := Disks.GetLogPartByName(md.Path, md.Name)
		sb = Structs.GetSuperBlock(md.Path, e.Part_start+int32(unsafe.Sizeof(Structs.EBR{}))+1)
	}

	graph += createRowSBRpt("sb_nombre_hd", Tools.GetFileName(md.Path))
	graph += createRowSBRpt("sb_arbol_virtual_count", "0")
	graph += createRowSBRpt("sb_detalle_directorio_count", "0")
	graph += createRowSBRpt("sb_inodos_count", string(sb.S_inodes_count))
	graph += createRowSBRpt("sb_bloques_count", string(sb.S_blocks_count))
	graph += createRowSBRpt("sb_arbol_virtual_free", "0")
	graph += createRowSBRpt("sb_detalle_directorio_free", "0")
	graph += createRowSBRpt("sb_inodos_free", string(sb.S_free_inodes_count))
	graph += createRowSBRpt("sb_bloques_free", string(sb.S_free_blocks_count))
	fecha := string(sb.S_mtime.Year) + "-" + string(sb.S_mtime.Month) + "-" + string(sb.S_mtime.Day) + " " + string(sb.S_mtime.Hour) + ":" + string(sb.S_mtime.Min) + ":" + string(sb.S_mtime.Sec)
	graph += createRowSBRpt("sb_date_creacion", fecha)
	fecha = string(sb.S_mtime.Year) + "-" + string(sb.S_mtime.Month) + "-" + string(sb.S_mtime.Day) + " " + string(sb.S_mtime.Hour) + ":" + string(sb.S_mtime.Min) + ":" + string(sb.S_mtime.Sec)
	graph += createRowSBRpt("sb_date_ultimo_montaje", fecha)
	graph += createRowSBRpt("sb_montajes_count", "0")
	graph += createRowSBRpt("sb_ap_bitmap_arbol_directorio", string(sb.S_bm_inode_start))
	graph += createRowSBRpt("sb_ap_arbol_directorio", string(sb.S_inode_start))
	graph += createRowSBRpt("sb_ap_bitmap_detalle_directorio", string(sb.S_bm_block_start))
	graph += createRowSBRpt("sb_ap_detalle_directorio", string(sb.S_block_start))
	graph += createRowSBRpt("sb_ap_bitmap_tabla_inodo", string(sb.S_bm_inode_start))
	graph += createRowSBRpt("sb_ap_tabla_inodo", string(sb.S_inode_start))
	graph += createRowSBRpt("sb_ap_bitmap_bloques", string(sb.S_bm_block_start))
	graph += createRowSBRpt("sb_ap_bloques", string(sb.S_block_start))
	graph += createRowSBRpt("sb_ap_log", "0")
	graph += createRowSBRpt("sb_size_struct_arbol_directorio", "0")
	graph += createRowSBRpt("sb_size_struct_detalle_directorio", "0")
	graph += createRowSBRpt("sb_size_struct_inodo", "0")
	graph += createRowSBRpt("sb_size_struct_bloque", "0")
	graph += createRowSBRpt("sb_first_free_bit_arbol_directorio", "0")
	graph += createRowSBRpt("sb_first_free_bit_detalle_directorio", "0")
	graph += createRowSBRpt("sb_first_free_bit_tabla_inodo", "0")
	graph += createRowSBRpt("sb_first_free_bit_bloques", "0")
	graph += createRowSBRpt("sb_magic_num", string(sb.S_magic))

	return graph
}

func getSBReport(id string) string {
	graph := ""
	start := "digraph {\n\tnode [ shape=none fontname=Arial fontsize=12];\n\n\tn1 [ label = <\n"
	st_table := "\t\t<table border=\"2\" cellspacing=\"0\" cellpadding=\"10\">\n"
	end := "\t\t</table>\n\t> ];\n\n\t{rank=same n1};\n}"

	graph += start
	graph += st_table
	graph += "\t\t\t<tr><td colspan=\"2\" bgcolor=\"royalblue\" border=\"1\" align=\"left\" width=\"500\" color=\"white\"><b><font point-size=\"16\" color=\"white\">REPORTE DE SUPERBLOQUE</font></b></td></tr>\n"
	graph += getTableSBInfo(id)
	graph += end
	return graph
}

// Crea el archivo dot y la imagen del reporte del Super Bloque
// Recibe como parametro la ruta donde se guardará la imagen y el id del disco
func CreateSBReport(id string) string {
	text := ""
	if Disks.IdExists(id) {
		text = getSBReport(id)
	} else {
		fmt.Println("ERROR: La particion no está montada")
	}
	return text
}
