package pedidos

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sistema-gestion-ecommerce/almacenamiento"
	"sistema-gestion-ecommerce/productos"
)

type Pedido struct {
	ID       int
	Cliente  string
	Producto string
	Cantidad int
	Precio   float64
	Total    float64
	Estado   string
	Fecha    string
}

var pedidosRegistrados []Pedido

func init() {
	err := almacenamiento.CargarJSON("datos/pedidos.json", &pedidosRegistrados)
	if err != nil {
		fmt.Println("Error al cargar pedidos:", err)
	}
}

func GestionarPedidos(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== GESTION DE PEDIDOS =====")
		fmt.Println("1. Registrar pedido")
		fmt.Println("2. Mostrar pedidos")
		fmt.Println("3. Actualizar estado del pedido")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			registrarPedido(lector)
		case "2":
			mostrarPedidos()
		case "3":
			actualizarEstado(lector)
		case "0":
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func registrarPedido(lector *bufio.Reader) {
	fmt.Print("Nombre del cliente: ")
	cliente, _ := lector.ReadString('\n')
	cliente = strings.TrimSpace(cliente)

	fmt.Print("Nombre del producto: ")
	nombre, _ := lector.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	fmt.Print("Cantidad: ")
	textoCantidad, _ := lector.ReadString('\n')
	cantidad, err := strconv.Atoi(strings.TrimSpace(textoCantidad))

	if err != nil || cantidad <= 0 {
		fmt.Println("Cantidad no valida.")
		return
	}

	fmt.Print("Precio unitario: ")
	textoPrecio, _ := lector.ReadString('\n')
	precio, err := strconv.ParseFloat(strings.TrimSpace(textoPrecio), 64)

	if err != nil || precio <= 0 {
		fmt.Println("Precio no valido.")
		return
	}

	err = CrearPedido(cliente, nombre, cantidad, precio)
	if err != nil {
		fmt.Println("No se pudo registrar el pedido:", err)
		return
	}

	fmt.Println("Pedido registrado correctamente.")
	fmt.Printf("Total del pedido: $%.2f\n", float64(cantidad)*precio)
}

func CrearPedido(cliente, nombreProducto string, cantidad int, precio float64) error {
	if strings.TrimSpace(cliente) == "" || strings.TrimSpace(nombreProducto) == "" {
		return errors.New("el cliente y el producto son obligatorios")
	}

	if cantidad <= 0 || precio <= 0 {
		return errors.New("la cantidad y el precio deben ser mayores que cero")
	}

	producto, encontrado := productos.BuscarProducto(nombreProducto)
	if !encontrado {
		return errors.New("el producto no existe")
	}

	// Se utiliza el precio registrado en el catálogo.
	precio = producto.Precio

	err := productos.ActualizarStock(producto.Nombre, -cantidad)
	if err != nil {
		return err
	}

	pedido := Pedido{
		ID:       len(pedidosRegistrados) + 1,
		Cliente:  strings.TrimSpace(cliente),
		Producto: producto.Nombre,
		Cantidad: cantidad,
		Precio:   precio,
		Total:    float64(cantidad) * precio,
		Estado:   "Pendiente",
		Fecha:    time.Now().Format("02/01/2006"),
	}

	pedidosRegistrados = append(pedidosRegistrados, pedido)

	err = almacenamiento.GuardarJSON("datos/pedidos.json", pedidosRegistrados)
	if err != nil {
		pedidosRegistrados = pedidosRegistrados[:len(pedidosRegistrados)-1]
		productos.ActualizarStock(producto.Nombre, cantidad)
		return err
	}

	return nil
}

func mostrarPedidos() {
	fmt.Println("\n===== LISTA DE PEDIDOS =====")

	if len(pedidosRegistrados) == 0 {
		fmt.Println("No existen pedidos registrados.")
		return
	}

	for _, pedido := range pedidosRegistrados {
		fmt.Printf("\nID: %d\n", pedido.ID)
		fmt.Printf("Cliente: %s\n", pedido.Cliente)
		fmt.Printf("Producto: %s\n", pedido.Producto)
		fmt.Printf("Cantidad: %d\n", pedido.Cantidad)
		fmt.Printf("Total: $%.2f\n", pedido.Total)
		fmt.Printf("Estado: %s\n", pedido.Estado)
		fmt.Printf("Fecha: %s\n", pedido.Fecha)
	}
}

func actualizarEstado(lector *bufio.Reader) {
	if len(pedidosRegistrados) == 0 {
		fmt.Println("No hay pedidos para actualizar.")
		return
	}

	mostrarPedidos()

	fmt.Print("\nIngrese el ID del pedido: ")
	texto, _ := lector.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(texto))

	if err != nil || id <= 0 || id > len(pedidosRegistrados) {
		fmt.Println("ID de pedido no valido.")
		return
	}

	fmt.Println("1. Pendiente")
	fmt.Println("2. Enviado")
	fmt.Println("3. Entregado")
	fmt.Print("Seleccione el nuevo estado: ")

	estadoTexto, _ := lector.ReadString('\n')

	switch strings.TrimSpace(estadoTexto) {
	case "1":
		pedidosRegistrados[id-1].Estado = "Pendiente"
	case "2":
		pedidosRegistrados[id-1].Estado = "Enviado"
	case "3":
		pedidosRegistrados[id-1].Estado = "Entregado"
	default:
		fmt.Println("Estado no valido.")
		return
	}

	err = almacenamiento.GuardarJSON("datos/pedidos.json", pedidosRegistrados)
	if err != nil {
		fmt.Println("Error al actualizar el pedido:", err)
		return
	}

	fmt.Println("Estado del pedido actualizado correctamente.")
}

func ObtenerTotalPedidos() int {
	return len(pedidosRegistrados)
}

func ObtenerTotalVentas() float64 {
	var total float64

	for _, pedido := range pedidosRegistrados {
		total += pedido.Total
	}

	return total
}