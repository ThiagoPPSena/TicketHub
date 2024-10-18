package passagesService

func GetAllRoutes(origin string, destination string) ([]string, error) {
	return []string{"Salvador", "Sao Paulo", "Recife"}, nil
}

func Buy(routes []string) (bool, error) {
	return false, nil
}
