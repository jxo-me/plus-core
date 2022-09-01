package tus

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PostFileReq struct {
	g.Meta `path:"/upload" method:"post" tags:"TusUpload" summary:"断点续传上传" required:"tags_name"`
}

type HeadFileReq struct {
	g.Meta `path:"/upload/{id}" method:"head" tags:"TusUpload" summary:"Head文件详情" required:"id"`
	GetByIdInput
}

type GetFileReq struct {
	g.Meta `path:"/upload/{id}" method:"get" tags:"TusUpload" summary:"续传文件下载" required:"id"`
	GetByIdInput
}

type PatchFileReq struct {
	g.Meta `path:"/upload" method:"patch" tags:"TusUpload" summary:"切片文件续传" required:"tags_name"`
}

type DelFileReq struct {
	g.Meta `path:"/upload/{id}" method:"delete" tags:"TusUpload" summary:"删除续传文件" required:"id"`
	GetByIdInput
}
