package Commands

import (
	"fmt"
	Disks "pack/packages/Disks"
	Filesystem "pack/packages/FileSystem"
	"pack/packages/Structs"
	"strconv"
	"strings"
)

type PDM struct {
	s    int
	f    string
	u    string
	path string
	t    string
	del  string
	name string
	add  int
	id   string
	ruta string
	fs   string
	usr  string
	pwd  string
	grp  string
}

func (p *PDM) ResetPDM() {
	p.s = -1
	p.f = "-1"
	p.u = "-1"
	p.path = "-1"
	p.t = "-1"
	p.del = "-1"
	p.name = "-1"
	p.add = 0
	p.id = "-1"
	p.ruta = "-1"
	p.fs = "-1"
	p.usr = "-1"
	p.pwd = "-1"
	p.grp = "-1"
}

var pdm PDM

// Comprueba si el parametro size es valido
func check_param_s() bool {
	if pdm.s > 0 {
		return true
	} else {
		fmt.Println("ERROR: Tamanio no valido")
	}
	return false
}

func check_param_path() bool {
	return pdm.path != "-1"
}

func check_param_f() bool {
	if pdm.f == "-1" {
		pdm.f = "ff"
		return true
	}
	if pdm.f == "bf" || pdm.f == "ff" || pdm.f == "wf" {
		return true
	} else {
		fmt.Println("ERROR: Fit no valido")
	}
	return false
}

func check_param_f_fdisk() bool {
	if pdm.f == "-1" {
		pdm.f = "wf"
		return true
	}
	if pdm.f == "bf" || pdm.f == "ff" || pdm.f == "wf" {
		return true
	} else {
		fmt.Println("ERROR: Fit no valido")
	}
	return false
}

func check_param_u() bool {
	if pdm.u == "-1" {
		pdm.u = "m"
		return true
	}
	if pdm.u == "m" || pdm.u == "k" {
		return true
	} else {
		fmt.Println("ERROR: Unit no valido")
	}
	return false
}

func check_param_u_fdisk() bool {
	if pdm.u == "-1" {
		pdm.u = "k"
		return true
	}
	if pdm.u == "m" || pdm.u == "k" || pdm.u == "b" {
		return true
	} else {
		fmt.Println("ERROR: Unit no valido")
	}
	return false
}

func check_param_t() bool {
	if pdm.t == "-1" {
		pdm.t = "p"
		return true
	}
	if pdm.t == "p" || pdm.t == "e" || pdm.t == "l" {
		return true
	} else {
		fmt.Println("ERROR: Type no valido")
	}
	return false
}

func check_param_name() bool {
	return pdm.name != "-1"
}

func check_param_id() bool {
	return pdm.id != "-1"
}

func check_param_t_mkfs() bool {
	if pdm.t == "full" || pdm.t == "-1" {
		return true
	}
	fmt.Println("ERROR: Type no valido")
	return false
}

