package dao

import "blog-service/internal/model"

// GetAuth 新增针对获取认证信息的方法
func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
	auth := model.Auth{AppKey: appKey, AppSecret: appSecret}
	return auth.Get(d.engine)
}
