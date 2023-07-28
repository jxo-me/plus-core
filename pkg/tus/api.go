package tus

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"net/http"
	"strconv"
)

// PostFile creates a new file upload using the datastore after validating the
// length and parsing the metadata.
func (u *Uploader) PostFile(r *ghttp.Request) {
	ctx := context.Background()

	// Check for the presence of application/offset+octet-stream.
	// If another content
	// type is defined, it will be ignored and treated as none was set because
	// some HTTP clients may enforce a default value for this header.
	containsChunk := r.Header.Get("Content-Type") == "application/offset+octet-stream"

	// Only use the proper Upload-Concat header if the data store
	// even supports the concatenation extension.
	var concatHeader string
	if u.composer.UsesConcater {
		concatHeader = r.Header.Get("Upload-Concat")
	}

	// Parse Upload-Concat header
	isPartial, isFinal, partialUploadIDs, err := parseConcat(concatHeader)
	if err != nil {
		u.sendError(r, err)
		return
	}

	// If the upload is a final upload created by concatenation multiple partial
	// uploads, the size is sum of all sizes of these files (no need for
	// Upload-Length header)
	var size int64
	var sizeIsDeferred bool
	var partialUploads []Upload
	if isFinal {
		// A final upload must not contain a chunk within the creation request
		if containsChunk {
			u.sendError(r, ErrModifyFinal)
			return
		}

		partialUploads, size, err = u.sizeOfUploads(ctx, partialUploadIDs)
		if err != nil {
			u.sendError(r, err)
			return
		}
	} else {
		uploadLengthHeader := r.Header.Get("Upload-Length")
		uploadDeferLengthHeader := r.Header.Get("Upload-Defer-Length")
		size, sizeIsDeferred, err = u.validateNewUploadLengthHeaders(uploadLengthHeader, uploadDeferLengthHeader)
		if err != nil {
			u.sendError(r, err)
			return
		}
	}

	// Test whether the size is still allowed
	if u.config.MaxSize > 0 && size > u.config.MaxSize {
		u.sendError(r, ErrMaxSizeExceeded)
		return
	}

	// Parse metadata
	meta := ParseMetadataHeader(r.Header.Get("Upload-Metadata"))

	info := FileInfo{
		Size:           size,
		SizeIsDeferred: sizeIsDeferred,
		MetaData:       meta,
		IsPartial:      isPartial,
		IsFinal:        isFinal,
		PartialUploads: partialUploadIDs,
	}

	if u.config.PreUploadCreateCallback != nil {
		if err := u.config.PreUploadCreateCallback(newHookEvent(info, r)); err != nil {
			u.sendError(r, err)
			return
		}
	}

	upload, err := u.composer.Core.NewUpload(ctx, info)
	if err != nil {
		u.sendError(r, err)
		return
	}

	info, err = upload.GetInfo(ctx)
	if err != nil {
		u.sendError(r, err)
		return
	}

	id := info.ID

	// Add the Location header directly after creating the new resource to even
	// include it in cases of failure when an error is returned
	url := u.absFileURL(r, id)
	r.Response.Header().Set("Location", url)

	u.Metrics.incUploadsCreated()
	u.log(ctx, "UploadCreated", "id", id, "size", i64toa(size), "url", url)

	if u.config.NotifyCreatedUploads {
		u.CreatedUploads <- newHookEvent(info, r)
	}

	if isFinal {
		concatableUpload := u.composer.Concater.AsConcatableUpload(upload)
		if err := concatableUpload.ConcatUploads(ctx, partialUploads); err != nil {
			u.sendError(r, err)
			return
		}
		info.Offset = size

		if u.config.NotifyCompleteUploads {
			u.CompleteUploads <- newHookEvent(info, r)
		}
	}

	if containsChunk {
		if u.composer.UsesLocker {
			lock, err := u.lockUpload(id)
			if err != nil {
				u.sendError(r, err)
				return
			}

			defer lock.Unlock()
		}

		if err := u.writeChunk(ctx, upload, info, r); err != nil {
			u.sendError(r, err)
			return
		}
	} else if !sizeIsDeferred && size == 0 {
		// Directly finish the upload if the upload is empty (i.e. has a size of 0).
		// This statement is in an else-if block to avoid causing duplicate calls
		// to finishUploadIfComplete if an upload is empty and contains a chunk.
		if err := u.finishUploadIfComplete(ctx, upload, info, r); err != nil {
			u.sendError(r, err)
			return
		}
	}

	u.sendResp(r, http.StatusCreated)
}

// HeadFile returns the length and offset for the HEAD request
func (u *Uploader) HeadFile(r *ghttp.Request) {
	ctx := context.Background()

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		u.sendError(r, err)
		return
	}

	if u.composer.UsesLocker {
		lock, err := u.lockUpload(id)
		if err != nil {
			u.sendError(r, err)
			return
		}

		defer lock.Unlock()
	}

	upload, err := u.composer.Core.GetUpload(ctx, id)
	if err != nil {
		u.sendError(r, err)
		return
	}

	info, err := upload.GetInfo(ctx)
	if err != nil {
		u.sendError(r, err)
		return
	}

	// Add Upload-Concat header if possible
	if info.IsPartial {
		r.Response.Header().Set("Upload-Concat", "partial")
	}

	if info.IsFinal {
		v := "final;"
		for _, uploadID := range info.PartialUploads {
			v += u.absFileURL(r, uploadID) + " "
		}
		// Remove trailing space
		v = v[:len(v)-1]

		r.Response.Header().Set("Upload-Concat", v)
	}

	if len(info.MetaData) != 0 {
		r.Response.Header().Set("Upload-Metadata", SerializeMetadataHeader(info.MetaData))
	}

	if info.SizeIsDeferred {
		r.Response.Header().Set("Upload-Defer-Length", UploadLengthDeferred)
	} else {
		r.Response.Header().Set("Upload-Length", strconv.FormatInt(info.Size, 10))
		r.Response.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	}

	r.Response.Header().Set("Cache-Control", "no-store")
	r.Response.Header().Set("Upload-Offset", strconv.FormatInt(info.Offset, 10))
	u.sendResp(r, http.StatusOK)
}