func MKDISK(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">size" {
			pdm.s, _ = strconv.Atoi(param[1])
		} else if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">unit" {
			pdm.u = strings.ToLower(param[1])
		} else if strings.ToLower(param[0]) == ">fit" {
			pdm.f = strings.ToLower(param[1])
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	//Comprobar si son valores validos
	if check_param_s() && check_param_f() && check_param_u() && check_param_path() {
		//Convertir el tamaño a bytes
		if pdm.u == "k" {
			pdm.s *= 1024
		} else {
			pdm.s *= 1024 * 1024
		}
		//Crear el disco y comprobar su creacion
		flag, message := Disks.CreateDisk(pdm.path, pdm.s)
		if flag {
			t := Structs.Time{}
			t.SetTime()
			m := Structs.MBR{Mbr_tamano: int32(pdm.s), Mbr_fecha_creacion: t, Mbr_dsk_signature: 0, Dsk_fit: byte(pdm.f[0]), Mbr_partition_1: Structs.RPV(), Mbr_partition_2: Structs.RPV(), Mbr_partition_3: Structs.RPV(), Mbr_partition_4: Structs.RPV()} //Crear el MBR
			//Agregar el MBR al disco
			Structs.AddMBR(pdm.path, m)
			//Imprimir el MBR
			Structs.ReadMBR(pdm.path)
		}
		messages += message
	}
	return messages
}

func RMDISK(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if check_param_path() {
		_, message := Disks.DeleteDisk(pdm.path)
		messages += message
	}
	return messages
}

func FDISK(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">size" {
			pdm.s, _ = strconv.Atoi(param[1])
		} else if strings.ToLower(param[0]) == ">unit" {
			pdm.u = strings.ToLower(param[1])
		} else if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">type" {
			pdm.t = strings.ToLower(param[1])
		} else if strings.ToLower(param[0]) == ">fit" {
			pdm.f = strings.ToLower(param[1])
		} else if strings.ToLower(param[0]) == ">name" {
			pdm.name = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	//Comprobar los parametros
	if check_param_s() && check_param_u_fdisk() && check_param_path() && check_param_t() && check_param_f_fdisk() && check_param_name() {
		//Convertir el tamaño a bytes
		if pdm.u == "k" {
			pdm.s *= 1024
		} else if pdm.u == "m" {
			pdm.s *= 1024 * 1024
		}

		//Crear una variable particion para guardar los datos, excepto el inicio de particion y el estado
		p := Structs.Partition{Part_status: '0', Part_type: byte(pdm.t[0]), Part_fit: byte(pdm.f[0]), Part_start: int32(-1), Part_s: int32(pdm.s)}
		copy(p.Part_name[:], pdm.name)
		flag, message := Disks.CreatePart(pdm.path, p)
		if flag { //Creo la particion y compruebo que todo haya salido bien
			Structs.ReadMBR(pdm.path)
		}
		messages += message
	}
	return messages
}

func MOUNT(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">name" {
			pdm.name = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if check_param_path() && check_param_name() {
		_, message := Disks.MountDisk(pdm.path, pdm.name)
		messages += message
		/* mds := Disks.GetDisksMounted()
		for _, md := range mds {
			fmt.Println("->:", md.Id)
		} */
	}
	return messages
}

func MKFS(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">type" {
			pdm.t = strings.ToLower(param[1])
		} else if strings.ToLower(param[0]) == ">id" {
			pdm.id = strings.ToLower(param[1])
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if check_param_t_mkfs() && check_param_id() {
		flag, message := Disks.MakeFileSystem(pdm.id)
		messages += message
		if flag {
			messages += "Usuario y Grupo Root creado" + "\n"
			Filesystem.InitGrpNUsr(pdm.id)
		}
	}
	return messages
}

func LOGIN(params []string) {
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">user" {
			pdm.usr = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">pwd" {
			pdm.pwd = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">id" {
			pdm.id = strings.ToLower(param[1])
		} else {
			fmt.Println("ERROR: el parametro \"" + param[0] + "\" no es valido.")
		}
	}

	if check_param_id() && pdm.usr != "-1" && pdm.pwd != "-1" {
		Filesystem.Login(pdm.usr, pdm.pwd, pdm.id)
	}
}

func LOGOUT() {
	pdm.ResetPDM()
	Filesystem.Logout()
}

func MKGRP(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">name" {
			pdm.name = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if check_param_name() {
		_, message := Filesystem.MakeGroup(pdm.name)
		messages += message
	}

	return messages
}

func RMGRP(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">name" {
			pdm.name = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if check_param_name() {
		_, message := Filesystem.RemoveGroup(pdm.name)
		messages += message
	}

	return messages
}

func MKUSR(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">user" {
			pdm.usr = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">pwd" {
			pdm.pwd = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">grp" {
			pdm.grp = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	if pdm.usr != "-1" && pdm.pwd != "-1" && pdm.grp != "-1" {
		_, message := Filesystem.MakeUser(pdm.usr, pdm.pwd, pdm.grp)
		messages += message
	}

	return messages
}

func RMUSR(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">user" {
			pdm.usr = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}
	//!Falta poner lo de los usuarios en el analyzer
	if pdm.usr != "-1" {
		_, message := Filesystem.RemoveUser(pdm.usr)
		messages += message
	}

	return messages
}

func REP(params []string) string {
	messages := ""
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">name" {
			pdm.name = strings.ToLower(strings.Trim(param[1], "\""))
		} else if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else if strings.ToLower(param[0]) == ">id" {
			pdm.id = param[1]
		} else if strings.ToLower(param[0]) == ">ruta" {
			pdm.ruta = strings.Trim(param[1], "\"")
		} else {
			messages += "ERROR: el parametro \"" + param[0] + "\" no es valido." + "\n"
		}
	}

	/* if pdm.name == "disk" && check_param_path() && check_param_id() {
		messages += Graphviz.GetDiskGraph(pdm.path, pdm.id)
	} */
	return messages
}
