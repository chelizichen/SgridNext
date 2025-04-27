package patchserver

import (
	"gorm.io/gorm"
	"sgridnext.com/src/domain/db"
)

type T_PatchServer_Mapper struct {
	db gorm.DB
}

func (t *T_PatchServer_Mapper) CeateServer(req *db.Server) (int, error) {
	res := t.db.Debug().Create(*req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) CreateNode(req *db.Node) (int, error) {
	res := t.db.Debug().Create(*req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) CreateGroup(req *db.ServerGroup) (int, error) {
	res := t.db.Debug().Create(*req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) CreateServerPackage(req *db.ServerPackage) (int, error) {
	res := t.db.Debug().Create(*req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}
