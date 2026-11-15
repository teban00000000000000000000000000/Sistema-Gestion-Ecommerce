package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"sistema-gestion-ecommerce/usuarios"
	"sistema-gestion-ecommerce/productos"
	"sistema-gestion-ecommerce/inventario"
	"sistema-gestion-ecommerce/carrito"
	"sistema-gestion-ecommerce/pedidos"
	"sistema-gestion-ecommerce/reportes"
)

func mostrarMenu() {
	fmt.Println("\n======================================")
	fmt.Println("       SISTEMA DE GESTION E-COMMERCE")
	fmt.Println("======================================")
	fmt.Println("1. Gestion de usuarios")
	fmt.Println("2. Gestion de productos")
	fmt.Println("3. Gestion de inventario")
	fmt.Println("4. Carrito de compras")
	fmt.Println("5. Gestion de pedidos")
	fmt.Println("6. Reportes")
	fmt.Println("0. Salir")
	fmt.Println("======================================")
	fmt.Print("Seleccione una opcion: ")
}

func gestionarReportes(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== MODULO DE REPORTES =====")
		fmt.Println("1. Generar reporte general")
		fmt.Println("2. Mostrar historial de reportes")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			reportes.GenerarReporte(
				usuarios.ObtenerTotalUsuarios(),
				productos.ObtenerTotalProductos(),
				inventario.ObtenerTotalMovimientos(),
				pedidos.ObtenerTotalPedidos(),
				pedidos.ObtenerTotalVentas(),
			)

		case "2":
			reportes.MostrarHistorialReportes()

		case "0":
			return

		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func main() {
	lector := bufio.NewReader(os.Stdin)

	for {
		mostrarMenu()

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			usuarios.GestionarUsuarios(lector)

		case "2":
			productos.GestionarProductos(lector)

		case "3":
			inventario.GestionarInventario(lector)

		case "4":
			carrito.GestionarCarrito(lector)

		case "5":
			pedidos.GestionarPedidos(lector)

		case "6":
			gestionarReportes(lector)

		case "0":
			fmt.Println("Gracias por utilizar el sistema.")
			return

		default:
			fmt.Println("Opcion no valida. Intente nuevamente.")
		}
	}
}