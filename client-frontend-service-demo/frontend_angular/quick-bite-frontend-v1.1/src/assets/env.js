(function (window) {
  window.__env = window.__env || {};

  window.__env.apiBaseUrl = "${API_BASE_URL}";
  window.__env.restaurantsServiceUrl = "${RESTAURANTS_SERVICE_URL}";
  window.__env.ordersServiceUrl = "${ORDERS_SERVICE_URL}";
  window.__env.paymentsServiceUrl = "${PAYMENTS_SERVICE_URL}";
  window.__env.deliveriesServiceUrl = "${DELIVERIES_SERVICE_URL}";
  window.__env.notificationsServiceUrl = "${NOTIFICATIONS_SERVICE_URL}";
  window.__env.frontendPort = "${FRONTEND_PORT}";
})(this);
