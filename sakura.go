package sakura

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
)

type Sakura struct {
}

func NewSakura() *Sakura {
	return &Sakura{}
}

func (s Sakura) Version() string {
	return "0.1.x"
}

func (s Sakura) Author() string {
	return "zhenhaowang"
}

func (s Sakura) Email() string {
	return "hit.zhenhao@gmail.com"
}

func LocalIP() (ips []string, err error) {
	var addrs []net.Addr

	if addrs, err = net.InterfaceAddrs(); nil != err {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return
}

func HttpBuildQuery(params map[string][]string) string {
	var seg []string

	for k, v := range params {
		switch len(v) {
		case 0:
			seg = append(seg, url.QueryEscape(k)+"=")
		case 1:
			seg = append(seg, strings.Join([]string{url.QueryEscape(k), url.QueryEscape(v[0])}, "="))
		default:
			for idx, vv := range v {
				kk := fmt.Sprintf("%s[%d]", k, idx)
				seg = append(seg, strings.Join([]string{url.QueryEscape(kk), url.QueryEscape(vv)}, "="))
			}
		}
	}

	return strings.Join(seg, "&")
}

func MD5(str string) string {
	bs := md5.Sum([]byte(str))
	return hex.EncodeToString(bs[:])
}

func SliceInterface(s interface{}) (r []interface{}) {
	rs := reflect.ValueOf(s)
	kind := rs.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic(&reflect.ValueError{Method: "gohelper.utils.SliceInterface", Kind: kind})
	}

	for i := 0; i < rs.Len(); i++ {
		r = append(r, rs.Index(i).Interface())
	}

	return
}

func ArrayIndex(niddle, s interface{}) int {
	slice := SliceInterface(s)

	for k, v := range slice {
		if reflect.DeepEqual(niddle, v) {
			return k
		}
	}

	return -1
}
