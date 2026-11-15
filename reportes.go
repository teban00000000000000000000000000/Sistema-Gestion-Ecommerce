package reportes

import (
	"fmt"
	"time"

	"sistema-gestion-ecommerce/almacenamiento"
)

type Reporte struct {
	Fecha            string
	TotalUsuarios    int
	TotalProductos   int
	TotalMovimientos int
	TotalPedidos     int
	TotalVentas      float64
}

var reportesGenerados []Reporte

func init() {
	err := almacenamiento.CargarJSON("datos/reportes.json", &reportesGenerados)
	if err != nil {
		fmt.Println("Error al cargar reportes:", err)
	}
}

func GenerarReporte(usuarios, productos, movimientos, pedidos int, ventas float64) {
	reporte := Reporte{
		Fecha:            time.Now().Format("02/01/2006 15:04"),
		TotalUsuarios:    usuarios,
		TotalProductos:   productos,
		TotalMovimientos: movimientos,
		TotalPedidos:     pedidos,
		TotalVentas:      ventas,
	}

	reportesGenerados = append(reportesGenerados, reporte)

	err := almacenamiento.GuardarJSON("datos/reportes.json", reportesGenerados)
	if err != nil {
		fmt.Println("Error al guardar el reporte:", err)
		reportesGenerados = reportesGenerados[:len(reportesGenerados)-1]
		return
	}

	fmt.Println("\n========== REPORTE GENERAL ==========")
	fmt.Println("Fecha:", reporte.Fecha)
	fmt.Println("-------------------------------------")
	fmt.Println("Usuarios registrados:", reporte.TotalUsuarios)
	fmt.Println("Productos registrados:", reporte.TotalProductos)
	fmt.Println("Movimientos de inventario:", reporte.TotalMovimientos)
	fmt.Println("Pedidos registrados:", reporte.TotalPedidos)
	fmt.Printf("Total de ventas: $%.2f\n", reporte.TotalVentas)
	fmt.Println("=====================================")
}

func MostrarHistorialReportes() {
	fmt.Println("\n===== HISTORIAL DE REPORTES =====")

	if len(reportesGenerados) == 0 {
		fmt.Println("No existen reportes registrados.")
		return
	}

	for i, reporte := range reportesGenerados {
		fmt.Printf("\nReporte #%d\n", i+1)
		fmt.Println("Fecha:", reporte.Fecha)
		fmt.Println("Usuarios:", reporte.TotalUsuarios)
		fmt.Println("Productos:", reporte.TotalProductos)
		fmt.Println("Movimientos:", reporte.TotalMovimientos)
		fmt.Println("Pedidos:", reporte.TotalPedidos)
		fmt.Printf("Ventas: $%.2f\n", reporte.TotalVentas)
	}
}