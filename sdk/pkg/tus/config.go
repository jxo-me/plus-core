package tus

import (
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"net/url"
)

// Config provides a way to configure the Handler depending on your needs.
type Config struct {
	// StoreComposer points to the store composer from which the core data store
	// and optional dependencies should be taken. May only be nil if DataStore is
	// set.
	// TODO: Remove pointer?
	StoreComposer *StoreComposer `yaml:"-" json:"-"`
	// MaxSize defines how many bytes may be stored in one single upload. If its
	// value is is 0 or smaller no limit will be enforced.
	MaxSize int64 `yaml:"maxSize" json:"max_size"`
	// BasePath defines the URL path used for handling uploads, e.g. "/files/".
	// If no trailing slash is presented it will be added. You may specify an
	// absolute URL containing a scheme, e.g. "http://tus.io"
	BasePath string `yaml:"basePath" json:"base_path"`
	IsAbs    bool   `yaml:"isAbs" json:"is_abs"`
	// DisableDownload indicates whether the server will refuse downloads of the
	// uploaded file, by not mounting the GET handler.
	DisableDownload bool `yaml:"disableDownload" json:"disable_download"`
	// DisableTermination indicates whether the server will refuse termination
	// requests of the uploaded file, by not mounting the DELETE handler.
	DisableTermination bool `yaml:"disableTermination" json:"disable_termination"`
	// NotifyCompleteUploads indicates whether sending notifications about
	// completed uploads using the CompleteUploads channel should be enabled.
	NotifyCompleteUploads bool `yaml:"notifyCompleteUploads" json:"notify_complete_uploads"`
	// NotifyTerminatedUploads indicates whether sending notifications about
	// terminated uploads using the TerminatedUploads channel should be enabled.
	NotifyTerminatedUploads bool `yaml:"notifyTerminatedUploads" json:"notify_terminated_uploads"`
	// NotifyUploadProgress indicates whether sending notifications about
	// the upload progress using the UploadProgress channel should be enabled.
	NotifyUploadProgress bool `yaml:"notifyUploadProgress" json:"notify_upload_progress"`
	// NotifyCreatedUploads indicates whether sending notifications about
	// the upload having been created using the CreatedUploads channel should be enabled.
	NotifyCreatedUploads bool `yaml:"notifyCreatedUploads" json:"notify_created_uploads"`
	// Logger is the Logger to use internally, mostly for printing requests.
	Logger *glog.Logger `yaml:"-" json:"-"`
	// Respect the X-Forwarded-Host, X-Forwarded-Proto and Forwarded headers
	// potentially set by proxies when generating an absolute URL in the
	// response to POST requests.
	RespectForwardedHeaders bool `yaml:"respectForwardedHeaders" json:"respect_forwarded_headers"`
	// PreUploadCreateCallback will be invoked before a new upload is created, if the
	// property is supplied. If the callback returns nil, the upload will be created.
	// Otherwise the HTTP request will be aborted. This can be used to implement
	// validation of upload metadata etc.
	PreUploadCreateCallback func(hook HookEvent) error `yaml:"-" json:"-"`
	// PreFinishResponseCallback will be invoked after an upload is completed but before
	// a response is returned to the client. Error responses from the callback will be passed
	// back to the client. This can be used to implement post-processing validation.
	PreFinishResponseCallback func(hook HookEvent) error `yaml:"-" json:"-"`

	Path string `yaml:"path" json:"path"`
	// log
	LogPath   string `yaml:"logPath" json:"log_path"`
	LogFile   string `yaml:"logFile" json:"log_file"`
	LogLevel  string `yaml:"logLevel" json:"log_level"`
	LogStdout bool   `yaml:"logStdout" json:"log_stdout"`
}

func (config *Config) validate() error {
	if config.Logger == nil {
		config.Logger = g.Log()
	}

	base := config.BasePath
	uri, err := url.Parse(base)
	if err != nil {
		return err
	}

	// Ensure base path ends with slash to remove logic from absFileURL
	if base != "" && string(base[len(base)-1]) != "/" {
		base += "/"
	}

	// Ensure base path begins with slash if not absolute (starts with scheme)
	if !uri.IsAbs() && len(base) > 0 && string(base[0]) != "/" {
		base = "/" + base
	}
	config.BasePath = base
	config.IsAbs = uri.IsAbs()

	if config.StoreComposer == nil {
		return errors.New("tus: StoreComposer must no be nil")
	}

	if config.StoreComposer.Core == nil {
		return errors.New("tus: StoreComposer in Config needs to contain a non-nil core")
	}

	return nil
}
