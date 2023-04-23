package Disks

import (
	"fmt"
	"math"
	"os"
	"pack/packages/Structs"
	"pack/packages/Tools"
	"sort"
	"strconv"
	"unsafe"
)

var mds []Structs.MountedDisk

func CreateDisk(path string, tam int) bool {
	if !Tools.Exists(path) {
		//Creo los directorios si no existen
		if Tools.CreateDir(path) {
			myfile, err := os.Create(path)
			if err != nil {
				fmt.Println("ERROR: No se pudo crear el disco")
			}
			// defer myfile.Close()
			/* Lleno el archivo con caracteres nulos para simular el tamaño */
			var buffer [1024]byte
			for i := 0; i < tam/1024; i++ {
				myfile.Write(buffer[:])
			}
			myfile.Close()
			fmt.Println("Disco creado exitosamente")
			return true
		}
	} else {
		fmt.Println("El disco \"" + Tools.GetFileName(path) + "\" ya existe")
	}
	return false
}

func DeleteDisk(path string) bool {
	if Tools.Exists(path) {
		if err := os.Remove(path); err != nil {
			fmt.Println("ERROR: No se pudo eliminar el disco")
		} else {
			fmt.Println("Disco eliminado exitosamente")
			return true
		}
	} else {
		fmt.Println("El disco \"" + Tools.GetFileName(path) + "\" no existe")
	}
	return false
}

// Retorna la particion extendida del disco
func GetExtPart(path string) Structs.Partition {
	m := Structs.GetMBR(path)
	if m.Mbr_partition_1.Part_type == 'e' {
		return m.Mbr_partition_1
	} else if m.Mbr_partition_2.Part_type == 'e' {
		return m.Mbr_partition_2
	} else if m.Mbr_partition_3.Part_type == 'e' {
		return m.Mbr_partition_3
	} else if m.Mbr_partition_4.Part_type == 'e' {
		return m.Mbr_partition_4
	}
	return Structs.RPV()
}

// Comruebo si una particion dada es primaria
func isPrimPart(m Structs.MBR, name string) bool {
	if string(m.Mbr_partition_1.Part_name[:]) == name && m.Mbr_partition_1.Part_type == 'p' {
		return true
	} else if string(m.Mbr_partition_2.Part_name[:]) == name && m.Mbr_partition_2.Part_type == 'p' {
		return true
	} else if string(m.Mbr_partition_3.Part_name[:]) == name && m.Mbr_partition_3.Part_type == 'p' {
		return true
	} else if string(m.Mbr_partition_4.Part_name[:]) == name && m.Mbr_partition_4.Part_type == 'p' {
		return true
	}
	return false
}

// Compruebo si una particion dada es extendida
func isExtPart(m Structs.MBR, name string) bool {
	if string(m.Mbr_partition_1.Part_name[:]) == name && m.Mbr_partition_1.Part_type == 'e' {
		return true
	} else if string(m.Mbr_partition_2.Part_name[:]) == name && m.Mbr_partition_2.Part_type == 'e' {
		return true
	} else if string(m.Mbr_partition_3.Part_name[:]) == name && m.Mbr_partition_3.Part_type == 'e' {
		return true
	} else if string(m.Mbr_partition_4.Part_name[:]) == name && m.Mbr_partition_4.Part_type == 'e' {
		return true
	}
	return false
}

// Compruebo si una particion dada es logica
func isLogPart(path string, name string) bool {
	var ep = GetExtPart(path)
	if ep.Part_start > 0 {
		var start = Structs.GetEBR(path, ep.Part_start)
		if start.Part_next > 0 {
			actual := Structs.GetEBR(path, start.Part_next)
			for actual.Part_next > 0 {
				if string(actual.Part_name[:]) == name {
					return true
				}
				actual = Structs.GetEBR(path, actual.Part_next)
			}
			if string(actual.Part_name[:]) == name {
				return true
			}
		}

	}
	return false
}

