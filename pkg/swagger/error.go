package swagger

type InternalServerError struct {
	Status   string `json:"status" example:"InternalServerError"`
	Response struct {
		Code    int    `json:"code" example:"500"`
		Message string `json:"message" example:"Внутренняя ошибка сервера"`
	} `json:"response"`
}

type IncorrectFormatError struct {
	Status   string `json:"status" example:"IncorrectFormatError"`
	Response struct {
		Code    int    `json:"code" example:"400"`
		Message string `json:"message" example:"Неверный формат запроса"`
	} `json:"response"`
}

type IncorrectDataError struct {
	Status   string `json:"status" example:"IncorrectDataError"`
	Response struct {
		Code    int    `json:"code" example:"400"`
		Message string `json:"message" example:"Некорректные данные"`
	} `json:"response"`
}

type NotFoundError struct {
	Status   string `json:"status" example:"NotFoundError"`
	Response struct {
		Code    int    `json:"code" example:"404"`
		Message string `json:"message" example:"Данные не найдены"`
	} `json:"response"`
}
