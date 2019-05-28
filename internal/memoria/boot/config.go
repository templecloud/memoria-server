package boot

import (
	"github.com/templecloud/memoria-server/internal/memoria/boot/logging"
	"github.com/templecloud/memoria-server/internal/memoria/boot/persistence"
	"github.com/templecloud/memoria-server/internal/memoria/controller"
)

//-------------------------------------------------------------------------------------------------
// Models

// Config denotes the configuration of the server controller.
type Config struct {
	Controller *controller.Config `json:"controller" binding:"optional"`
	Logging *logging.Config `json:"logging" binding:"optional"`
	Persistence *persistence.Config `json:"persistence" binding:"optional"`
}

//-------------------------------------------------------------------------------------------------
// Public Functions

// NewDefaultConfig returns a default Memoria configuration.
func NewDefaultConfig() *Config {
	return &Config{
		Controller: controller.NewDefaultConfig(),
		Logging: logging.NewDefaultConfig(),
		Persistence: persistence.NewDefaultConfig(),
	}
}
