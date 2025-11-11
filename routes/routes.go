package routes

import (
	notificationController "github.com/Danny19977/certikiosk.git/controller/Notification"
	"github.com/Danny19977/certikiosk.git/controller/auth"
	certificationController "github.com/Danny19977/certikiosk.git/controller/certification"
	citizensController "github.com/Danny19977/certikiosk.git/controller/citizens"
	documentsController "github.com/Danny19977/certikiosk.git/controller/documents"
	fingerprintController "github.com/Danny19977/certikiosk.git/controller/fingerprint"
	"github.com/Danny19977/certikiosk.git/controller/user"
	"github.com/Danny19977/certikiosk.git/controller/userlog"
	"github.com/Danny19977/certikiosk.git/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	api := app.Group("/api")

	// ============================================
	// PUBLIC ROUTES (No Authentication Required)
	// ============================================
	public := api.Group("/public")

	// Public citizen registration
	public.Post("/citizens/register", citizensController.CreateCitizen)

	// Public fingerprint operations for kiosk
	public.Post("/fingerprint/enroll", fingerprintController.EnrollFingerprint)
	public.Post("/fingerprint/verify", fingerprintController.VerifyFingerprint)

	// Public document access for kiosk (read-only operations)
	public.Get("/documents/national-id/:national_id", documentsController.GetDocumentsByNationalID)
	public.Get("/documents/active", documentsController.GetActiveDocuments)
	public.Get("/documents/:uuid", documentsController.GetDocument)
	public.Post("/documents/send-email", documentsController.SendDocumentEmail)

	// Public Google Drive document operations for kiosk
	public.Post("/documents/send-email-gdrive", documentsController.SendDocumentEmailFromGDrive)
	public.Get("/documents/generate-stamped-pdf", documentsController.GenerateStampedPDF)
	public.Post("/documents/generate-stamped-pdf", documentsController.GenerateStampedPDF)
	public.Get("/documents/stamped-pdf-metadata", documentsController.GenerateStampedPDFMetadata)

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

	// Citizens controller - Protected routes (admin operations only)
	// Note: Public citizen registration is available at /api/public/citizens/register
	citizens := api.Group("/citizens")
	citizens.Use(middlewares.IsAuthenticated)
	citizens.Get("/all", citizensController.GetAllCitizens)
	citizens.Get("/all/paginate", citizensController.GetPaginatedCitizens)
	citizens.Get("/get/:uuid", citizensController.GetCitizen)
	citizens.Get("/national-id/:national_id", citizensController.GetCitizenByNationalID)
	citizens.Put("/update/:uuid", citizensController.UpdateCitizen)
	citizens.Delete("/delete/:uuid", citizensController.DeleteCitizen)

	// Fingerprint controller - Protected routes
	// Note: Public fingerprint enrollment and verification are available at /api/public/fingerprint/*
	fingerprint := api.Group("/fingerprint")
	fingerprint.Use(middlewares.IsAuthenticated)
	fingerprint.Get("/all/paginate", fingerprintController.GetPaginatedFingerprints)
	fingerprint.Get("/citizen/:citizen_uuid", fingerprintController.GetFingerprintByCitizen)
	fingerprint.Put("/update/:citizen_uuid", fingerprintController.UpdateFingerprint)
	fingerprint.Delete("/delete/:citizen_uuid", fingerprintController.DeleteFingerprint)

	// Documents controller - Protected routes (admin operations)
	// Note: Public read-only document access and email sending are available at /api/public/documents/*
	documents := api.Group("/documents")
	documents.Use(middlewares.IsAuthenticated)
	documents.Get("/all", documentsController.GetAllDocuments)
	documents.Get("/all/paginate", documentsController.GetPaginatedDocuments)
	documents.Get("/user/:user_uuid", documentsController.GetDocumentsByUserUUID)
	documents.Post("/create", documentsController.CreateDocument)
	documents.Post("/fetch-external", documentsController.FetchDocumentFromExternalSource)
	documents.Put("/update/:uuid", documentsController.UpdateDocument)
	documents.Put("/toggle-status/:uuid", documentsController.ToggleDocumentStatus)
	documents.Delete("/delete/:uuid", documentsController.DeleteDocument)
	documents.Post("/send-email", documentsController.SendDocumentEmail)
	documents.Post("/send-email-gdrive", documentsController.SendDocumentEmailFromGDrive)
	documents.Get("/generate-stamped-pdf", documentsController.GenerateStampedPDF)
	documents.Post("/generate-stamped-pdf", documentsController.GenerateStampedPDF)
	documents.Get("/stamped-pdf-metadata", documentsController.GenerateStampedPDFMetadata)

	// Certification controller - Protected routes
	certification := api.Group("/certification")
	certification.Use(middlewares.IsAuthenticated)
	certification.Get("/all", certificationController.GetAllCertifications)
	certification.Get("/all/paginate", certificationController.GetPaginatedCertifications)
	certification.Get("/get/:uuid", certificationController.GetCertification)
	certification.Get("/citizen/:citizen_uuid", certificationController.GetCertificationsByCitizen)
	certification.Get("/document/:document_uuid", certificationController.GetCertificationsByDocument)
	certification.Get("/download/:uuid", certificationController.DownloadCertifiedDocument)
	certification.Get("/print/:uuid", certificationController.PrintCertifiedDocument)
	certification.Post("/certify", certificationController.CertifyDocument)
	certification.Put("/revoke/:uuid", certificationController.RevokeCertification)
	certification.Delete("/delete/:uuid", certificationController.DeleteCertification)

}
