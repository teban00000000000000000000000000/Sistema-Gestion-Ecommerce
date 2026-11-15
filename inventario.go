package inventario

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"sistema-gestion-ecommerce/almacenamiento"
	"sistema-gestion-ecommerce/productos"
)

type Movimiento struct {
	ID       int
	Producto string
	Tipo     string
	Cantidad int
}

var movimientos []Movimiento

func init() {
	err := almacenamiento.CargarJSON("datos/inventario.json", &movimientos)
	if err != nil {
		fmt.Println("Error al cargar inventario:", err)
	}
}

func GestionarInventario(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== GESTION DE INVENTARIO =====")
		fmt.Println("1. Registrar entrada de productos")
		fmt.Println("2. Registrar salida de productos")
		fmt.Println("3. Mostrar movimientos")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			registrarMovimiento(lector, "Entrada")
		case "2":
			registrarMovimiento(lector, "Salida")
		case "3":
			mostrarMovimientos()
		case "0":
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func registrarMovimiento(lector *bufio.Reader, tipo string) {
	var movimiento Movimiento

	fmt.Print("Ingrese el nombre del producto: ")
	movimiento.Producto, _ = lector.ReadString('\n')
	movimiento.Producto = strings.TrimSpace(movimiento.Producto)

	if movimiento.Producto == "" {
		fmt.Println("El nombre del producto es obligatorio.")
		return
	}

	fmt.Print("Ingrese la cantidad: ")
	texto, _ := lector.ReadString('\n')
	cantidad, err := strconv.Atoi(strings.TrimSpace(texto))

	if err != nil || cantidad <= 0 {
		fmt.Println("Ingrese una cantidad valida.")
		return
	}

	cambio := cantidad
	if tipo == "Salida" {
		cambio = -cantidad
	}

	err = productos.ActualizarStock(movimiento.Producto, cambio)
	if err != nil {
		fmt.Println("No se pudo actualizar el inventario:", err)
		return
	}

	movimiento.ID = len(movimientos) + 1
	movimiento.Tipo = tipo
	movimiento.Cantidad = cantidad

	movimientos = append(movimientos, movimiento)

	err = almacenamiento.GuardarJSON("datos/inventario.json", movimientos)
	if err != nil {
		fmt.Println("Error al guardar el movimiento:", err)
		movimientos = movimientos[:len(movimientos)-1]

		// Revertir el cambio de existencias.
		productos.ActualizarStock(movimiento.Producto, -cambio)
		return
	}

	fmt.Println(tipo, "registrada correctamente.")
	fmt.Println("Las existencias del producto se actualizaron.")
}

func mostrarMovimientos() {
	fmt.Println("\n===== MOVIMIENTOS DE INVENTARIO =====")

	if len(movimientos) == 0 {
		fmt.Println("No existen movimientos registrados.")
		return
	}

	for _, movimiento := range movimientos {
		fmt.Printf("\nID: %d\n", movimiento.ID)
		fmt.Printf("Producto: %s\n", movimiento.Producto)
		fmt.Printf("Tipo: %s\n", movimiento.Tipo)
		fmt.Printf("Cantidad: %d\n", movimiento.Cantidad)
	}
}

func ObtenerTotalMovimientos() int {
	return len(movimientos)
}