func getPartByName(path string, name string) Structs.Partition {
	m := Structs.GetMBR(path)
	if string(m.Mbr_partition_1.Part_name[:]) == name {
		return m.Mbr_partition_1
	} else if string(m.Mbr_partition_2.Part_name[:]) == name {
		return m.Mbr_partition_2
	} else if string(m.Mbr_partition_3.Part_name[:]) == name {
		return m.Mbr_partition_3
	} else if string(m.Mbr_partition_4.Part_name[:]) == name {
		return m.Mbr_partition_4
	}
	return Structs.RPV()
}

func getLogPartByName(path string, name string) Structs.EBR {
	var ep = GetExtPart(path)
	if ep.Part_s > 0 {
		var start = Structs.GetEBR(path, ep.Part_start)
		if start.Part_next > 0 {
			actual := Structs.GetEBR(path, start.Part_next)
			for actual.Part_next != -1 {
				if string(actual.Part_name[:]) == name {
					return actual
				}
				actual = Structs.GetEBR(path, actual.Part_next)
			}
			if string(actual.Part_name[:]) == name {
				return actual
			}
		}
	}
	return Structs.REBRV()
}

// Comprueba si existe una particion
// Recibe la ruta del disco y el nombre de la particion
func partExists(path string, name string) bool {
	var m = Structs.GetMBR(path)
	if isPrimPart(m, name) { //Busco si existe una particion primaria con el nombre dado
		return true
	} else if isExtPart(m, name) { //Busco si existe una particion extendida con el nombre dado
		return true
	} else if isLogPart(path, name) { //Busco si existe una particion logica con el nombre dado
		return true
	}
	return false
}

// Retorna el numero de particiones primarias
func numPrimPart(m Structs.MBR) int {
	var cont = 0
	if m.Mbr_partition_1.Part_type == 'p' {
		cont++
	}
	if m.Mbr_partition_2.Part_type == 'p' {
		cont++
	}
	if m.Mbr_partition_3.Part_type == 'p' {
		cont++
	}
	if m.Mbr_partition_4.Part_type == 'p' {
		cont++
	}
	return cont
}

// Retona el numero de particiones extendidas
func numExtPart(m Structs.MBR) int {
	var cont = 0
	if m.Mbr_partition_1.Part_type == 'e' {
		cont++
	}
	if m.Mbr_partition_2.Part_type == 'e' {
		cont++
	}
	if m.Mbr_partition_3.Part_type == 'e' {
		cont++
	}
	if m.Mbr_partition_4.Part_type == 'e' {
		cont++
	}
	return cont
}

// Comprueba si el disco está vacio
// Recibe la ruta del disco como parametro
func isDiskEmpty(path string) bool {
	var m = Structs.GetMBR(path)
	//Compruebo si algún slot de particion no está vacio
	if m.Mbr_partition_1.Part_s > 0 || m.Mbr_partition_2.Part_s > 0 || m.Mbr_partition_3.Part_s > 0 || m.Mbr_partition_4.Part_s > 0 {
		return false
	}
	return true
}

func addPartition(m *Structs.MBR, p Structs.Partition) {
	if m.Mbr_partition_1.Part_s < 0 {
		m.Mbr_partition_1 = p
	} else if m.Mbr_partition_2.Part_s < 0 {
		m.Mbr_partition_2 = p
	} else if m.Mbr_partition_3.Part_s < 0 {
		m.Mbr_partition_3 = p
	} else if m.Mbr_partition_4.Part_s < 0 {
		m.Mbr_partition_4 = p
	}
}