// PatchFile adds a chunk to an upload. This operation is only allowed
// if enough space in the upload is left.
func (u *Uploader) PatchFile(r *ghttp.Request) {
	ctx := context.Background()

	// Check for presence of application/offset+octet-stream
	if r.Header.Get("Content-Type") != "application/offset+octet-stream" {
		u.sendError(r, ErrInvalidContentType)
		return
	}

	// Check for presence of a valid Upload-Offset Header
	offset, err := strconv.ParseInt(r.Header.Get("Upload-Offset"), 10, 64)
	if err != nil || offset < 0 {
		u.sendError(r, ErrInvalidOffset)
		return
	}

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		u.sendError(r, err)
		return
	}

	if u.composer.UsesLocker {
		lock, err := u.lockUpload(id)
		if err != nil {
			u.sendError(r, err)
			return
		}

		defer lock.Unlock()
	}

	upload, err := u.composer.Core.GetUpload(ctx, id)
	if err != nil {
		u.sendError(r, err)
		return
	}

	info, err := upload.GetInfo(ctx)
	if err != nil {
		u.sendError(r, err)
		return
	}

	// Modifying a final upload is not allowed
	if info.IsFinal {
		u.sendError(r, ErrModifyFinal)
		return
	}

	if offset != info.Offset {
		u.sendError(r, ErrMismatchOffset)
		return
	}

	// Do not proxy the call to the data store if the upload is already completed
	if !info.SizeIsDeferred && info.Offset == info.Size {
		r.Response.Header().Set("Upload-Offset", strconv.FormatInt(offset, 10))
		u.sendResp(r, http.StatusNoContent)
		return
	}

	if r.Header.Get("Upload-Length") != "" {
		if !u.composer.UsesLengthDeferrer {
			u.sendError(r, ErrNotImplemented)
			return
		}
		if !info.SizeIsDeferred {
			u.sendError(r, ErrInvalidUploadLength)
			return
		}
		uploadLength, err := strconv.ParseInt(r.Header.Get("Upload-Length"), 10, 64)
		if err != nil || uploadLength < 0 || uploadLength < info.Offset || (u.config.MaxSize > 0 && uploadLength > u.config.MaxSize) {
			u.sendError(r, ErrInvalidUploadLength)
			return
		}

		lengthDeclarableUpload := u.composer.LengthDeferrer.AsLengthDeclarableUpload(upload)
		if err := lengthDeclarableUpload.DeclareLength(ctx, uploadLength); err != nil {
			u.sendError(r, err)
			return
		}

		info.Size = uploadLength
		info.SizeIsDeferred = false

	}

	if err := u.writeChunk(ctx, upload, info, r); err != nil {
		u.sendError(r, err)
		return
	}

	u.sendResp(r, http.StatusNoContent)
}

// GetFile handles requests to download a file using a GET request. This is not
// part of the specification.
func (u *Uploader) GetFile(r *ghttp.Request) {
	ctx := context.Background()

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		u.sendError(r, err)
		return
	}

	if u.composer.UsesLocker {
		lock, err := u.lockUpload(id)
		if err != nil {
			u.sendError(r, err)
			return
		}

		defer lock.Unlock()
	}

	upload, err := u.composer.Core.GetUpload(ctx, id)
	if err != nil {
		u.sendError(r, err)
		return
	}

	info, err := upload.GetInfo(ctx)
	if err != nil {
		u.sendError(r, err)
		return
	}

	// Set headers before sending responses
	r.Response.Header().Set("Content-Length", strconv.FormatInt(info.Offset, 10))

	contentType, contentDisposition := filterContentType(info)
	r.Response.Header().Set("Content-Type", contentType)
	r.Response.Header().Set("Content-Disposition", contentDisposition)

	// If no data has been uploaded yet, respond with an empty "204 No Content" status.
	if info.Offset == 0 {
		u.sendResp(r, http.StatusNoContent)
		return
	}

	src, err := upload.GetReader(ctx)
	if err != nil {
		u.sendError(r, err)
		return
	}

	u.sendResp(r, http.StatusOK)
	_, _ = io.Copy(r.Response.ResponseWriter, src)

	// Try to close the reader if the io.Closer interface is implemented
	if closer, ok := src.(io.Closer); ok {
		_ = closer.Close()
	}
}

// DelFile terminates an upload permanently.
func (u *Uploader) DelFile(r *ghttp.Request) {
	ctx := context.Background()

	// Abort the request handling if the required interface is not implemented
	if !u.composer.UsesTerminater {
		u.sendError(r, ErrNotImplemented)
		return
	}

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		u.sendError(r, err)
		return
	}

	if u.composer.UsesLocker {
		lock, err := u.lockUpload(id)
		if err != nil {
			u.sendError(r, err)
			return
		}

		defer lock.Unlock()
	}

	upload, err := u.composer.Core.GetUpload(ctx, id)
	if err != nil {
		u.sendError(r, err)
		return
	}

	var info FileInfo
	if u.config.NotifyTerminatedUploads {
		info, err = upload.GetInfo(ctx)
		if err != nil {
			u.sendError(r, err)
			return
		}
	}

	err = u.terminateUpload(ctx, upload, info, r)
	if err != nil {
		u.sendError(r, err)
		return
	}

	u.sendResp(r, http.StatusNoContent)
}
