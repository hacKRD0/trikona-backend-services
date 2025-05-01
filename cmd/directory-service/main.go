package main

import (
	"os"

	"github.com/gin-gonic/gin"

	// user-management imports
	user_api "github.com/hacKRD0/trikona_go/api/user-management-service"
	user_repo "github.com/hacKRD0/trikona_go/internal/user-management-service/repository"
	user_uc "github.com/hacKRD0/trikona_go/internal/user-management-service/usecase"

	// directory-service imports
	dir_api "github.com/hacKRD0/trikona_go/api/directory-service/student"
	dir_repo "github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	dir_uc "github.com/hacKRD0/trikona_go/internal/directory-service/usecase"

	"github.com/hacKRD0/trikona_go/pkg/auth"
	"github.com/hacKRD0/trikona_go/pkg/config"
	"github.com/hacKRD0/trikona_go/pkg/database"
	"github.com/hacKRD0/trikona_go/pkg/logger"
	"github.com/hacKRD0/trikona_go/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	// 1. Init logger & config
	if err := logger.InitLogger(); err != nil {
		panic("Logger init failed: " + err.Error())
	}
	defer logger.Log.Sync()

	if err := config.LoadEnv(); err != nil {
		logger.Fatal("Env load failed", err)
	}

	// 2. Init DB
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("DB init failed", err)
	}

	// 3. JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtSvc := auth.NewJWTService(jwtSecret)

	// 4. Wire up user-management
	usrRepo := user_repo.NewUserRepository(db)
	usrUse := user_uc.NewUserUsecase(usrRepo, jwtSvc)
	usrHnd := user_api.NewHandler(usrUse)

	// 5. Wire up directory-service
	stuRepo := dir_repo.NewStudentRepository(db)
	stuUse := dir_uc.NewStudentUsecase(stuRepo)
	stuHnd := dir_api.NewStudentHandler(stuUse) // adjust if your constructor differs

	// 6. Setup Gin router
	router := gin.Default()
	router.Use(middleware.RequestLogger())

	// 7. CORS (copy your existing corsConfig from user main)
	// Configure CORS
	corsConfig := &middleware.CorsConfig{
		AllowOrigins: []string{
			"*",                                // Allow all origins
			"http://localhost:5172",            // Development frontend
			os.Getenv("FRONTEND_URL"),          // Production frontend
			os.Getenv("LINKEDIN_REDIRECT_URI"), // LinkedIn callback
			os.Getenv("BASE_URL"),
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"Access-Control-Allow-Origin",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}
	router.Use(middleware.Cors(corsConfig))

	// 8. Mount everything under /api/v1
	api := router.Group("/api/v1")
	{
		// --- user-management public ---
		api.POST("/register", usrHnd.Register)
		api.POST("/login", usrHnd.Login)
		api.POST("/verify", usrHnd.RequestVerification)
		api.POST("/forgot-password", usrHnd.RequestPasswordReset)
		api.GET("/linkedin/auth-url", usrHnd.GetLinkedInAuthURL)
		api.POST("/linkedin/callback", usrHnd.GetLinkedInProfileInfo)

		// --- protected user routes ---
		apiAuth := api.Group("")
		apiAuth.Use(user_api.AuthMiddleware(jwtSvc))
		{
			apiAuth.PUT("/users", usrHnd.UpdateUser)
			apiAuth.DELETE("/users", usrHnd.DeleteUser)
			apiAuth.POST("/reset-password", usrHnd.ResetPassword)
			// --- directory-service student routes ---
			apiStu := apiAuth.Group("/directory")
			{
				apiStu.GET("/students", stuHnd.GetStudents)
				apiStu.GET("/students/:id", stuHnd.GetStudent)
				apiStu.POST("/students", stuHnd.CreateStudent)
				apiStu.PUT("/students/:id", stuHnd.UpdateStudent)
				apiStu.DELETE("/students/:id", stuHnd.DeleteStudent)
			}

			// --- directory-service college routes ---
			// apiCol := apiAuth.Group("/directory/colleges")
			// {
			// 	apiCol.GET("/", colHnd.GetColleges)
			// 	apiCol.GET("/:id", colHnd.GetCollege)
			// 	apiCol.POST("/", colHnd.CreateCollege)
			// 	apiCol.PUT("/:id", colHnd.UpdateCollege)
			// 	apiCol.DELETE("/:id", colHnd.DeleteCollege)
			// }
		}
	}

	// 9. Start
	port := os.Getenv("PORT")
	logger.Info("Starting combined server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Server failed", err)
	}
}
