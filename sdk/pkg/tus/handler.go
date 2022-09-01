package tus

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	Handler = hTus{}
)

type hTus struct{}

// PostFile creates a new file upload using the datastore after validating the
// length and parsing the metadata.
func (h *hTus) PostFile(ctx context.Context, req *PostFileReq) (result *NullStructRes, err error) {
	r := g.RequestFromCtx(ctx)
	//sdk.Runtime.Config()
	s, err := NewTus(Config{})
	if err != nil {
		return nil, err
	}
	s.PostFile(ctx, r.Response.RawWriter(), r.Request)
	return nil, nil
}

// HeadFile returns the length and offset for the HEAD request
func (h *hTus) HeadFile(ctx context.Context, req *HeadFileReq) (result *NullStructRes, err error) {
	r := g.RequestFromCtx(ctx)
	s, err := NewTus(Config{})
	if err != nil {
		return nil, err
	}
	s.HeadFile(ctx, r.Response.RawWriter(), r.Request)
	return nil, nil
}

// GetFile handles requests to download a file using a GET request. This is not
// part of the specification.
func (h *hTus) GetFile(ctx context.Context, req *GetFileReq) (result *NullStructRes, err error) {
	r := g.RequestFromCtx(ctx)
	s, err := NewTus(Config{})
	if err != nil {
		return nil, err
	}
	s.GetFile(ctx, r.Response.RawWriter(), r.Request)
	return nil, nil
}

// PatchFile adds a chunk to an upload. This operation is only allowed
// if enough space in the upload is left.
func (h *hTus) PatchFile(ctx context.Context, req *PatchFileReq) (result *NullStructRes, err error) {
	r := g.RequestFromCtx(ctx)
	s, err := NewTus(Config{})
	if err != nil {
		return nil, err
	}
	s.PatchFile(ctx, r.Response.RawWriter(), r.Request)
	return nil, nil
}

// DelFile terminates an upload permanently.
func (h *hTus) DelFile(ctx context.Context, req *DelFileReq) (result *NullStructRes, err error) {
	r := g.RequestFromCtx(ctx)
	s, err := NewTus(Config{})
	if err != nil {
		return nil, err
	}
	s.DelFile(ctx, r.Response.RawWriter(), r.Request)
	return nil, nil
}