func addLogPartition(path string, prevSpace int, nuevo *Structs.EBR) {
	ep := GetExtPart(path)
	start := Structs.GetEBR(path, ep.Part_start)
	prev := Structs.REBRV()

	if prevSpace != -1 {
		prev = Structs.GetEBR(path, int32(prevSpace))
		if prev.Part_next != -1 {
			nxt := Structs.GetEBR(path, prev.Part_next)
			prev.Part_next = nuevo.Part_start
			nuevo.Part_next = nxt.Part_start
			Structs.AddEBR(path, prev.Part_start, prev)
			Structs.AddEBR(path, nuevo.Part_start, *nuevo)
		} else {
			prev.Part_next = nuevo.Part_start
			Structs.AddEBR(path, prev.Part_start, prev)
			Structs.AddEBR(path, nuevo.Part_start, *nuevo)
		}
	} else {
		if start.Part_next != -1 {
			nuevo.Part_next = Structs.GetEBR(path, start.Part_next).Part_start
		}
		start.Part_next = nuevo.Part_start
		Structs.AddEBR(path, start.Part_start, start)
		Structs.AddEBR(path, nuevo.Part_start, *nuevo)
	}
}

func sortPartitions(m *Structs.MBR) {
	var sizes []int
	var parts []Structs.Partition
	var nuevo []Structs.Partition

	sizes = append(sizes, int(m.Mbr_partition_1.Part_start))
	sizes = append(sizes, int(m.Mbr_partition_2.Part_start))
	sizes = append(sizes, int(m.Mbr_partition_3.Part_start))
	sizes = append(sizes, int(m.Mbr_partition_4.Part_start))

	parts = append(parts, m.Mbr_partition_1)
	parts = append(parts, m.Mbr_partition_2)
	parts = append(parts, m.Mbr_partition_3)
	parts = append(parts, m.Mbr_partition_4)

	sort.Ints(sizes)

	for i := 0; i < len(sizes); i++ {
		for j := 0; j < len(parts); j++ {
			if parts[j].Part_start == int32(sizes[i]) {
				nuevo = append(nuevo, parts[i])
				break
			}
		}
	}

	m.Mbr_partition_1 = nuevo[0]
	m.Mbr_partition_2 = nuevo[1]
	m.Mbr_partition_3 = nuevo[2]
	m.Mbr_partition_4 = nuevo[3]
}

