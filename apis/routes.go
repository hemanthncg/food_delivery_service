package apis

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/restaurant/enroll", RestaurantEnrollHandler)
	http.HandleFunc("/restaurants", RestaurantsHandler)
	http.HandleFunc("/menu", MenuHandler)
	http.HandleFunc("/driver/enroll", DriverEnrollHandler)
	http.HandleFunc("/drivers", DriversHandler)
	http.HandleFunc("/driver/location", DriverLocationHandler)
	http.HandleFunc("/order", OrderHandler)
	http.HandleFunc("/order/delivered", DeliveredHandler)
	http.HandleFunc("/order/get", GetOrderHandler)
	http.HandleFunc("/orders", ListOrdersHandler)

	// OpenAPI documentation endpoints
	http.HandleFunc("/openapi.json", ServeOpenAPIDoc)
	http.HandleFunc("/docs", ServeOpenAPISwagger)
}
