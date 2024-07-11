# Proyecto de Manejo e Implementación de Archivos

Este proyecto implementa una aplicación para la administración de archivos y particiones utilizando comandos específicos. A continuación, se describen los comandos disponibles.

## Comandos

### MKDISK

Crea un archivo binario que simula un disco duro.

**Sintaxis**:

```bash
mkdisk >size=<tamaño> >path=<ruta> [>fit=<ajuste>] [>unit=<unidad>]
```

**Parámetros**:

- `>size`: (Obligatorio) Tamaño del disco a crear. Debe ser un número positivo mayor a cero.
- `>path`: (Obligatorio) Ruta donde se creará el archivo que representará el disco duro. Se crearán las carpetas necesarias si no existen.
- `>fit`: (Opcional) Ajuste que utilizará el disco para crear las particiones. Valores posibles: `BF` (Best Fit), `FF` (First Fit), `WF` (Worst Fit). Por defecto, `FF`.
- `>unit`: (Opcional) Unidades para el tamaño del disco. Valores posibles: `K` (Kilobytes), `M` (Megabytes). Por defecto, `M`.

**Ejemplo**:

```bash
# Crea un disco de 100 Megabytes en la ruta especificada utilizando el ajuste Best Fit
mkdisk >size=100 >path=/home/user/discos/Disco1.dsk >fit=BF >unit=M
```

### RMDISK

Este parámetro elimina un archivo que representa a un disco duro

**Sintaxis**:

```bash
rmdisk >path=<ruta>
```

**Parámetros**:

- `>path`: (Obligatorio) Ruta donde se creará el archivo que representará el disco duro. Se crearán las carpetas necesarias si no existen.

**Ejemplo**:

```bash
#Elimina con rmdisk Disco4.dsk
rmdisk >path="/home/user/discos/Disco4.dsk”
```

### FDISK

Administra las particiones en el archivo que representa el disco duro.

**Sintaxis**:

```bash
fdisk >size=<tamaño> >path=<ruta> >name=<nombre> [>unit=<unidad>] [>type=<tipo>] [>fit=<ajuste>]
```

**Parámetros**:

- `>size`: (Obligatorio al crear) Tamaño de la partición a crear. Debe ser un número positivo mayor a cero.
- `>path`: (Obligatorio) Ruta del disco en el que se creará la partición. El archivo debe existir.
- `>name`: (Obligatorio) Nombre de la partición. No debe repetirse dentro de las particiones del mismo disco.
- `>unit` (Opcional): Indica las unidades que utilizará el parámetro size. Podrá tener los siguientes valores:
  - `K`: Kilobytes (1024 bytes)
  - `M`: Megabytes (1024 \* 1024 bytes)
  - Por defecto, se utilizarán Megabytes.
- - `>type` (Opcional): Indica el tipo de partición. Podrá tener los siguientes valores:
  - `P`: Primaria
  - `E`: Extendida
  - `L`: Lógica
  - Por defecto, se creará una partición primaria.
- `>fit` (Opcional): Indica el ajuste que utilizará el disco para crear las particiones. Podrá tener los siguientes valores:
  - `BF`: Mejor ajuste (Best Fit)
  - `FF`: Primer ajuste (First Fit)
  - `WF`: Peor ajuste (Worst Fit)
  - Por defecto, se tomará el primer ajuste (FF).
- `>delete` (Opcional): Indica que se eliminará una partición. Requiere el parámetro `>name`.
- `>name` (Obligatorio para eliminar): Indica el nombre de la partición a eliminar.
- `>add` (Opcional): Indica el tamaño adicional que se añadirá a una partición existente. Requiere el parámetro `>name`.

**Ejemplos**:

```bash
# Crear una partición primaria de 500 MB
fdisk >size=500 >path=/home/disco1.dsk

# Crear una partición extendida de 1 GB con mejor ajuste
fdisk >size=1024 >unit=M >path=/home/disco1.dsk >type=E >fit=BF

# Eliminar una partición
fdisk >delete >name=Part1 >path=/home/disco1.dsk

# Añadir 200 MB a una partición existente
fdisk >add=200 >unit=M >name=Part1 >path=/home/disco1.dsk
```

### MOUNT

Monta una partición del disco duro en el sistema con un ID con la siguiente estructura: últimos dos dígitos del Carnet + Número Partición + NombreDisco.

