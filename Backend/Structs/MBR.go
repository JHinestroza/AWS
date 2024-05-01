package Structs

type MBR struct {
	MBR_tamanio        int64
	MBR_fecha_creacion [16]byte
	MBR_DISK_SIGNATURE int64
	DSK_fit            [1]byte
	MBR_particion1     Particion
	MBR_particion2     Particion
	MBR_particion3     Particion
	MBR_particion4     Particion
}

func NewMBR() MBR {
	var mb MBR
	return mb
}
