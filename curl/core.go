// libcurl go bingding
package curl

/*
#cgo linux pkg-config: libcurl
#cgo windows LDFLAGS: -lcurl
#include <stdlib.h>
#include <curl/curl.h>

#ifndef CURL_VERSION_SSL
#define CURL_VERSION_SSL       (1<<2)
#endif
#ifndef CURL_VERSION_LIBZ
#define CURL_VERSION_LIBZ      (1<<3)
#endif
#ifndef CURL_VERSION_NTLM
#define CURL_VERSION_NTLM      (1<<4)
#endif
#ifndef CURL_VERSION_GSSNEGOTIATE
#define CURL_VERSION_GSSNEGOTIATE (1<<5)
#endif
#ifndef CURL_VERSION_DEBUG
#define CURL_VERSION_DEBUG     (1<<6)
#endif
#ifndef CURL_VERSION_ASYNCHDNS
#define CURL_VERSION_ASYNCHDNS (1<<7)
#endif
#ifndef CURL_VERSION_SPNEGO
#define CURL_VERSION_SPNEGO    (1<<8)
#endif
#ifndef CURL_VERSION_LARGEFILE
#define CURL_VERSION_LARGEFILE (1<<9)
#endif
#ifndef CURL_VERSION_IDN
#define CURL_VERSION_IDN       (1<<10)
#endif
#ifndef CURL_VERSION_SSPI
#define CURL_VERSION_SSPI      (1<<11)
#endif
#ifndef CURL_VERSION_CONV
#define CURL_VERSION_CONV      (1<<12)
#endif
#ifndef CURL_VERSION_CURLDEBUG
#define CURL_VERSION_CURLDEBUG (1<<13)
#endif
#ifndef CURL_VERSION_TLSAUTH_SRP
#define CURL_VERSION_TLSAUTH_SRP (1<<14)
#endif
#ifndef CURL_VERSION_NTLM_WB
#define CURL_VERSION_NTLM_WB   (1<<15)
#endif

static char *string_array_index(char **p, int i) {
  return p[i];
}
*/
import "C"

import (
	"time"
	"unsafe"
)


// curl_global_init - Global libcurl initialisation
func GlobalInit(flags int) error {
	return newCurlError(C.curl_global_init(C.long(flags)))
}

// curl_global_cleanup - global libcurl cleanup
func GlobalCleanup() {
	C.curl_global_cleanup()
}

type VersionInfoData struct {
	Age C.CURLversion
	// age >= 0
	Version       string
	VersionNum    uint
	Host          string
	Features      int
	SslVersion    string
	SslVersionNum int
	LibzVersion   string
	Protocols     []string
	// age >= 1
	Ares    string
	AresNum int
	// age >= 2
	Libidn string
	// age >= 3
	IconvVerNum   int
	LibsshVersion string
}

// curl_version - returns the libcurl version string
func Version() string {
	return C.GoString(C.curl_version())
}

// curl_version_info - returns run-time libcurl version info
func VersionInfo(ver C.CURLversion) *VersionInfoData {
	data := C.curl_version_info(ver)
	ret := new(VersionInfoData)
	ret.Age = data.age
	switch age := ret.Age; {
	case age >= 0:
		ret.Version = string(C.GoString(data.version))
		ret.VersionNum = uint(data.version_num)
		ret.Host = C.GoString(data.host)
		ret.Features = int(data.features)
		ret.SslVersion = C.GoString(data.ssl_version)
		ret.SslVersionNum = int(data.ssl_version_num)
		ret.LibzVersion = C.GoString(data.libz_version)
		// ugly but works
		ret.Protocols = []string{}
		for i := C.int(0); C.string_array_index(data.protocols, i) != nil; i++ {
			p := C.string_array_index(data.protocols, i)
			ret.Protocols = append(ret.Protocols, C.GoString(p))
		}
		fallthrough
	case age >= 1:
		ret.Ares = C.GoString(data.ares)
		ret.AresNum = int(data.ares_num)
		fallthrough
	case age >= 2:
		ret.Libidn = C.GoString(data.libidn)
		fallthrough
	case age >= 3:
		ret.IconvVerNum = int(data.iconv_ver_num)
		ret.LibsshVersion = C.GoString(data.libssh_version)
	}
	return ret
}

// curl_getdate - Convert a date string to number of seconds since January 1, 1970
// In golang, we convert it to a *time.Time
func Getdate(date string) *time.Time {
	datestr := C.CString(date)
	defer C.free(unsafe.Pointer(datestr))
	t := C.curl_getdate(datestr, nil)
	if t == -1 {
		return nil
	}
    unix := time.Unix(int64(t), 0).UTC()
	return &unix

	/*
		// curl_getenv - return value for environment name
		func Getenv(name string) string {
			namestr := C.CString(name)
			defer C.free(unsafe.Pointer(namestr))
			ret := C.curl_getenv(unsafe.Pointer(namestr))
			defer C.free(unsafe.Pointer(ret))

			return C.GoString(ret)
		}
	*/
}

// TODO: curl_global_init_mem
