package Filesystem

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"pack/packages/Disks"
	"pack/packages/Structs"
	"strconv"
	"strings"
	"unsafe"
)

func createInode(uid int32, gid int32, size int32, tipo int32) Structs.Inodo {
	i := Structs.RIV()
	i.I_uid = 0
	t := Structs.Time{}
	t.SetTime()
	i.I_atime = t
	i.I_ctime = t
	i.I_mtime = t
	i.I_type = tipo

	return i
}

func createFB(parent string, iparent int32, child string, ichild int32) Structs.FolderBlock {
	fb := Structs.RFBV()
	copy(fb.B_content[1].B_name[:], parent)
	fb.B_content[1].B_inodo = iparent
	copy(fb.B_content[0].B_name[:], child)
	fb.B_content[0].B_inodo = ichild

	return fb
}

// Retorna el inicio de la particion indicada
// Recibe la ruta del disco y el nombre de la partición
func getPartStart(path string, name string) int32 {
	m := Structs.GetMBR(path)
	if Disks.IsPrimPart(m, name) || Disks.IsExtPart(m, name) {
		return Disks.GetPartByName(path, name).Part_start
	} else if Disks.IsLogPart(path, name) {
		return (Disks.GetLogPartByName(path, name).Part_start + int32(unsafe.Sizeof(Structs.EBR{})) + 1)
	}
	return -1
}

