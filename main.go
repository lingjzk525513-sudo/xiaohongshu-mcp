package main

import (
	"flag"
		"net/http"   // 新增
			"os"
				"strings"    // 新增

					"github.com/sirupsen/logrus"
						"github.com/xpzouying/xiaohongshu-mcp/browser"
							"github.com/xpzouying/xiaohongshu-mcp/configs"
							)

							// ========== 鉴权中间件（Bearer Token）==========
							// 从环境变量 MCP_API_TOKEN 读取期望的 Token，检查请求头 Authorization: Bearer <token>
							func authMiddleware(next http.Handler) http.Handler {
								return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
										expectedToken := os.Getenv("MCP_API_TOKEN")
												if expectedToken == "" {
															// 未设置 Token 时直接放行（本地开发）
																		next.ServeHTTP(w, r)
																					return
																							}
																									authHeader := r.Header.Get("Authorization")
																											if !strings.HasPrefix(authHeader, "Bearer ") {
																														http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
																																	return
																																			}
																																					token := strings.TrimPrefix(authHeader, "Bearer ")
																																							if token != expectedToken {
																																										http.Error(w, "Invalid token", http.StatusUnauthorized)
																																													return
																																															}
																																																	next.ServeHTTP(w, r)
																																																		})
																																																		}
																																																		// ===============================================

																																																		func main() {
																																																			var (
																																																					headless bool
																																																							binPath  string
																																																									port     string
																																																										)
																																																											flag.BoolVar(&headless, "headless", true, "是否无头模式")
																																																												flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
																																																													flag.StringVar(&port, "port", ":18060", "端口")
																																																														flag.Parse()

																																																															if len(binPath) == 0 {
																																																																	binPath = os.Getenv("ROD_BROWSER_BIN")
																																																																		}
																																																																			if binPath == "" {
																																																																					bin, err := browser.EnsureBrowser()
																																																																							if err != nil {
																																																																										logrus.Fatalf("%v", err)
																																																																												}
																																																																														binPath = bin
																																																																															}
																																																																																logrus.Infof("using browser binary: %s", binPath)

																																																																																	configs.InitHeadless(headless)
																																																																																		configs.SetBinPath(binPath)
																																																																																			configs.SetFingerprintSeed(configs.FingerprintSeedFromEnv())
																																																																																				configs.SetProxy(configs.ProxyFromEnv())

																																																																																					xiaohongshuService := NewXiaohongshuService()

																																																																																						appServer := NewAppServer(xiaohongshuService)

																																																																																							// ⚠️ 重要：此处 appServer.Start 内部注册了 /mcp 路由，
																																																																																								// 我们需要在 appServer 的 Start 方法中将原始 handler 用 authMiddleware 包裹。
																																																																																									// 请打开 server.go（或 app.go），找到 Start 方法，将路由注册改成：
																																																																																										//   http.Handle("/mcp", authMiddleware(原始handler))
																																																																																											if err := appServer.Start(port); err != nil {
																																																																																													logrus.Fatalf("failed to run server: %v", err)
																																																																																														}
																																																																																														}