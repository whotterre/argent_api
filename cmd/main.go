//	@title			Argent Wallet API
//	@version		1.0
//	@description	A wallet API for managing deposits, transfers, and balances with Paystack integration.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host			localhost:9000
//	@BasePath		/

//	@securityDefinitions.apikey	BearerAuth
//	@type						apiKey
//	@in							header
//	@name						Authorization
//	@description				Enter the token in the format: Bearer {token}

package main

import (
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/initializers"
	"whotterre/argent/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	// Connect to database
	initializers.ConnectToDB(cfg.DatabaseURL)
	db := initializers.DB
	routes.SetupRoutes(app, cfg, db)

	port := ":" + cfg.Port
	app.Run(port)
}
