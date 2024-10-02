package handler

import (
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type RoleHandler struct {
	crud services.Repository
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler(c *config.Configuration, db *mongo.Client) *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.crud = *services.NewRepository(c, db, global.ColNames.Roles)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo một vai trò mới
// func (h *RoleHandler) Create(ctx *fasthttp.RequestCtx) {
// 	var response map[string]interface{} = nil

// 	// Lấy dữ liệu từ yêu cầu
// 	postValues := ctx.PostBody()
// 	inputStruct := new(models.RoleCreateInput)
// 	response = utility.Convert2Struct(postValues, inputStruct)
// 	if response == nil { // Kiểm tra dữ liệu đầu vào
// 		response = utility.ValidateStruct(inputStruct)

// 		// Chuyển đổi chuỗi thành ObjectId
// 		newRole := new(models.Role)
// 		newRole.Name = inputStruct.Name
// 		newRole.Describe = inputStruct.Describe
// 		var listObjectID []models.PermissionDetail
// 		for _, s := range inputStruct.Permissions {
// 			newPermissonDetail := new(models.PermissionDetail)
// 			newPermissonDetail.ID = permissionMap["_id"].(primitive.ObjectID)
// 			newPermissonDetail.Name = permissionMap["name"].(string)

// 			listObjectID = append(listObjectID, *newPermissonDetail)
// 		}
// 		newRole.Permissions = listObjectID

// 		if response == nil { // Gọi hàm xử lý logic
// 			response = utility.FinalResponse(h.crud.InsertOne(ctx, newRole))
// 		}
// 	}

// 	utility.JSON(ctx, response)
// }

// Tìm một vai trò theo ID
func (h *RoleHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.FindOneById(ctx, id, nil))

	utility.JSON(ctx, response)
}

// Tìm tất cả các vai trò với phân trang
func (h *RoleHandler) FindAll(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	buf := string(ctx.FormValue("limit"))
	limit, err := strconv.ParseInt(buf, 10, 64)
	if err != nil {
		limit = 10
	}

	buf = string(ctx.FormValue("page"))
	page, err := strconv.ParseInt(buf, 10, 64)
	if err != nil {
		page = 0
	}

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	response = utility.FinalResponse(h.crud.FindAllWithPaginate(ctx, bson.D{}, opts))

	utility.JSON(ctx, response)
}

// Cập nhật một vai trò theo ID
// func (h *RoleHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
// 	var response map[string]interface{} = nil

// 	// Lấy ID từ yêu cầu
// 	id := ctx.UserValue("id").(string)
// 	response = utility.FinalResponse(h.crud.FindOneById(ctx, id))

// 	// Lấy dữ liệu từ yêu cầu
// 	postValues := ctx.PostBody()
// 	inputStruct := new(models.RoleUpdateInput)
// 	response = utility.Convert2Struct(postValues, inputStruct)
// 	if response == nil { // Kiểm tra dữ liệu đầu vào
// 		response = utility.ValidateStruct(inputStruct)

// 		// Chuyển đổi chuỗi thành ObjectId
// 		newRole := new(models.Role)
// 		newRole.Name = inputStruct.Name
// 		newRole.Describe = inputStruct.Describe
// 		var listObjectID []primitive.ObjectID
// 		for _, s := range inputStruct.Permissions {
// 			listObjectID = append(listObjectID, utility.String2ObjectID(s))
// 		}
// 		newRole.Permissions = listObjectID

// 		if response == nil { // Gọi hàm tạo json changes
// 			var change map[string]interface{}
// 			response = utility.CreateChangeMap(newRole, &change)
// 			if response == nil { // Gọi hàm xử lý logic
// 				response = utility.FinalResponse(h.crud.UpdateOneById(ctx, id, change))
// 			}
// 		}
// 	}

// 	utility.JSON(ctx, response)
// }

// Xóa một vai trò theo ID
// func (h *RoleHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
// 	var response map[string]interface{} = nil

// 	// Lấy ID từ yêu cầu
// 	id := ctx.UserValue("id").(string)
// 	response = utility.FinalResponse(h.crud.DeleteOneById(ctx, id))

// 	utility.JSON(ctx, response)
// }

// Other functions =========================================================================
