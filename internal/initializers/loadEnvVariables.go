package initializers

import (
	"github.com/joho/godotenv"
)

// esto solo sirve en local que recibe el archivo .env
// en eks se le pasaran las variables a traves de un config map
// y seran directamente accesibles haciendo os.Getenv("var")
func LoadEnvVariables() {
	_ = godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}
