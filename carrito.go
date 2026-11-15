package carrito

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"sistema-gestion-ecommerce/almacenamiento"
	"sistema-gestion-ecommerce/productos"
	"sistema-gestion-ecommerce/pedidos"
)

type ItemCarrito struct {
	Producto string
	Precio   float64
	Cantidad int
	Subtotal float64
}

var items []ItemCarrito

func init() {
	err := almacenamiento.CargarJSON("datos/carrito.json", &items)
	if err != nil {
		fmt.Println("Error al cargar carrito:", err)
	}
}

func GestionarCarrito(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== CARRITO DE COMPRAS =====")
		fmt.Println("1. Agregar producto")
		fmt.Println("2. Mostrar carrito")
		fmt.Println("3. Calcular total")
		fmt.Println("4. Finalizar compra")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			agregarProducto(lector)
		case "2":
			mostrarCarrito()
		case "3":
			calcularTotal()
		case "4":
			finalizarCompra(lector)
		case "0":
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func agregarProducto(lector *bufio.Reader) {
	fmt.Println("\n===== AGREGAR PRODUCTO =====")

	fmt.Print("Ingrese el nombre del producto: ")
	nombre, _ := lector.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	producto, encontrado := productos.BuscarProducto(nombre)

	if !encontrado {
		fmt.Println("El producto no existe en el catalogo.")
		return
	}

	fmt.Printf("\nProducto: %s\n", producto.Nombre)
	fmt.Printf("Precio: $%.2f\n", producto.Precio)
	fmt.Printf("Stock disponible: %d\n", producto.Cantidad)

	if producto.Cantidad <= 0 {
		fmt.Println("Este producto no tiene existencias.")
		return
	}

	fmt.Print("Ingrese la cantidad: ")
	texto, _ := lector.ReadString('\n')

	cantidad, err := strconv.Atoi(strings.TrimSpace(texto))

	if err != nil || cantidad <= 0 {
		fmt.Println("Ingrese una cantidad valida.")
		return
	}

	cantidadEnCarrito := 0

	for _, item := range items {
		if strings.EqualFold(item.Producto, producto.Nombre) {
			cantidadEnCarrito += item.Cantidad
		}
	}

	if cantidadEnCarrito+cantidad > producto.Cantidad {
		fmt.Println("No hay suficientes unidades disponibles.")
		return
	}

	item := ItemCarrito{
		Producto: producto.Nombre,
		Precio:   producto.Precio,
		Cantidad: cantidad,
		Subtotal: producto.Precio * float64(cantidad),
	}

	items = append(items, item)

	err = almacenamiento.GuardarJSON("datos/carrito.json", items)

	if err != nil {
		items = items[:len(items)-1]
		fmt.Println("Error al guardar el carrito:", err)
		return
	}

	fmt.Println("\nProducto agregado correctamente.")
	fmt.Printf("Subtotal: $%.2f\n", item.Subtotal)
}

func mostrarCarrito() {
	fmt.Println("\n===== PRODUCTOS EN EL CARRITO =====")

	if len(items) == 0 {
		fmt.Println("El carrito esta vacio.")
		return
	}

	for _, item := range items {
		fmt.Println("------------------------------")
		fmt.Printf("Producto: %s\n", item.Producto)
		fmt.Printf("Precio: $%.2f\n", item.Precio)
		fmt.Printf("Cantidad: %d\n", item.Cantidad)
		fmt.Printf("Subtotal: $%.2f\n", item.Subtotal)
	}

	fmt.Println("------------------------------")
}

func calcularTotal() {
	var total float64

	for _, item := range items {
		total += item.Subtotal
	}

	fmt.Println("\n===== TOTAL DE COMPRA =====")
	fmt.Printf("TOTAL: $%.2f\n", total)
}

func finalizarCompra(lector *bufio.Reader) {
	if len(items) == 0 {
		fmt.Println("El carrito esta vacio.")
		return
	}

	mostrarCarrito()
	calcularTotal()

	fmt.Print("\nIngrese el nombre del cliente: ")
	cliente, _ := lector.ReadString('\n')
	cliente = strings.TrimSpace(cliente)

	if cliente == "" {
		fmt.Println("El nombre del cliente es obligatorio.")
		return
	}

	fmt.Print("¿Desea confirmar la compra? (S/N): ")
	confirmacion, _ := lector.ReadString('\n')

	if strings.ToUpper(strings.TrimSpace(confirmacion)) != "S" {
		fmt.Println("Compra cancelada.")
		return
	}

	for len(items) > 0 {
		item := items[0]

		err := pedidos.CrearPedido(
			cliente,
			item.Producto,
			item.Cantidad,
			item.Precio,
		)

		if err != nil {
			fmt.Println("No se pudo completar la compra:", err)
			fmt.Println("Los productos restantes permanecen en el carrito.")
			return
		}

		items = items[1:]

		err = almacenamiento.GuardarJSON("datos/carrito.json", items)

		if err != nil {
			fmt.Println("Pedido registrado, pero hubo un error al actualizar el carrito:", err)
			return
		}
	}

	fmt.Println("\n¡COMPRA REALIZADA CORRECTAMENTE!")
	fmt.Println("Los pedidos fueron registrados.")
	fmt.Println("El inventario fue actualizado.")
	fmt.Println("El carrito esta vacio.")
}