package middleware

import (
	"context"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"atk-go-server/app/models"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
)

// JwtToken , basic jwt model
// Cấu trúc JwtToken, mô hình jwt cơ bản
type JwtToken struct {
	C              *config.Configuration
	UserCRUD       services.Repository
	RoleCRUD       services.Repository
	PermissionCRUD services.Repository
	MtServiceCRUD  services.Repository
}

// NewJwtToken , khởi tạo một JwtToken mới
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {

	newHandler := new(JwtToken)
	newHandler.C = c
	newHandler.UserCRUD = *services.NewRepository(c, db, global.ColNames.Users)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.ColNames.Roles)
	newHandler.PermissionCRUD = *services.NewRepository(c, db, global.ColNames.Permissions)
	newHandler.MtServiceCRUD = *services.NewRepository(c, db, global.ColNames.MtServices)

	return newHandler
}

// CheckUserAuth , kiểm tra xác thực người dùng
// Dành cho user
func (jt *JwtToken) CheckUserAuth(requirePermissions []string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		notAuthError := "An unauthorized access!"
		notPermissionError := "You do not have permission to perform the action!"

		tokenString := string(ctx.Request.Header.Peek("Authorization"))
		if tokenString != "" {
			splitToken := strings.Split(tokenString, "Bearer ")
			if len(splitToken) > 1 {
				tokenString = splitToken[1]

				// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
				t := models.JwtToken{}
				token, err := jwt.ParseWithClaims(tokenString, &t, func(token *jwt.Token) (interface{}, error) {
					return []byte(jt.C.JwtSecret), nil
				})

				if err != nil || !token.Valid {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
				} else {
					findUser, err := jt.UserCRUD.FindOneById(context.TODO(), t.ID, nil)
					if findUser == nil {
						utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					} else {
						var user models.User
						bsonBytes, err := bson.Marshal(findUser)
						if err != nil {
							utility.JSON(ctx, utility.Payload(false, err, notAuthError))
						} else {
							err = bson.Unmarshal(bsonBytes, &user)
							if err != nil {
								utility.JSON(ctx, utility.Payload(false, err, notAuthError))
							} else {
								if user.IsBlock {
									utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
								} else {
									isRightToken := false
									for _, _token := range user.Tokens {
										if _token.Token == tokenString {
											ctx.SetUserValue("userId", t.ID)                               // set loggedIn user id in context
											ctx.SetUserValue("roleId", utility.ObjectID2String(user.Role)) // set loggedIn user id in context
											isRightToken = true
											break
										}
									}

									if isRightToken == false {
										utility.JSON(ctx, utility.Payload(false, nil, notAuthError))

									} else {
										if len(requirePermissions) == 0 {
											next(ctx)
										} else {
											strRoleID := utility.ObjectID2String(user.Role)
											findRole, err := jt.RoleCRUD.FindOneById(context.TODO(), strRoleID, nil)
											if findRole == nil {
												utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
											} else {
												var result_findRole models.Role
												bsonBytes, err := bson.Marshal(findRole)
												if err != nil {
													utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
												} else {
													err = bson.Unmarshal(bsonBytes, &result_findRole)
													if err != nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													} else {
														totalCheck := true
														for _, requirePermisson := range requirePermissions {
															checkPermission := false
															for _, s := range result_findRole.Permissions {
																if requirePermisson == s.Name {
																	checkPermission = true
																	break
																}
															}

															if checkPermission == false {
																totalCheck = false
																break
															}
														}

														if totalCheck == true {
															next(ctx)
														} else {
															utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
														}

													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			} else {
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}

// CheckMtServiceAuth , kiểm tra xác thực dịch vụ
// Dành cho các service giao tiếp với nhau
func (jt *JwtToken) CheckMtServiceAuth(requirePermissions []string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		notAuthError := "An unauthorized access!"
		notPermissionError := "You do not have permission to perform the action!"

		tokenString := string(ctx.Request.Header.Peek("Authorization"))
		if tokenString != "" {
			splitToken := strings.Split(tokenString, "Bearer ")
			if len(splitToken) > 1 {
				tokenString = splitToken[1]

				// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
				t := models.JwtToken{}
				token, err := jwt.ParseWithClaims(tokenString, &t, func(token *jwt.Token) (interface{}, error) {
					return []byte(jt.C.JwtSecret), nil
				})

				if err != nil || !token.Valid {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
				} else {
					findUser, err := jt.MtServiceCRUD.FindOneById(context.TODO(), t.ID, nil)
					if findUser == nil {
						utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					} else {
						var user models.MtService
						bsonBytes, err := bson.Marshal(findUser)
						if err != nil {
							utility.JSON(ctx, utility.Payload(false, err, notAuthError))
						} else {
							err = bson.Unmarshal(bsonBytes, &user)
							if err != nil {
								utility.JSON(ctx, utility.Payload(false, err, notAuthError))
							} else {
								if user.IsBlock {
									utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
								} else {
									isRightToken := false
									for _, _token := range user.Tokens {
										if _token.Token == tokenString {
											ctx.SetUserValue("serviceId", t.ID)                            // set loggedIn user id in context
											ctx.SetUserValue("roleId", utility.ObjectID2String(user.Role)) // set loggedIn user id in context
											isRightToken = true
											break
										}
									}

									if isRightToken == false {
										utility.JSON(ctx, utility.Payload(false, nil, notAuthError))

									} else {
										if len(requirePermissions) == 0 {
											next(ctx)
										} else {
											strRoleID := utility.ObjectID2String(user.Role)
											findRole, err := jt.RoleCRUD.FindOneById(context.TODO(), strRoleID, nil)
											if findRole == nil {
												utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
											} else {
												var result_findRole models.Role
												bsonBytes, err := bson.Marshal(findRole)
												if err != nil {
													utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
												} else {
													err = bson.Unmarshal(bsonBytes, &result_findRole)
													if err != nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													} else {
														totalCheck := true
														for _, requirePermisson := range requirePermissions {
															checkPermission := false
															for _, s := range result_findRole.Permissions {
																if requirePermisson == s.Name {
																	checkPermission = true
																	break
																}
															}

															if checkPermission == false {
																totalCheck = false
																break
															}
														}

														if totalCheck == true {
															next(ctx)
														} else {
															utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
														}

													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			} else {
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}
