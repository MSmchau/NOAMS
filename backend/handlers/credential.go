package handlers

import (
	"strconv"

	"noams/models"
	"noams/services"
	"noams/utils"

	"github.com/gin-gonic/gin"
)

type CredentialHandler struct {
	svc *services.CredentialService
}

func NewCredentialHandler(svc *services.CredentialService) *CredentialHandler {
	return &CredentialHandler{svc: svc}
}

// createCredentialRequest 创建凭证的请求结构（Password 需要显式接收）
type createCredentialRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	EnablePW    string `json:"enable_pw"`
	AuthMethod  string `json:"auth_method"`
	Description string `json:"description"`
}

func (h *CredentialHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		utils.ServerError(c, "查询凭证列表失败")
		return
	}
	utils.Success(c, list)
}

func (h *CredentialHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的凭证 ID")
		return
	}
	cred, err := h.svc.GetByID(uint(id))
	if err != nil {
		utils.BadRequest(c, "凭证不存在")
		return
	}
	utils.Success(c, cred)
}

func (h *CredentialHandler) Create(c *gin.Context) {
	var req createCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}
	if req.Name == "" || req.Username == "" || req.Password == "" {
		utils.BadRequest(c, "名称、用户名、密码为必填项")
		return
	}
	cred := models.Credential{
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		EnablePW:    req.EnablePW,
		AuthMethod:  req.AuthMethod,
		Description: req.Description,
	}
	if cred.AuthMethod == "" {
		cred.AuthMethod = "password"
	}
	if err := h.svc.Create(&cred); err != nil {
		utils.ServerError(c, "创建凭证失败: "+err.Error())
		return
	}
	utils.Success(c, cred)
}

func (h *CredentialHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的凭证 ID")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "请求参数无效")
		return
	}
	if name, ok := updates["name"]; ok && name == "" {
		utils.BadRequest(c, "名称不能为空")
		return
	}
	if err := h.svc.Update(uint(id), updates); err != nil {
		utils.ServerError(c, "更新凭证失败")
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func (h *CredentialHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的凭证 ID")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除凭证失败")
		return
	}
	utils.Success(c, gin.H{"id": id})
}
