package routes

import (
	notificationController "github.com/Danny19977/certikiosk.git/controller/Notification"
	"github.com/Danny19977/certikiosk.git/controller/auth"
	"github.com/Danny19977/certikiosk.git/controller/user"
	"github.com/Danny19977/certikiosk.git/controller/userlog"
	"github.com/Danny19977/certikiosk.git/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	api := app.Group("/api")

	// Authentification controller - Public routes (no authentication required)
	a := api.Group("/auth")
	a.Post("/register", auth.Register)
	a.Post("/login", auth.Login)
	a.Post("/forgot-password", auth.Forgot)
	a.Post("/reset/:token", auth.ResetPassword)

	// Protected routes (authentication required)
	protected := api.Group("/auth")
	protected.Use(middlewares.IsAuthenticated)
	protected.Get("/user", auth.AuthUser)
	protected.Put("/profil/info", auth.UpdateInfo)
	protected.Put("/change-password", auth.ChangePassword)
	protected.Post("/logout", auth.Logout)

	// Users controller - Protected routes
	u := api.Group("/users")
	u.Use(middlewares.IsAuthenticated)
	u.Get("/all", user.GetAllUsers)
	u.Get("/all/paginate", user.GetPaginatedUsers)
	u.Get("/all/paginate/nosearch", user.GetPaginatedNoSerach)

	u.Get("/get/:uuid", user.GetUser)
	u.Post("/create", user.CreateUser)
	u.Put("/update/:uuid", user.UpdateUser)
	u.Delete("/delete/:uuid", user.DeleteUser)

	// UserLogs controller - Protected routes
	log := api.Group("/users-logs")
	log.Use(middlewares.IsAuthenticated)
	log.Get("/all", userlog.GetUserLogs)
	log.Get("/all/paginate", userlog.GetPaginatedUserLogs)
	log.Get("/all/paginate/:user_uuid", userlog.GetUserLogByID)
	log.Get("/get/:uuid", userlog.GetUserLog)
	log.Post("/create", userlog.CreateUserLog)
	log.Put("/update/:uuid", userlog.UpdateUserLog)
	log.Delete("/delete/:uuid", userlog.DeleteUserLog)

	// Notification controller - Protected routes
	notificationGroup := api.Group("/notifications")
	notificationGroup.Use(middlewares.IsAuthenticated)
	notificationGroup.Get("/all", notificationController.GetAllNotifications)
	notificationGroup.Get("/all/paginate", notificationController.GetPaginatedNotification)
	notificationGroup.Get("/get/:uuid", notificationController.GetNotification)
	notificationGroup.Get("/get/title/:title", notificationController.GetNotificationByTitleString)
	notificationGroup.Post("/create", notificationController.CreateNotification)
	notificationGroup.Put("/update/:uuid", notificationController.UpdateNotification)
	notificationGroup.Delete("/delete/:uuid", notificationController.DeleteNotification)

}
