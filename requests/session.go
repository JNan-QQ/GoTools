package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

// Session http.client 简单功能封装
type Session struct {
	session *http.Client
	header  map[string]string
}

// Response 请求结果
type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200

	// Header 将标头键映射到值。如果响应具有具有相同键的多个标头，则可以使用逗号分隔符将它们连接起来。
	//（RFC 7230 第 3.2.2 节要求多个标头在语义上等效于逗号分隔的序列。
	// 当标头值被此结构中的其他字段（例如，ContentLength、TransferEncoding、Trailer）复制时，字段值是权威的。
	//
	// 映射中的键是规范化的（请参阅 CanonicalHeaderKey）
	Header http.Header

	Body []byte

	Cookies []*http.Cookie
}

// NewSession 创建 Session 对象。可以自定义请求头
func NewSession(defaultHeader map[string]string) *Session {
	// 设置默认请求头
	if defaultHeader == nil {
		defaultHeader = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
				"Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.62",
		}
	}

	// 创建session缓存对象
	jar, _ := cookiejar.New(nil)
	return &Session{session: &http.Client{Jar: jar}, header: defaultHeader}
}

// Proxy 设置代理，例：http://127.0.0.1:8888
func (s *Session) Proxy(proxy string) {
	p, _ := url.Parse(proxy)
	s.session.Transport = &http.Transport{
		// 设置代理
		Proxy: http.ProxyURL(p),
	}
}

// 设置请求cookies
func (s *Session) cookies(cookies string) {
	s.header["Cookie"] = cookies
}

// Do 发送请求
//
//	method:	请求方法 MethodGet , MethodPost , MethodPut , MethodDelete
//	params: GET请求参数，可选 字符串（string）/ 键值对（map[string]any）/ nil 例：a=5&b=6&c="acc" / map[string]any{ "a":5, "b":6, "c":"acc" }
//	data:	用于提交 x-www-form-urlencoded 数据, 可选 string / []byte / map[string]any / io.Reader / nil
//	_json:	用于提交 json 数据, 可选 string / []byte / map[string]any / io.Reader / nil
func (s *Session) Do(method, url string, params, _data, _json any) (*Response, error) {

	var body io.Reader
	cType := CJson

	// 判断格式类型
	switch method {

	case MethodGet:
		// 格式化GET请求参数
		if params != nil {
			switch v := params.(type) {
			case string:
				if strings.HasPrefix(v, "?") {
					url += v
				} else {
					url += "?" + v
				}
			case map[string]any:
				var p []string
				for key, value := range v {
					p = append(p, fmt.Sprintf("%s=%v", key, value))
				}

				if len(p) > 0 {
					url += "?" + strings.Join(p, "&")
				}
			default:
				return nil, fmt.Errorf("不支持的参数格式")
			}
		}
		cType = CFormData
		body = nil

	case MethodPost:
		if _data != nil {
			switch v := _data.(type) {
			case string:
				body = strings.NewReader(v)

			case map[string]any:
				body = strings.NewReader(FormatUrlValues(v).Encode())

			case []byte:
				body = strings.NewReader(string(v))

			case io.Reader:
				body = v

			default:
				return nil, fmt.Errorf("不支持的参数格式")
			}
			cType = CUrlencoded

		} else if _json != nil {
			switch v := _json.(type) {
			case string:
				body = strings.NewReader(v)

			case map[string]string, map[string]int, map[string]float64:
				bytes, err := json.MarshalIndent(v, "", "\t")
				if err != nil {
					return nil, err
				}
				body = strings.NewReader(string(bytes))

			case []byte:
				body = strings.NewReader(string(v))

			case io.Reader:
				body = v

			default:
				return nil, fmt.Errorf("不支持的参数格式")
			}
			cType = CJson

		} else {
			body = nil
			cType = CJson
		}
	}

	// 创建请求
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	s.header["Content-Type"] = cType
	for key, value := range s.header {
		request.Header.Set(key, value)
	}

	// 发送请求
	response, err := s.session.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 获取结果
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// 提取返回值
	res := new(Response)
	res.Status = response.Status
	res.StatusCode = response.StatusCode
	res.Header = response.Header
	res.Body = b
	// 获取cookies
	cookies := response.Cookies()
	res.Cookies = make([]*http.Cookie, len(cookies))
	copy(res.Cookies, cookies)

	return res, nil

}

