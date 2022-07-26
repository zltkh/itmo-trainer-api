package itmoTrainerApi

import "os"

func GetAPIKEY() string {
	return os.Getenv("ITMO_TRAINER_API_APIKEY")
}

func GetDBNAME() string {
	return os.Getenv("ITMO_TRAINER_API_DBNAME")
}

func GetDBUSER() string {
	return os.Getenv("ITMO_TRAINER_API_DBUSER")
}

func GetDBPASSWORD() string {
	return os.Getenv("ITMO_TRAINER_API_DBPASSWORD")
}

func GetDBHOST() string {
	return os.Getenv("ITMO_TRAINER_API_DBHOST")
}
