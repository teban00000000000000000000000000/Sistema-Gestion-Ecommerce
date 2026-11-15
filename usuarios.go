package usuarios

import (
	"bufio"
	"fmt"
	"strings"

	"sistema-gestion-ecommerce/almacenamiento"
)

type Usuario struct {
	ID       int
	Nombre   string
	Correo   string
	Telefono string
}

var usuariosRegistrados []Usuario

func init() {
	err := almacenamiento.CargarJSON("datos/usuarios.json", &usuariosRegistrados)
	if err != nil {
		fmt.Println("Error al cargar usuarios:", err)
	}
}

func GestionarUsuarios(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== GESTION DE USUARIOS =====")
		fmt.Println("1. Registrar usuario")
		fmt.Println("2. Mostrar usuarios")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			registrarUsuario(lector)
		case "2":
			mostrarUsuarios()
		case "0":
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func registrarUsuario(lector *bufio.Reader) {
	var usuario Usuario
	usuario.ID = len(usuariosRegistrados) + 1

	fmt.Print("Ingrese el nombre: ")
	usuario.Nombre, _ = lector.ReadString('\n')
	usuario.Nombre = strings.TrimSpace(usuario.Nombre)

	fmt.Print("Ingrese el correo: ")
	usuario.Correo, _ = lector.ReadString('\n')
	usuario.Correo = strings.TrimSpace(usuario.Correo)

	fmt.Print("Ingrese el telefono: ")
	usuario.Telefono, _ = lector.ReadString('\n')
	usuario.Telefono = strings.TrimSpace(usuario.Telefono)

	if usuario.Nombre == "" || usuario.Correo == "" || usuario.Telefono == "" {
		fmt.Println("Todos los campos son obligatorios.")
		return
	}

	usuariosRegistrados = append(usuariosRegistrados, usuario)

	err := almacenamiento.GuardarJSON("datos/usuarios.json", usuariosRegistrados)
	if err != nil {
		fmt.Println("Error al guardar el usuario:", err)
		usuariosRegistrados = usuariosRegistrados[:len(usuariosRegistrados)-1]
		return
	}

	fmt.Println("Usuario registrado y guardado correctamente.")
}

func mostrarUsuarios() {
	fmt.Println("\n===== LISTA DE USUARIOS =====")

	if len(usuariosRegistrados) == 0 {
		fmt.Println("No existen usuarios registrados.")
		return
	}

	for _, usuario := range usuariosRegistrados {
		fmt.Printf("\nID: %d\n", usuario.ID)
		fmt.Printf("Nombre: %s\n", usuario.Nombre)
		fmt.Printf("Correo: %s\n", usuario.Correo)
		fmt.Printf("Telefono: %s\n", usuario.Telefono)
	}
}

// ObtenerTotalUsuarios devuelve la cantidad de usuarios registrados.
func ObtenerTotalUsuarios() int {
	return len(usuariosRegistrados)
}