**Parámetros**:

- `>path` (Obligatorio): Es la ruta en la que se encuentra el disco que se montará en el sistema. Este archivo ya debe existir.
- `>name` (Obligatorio): Indica el nombre de la partición a cargar. Si no existe, debe mostrar un error.

**Ejemplo**:

```bash
# Monta la partición Part1 del disco Disco1.dsk
mount >path=/home/Disco1.dsk >name=Part1

# Monta la partición Part2 del disco Disco2.dsk
mount >path=/home/Disco2.dsk >name=Part2
```

### UNMOUNT

Desmonta una partición del sistema. Se utilizará el ID que se le asignó a la partición al momento de cargarla.

**Parámetros**:

- `>id` (Obligatorio): Especifica el ID de la partición que se desmontará. Si no existe el ID, deberá mostrar un error.

**Ejemplo**:

```bash
# Desmonta la partición con ID 061Disco1
unmount >id=061Disco1

# Si no existe, se debe mostrar error
unmount >id=061XXX
```

### MKFS

Este comando realiza un formateo completo de la partición. Por defecto, se formateará como ext2 si no se especifica el parámetro `fs`. También creará un archivo en la raíz llamado `users.txt` que contendrá los usuarios y contraseñas del sistema de archivos.

**Parámetros**:

- `>id` (Obligatorio): Indica el ID de la partición que se formateará.
- `>type` (Opcional): Indica el tipo de formateo. Podrá tener los siguientes valores:
  - `full`: Realiza un formateo completo. Por defecto, se realizará un formateo completo si no se especifica esta opción.
- `>fs` (Opcional): Indica el sistema de archivos a formatear. Podrá tener los siguientes valores:
  - `2fs`: Para el sistema EXT2
  - `3fs`: Para el sistema EXT3
  - Por defecto, se formateará como ext2.

**Ejemplo**:

```bash
# Realiza un formateo completo de la partición en el ID 061Disco1 en ext2
mkfs >type=full >id=061Disco1

# Realiza un formateo completo de la partición en el ID 062Disco1 en ext3
mkfs >type=full >id=062Disco1 >fs=3fs
```

### LOGIN

Este comando se utiliza para iniciar sesión en el sistema.

**Parámetros**:

- `>user` (Obligatorio): Indica el nombre del usuario que inicia sesión.
- `>pass` (Obligatorio): Indica la contraseña del usuario que inicia sesión.

**Ejemplo**:

```bash
# Iniciar sesión con el usuario admin y la contraseña admin123
login >user=admin >pass=admin123

# Iniciar sesión con el usuario user1 y la contraseña pass123
login >user=user1 >pass=pass123
```

### LOGOUT

Este comando se utiliza para cerrar la sesión actual del usuario.

**Ejemplo**:

```bash
# Cerrar sesión del usuario actual
logout
```

### MKGRP

Este comando crea un nuevo grupo en el sistema.

**Parámetros**:

- `>name` (Obligatorio): Indica el nombre del grupo a crear. Si el grupo ya existe, debe mostrar un mensaje de error.

**Ejemplo**:

```bash
# Crear un nuevo grupo llamado grupo1
mkgrp >name=grupo1

# Crear un nuevo grupo llamado admin
mkgrp >name=admin
```

### RMGRP

Este comando elimina un grupo del sistema.

**Parámetros**:

- `>name` (Obligatorio): Indica el nombre del grupo a eliminar. Si el grupo no existe, debe mostrar un mensaje de error.

**Ejemplo**:

```bash
# Eliminar el grupo llamado grupo1
rmgrp >name=grupo1

# Eliminar el grupo llamado admin
rmgrp >name=admin
```

### MKUSR

Este comando crea un nuevo usuario en el sistema.

**Parámetros**:

- `>user` (Obligatorio): Indica el nombre del usuario a crear. Si el usuario ya existe, debe mostrar un mensaje de error.
- `>pass` (Obligatorio): Indica la contraseña del nuevo usuario.
- `>grp` (Obligatorio): Indica el grupo al que pertenecerá el nuevo usuario. Si el grupo no existe, debe mostrar un mensaje de error.

**Ejemplo**:

```bash
# Crear un nuevo usuario llamado user1 con contraseña pass123 en el grupo grupo1
mkusr >user=user1 >pass=pass123 >grp=grupo1

# Crear un nuevo usuario llamado admin con contraseña admin123 en el grupo admin
mkusr >user=admin >pass=admin123 >grp=admin
```

