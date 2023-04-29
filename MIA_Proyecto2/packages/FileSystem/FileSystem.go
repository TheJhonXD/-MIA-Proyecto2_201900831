package Filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"pack/packages/Disks"
	"pack/packages/Structs"
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

func getNxtFreeApInodo(content Structs.FolderBlock) int {
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

func linkInodeToFile(path string, start int32, id_bloque int32, i Structs.Inodo) {
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

func CreateUsersFile(id string) bool {
	md := Disks.GetDiskMtd(id)
	start := getPartStart(md.Path, md.Name)
	fb := Structs.GetFolderBlock(md.Path, start)
	indice := getNxtFreeApInodo(fb)
	i := createInode(1, 1, 0, 0)
}
