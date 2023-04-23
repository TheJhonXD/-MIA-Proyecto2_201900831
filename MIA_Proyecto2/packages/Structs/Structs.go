package Structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"
)

type SpaceSize struct {
	Part_start int32
	Part_s     int32
	In_use     byte
	Type       byte
}

type MountedDisk struct {
	Path string
	Name string
	Id   string
}

type Time struct {
	Year  int32
	Month int32
	Day   int32
	Hour  int32
	Min   int32
	Sec   int32
}

func (t *Time) SetTime() {
	t.Year = int32(time.Now().Year())
	t.Month = int32(time.Now().Month())
	t.Day = int32(time.Now().Day())
	t.Hour = int32(time.Now().Hour())
	t.Min = int32(time.Now().Minute())
	t.Sec = int32(time.Now().Second())
}

type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int32
	Part_s      int32
	Part_name   [16]byte
}

type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion Time
	Mbr_dsk_signature  int32
	Dsk_fit            byte
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int32
	Part_s      int32
	Part_next   int32
	Part_name   [16]byte
}

// Reset Partition Variable
// Limpia la variable de tipo Partition o inicializa
func RPV() Partition {
	return Partition{'0', '0', '0', int32(-1), int32(-1), [16]byte(bytes.Repeat([]byte("-1"), 16))}
}

// Reset EBR Variable
// Limpia la variable de tipo EBR o inicializa
func REBRV() EBR {
	return EBR{'0', 0, int32(-1), int32(-1), int32(-1), [16]byte(bytes.Repeat([]byte("-1"), 16))}
}

// Añade el MBR a un disco especificado
// Recibe el path del disco y el MBR
func AddMBR(path string, m MBR) {
	myfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(0, 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	err = binary.Write(myfile, binary.LittleEndian, &m)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
}

// Retorna el MBR de un disco especificado
// Recibe el path del disco
func GetMBR(path string) MBR {
	myfile, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(0, 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	var mbr MBR
	err = binary.Read(myfile, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
	return mbr
}

// Lee el MBR de un disco especificado
// Recibe el path del disco
func ReadMBR(path string) {
	myfile, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(0, 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	var mbr MBR
	err = binary.Read(myfile, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
	fmt.Println("Size: ", mbr.Mbr_tamano)
	fmt.Println("Fecha: ", mbr.Mbr_fecha_creacion.Day, "/", mbr.Mbr_fecha_creacion.Month, "/", mbr.Mbr_fecha_creacion.Year)
	fmt.Println("Hora: ", mbr.Mbr_fecha_creacion.Hour, ":", mbr.Mbr_fecha_creacion.Min, ":", mbr.Mbr_fecha_creacion.Sec)
	fmt.Println("ID: ", mbr.Mbr_dsk_signature)
	fmt.Println("Type: ", strconv.Quote(string(mbr.Dsk_fit)))
	fmt.Println("Part 1: ", string(mbr.Mbr_partition_1.Part_name[:]))
	fmt.Println("Part 2: ", string(mbr.Mbr_partition_2.Part_name[:]))
	fmt.Println("Part 3: ", string(mbr.Mbr_partition_3.Part_name[:]))
	fmt.Println("Part 4: ", string(mbr.Mbr_partition_4.Part_name[:]))
}

func AddEBR(path string, start int32, e EBR) {
	myfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(int64(start), 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	err = binary.Write(myfile, binary.LittleEndian, &e)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
}

func GetEBR(path string, start int32) EBR {
	myfile, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	_, err = myfile.Seek(int64(start), 0)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	var e EBR
	err = binary.Read(myfile, binary.LittleEndian, &e)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	myfile.Close()
	return e
}

// Retorna el EBR de la particion, la busca por nombre
// Reciba la ruta del disco, la particion extendida y el nombre
func GetEBRByName(path string, part Partition, name string) EBR {
	start := GetEBR(path, part.Part_start)
	if start.Part_next != -1 {
		actual := GetEBR(path, start.Part_next)
		for actual.Part_next != -1 {
			if string(actual.Part_name[:]) == name {
				return actual
			}
			actual = GetEBR(path, actual.Part_next)
		}
		if string(actual.Part_name[:]) == name {
			return actual
		}
	}
	return REBRV()
}

// Lee todos los ebrs de una particion extendida
// Recibe la ruta del disco, la particion extendida y el nombre
func ReadEBRs(path string, ep Partition, name string) {
	start := GetEBR(path, ep.Part_start)
	if start.Part_next != -1 {
		actual := GetEBR(path, start.Part_next)
		for actual.Part_next != -1 {
			fmt.Println("Nombre: ", string(actual.Part_name[:]))
			fmt.Println("Status: ", string(actual.Part_status))
			fmt.Println("Fit: ", string(actual.Part_fit))
			fmt.Println("Start: ", actual.Part_start)
			fmt.Println("Size: ", actual.Part_s)
			fmt.Println("Next: ", actual.Part_next)
			actual = GetEBR(path, actual.Part_next)
			fmt.Println("+++++++++++++++++++++++++++")
		}
		if string(actual.Part_name[:]) == name {
			fmt.Println("Nombre: ", string(actual.Part_name[:]))
			fmt.Println("Status: ", string(actual.Part_status))
			fmt.Println("Fit: ", string(actual.Part_fit))
			fmt.Println("Start: ", actual.Part_start)
			fmt.Println("Size: ", actual.Part_s)
			fmt.Println("Next: ", actual.Part_next)
		}
	}
}