### RMUSR

Este comando elimina un usuario del sistema.

**Parámetros**:

- `>user` (Obligatorio): Indica el nombre del usuario a eliminar. Si el usuario no existe, debe mostrar un mensaje de error.

**Ejemplo**:

```bash
# Eliminar el usuario llamado user1
rmusr >user=user1

# Eliminar el usuario llamado admin
rmusr >user=admin
```

### MKFILE

Crea un archivo en el sistema de archivos simulado.

**Sintaxis**:

```bash
mkfile >path=<ruta> [>size=<tamaño>] [>cont=<contenido>] [>r]
```

**Parámetros**:

- `>path`: (Obligatorio) Ruta del archivo a crear. Si incluye espacios en blanco, debe encerrarse entre comillas.
- `>size`: (Opcional) Tamaño del archivo en bytes. El contenido será números del 0 al 9 repetidos hasta alcanzar el tamaño especificado.
- `>cont`: (Opcional) Ruta de un archivo cuyo contenido se copiará en el nuevo archivo.
- `>r`: (Opcional) Crea las carpetas padres especificadas en la ruta si no existen. No recibe ningún valor adicional.

**Ejemplo**:

```bash
mkfile >path="/home/user/docs/archivo.txt" >size=100 >r
```

### MKDIR

Crea una carpeta en el sistema de archivos simulado.

**Sintaxis**:

```bash
mkdir >path=<ruta> [>r]
```

**Parámetros**:

- `>path`: (Obligatorio) Ruta de la carpeta a crear. Si incluye espacios en blanco, debe encerrarse entre comillas.
- `>r`: (Opcional) Crea las carpetas padres especificadas en la ruta si no existen. No recibe ningún valor adicional.

**Ejemplo**:

```bash
mkdir >path="/home/user/docs/nueva_carpeta" >r
```

### CAT

Muestra el contenido de un archivo.

**Sintaxis**:

```bash
cat >fileN=<archivo1> [>fileN=<archivo2>] ...
```

**Parámetros**:

- `>fileN`: (Obligatorio) Lista de archivos cuyos contenidos se mostrarán.

**Ejemplo**:

```bash
cat >fileN="/home/user/docs/archivo.txt"
```

### CHMOD

Cambia los permisos de un archivo o carpeta.

**Sintaxis**:

```bash
chmod >path=<ruta> >ugo=<permisos> [>r]
```

**Parámetros**:

- `>path`: (Obligatorio) Ruta del archivo o carpeta cuyos permisos se cambiarán.
- `>ugo`: (Obligatorio) Permisos que se asignarán. Formato: `U G O` (Usuario, Grupo, Otros), cada uno con valores de `0` a `7`.
- `>r`: (Opcional) Aplica los cambios recursivamente a todos los archivos y carpetas dentro de la ruta especificada.

**Ejemplo**:

```bash
chmod >path="/home/user/docs" >ugo=764 >r
```

### PAUSE

Pone la aplicación en pausa hasta que se presione una tecla.

**Sintaxis**:

```bash
pause
```

### REP

Genera reportes gráficos utilizando Graphviz.

**Sintaxis**:

```bash
rep >id=<id_partición> >path=<ruta_reporte> >name=<nombre_reporte> [>ruta=<ruta_archivo>]
```

**Parámetros**:

- `>id`: (Obligatorio) ID de la partición sobre la que se generará el reporte.
- `>path`: (Obligatorio) Ruta donde se guardará el reporte.
- `>name`: (Obligatorio) Tipo de reporte a generar. Valores posibles: `mbr`, `disk`, `file`, `ls`.
- `>ruta`: (Opcional) Ruta del archivo o carpeta para los reportes `file` y `ls`.

**Ejemplo**:

```bash
rep >id=061Disco1 >path=/home/user/reports/reporte1.jpg >name=mbr
```

## Interfaz

### Pagina principal

Carga de archivo con la lista de comandos o entrada manual.

![Interfaz principal](/img/interfaz_principal.png)

### Login

Inicia la sesión de un usuario en el disco montado en la partición determinada.

![Login y logout](/img/login.png)

### Reportes

Genera el reporte que se solicite de la lista de reportes disponibles.

![Reportes](/img//reportes.png)
