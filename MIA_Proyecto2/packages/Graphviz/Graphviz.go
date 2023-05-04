package Graphviz

import (
	"fmt"
	"pack/packages/Disks"
	Filesystem "pack/packages/FileSystem"
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
	graph += createRowSBRpt("sb_arbol_virtual_count", strconv.Itoa(int(sb.S_inodes_count+sb.S_blocks_count)))
	graph += createRowSBRpt("sb_detalle_directorio_count", strconv.Itoa(int(sb.S_inodes_count+sb.S_blocks_count)))
	graph += createRowSBRpt("sb_inodos_count", strconv.Itoa(int(sb.S_inodes_count-1)))
	graph += createRowSBRpt("sb_bloques_count", strconv.Itoa(int(sb.S_blocks_count+2)))
	graph += createRowSBRpt("sb_arbol_virtual_free", strconv.Itoa(int(sb.S_bm_block_start/2)))
	graph += createRowSBRpt("sb_detalle_directorio_free", strconv.Itoa(int(sb.S_bm_block_start/2+sb.S_bm_inode_start)))
	graph += createRowSBRpt("sb_inodos_free", strconv.Itoa(int(sb.S_bm_block_start)))
	graph += createRowSBRpt("sb_bloques_free", strconv.Itoa(int(sb.S_bm_block_start+sb.S_bm_block_start)))
	fecha := strconv.Itoa(int(sb.S_mtime.Year)) + "-0" + strconv.Itoa(int(sb.S_mtime.Month)) + "-" + strconv.Itoa(int(sb.S_mtime.Day)) + " " + strconv.Itoa(int(sb.S_mtime.Hour)) + ":" + strconv.Itoa(int(sb.S_mtime.Min)) + ":" + strconv.Itoa(int(sb.S_mtime.Sec))
	graph += createRowSBRpt("sb_date_creacion", fecha)
	fecha = strconv.Itoa(int(sb.S_mtime.Year)) + "-0" + strconv.Itoa(int(sb.S_mtime.Month)) + "-" + strconv.Itoa(int(sb.S_mtime.Day)) + " " + strconv.Itoa(int(sb.S_mtime.Hour)) + ":" + strconv.Itoa(int(sb.S_mtime.Min)) + ":" + strconv.Itoa(int(sb.S_mtime.Sec))
	// fecha = string(sb.S_mtime.Year) + "-" + string(sb.S_mtime.Month) + "-" + string(sb.S_mtime.Day) + " " + string(sb.S_mtime.Hour) + ":" + string(sb.S_mtime.Min) + ":" + string(sb.S_mtime.Sec)
	graph += createRowSBRpt("sb_date_ultimo_montaje", fecha)
	graph += createRowSBRpt("sb_montajes_count", "1")
	graph += createRowSBRpt("sb_ap_bitmap_arbol_directorio", strconv.Itoa(int(sb.S_bm_inode_start)))
	graph += createRowSBRpt("sb_ap_arbol_directorio", strconv.Itoa(int(sb.S_inode_start)))
	graph += createRowSBRpt("sb_ap_bitmap_detalle_directorio", strconv.Itoa(int(sb.S_bm_block_start)))
	graph += createRowSBRpt("sb_ap_detalle_directorio", strconv.Itoa(int(sb.S_block_start)))
	graph += createRowSBRpt("sb_ap_bitmap_tabla_inodo", strconv.Itoa(int(sb.S_bm_inode_start)))
	graph += createRowSBRpt("sb_ap_tabla_inodo", strconv.Itoa(int(sb.S_inode_start)))
	graph += createRowSBRpt("sb_ap_bitmap_bloques", strconv.Itoa(int(sb.S_bm_block_start)))
	graph += createRowSBRpt("sb_ap_bloques", strconv.Itoa(int(sb.S_block_start)))
	graph += createRowSBRpt("sb_ap_log", "0")
	graph += createRowSBRpt("sb_size_struct_arbol_directorio", strconv.Itoa(int(sb.S_block_size)))
	graph += createRowSBRpt("sb_size_struct_detalle_directorio", strconv.Itoa(int(sb.S_inode_size)))
	graph += createRowSBRpt("sb_size_struct_inodo", strconv.Itoa(int(sb.S_inode_size)))
	graph += createRowSBRpt("sb_size_struct_bloque", strconv.Itoa(int(sb.S_block_size)))
	graph += createRowSBRpt("sb_first_free_bit_arbol_directorio", strconv.Itoa(int(sb.S_bm_block_start)))
	graph += createRowSBRpt("sb_first_free_bit_detalle_directorio", "0")
	graph += createRowSBRpt("sb_first_free_bit_tabla_inodo", strconv.Itoa(int(sb.S_bm_block_start)))
	graph += createRowSBRpt("sb_first_free_bit_bloques", strconv.Itoa(int(sb.S_block_start+(sb.S_bm_inode_start*5))))
	graph += createRowSBRpt("sb_magic_num", string("0xEF53"))

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

func createRowTreeRpt(param string, valor string, port string) string {
	graph := "<tr><td>" + param + "</td><td port='" + port + "'>" + valor + "</td></tr>"

	return graph
}

func createFileBlockTreeRpt(id string, valor string, port string) string {
	graph := "fileblock" + id + " [shape=plaintext label=< <table border='0' cellborder='1' cellspacing='0' bgcolor='lightcoral'>\n\t<tr><td port='" + port + "'>Bloque " + id + "</td></tr>\n\t"
	graph += "<tr><td>" + valor + "</td></tr>"
	graph += "</table> >];"

	return graph
}

func createInodeRootTreeRpt(fecha string) string {
	graph := "root [shape=plaintext label=< <table border='0' cellborder='1' cellspacing='0' bgcolor='lightblue'>\n\t<tr><td colspan='2'>INODO 0</td></tr>\n\t"
	graph += createRowTreeRpt("UID", "1", "")
	graph += createRowTreeRpt("GID", "1", "")
	graph += createRowTreeRpt("Fecha creacion", fecha, "")
	graph += createRowTreeRpt("Tipo", "1", "")
	graph += createRowTreeRpt("Size", "0", "")
	graph += createRowTreeRpt("Ap0", "0", "f0")
	graph += createRowTreeRpt("Ap1", "-1", "")
	graph += createRowTreeRpt("Ap2", "-1", "")
	graph += createRowTreeRpt("Ap3", "-1", "")
	graph += createRowTreeRpt("Ap4", "-1", "")
	graph += createRowTreeRpt("Ap5", "-1", "")
	graph += createRowTreeRpt("Ap6", "-1", "")
	graph += createRowTreeRpt("Ap7", "-1", "")
	graph += createRowTreeRpt("Ap8", "-1", "")
	graph += createRowTreeRpt("Ap9", "-1", "")
	graph += createRowTreeRpt("Ap10", "-1", "")
	graph += createRowTreeRpt("Ap11", "-1", "")
	graph += createRowTreeRpt("Ap12", "-1", "")
	graph += createRowTreeRpt("Ap13", "-1", "")
	graph += createRowTreeRpt("Ap14", "-1", "")
	graph += createRowTreeRpt("Ap15", "-1", "")

	graph += "</table> >];"
	return graph
}

func unirNodePort(nodo1 string, nodo2 string, port1 string, port2 string) string {
	text := ""
	text += nodo1 + ":" + port1 + " -> " + nodo2 + ":" + port2 + ";\n"
	return text
}

func createInodeTreeRpt(uid string, gid string, fecha string, tipo string, contenido string) (string, string, string) {
	graph := "inodeUser [shape=plaintext label=< <table border='0' cellborder='1' cellspacing='0' bgcolor='lightblue'>\n\t<tr><td port='b1' colspan='2'>INODO 0</td></tr>\n\t"
	graph += createRowTreeRpt("UID", uid, "")
	graph += createRowTreeRpt("GID", gid, "")
	graph += createRowTreeRpt("Fecha creacion", fecha, "")
	graph += createRowTreeRpt("Tipo", tipo, "")
	graph += createRowTreeRpt("Size", strconv.Itoa(len(contenido))+" bytes", "")
	// tam := int(math.Round(float64(len(contenido)) / float64(64)))
	graphFB := ""
	unirNodo := ""
	c := contenido
	// fmt.Println("CADENA:", c)
	for i := 0; i < 16; i++ {
		if len(c) >= 64 {
			graph += createRowTreeRpt("Ap"+strconv.Itoa(i), strconv.Itoa(i+1), "fb"+strconv.Itoa(i))
			graphFB += createFileBlockTreeRpt(strconv.Itoa(i+1), c[:64], "fb"+strconv.Itoa(i))
			unirNodo += unirNodePort("inodeUser", "fileblock"+strconv.Itoa(i+1), "fb"+strconv.Itoa(i), "fb"+strconv.Itoa(i))
			c = c[64:]
		} else if len(c) > 0 {
			graph += createRowTreeRpt("Ap"+strconv.Itoa(i), strconv.Itoa(i+1), "fb"+strconv.Itoa(i))
			graphFB += createFileBlockTreeRpt(strconv.Itoa(i+1), c, "fb"+strconv.Itoa(i))
			unirNodo += unirNodePort("inodeUser", "fileblock"+strconv.Itoa(i+1), "fb"+strconv.Itoa(i), "fb"+strconv.Itoa(i))
			c = ""
		} else {
			graph += createRowTreeRpt("Ap"+strconv.Itoa(i), "-1", "")
		}
	}

	/* graph += createRowTreeRpt("Ap0", "-1", "fb0")
	graph += createRowTreeRpt("Ap1", "-1", "fb1")
	graph += createRowTreeRpt("Ap2", "-1", "fb2")
	graph += createRowTreeRpt("Ap3", "-1", "fb3")
	graph += createRowTreeRpt("Ap4", "-1", "fb4")
	graph += createRowTreeRpt("Ap5", "-1", "fb5")
	graph += createRowTreeRpt("Ap6", "-1", "fb6")
	graph += createRowTreeRpt("Ap7", "-1", "fb7")
	graph += createRowTreeRpt("Ap8", "-1", "fb8")
	graph += createRowTreeRpt("Ap9", "-1", "fb9")
	graph += createRowTreeRpt("Ap10", "-1", "fb10")
	graph += createRowTreeRpt("Ap11", "-1", "fb11")
	graph += createRowTreeRpt("Ap12", "-1", "fb12")
	graph += createRowTreeRpt("Ap13", "-1", "fb13")
	graph += createRowTreeRpt("Ap14", "-1", "fb14")
	graph += createRowTreeRpt("Ap15", "-1", "fb15") */

	graph += "</table> >];"
	return graph, graphFB, unirNodo
}

func createFolderBlockRootTreeRpt() string {
	graph := "folderblock1 [shape=plaintext label=< <table border='0' cellborder='1' cellspacing='0' bgcolor='lightcoral'>\n\t<tr><td port='b0' colspan='2'>Bloque 0</td></tr>\n\t"
	graph += createRowTreeRpt(".", "0", "")
	graph += createRowTreeRpt("..", "0", "")
	graph += createRowTreeRpt("users.txt", "8", "f1")
	graph += createRowTreeRpt("", "-1", "")
	graph += "</table> >];"

	return graph
}

func CreateTreeReport(id string) string {
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

	grupos := Filesystem.GetAllGroupsGraph(id)
	usuarios := Filesystem.GetAllUsersGraph(id)

	graph := "digraph G {\n"
	fecha := strconv.Itoa(int(sb.S_mtime.Year)) + "/" + strconv.Itoa(int(sb.S_mtime.Month)) + "/" + strconv.Itoa(int(sb.S_mtime.Day)) + " " + strconv.Itoa(int(sb.S_mtime.Hour)) + ":" + strconv.Itoa(int(sb.S_mtime.Min)) + ":" + strconv.Itoa(int(sb.S_mtime.Sec))
	graph += createInodeRootTreeRpt(fecha)
	graph += createFolderBlockRootTreeRpt()
	graphAux, graphFB, unirNodo := createInodeTreeRpt("1", "1", fecha, "0", grupos+usuarios)
	graph += graphAux
	graph += "\n\troot:f0   -> folderblock1:b0\n"
	graph += "\n\tfolderblock1:f1  -> inodeUser:b1\n"
	graph += graphFB
	graph += unirNodo
	graph += "\n\trankdir=LR;\n}"

	return graph
}