// Retorna un vector de structs que contienen el inicio y el final de cada espacio libre y ocupado en el disco
// Recibe la ruta del disco y el mbr del mismo
func BlockSize(path string, m Structs.MBR) []Structs.SpaceSize {
	var v []Structs.SpaceSize
	inicio := int(unsafe.Sizeof(Structs.MBR{})) + 1
	prevSize := inicio - 1
	prevStart := 0
	endSpace := prevSize + prevStart

	if !isDiskEmpty(path) { //Compruebo si no está vacio el disco
		if m.Mbr_partition_1.Part_s > 0 { //Compruebo si la particion 1 no está vacia
			//Compruebo si hay espacio libre esntre el mbr y la particion 1
			if math.Abs(float64(inicio)-float64(m.Mbr_partition_1.Part_start)) > 0 {
				/* Agrego los datos del espacio a la lista de tamaño d espacio */
				v = append(v, Structs.SpaceSize{Part_start: int32(inicio), Part_s: int32(math.Abs(float64(m.Mbr_partition_1.Part_start) - float64(inicio))), In_use: 'n', Type: '0'})
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_1.Part_start, Part_s: m.Mbr_partition_1.Part_s, In_use: 's', Type: m.Mbr_partition_1.Part_type})
				prevSize = int(m.Mbr_partition_1.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_1.Part_start) //Guardo el inicio de la particion 1
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			} else if math.Abs(float64(inicio)-float64(m.Mbr_partition_1.Part_start)) == 0 {
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_1.Part_start, Part_s: m.Mbr_partition_1.Part_s, In_use: 's', Type: m.Mbr_partition_1.Part_type})
				prevSize = int(m.Mbr_partition_1.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_1.Part_start) //Guardo el inicio de la particion 1
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			}
		}
		if m.Mbr_partition_2.Part_s > 0 { //Compruebo si la particion 2 no está vacia
			//Compruebo si hay espacio libre esntre la particion 1 y la particion 2
			if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_2.Part_start)) > 0 {
				/* Agrego los datos del espacio a la lista de tamaño d espacio */
				v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(m.Mbr_partition_2.Part_start) - float64(endSpace-1))), In_use: 'n', Type: '0'})
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_2.Part_start, Part_s: m.Mbr_partition_2.Part_s, In_use: 's', Type: m.Mbr_partition_2.Part_type})
				prevSize = int(m.Mbr_partition_2.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_2.Part_start) //Guardo el inicio de la particion 2
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			} else if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_2.Part_start)) == 0 {
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_2.Part_start, Part_s: m.Mbr_partition_2.Part_s, In_use: 's', Type: m.Mbr_partition_2.Part_type})
				prevSize = int(m.Mbr_partition_2.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_2.Part_start) //Guardo el inicio de la particion 2
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			}
		}
		if m.Mbr_partition_3.Part_s > 0 { //Compruebo si la particion 3 no está vacia
			//Compruebo si hay espacio libre esntre la particion 2 y la particion 3
			if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_3.Part_start)) > 0 {
				/* Agrego los datos del espacio a la lista de tamaño d espacio */
				v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(m.Mbr_partition_3.Part_start) - float64(endSpace-1))), In_use: 'n', Type: '0'})
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_3.Part_start, Part_s: m.Mbr_partition_3.Part_s, In_use: 's', Type: m.Mbr_partition_3.Part_type})
				prevSize = int(m.Mbr_partition_3.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_3.Part_start) //Guardo el inicio de la particion 3
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			} else if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_3.Part_start)) == 0 {
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_3.Part_start, Part_s: m.Mbr_partition_3.Part_s, In_use: 's', Type: m.Mbr_partition_3.Part_type})
				prevSize = int(m.Mbr_partition_3.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_3.Part_start) //Guardo el inicio de la particion 3
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			}
		}
		if m.Mbr_partition_4.Part_s > 0 { //Compruebo si la particion 4 no está vacia
			//Compruebo si hay espacio libre esntre la particion 3 y la particion 4
			if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_4.Part_start)) > 0 {
				/* Agrego los datos del espacio a la lista de tamaño d espacio */
				v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(m.Mbr_partition_4.Part_start) - float64(endSpace-1))), In_use: 'n', Type: '0'})
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_4.Part_start, Part_s: m.Mbr_partition_4.Part_s, In_use: 's', Type: m.Mbr_partition_4.Part_type})
				prevSize = int(m.Mbr_partition_4.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_4.Part_start) //Guardo el inicio de la particion 4
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			} else if math.Abs(float64(endSpace+1)-float64(m.Mbr_partition_4.Part_start)) == 0 {
				v = append(v, Structs.SpaceSize{Part_start: m.Mbr_partition_4.Part_start, Part_s: m.Mbr_partition_4.Part_s, In_use: 's', Type: m.Mbr_partition_4.Part_type})
				prevSize = int(m.Mbr_partition_4.Part_s)      //Guardo el tamaño de la particion actual
				prevStart = int(m.Mbr_partition_4.Part_start) //Guardo el inicio de la particion 4
				endSpace = prevSize + prevStart               //Guardo el calculo de posicion final de la particion actual
			}
		}

		//Compruebo si hay espacio libre entre la particion 4 y el final del disco
		if (endSpace + 1) < int(m.Mbr_tamano) {
			/* Agrego los datos del espacio a la lista de tamaño de espacio */
			v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(m.Mbr_tamano) - float64(endSpace+1))), In_use: 'n', Type: '0'})
		}
	} else {
		v = append(v, Structs.SpaceSize{Part_start: int32(inicio), Part_s: m.Mbr_tamano - int32(inicio), In_use: 'n', Type: '0'})
	}
	return v
}

