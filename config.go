package itmo_trainer_api

import "os"

func getAPIKEY() string {
	return os.Getenv("ITMO_TRAINER_API_APIKEY")
}

func getDBNAME() string {
	return os.Getenv("ITMO_TRAINER_API_DBNAME")
}

func getDBUSER() string {
	return os.Getenv("ITMO_TRAINER_API_DBUSER")
}

func getDBPASSWORD() string {
	return os.Getenv("ITMO_TRAINER_API_DBPASSWORD")
}

func getDBHOST() string {
	return os.Getenv("ITMO_TRAINER_API_DBHOST")
}
