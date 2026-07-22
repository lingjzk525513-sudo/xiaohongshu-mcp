package main

import (
	"net/http"
		"os"
			"strings"

				"github.com/gin-gonic/gin"
					"github.com/modelcontextprotocol/go-sdk/mcp"
					)

					// tokenAuthMiddleware 验证 Bearer Token（Gin 中间件）
					func tokenAuthMiddleware() gin.HandlerFunc {
						return func(c *gin.Context) {
								expectedToken := os.Getenv("MCP_API_TOKEN")
										if expectedToken == "" {
													// 未设置 Token，直接放行（本地开发）
																c.Next()
																			return
																					}
																							authHeader := c.GetHeader("Authorization")
																									if !strings.HasPrefix(authHeader, "Bearer ") {
																												c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
																															return
																																	}
																																			token := strings.TrimPrefix(authHeader, "Bearer ")
																																					if token != expectedToken {
																																								c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
																																											return
																																													}
																																															c.Next()
																																																}
																																																}

																																																// setupRoutes 设置路由配置
																																																func setupRoutes(appServer *AppServer) *gin.Engine {
																																																	// 设置 Gin 模式
																																																		gin.SetMode(gin.ReleaseMode)

																																																			router := gin.New()
																																																				router.Use(gin.Logger())
																																																					router.Use(gin.Recovery())

																																																						// 添加中间件
																																																							router.Use(errorHandlingMiddleware())
																																																								router.Use(corsMiddleware())

																																																									// 健康检查
																																																										router.GET("/health", healthHandler)

																																																											// MCP 端点 - 使用官方 SDK 的 Streamable HTTP Handler
																																																												mcpHandler := mcp.NewStreamableHTTPHandler(
																																																														func(r *http.Request) *mcp.Server {
																																																																	return appServer.mcpServer
																																																																			},
																																																																					&mcp.StreamableHTTPOptions{
																																																																								JSONResponse: true, // 支持 JSON 响应
																																																																										},
																																																																											)

																																																																												// ========== 为 /mcp 路由添加 Token 鉴权 ==========
																																																																													mcpGroup := router.Group("/mcp")
																																																																														mcpGroup.Use(tokenAuthMiddleware())
																																																																															{
																																																																																	mcpGroup.Any("", gin.WrapH(mcpHandler))
																																																																																			mcpGroup.Any("/*path", gin.WrapH(mcpHandler))
																																																																																				}
																																																																																					// ================================================

																																																																																						// API 路由组
																																																																																							api := router.Group("/api/v1")
																																																																																								{
																																																																																										api.GET("/login/status", appServer.checkLoginStatusHandler)
																																																																																												api.GET("/login/qrcode", appServer.getLoginQrcodeHandler)
																																																																																														api.DELETE("/login/cookies", appServer.deleteCookiesHandler)
																																																																																																api.POST("/publish", appServer.publishHandler)
																																																																																																		api.POST("/publish_video", appServer.publishVideoHandler)
																																																																																																				api.GET("/feeds/list", appServer.listFeedsHandler)
																																																																																																						api.GET("/feeds/search", appServer.searchFeedsHandler)
																																																																																																								api.POST("/feeds/search", appServer.searchFeedsHandler)
																																																																																																										api.POST("/feeds/detail", appServer.getFeedDetailHandler)
																																																																																																												api.POST("/user/profile", appServer.userProfileHandler)
																																																																																																														api.POST("/feeds/comment", appServer.postCommentHandler)
																																																																																																																api.POST("/feeds/comment/reply", appServer.replyCommentHandler)
																																																																																																																		api.GET("/user/me", appServer.myProfileHandler)
																																																																																																																			}

																																																																																																																				return router
																																																																																																																				}