func ExtBlockSize(path string) []Structs.SpaceSize {
	var v []Structs.SpaceSize
	ep := GetExtPart(path) //Particion extendida
	inicio := ep.Part_start + int32(unsafe.Sizeof(Structs.EBR{})) + 1
	e := Structs.GetEBR(path, ep.Part_start) //EBR inicial
	endSpace := 0

	if e.Part_next != -1 { //Compruebo si no está vacio el disco
		actual := Structs.GetEBR(path, e.Part_next)
		var siguiente Structs.EBR

		//Primera parte
		if math.Abs(float64(inicio)-float64(actual.Part_start)) > 0 {
			v = append(v, Structs.SpaceSize{Part_start: inicio, Part_s: int32(math.Abs(float64(actual.Part_start) - float64(inicio))), In_use: 'n', Type: '0'})
		}

		//Segunda parte
		for actual.Part_next != -1 {
			endSpace = int(actual.Part_start + actual.Part_s)
			v = append(v, Structs.SpaceSize{Part_start: actual.Part_start, Part_s: actual.Part_s, In_use: 's', Type: 'l'})
			siguiente = Structs.GetEBR(path, actual.Part_next)
			if math.Abs(float64(endSpace+1)-float64(siguiente.Part_start)) > 0 {
				v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(siguiente.Part_start) - float64(endSpace+1))), In_use: 'n', Type: '0'})
			}

			actual = siguiente
		}
		if actual.Part_s > 0 {
			v = append(v, Structs.SpaceSize{Part_start: actual.Part_start, Part_s: actual.Part_s, In_use: 's', Type: 'l'})
		}

		//Tercera parte
		endSpace = int(actual.Part_start + actual.Part_s)
		if (endSpace + 1) < int(ep.Part_s) {
			v = append(v, Structs.SpaceSize{Part_start: int32(endSpace + 1), Part_s: int32(math.Abs(float64(ep.Part_s) - float64(endSpace+1))), In_use: 'n', Type: '0'})
		}
	} else {
		v = append(v, Structs.SpaceSize{Part_start: inicio, Part_s: ep.Part_s - inicio, In_use: 'n', Type: 'l'})
	}
	return v
}

// Asigna la particion en la memoria a bloques según el algoritmo de primer ajuste
// Recibe la ruta del disco y la particion
func firstFit(path string, p *Structs.Partition) bool {
	m := Structs.GetMBR(path)
	ss := BlockSize(path, m)
	for _, s := range ss {
		if (p.Part_s <= s.Part_s) && (s.In_use != 's') {
			p.Part_start = s.Part_start
			addPartition(&m, *p)
			sortPartitions(&m)
			Structs.AddMBR(path, m)
			return true
		}
	}
	return false
}

// Asigna la particion en la memoria a bloques según el algoritmo de mejor ajuste
// Recibe la ruta del disco y la particion
func bestFit(path string, p *Structs.Partition) bool {
	m := Structs.GetMBR(path) //Obtengo el MBR del disco
	ss := BlockSize(path, m)  //Obtengo el bloque de tamaños
	bestFitIdx := -1
	for i := 0; i < len(ss); i++ {
		if (ss[i].Part_s >= p.Part_s) && (ss[i].In_use != 's') {
			if bestFitIdx == -1 {
				bestFitIdx = i
			} else if ss[bestFitIdx].Part_s > ss[i].Part_s {
				bestFitIdx = i
			}
		}
	}
	if bestFitIdx != -1 {
		p.Part_start = ss[bestFitIdx].Part_start
		addPartition(&m, *p)
		sortPartitions(&m)
		Structs.AddMBR(path, m)
		// Structs.ReadMBR(path)
		return true
	}
	return false
}

// Asigna la particion en la memoria a bloques según el algoritmo de peor ajuste
// Recibe la ruta del disco y la particion
func worstFit(path string, p *Structs.Partition) bool {
	m := Structs.GetMBR(path) //Obtengo el MBR del disco
	ss := BlockSize(path, m)  //Obtengo el bloque de tamaños
	worstFitIdx := -1
	for i := 0; i < len(ss); i++ {
		if (ss[i].Part_s >= p.Part_s) && (ss[i].In_use != 's') {
			if worstFitIdx == -1 {
				worstFitIdx = i
			} else if ss[worstFitIdx].Part_s < ss[i].Part_s {
				worstFitIdx = i
			}
		}
	}
	if worstFitIdx != -1 {
		p.Part_start = ss[worstFitIdx].Part_start
		addPartition(&m, *p)
		sortPartitions(&m)
		Structs.AddMBR(path, m)
		return true
	}
	return false
}