func writeByteAtPosX(path string, pos int32, num byte) {
	myfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(int64(pos), 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	err = binary.Write(myfile, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
}

func getNxtFreeApD(apD [16]byte) int {
	for i := 0; i < 16; i++ {
		if int(apD[i]) < 0 {
			return i
		}
	}
	return -1
}

func getNxtFreeApFB(content Structs.FolderBlock) int {
	for i := 0; i < 4; i++ {
		if content.B_content[i].B_inodo < 0 {
			return i
		}
	}
	return -1
}

func linkInodeToBlock(path string, start int32, id_bloque int32, i Structs.Inodo) {
	indice := getNxtFreeApD(i.I_block)
	i.I_block[indice] = byte(id_bloque)
	Structs.AddInodo(path, start, i)
}

func linkInodeToFileBlock(path string, start int32, id_bloque int32, i Structs.Inodo) {
	indice := getNxtFreeApD(i.I_block)
	i.I_block[indice] = byte(id_bloque)
	Structs.AddInodo(path, start, i)
}

func CreateRoot(id string) bool {
	md := Disks.GetDiskMtd(id)
	i := createInode(0, 0, 0, 0)
	fb := createFB("..", 0, ".", 0)
	start := getPartStart(md.Path, md.Name)
	if start != -1 {
		sb := Structs.GetSuperBlock(md.Path, start)
		//Añado el inodo 0 (root)
		Structs.AddInodo(md.Path, sb.S_inode_start, i)
		//Añado el folderblock 0 (root)
		Structs.AddFolderBlock(md.Path, sb.S_block_start, fb)
		//Ocupo el espacio en el bitmap de inodos
		writeByteAtPosX(md.Path, sb.S_bm_inode_start, 1)
		//Ocupo el espacio en el bitmap de bloques
		writeByteAtPosX(md.Path, sb.S_bm_block_start, 1)
		//Enlazo el inodo 0 con el folderblock 0
		linkInodeToBlock(md.Path, sb.S_inode_start, 0, i)
	}
	return false
}

func getNxtFreePosBmInode(path string, startPart int32) int32 {
	sb := Structs.GetSuperBlock(path, startPart)
	myfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	var bPrev byte = 1
	var b byte
	for i := sb.S_bm_block_start - 1; i >= sb.S_bm_inode_start; i-- {
		_, _ = myfile.Seek(int64(i), 0)
		err = binary.Read(myfile, binary.LittleEndian, &b)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		if b == 1 && bPrev == 0 {
			return int32(bPrev)
		}
		bPrev = b
	}
	myfile.Close()
	return 0
}

func getNxtFreePosBmBlock(path string, startPart int32) int32 {
	sb := Structs.GetSuperBlock(path, startPart)
	myfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	var bPrev byte = 1
	var b byte
	for i := sb.S_inode_start - 1; i >= sb.S_bm_block_start; i-- {
		_, _ = myfile.Seek(int64(i), 0)
		err = binary.Read(myfile, binary.LittleEndian, &b)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		if b == 1 && bPrev == 0 {
			return int32(bPrev)
		}
		bPrev = b
	}
	myfile.Close()
	return 0
}

func getNxtPosInode(path string, name string) int32 {
	start := getPartStart(path, name)
	sb := Structs.GetSuperBlock(path, start)
	posRelative := getNxtFreePosBmInode(path, start) - sb.S_bm_inode_start
	return (posRelative * int32(unsafe.Sizeof(Structs.Inodo{}))) + sb.S_inode_start + 1
}

func getNxtPosBlock(path string, name string) int32 {
	start := getPartStart(path, name)
	sb := Structs.GetSuperBlock(path, start)
	posRelative := getNxtFreePosBmBlock(path, start) - sb.S_bm_block_start
	return (posRelative * int32(unsafe.Sizeof(Structs.Inodo{}))) + sb.S_block_start + 1
}

func splitStringInto64(cadena string) [][64]byte {
	var array [64]byte
	var result [][64]byte
	c := cadena
	flag := true
	for flag {
		if len(c) > 64 {
			copy(array[:], cadena[:64])
			c = cadena[64:]
			result = append(result, array)
		} else if len(c) > 0 {
			copy(array[:], bytes.Repeat([]byte{0}, len(array)))
			copy(array[:], c[:])
			result = append(result, array)
			c = ""
			flag = false
		} else {
			flag = false
		}
	}

	return result
}

func GetFileBlockById(path string, start int32, idBlock int32) Structs.FileBlock {
	sb := Structs.GetSuperBlock(path, start)
	// posActual := getNxtFreePosBmBlock(path, start) - sb.S_bm_block_start
	var pos int32
	if idBlock > 0 {
		pos = ((idBlock - 1) * int32(unsafe.Sizeof(Structs.FileBlock{}))) + sb.S_block_start + 1
	} else {
		pos = ((idBlock - 1) * int32(unsafe.Sizeof(Structs.FileBlock{}))) + sb.S_block_start
	}
	fb := Structs.GetFileBlock(path, pos)
	return fb
}

func GetFileBlockPosById(path string, start int32, idBlock int32) int32 {
	sb := Structs.GetSuperBlock(path, start)
	// posActual := getNxtFreePosBmBlock(path, start) - sb.S_bm_block_start
	var pos int32
	if idBlock > 0 {
		pos = ((idBlock - 1) * int32(unsafe.Sizeof(Structs.FileBlock{}))) + sb.S_block_start + 1
	} else {
		pos = ((idBlock - 1) * int32(unsafe.Sizeof(Structs.FileBlock{}))) + sb.S_block_start
	}
	return pos
}

func joinTextFileBlock(path string, start int32, inodo Structs.Inodo) string {
	// sb := Structs.GetSuperBlock(path, start)
	text := ""
	var fb Structs.FileBlock
	for _, inode := range inodo.I_block {
		if inode > 0 {
			fb = GetFileBlockById(path, start, int32(inode))
			text += strings.TrimRight(string(fb.B_content[:]), "\x00")
		}
	}
	return text
}

func addTextToFileBlocks(path string, start int32, matrix [][64]byte, inodo Structs.Inodo) {
	// sb := Structs.GetSuperBlock(path, start)
	var fb Structs.FileBlock
	for i := 0; i < len(matrix); i++ {
		fb = GetFileBlockById(path, start, int32(inodo.I_block[i]))
		fb.B_content = matrix[i]
		Structs.AddFileBlock(path, GetFileBlockPosById(path, start, int32(inodo.I_block[i])), fb)
	}
}

// ? Start es el inicio de la particion
func createGroup(path string, start int32, grupo Structs.Group, i Structs.Inodo) {
	textPrev := joinTextFileBlock(path, start, i)
	text := strconv.Itoa(int(grupo.GID)) + ", G, " + strings.TrimRight(string(grupo.Grp[:]), "\x00") + "\n"
	fmt.Println("TEXTOBLOQUE:", textPrev)
	textPrev += text
	arrays := splitStringInto64(textPrev)
	addTextToFileBlocks(path, start, arrays, i)
}

func crearBloquesDeArchivos(path string, name string, startInodo int32, inodo Structs.Inodo) {
	var b Structs.FileBlock
	start := getPartStart(path, name)
	sb := Structs.GetSuperBlock(path, start)
	for i := 0; i < 16; i++ {
		Structs.AddFileBlock(path, getNxtPosBlock(path, name), b)
		num := getNxtFreePosBmBlock(path, start) - sb.S_bm_block_start
		writeByteAtPosX(path, getNxtFreePosBmBlock(path, start), 1)
		//?Enlazo el inodo con el bloque
		linkInodeToFileBlock(path, startInodo, num, inodo)
	}
}

func CreateUsersFile(id string) bool {
	md := Disks.GetDiskMtd(id)
	start := getPartStart(md.Path, md.Name)
	sb := Structs.GetSuperBlock(md.Path, start)
	fb := Structs.GetFolderBlock(md.Path, sb.S_block_start)
	indice := getNxtFreeApFB(fb)
	var array [12]byte
	copy(array[:], "users.txt")
	fb.B_content[indice].B_name = array
	fb.B_content[indice].B_inodo = 1
	Structs.AddFolderBlock(md.Path, sb.S_block_start, fb)

	i := createInode(1, 1, 0, 0)
	nxtPosInodo := getNxtPosInode(md.Path, md.Name)
	nxtPosBmInodo := getNxtFreePosBmInode(md.Path, start)
	Structs.AddInodo(md.Path, nxtPosInodo, i)
	writeByteAtPosX(md.Path, nxtPosBmInodo, 1)
	crearBloquesDeArchivos(md.Path, md.Name, nxtPosInodo, i)
	var myarr [10]byte
	copy(myarr[:], "root")
	grp := Structs.Group{GID: 1, Grp: myarr}
	createGroup(md.Path, start, grp, i)
	return true
}
