package itmoTrainerApi

import "os"

func GetAPIKEY() string {
	return os.Getenv("ITMO-TRAINER-API-APIKEY")
}

func GetDBNAME() string {
	return os.Getenv("ITMO-TRAINER-API-DBNAME")
}

func GetDBUSER() string {
	return os.Getenv("ITMO-TRAINER-API-DBUSER")
}

func GetDBPASSWORD() string {
	return os.Getenv("ITMO-TRAINER-API-DBPASSWORD")
}

func GetDBHOST() string {
	return os.Getenv("ITMO-TRAINER-API-DBHOST")
}
