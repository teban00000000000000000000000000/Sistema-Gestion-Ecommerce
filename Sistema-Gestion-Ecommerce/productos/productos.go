package productos

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"sistema-gestion-ecommerce/almacenamiento"
)

type Producto struct {
	ID       int
	Nombre   string
	Precio   float64
	Cantidad int
}

var productosRegistrados []Producto

func init() {
	err := almacenamiento.CargarJSON("datos/productos.json", &productosRegistrados)
	if err != nil {
		fmt.Println("Error al cargar productos:", err)
	}
}

func GestionarProductos(lector *bufio.Reader) {
	for {
		fmt.Println("\n===== GESTION DE PRODUCTOS =====")
		fmt.Println("1. Registrar producto")
		fmt.Println("2. Mostrar productos")
		fmt.Println("0. Volver al menu principal")
		fmt.Print("Seleccione una opcion: ")

		entrada, _ := lector.ReadString('\n')
		opcion := strings.TrimSpace(entrada)

		switch opcion {
		case "1":
			registrarProducto(lector)
		case "2":
			mostrarProductos()
		case "0":
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func registrarProducto(lector *bufio.Reader) {
	var producto Producto

	producto.ID = len(productosRegistrados) + 1

	fmt.Print("Ingrese el nombre del producto: ")
	producto.Nombre, _ = lector.ReadString('\n')
	producto.Nombre = strings.TrimSpace(producto.Nombre)

	if producto.Nombre == "" {
		fmt.Println("El nombre es obligatorio.")
		return
	}

	fmt.Print("Ingrese el precio: ")
	precioTexto, _ := lector.ReadString('\n')
	precio, err := strconv.ParseFloat(strings.TrimSpace(precioTexto), 64)

	if err != nil || precio <= 0 {
		fmt.Println("Ingrese un precio valido.")
		return
	}

	fmt.Print("Ingrese la cantidad disponible: ")
	cantidadTexto, _ := lector.ReadString('\n')
	cantidad, err := strconv.Atoi(strings.TrimSpace(cantidadTexto))

	if err != nil || cantidad < 0 {
		fmt.Println("Ingrese una cantidad valida.")
		return
	}

	producto.Precio = precio
	producto.Cantidad = cantidad

	productosRegistrados = append(productosRegistrados, producto)

	err = almacenamiento.GuardarJSON("datos/productos.json", productosRegistrados)
	if err != nil {
		fmt.Println("Error al guardar el producto:", err)
		productosRegistrados = productosRegistrados[:len(productosRegistrados)-1]
		return
	}

	fmt.Println("Producto registrado y guardado correctamente.")
}

func mostrarProductos() {
	fmt.Println("\n===== LISTA DE PRODUCTOS =====")

	if len(productosRegistrados) == 0 {
		fmt.Println("No existen productos registrados.")
		return
	}

	for _, producto := range productosRegistrados {
		fmt.Printf("\nID: %d\n", producto.ID)
		fmt.Printf("Nombre: %s\n", producto.Nombre)
		fmt.Printf("Precio: $%.2f\n", producto.Precio)
		fmt.Printf("Cantidad disponible: %d\n", producto.Cantidad)
	}
}

func ObtenerTotalProductos() int {
	return len(productosRegistrados)
}

func ActualizarStock(nombre string, cambio int) error {
	for i := range productosRegistrados {
		if strings.EqualFold(productosRegistrados[i].Nombre, nombre) {
			nuevaCantidad := productosRegistrados[i].Cantidad + cambio

			if nuevaCantidad < 0 {
				return errors.New("no hay suficientes unidades disponibles")
			}

			cantidadAnterior := productosRegistrados[i].Cantidad
			productosRegistrados[i].Cantidad = nuevaCantidad

			err := almacenamiento.GuardarJSON("datos/productos.json", productosRegistrados)
			if err != nil {
				productosRegistrados[i].Cantidad = cantidadAnterior
				return err
			}

			return nil
		}
	}

	return errors.New("el producto no existe")
}

// BuscarProducto busca un producto por su nombre.
func BuscarProducto(nombre string) (Producto, bool) {
	for _, producto := range productosRegistrados {
		if strings.EqualFold(producto.Nombre, nombre) {
			return producto, true
		}
	}

	return Producto{}, false
}