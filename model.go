package turispro_user

type UserLevel int

const (
	Admin UserLevel = iota
	AdminTour
	VendedorTour
	Operador
	Chofer
	Guia
)

func (ul UserLevel) String() string {
	switch ul {
	case Admin:
		return "Admin"
	case AdminTour:
		return "AdminTour"
	case VendedorTour:
		return "VendedorTour"
	case Operador:
		return "Operador"
	case Chofer:
		return "Chofer"
	case Guia:
		return "Guia"
	default:
		return "Unknown"
	}
}
