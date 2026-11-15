package almacenamiento

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GuardarJSON almacena información en un archivo JSON.
func GuardarJSON(ruta string, datos interface{}) error {
	carpeta := filepath.Dir(ruta)

	err := os.MkdirAll(carpeta, 0755)
	if err != nil {
		return err
	}

	contenido, err := json.MarshalIndent(datos, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(ruta, contenido, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CargarJSON recupera información de un archivo JSON.
func CargarJSON(ruta string, destino interface{}) error {
	contenido, err := os.ReadFile(ruta)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	err = json.Unmarshal(contenido, destino)
	if err != nil {
		return err
	}

	return nil
}