// Post 会话请求
//
//	data:	用于提交 x-www-form-urlencoded 数据, 可选 string / []byte / map[string]any / io.Reader / nil
//	_json:	用于提交 json 数据, 可选 string / []byte / map[string]any / io.Reader / nil
func (s *Session) Post(url string, _data, _json any) (*Response, error) {
	do, err := s.Do(MethodPost, url, nil, _data, _json)
	if err != nil {
		return nil, err
	}
	return do, nil
}

// Get 会话请求
//
//	params: GET请求参数，可选 字符串（string）/ 键值对（map[string]any）/ nil
func (s *Session) Get(url string, params any) (*Response, error) {
	do, err := s.Do(MethodGet, url, params, nil, nil)
	if err != nil {
		return nil, err
	}
	return do, nil
}

// Post 会话请求 cookies 可选 string / map[string]string
func Post(url string, _data, _json any, header map[string]string, proxy string, cookies any) (*Response, error) {
	// 创建请求对象
	session := NewSession(header)

	// 设置代理
	if proxy != "" {
		session.Proxy(proxy)
	}

	// 设置cookies
	if cookies != nil {
		switch v := cookies.(type) {
		case string:
			session.cookies(v)
		case map[string]string:
			var cookie []string
			for key, value := range v {
				cookie = append(cookie, fmt.Sprintf("%s=%s", key, value))
			}
			session.cookies(strings.Join(cookie, ";"))
		}
	}

	do, err := session.Post(url, _data, _json)
	if err != nil {
		return nil, err
	}
	return do, nil
}

// Get 会话请求 cookies 可选 string / map[string]string
func Get(url string, params any, header map[string]string, proxy string, cookies any) (*Response, error) {
	// 创建请求对象
	session := NewSession(header)

	// 设置代理
	if proxy != "" {
		session.Proxy(proxy)
	}

	// 设置cookies
	if cookies != nil {
		switch v := cookies.(type) {
		case string:
			session.cookies(v)
		case map[string]string:
			var cookie []string
			for key, value := range v {
				cookie = append(cookie, fmt.Sprintf("%s=%s", key, value))
			}
			session.cookies(strings.Join(cookie, ";"))
		}
	}
	do, err := session.Get(url, params)
	if err != nil {
		return nil, err
	}
	return do, nil
}

const (
	CUrlencoded  = "application/x-www-form-urlencoded"
	CJson        = "application/json"
	CFormData    = "application/form-data"
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

// FormatUrlValues map to CUrlencoded 数据
func FormatUrlValues(b map[string]any) url.Values {
	formData := url.Values{}
	for key, value := range b {
		switch v := value.(type) {
		case string:
			formData.Set(key, v)
		case int:
			formData.Set(key, strconv.Itoa(v))
		case []string:
			for _, s := range v {
				formData.Add(key+"[]", s)
			}
		case []int:
			for _, s := range v {
				formData.Add(key+"[]", strconv.Itoa(s))
			}
		case []map[string]string:
			for i, s := range v {
				for k1, v1 := range s {
					formData.Set(fmt.Sprintf("%s[%d][%s]", key, i, k1), v1)
				}
			}
		case []map[string]int:
			for i, s := range v {
				for k1, v1 := range s {
					formData.Set(fmt.Sprintf("%s[%d][%s]", key, i, k1), strconv.Itoa(v1))
				}
			}
		}
	}

	return formData
}