// Asigna la particion en la memoria a bloques según el algoritmo de primer ajuste para la particion extendida
// Recibe la ruta del disco y la particion
func extFirstFit(path string, e *Structs.EBR) bool {
	// fmt.Println("First fit")
	prevSpace := -1
	ss := ExtBlockSize(path) //Obtengo el bloque de tamaños
	for _, s := range ss {
		if (e.Part_s <= s.Part_s) && (s.In_use != 's') {
			e.Part_start = s.Part_start
			addLogPartition(path, prevSpace, e)
			return true
		}
		prevSpace = int(s.Part_start) //Guardo el espacio anterior
	}
	return false
}

// Asigna la particion en la memoria a bloques según el algoritmo de mejor ajuste para la particion extendida
// Recibe la ruta del disco y la particion
func extBestFit(path string, e *Structs.EBR) bool {
	// fmt.Println("Best fit")
	ss := ExtBlockSize(path) //Obtengo el bloque de tamaños
	bestFitIdx := -1
	prevSpace := -1
	for i := 0; i < len(ss); i++ {
		if (ss[i].Part_s >= e.Part_s) && (ss[i].In_use != 's') {
			if bestFitIdx == -1 {
				bestFitIdx = i
			} else if ss[bestFitIdx].Part_s > ss[i].Part_s {
				bestFitIdx = i
			}
		}
	}
	if bestFitIdx != -1 {
		e.Part_start = ss[bestFitIdx].Part_start
		if bestFitIdx > 0 {
			prevSpace = int(ss[bestFitIdx-1].Part_start)
		}
		addLogPartition(path, prevSpace, e)
		return true
	}
	return false
}

// Asigna la particion en la memoria a bloques según el algoritmo de peor ajuste para la particion extendida
// Recibe la ruta del disco y la particion
func extWorstFit(path string, e *Structs.EBR) bool {
	// fmt.Println("Worst fit")
	ss := ExtBlockSize(path) //Obtengo el bloque de tamaños
	worstFitIdx := -1
	prevSpace := -1
	for i := 0; i < len(ss); i++ {
		if (ss[i].Part_s >= e.Part_s) && (ss[i].In_use != 's') {
			if worstFitIdx == -1 {
				worstFitIdx = i
			} else if ss[worstFitIdx].Part_s < ss[i].Part_s {
				worstFitIdx = i
			}
		}
	}
	if worstFitIdx != -1 {
		e.Part_start = ss[worstFitIdx].Part_start
		if worstFitIdx > 0 {
			prevSpace = int(ss[worstFitIdx-1].Part_start)
		}
		addLogPartition(path, prevSpace, e)
		return true
	}
	return false
}

func chooseFit(path string, fit byte, p *Structs.Partition) bool {
	if fit == 'f' {
		return firstFit(path, p)
	} else if fit == 'b' {
		return bestFit(path, p)
	} else if fit == 'w' {
		return worstFit(path, p)
	}
	return false
}

func chooseExtFit(path string, fit byte, e *Structs.EBR) bool {
	if fit == 'f' {
		return extFirstFit(path, e)
	} else if fit == 'b' {
		return extBestFit(path, e)
	} else if fit == 'w' {
		return extWorstFit(path, e)
	}
	return false
}

