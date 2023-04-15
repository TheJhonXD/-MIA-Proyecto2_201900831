package Commands

import (
	"fmt"
	Disks "pack/packages/Disks"
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

func MKDISK(params []string) {
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
			fmt.Println("ERROR: el parametro \"" + param[0] + "\" no es valido.")
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
		if Disks.CreateDisk(pdm.path, pdm.s) {
			t := Structs.Time{}
			t.SetTime()
			m := Structs.MBR{Mbr_tamano: int32(pdm.s), Mbr_fecha_creacion: t, Mbr_dsk_signature: 0, Dsk_fit: byte(pdm.f[0]), Mbr_partition_1: Structs.RPV(), Mbr_partition_2: Structs.RPV(), Mbr_partition_3: Structs.RPV(), Mbr_partition_4: Structs.RPV()} //Crear el MBR
			//Agregar el MBR al disco
			Structs.AddMBR(pdm.path, m)
			//Imprimir el MBR
			Structs.ReadMBR(pdm.path)
		}
	}
}

func RMDISK(params []string) {
	pdm.ResetPDM()
	var param []string
	for i := 1; i < len(params); i++ {
		param = strings.Split(params[i], "=")
		if strings.ToLower(param[0]) == ">path" {
			pdm.path = strings.Trim(param[1], "\"")
		} else {
			fmt.Println("ERROR: el parametro \"" + param[0] + "\" no es valido.")
		}
	}

	if check_param_path() {
		Disks.DeleteDisk(pdm.path)
	}
}

func FDISK(params []string) {
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
			fmt.Println("ERROR: el parametro \"" + param[0] + "\" no es valido.")
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
		if Disks.CreatePart(pdm.path, p) { //Creo la particion y compruebo que todo haya salido bien
			Structs.ReadMBR(pdm.path)
		}
	}

}
