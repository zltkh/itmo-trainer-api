package itmoTrainerApi

import "net/http"

func internalError(err error) APIGatewayResponse {
	return APIGatewayResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Internal server error: " + err.Error(),
	}
}

func ok() APIGatewayResponse {
	return APIGatewayResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}
}

func accessDenied() APIGatewayResponse {
	return APIGatewayResponse{
		StatusCode: http.StatusForbidden,
		Body:       "Access denied",
	}
}

func notFound() APIGatewayResponse {
	return APIGatewayResponse{
		StatusCode: http.StatusNotFound,
		Body:       "Not found",
	}
}

func notFoundCustomText(text string) APIGatewayResponse {
	return APIGatewayResponse{
		StatusCode: http.StatusNotFound,
		Body:       text,
	}
}