func CreatePart(path string, p Structs.Partition) bool {
	if Tools.Exists(path) { //Compruebo si existe el disco
		if !partExists(path, string(p.Part_name[:])) {
			m := Structs.GetMBR(path)
			if p.Part_type == 'p' { //Si es primaria
				if (numPrimPart(m) + numExtPart(m)) < 4 { //Compruebo si se exedio el numero de particiones permitidas
					if chooseFit(path, m.Dsk_fit, &p) { //Compruebo si se pudo asignar la particion
						fmt.Println("Particion primaria creada correctamente")
						return true
					} else {
						fmt.Println("ERROR: No se pudo asignar la particion")
					}
				} else {
					fmt.Println("ERROR: Se exedió el numero de particiones permitidas")
				}
			} else if p.Part_type == 'e' {
				if numExtPart(m) == 0 {
					if numPrimPart(m) < 4 {
						if chooseFit(path, m.Dsk_fit, &p) {
							var nuevo Structs.EBR = Structs.REBRV()
							nuevo.Part_start = p.Part_start
							Structs.AddEBR(path, p.Part_start, nuevo)
							fmt.Println("Particion extendida creada correctamente")
							return true
						} else {
							fmt.Println("ERROR: No se pudo asignar la particion")
						}
					} else {
						fmt.Println("ERROR: Se exedió el numero de particiones permitidas")
					}
				} else {
					fmt.Println("ERROR: Se exedió el numero de particiones extendidas permitidas")
				}
			} else if p.Part_type == 'l' {
				if numExtPart(m) > 0 {
					var e Structs.EBR = Structs.REBRV()
					e.Part_fit = p.Part_fit
					e.Part_start = p.Part_start
					e.Part_s = p.Part_s
					copy(e.Part_name[:], p.Part_name[:])
					if chooseExtFit(path, p.Part_fit, &e) {
						fmt.Println("Particion logica creada correctamente")
						ep := GetExtPart(path)
						Structs.ReadEBRs(path, ep, string(p.Part_name[:]))
					} else {
						fmt.Println("ERROR: No se pudo asignar la particion logica")
					}
				} else {
					fmt.Println("ERROR: No se puede crear una particion logica si no existe una extendida")
				}
			}
		} else {
			fmt.Println("ERROR: Ya existe una particion con el nombre \"" + string(p.Part_name[:]) + "\"")
		}
	} else {
		fmt.Println("ERROR: el disco \"" + Tools.GetFileName(path) + "\" no existe")
	}
	return false
}

func updatePart(m *Structs.MBR, p Structs.Partition, name string) {
	if string(m.Mbr_partition_1.Part_name[:]) == name {
		m.Mbr_partition_1 = p
	} else if string(m.Mbr_partition_2.Part_name[:]) == name {
		m.Mbr_partition_2 = p
	} else if string(m.Mbr_partition_3.Part_name[:]) == name {
		m.Mbr_partition_3 = p
	} else if string(m.Mbr_partition_4.Part_name[:]) == name {
		m.Mbr_partition_4 = p
	}
}

// Devuelve el siguiente numero correspondiente para montar una particion
func nextDiskNum(path string) int {
	for i, md := range mds {
		if md.Path == path {
			return i
		}
	}
	return -1
}

// Devuelve la letra correspondiente a la particion del mismo disco
func getPartLetter(path string, name string) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	cont := 0
	for _, md := range mds {
		if md.Path == path {
			cont++
		}
	}
	return string(letters[cont])
}

func getIdMtdDisk(path string, name string) string {
	lastNum := "31"
	return lastNum + strconv.Itoa(nextDiskNum(path)) + getPartLetter(path, name)
}

func MountDisk(path string, name string) bool {
	if Tools.Exists(path) {
		if partExists(path, name) {
			md := Structs.MountedDisk{Path: path, Name: name, Id: getIdMtdDisk(path, name)}
			mds = append(mds, md)
			m := Structs.GetMBR(path)
			if isPrimPart(m, name) || isExtPart(m, name) {
				p := getPartByName(path, name)
				p.Part_status = '1'
				updatePart(&m, p, name)
				Structs.AddMBR(path, m)
			} else if isLogPart(path, name) {
				e := getLogPartByName(path, name)
				if e.Part_s > 0 {
					e.Part_status = '1'
					Structs.AddEBR(path, e.Part_start, e)
				}
			}
			fmt.Println("Particion montada correctamente")
			return true
		} else {
			fmt.Println("ERROR: La particion \"" + name + "\" no existe")
		}
	} else {
		fmt.Println("ERROR: el disco \"" + Tools.GetFileName(path) + "\" no existe")
	}
	return false
}

func MakeFileSystem(id string) bool {
	return false
}