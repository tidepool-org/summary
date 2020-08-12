// Package Summary provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package Summary

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"strings"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9Ra63Ibt5J+FRQ2VZYUXoa0bMf8s2tbWa8qtuUy5WxVJJ4jcKaHRIIBxgCGEqOwKg+R",
	"J8yTnGpgrpwhJVk5lTp/JHIGaPTl668bAG9pqJJUSZDW0MktNeESEuY+TrMkYXpNNJhUSQP4jEURt1xJ",
	"Jj5qlYK2HAydxEwY6NEIYpYJSye3G/xiQs1THEwn9BUxuTQVE8GMJRFnc7BgSMQsI7RH05rAW8pCy1fc",
	"rvHzNzyiE/pfw2rEsHzdq6l0aiHp0OZitleb6yUPlySCFQ/BkGvQQDIDEWEyItdLkMQqkqVCsaipM+1R",
	"7he80y3NtTWkSluipBfPamKdEuSaGZIoY4mGEKQVa6dQy0Wh4CDR21srnC+B+HfEqNheMw3ELpklqVYr",
	"HkHbjqZcyRLolopv0GW2YwXnM6sI3FjNQuvGFOKZtaBRyD8uWP/XV/2fgv7L2cFv1Zf+7DboPR9taq8P",
	"Dy4vBw8YfvjtN7RHLbcC6IS+2dLO2dSjdp3iW2M1lwu66dFUMBsrnXTbW04vhvmgIV66XIBx05mkPQo3",
	"LEkFuvOCJiw8m5JRMBg9GxzT2W4dS1269NR8xawLSzH9PViWOzgfruY/Q2hpj970LVu4xd2A2aZHV6CN",
	"M6zLzvzl40J7EPx2Meq/nF1eRkeHl5eDvd8P+vWvv+GfPJj92UXQf1l8PnIwuO/Yw6PD/z64vPy2/vRb",
	"h6P6Axy1DyuFq1ph2FRzTooM2vZVysJf2AI6g1Kh4rZASH9UZRw95xGkSgny2dENaFoH6BaQaiGlo8H4",
	"6eCYbjaOd5FDHsPViIiciYrI58xxB3H4Se+ZzGIW2kyDNm28vZKEac0c9XrHEgQr4TLiIbPu+xJIUpNy",
	"YA4LaPolBpfyUp5KonQEGvVjK8UjEioZZw7IGkwmnKxYq4RoiEH7pRRhxHC5EM0lyDW3S5IoT5WSKOk5",
	"48/f/4iVJnm0eiTDuWSu7JI8ec8lTyB64irFk/cQWa0kD58QLi3ocMnkAthcrP/8/Y9rQJWs5nn2GLCV",
	"/aZ0s9MTP9RVM/igmCS4sVib5moFbl2QsdIhkCO4YaE9KlyaMBsuwZADLkORRfgoZKj7ofPdVUeorgg3",
	"hNWigz5Q2jslF3tJMVaXlMwhZJlBZ2G9ROAXBXTL9YyESgg2V5ohAMgc7DWA3HJ13d4eMVm4JMw4o8+Z",
	"jCAhb48JlyYTXJI0S1Ifrjdv3ztfL3LZB7sWu/JSrpzHrk7gJlTJ1eGgXsFbjJs/cN6gZV69VxGIDlAX",
	"DuIRSMvjdQljHN8BXsyxq5pI73wilexDktp1mRpYt0GGKgLTFJgLI/8PBGQEbKVy5NhwSYCFy4ZTnxhi",
	"LJMR0xFBQC/VNYrD5i7VYJDCvGRX4rkkFnRicB0Pmx72KhhskjJtTWGRG43RL8UgOI2L6bXSEeLHXqse",
	"ARsOuuqat2IKmjPxIUvmoPc41/mh8KyfiXa5yUS62dvOrUvOfbzLsw05WyEjH5TNmyhYgSQ8JqzJH0r6",
	"Ns2QiC+4RT4j+K8h1fSc0OkHYpYqExGZY2Yp7Z1WqqZhwXQkwBhnzmfJv2Qg6owZcxCRW8PlSFE3ctZ2",
	"HbULp+l1O+IoYesjXJ1JUsfbgJwvuUE/gcFocibE2uVwkmqVcAMTsmSrZgy2PIeTb6yGBFXmCfa6TFpy",
	"ACaF0At0lDpYDHpYPiUPmSDGZhEHc0jmmUW7xsHoObkGEjqEhVoZQ1SmScy1sQhKQERi0wXkKGRSKnuE",
	"/bJaIYPuD6cP5AIkaIaALfqZCfHUQN4+8z70bb/n5tLJdqlVtliSV2kqgPCz6RND/g+YsMsfuK3z0R7A",
	"n7v27CtLY5zJEGd0l8Ua+HGVq1xkBbiYcS3WxICI+3CTCiaZVXrtEFhT4B4rM0cGPMwE02Wm/G+9XOKk",
	"U+RtsOQskfyjikoej0DwFeg1MWtjISHLnPHdwlfzRUHX+fg+8v4V1m4XQCAfT94j1lwtZrJZHkIlrVZC",
	"gHYyfBEETLF5xoXtc0nmQqmILEQWKoO0KrlVulESmFyfxXRycUtBZgl20/MFNueh+1tXi842s95dJaTW",
	"Pfr+jpz4Vm27Vdz0KJKMfWwPZ3mSQz3/A/UNptuEP6y/Q4E79uPuVU0jSndqlEO2tsHN82ywtXEaB+Og",
	"Hzzvj4Pz0fFkHEyC8U+4f8J2mOEaEbPQz1cuXPsOzcr9e867N33+Qbcd+OpOOxw54FDvWXRnQRXOiBwu",
	"JpkvEC/zBapdN829mW0DAiJygpLPvRLt3YeGLxnXEKGIXFNn/qxlP7pmy/77bER8hHPPv+gHz8+D7ybj",
	"l87zlULeLPvrWRwbQCf1gxdBgNuPGsY/+dDmoXhVnde0wL5/BsmPoWi7JfMHKR1EOq1Od3KEFXnOpUdP",
	"Ts/3PcDZmWfFSmXiGMssN5aHpn1eo5I0sxi9TuiVr+8DP2ieRBTHWYjEUANWte5s2orpXdn0ETRXEXnj",
	"VPM99QmznUkFTAsOxp4wmyU7THTF+5+RG3FPKz1ZFMKJcQYVhB4VHZApMfJoAikimq9kLNO2y17B7B3W",
	"IsN+jbFe8teY+qIfHJ+PRpPxs8n4+MGmgow6j54cBnbYmL/sPapQeSHEtW5tTJuUSdM095ZmkmPi02uA",
	"X2iPrpjIgE7GWIabSZcP3KN8/YkffWe0cBhGKwFmMux0pXU9rS2NqVWBiGGwnKJbVQCfmTp55+lWaNEK",
	"RW7mbaVf0KVe1fM6SWSuMokImq+3ELS1Lnq6WtZ1sqAbjUuBmNLIFpkj/+3yt3/3VSf253VAlBS75XSy",
	"4CuQW8R3cUvZCjRbwFtfA2qooMliGIkKP8cvN76inspPTC4aQ1PQWEnYAs7iqZd+Kqf5ntqPLuV892Kz",
	"mc3+gvJSXVQUlwZxq0dr5MXfY2rr8qalxVc7wEsq67fjCM9XbvtbZVsXO7Rt3s0PRbomiRLDd7TnP+Bs",
	"lPIu/y/q6frabSFyK8nnzqztPoy/O4+3DXdT0GKf0sVVTtPy45c17TolVPp5kmi1lcUo75fd4robuQac",
	"vjbmH0v4uY1wFWydQ68z1g8E7d1AuIfAmnuw10Yd3YsdWLg/g7NEZdJlu83lOtu3qO27F/UKHwvFLN2l",
	"0Y+Pi35DWDv2HSXCE9a+fn/a4nO6Z8+8b/R9tjZV971re7PVwXa3j1uNXyXr2floXA2qmqZdnUpZKv82",
	"1u7ark2dd92JnHc4Zkd755UZ0Nz3hBpiV+ernxEM898QDKv7LAOanJ7swEn+64IHRbP6ecBF/SL8r7pI",
	"q1+jdd5qXdDXbA2azoqRH/zKr0wI0nBG3ihpVabJB7ix5Ez6Q4X2WTcdjZ8eP3veOPN55CZ8VtsYX/z7",
	"cV/fcPyH4N4jv4XO9jFCcRyf3234DbA7+GxN7tGE3bwDubBLOhkFPZpwWf/a/A1EHPRfzm5HwabB1y2R",
	"7RMgA2GmuV1PMcV8XsyBadDn6heQ5Q94cJZ/XklZWpu6Dh30CvQUQg0ObhyVWkKeJnkG3fRtrk7fj+8b",
	"P6HigpT/AGsv0GD6lBrcLdCN71s3oSUQzeQyVp60pWWhUxMSxgU6I0sR3P9TiBsovaiWefXxlEz9CCxj",
	"WuSGm8lwmM8c1GcO3VaeY9JCjUBeT0/IuP9GuDvOSzrlSSp4zCG6pORdPnpb/ILbZTYfhCoZlqbiCgXx",
	"DOdCzYcJMxb08N3pm+8/TL935RB0Ys7iKWjPOKXACFYgsEVpauzG91XczwwM27v5LT4LcAmVgmQppxP6",
	"dBAMAv+7jaVDz3A1GvpLGDO89R94tBmakvwnt3QBHb8y+gQ207Lcn9QO1tyujAlBUmaRmI1PGS+cOm38",
	"BckpktJbsD+O3ngN3D8eVZYUxcGpMQ6CAhPF6XiaCnc/oeTwZ6Mq/DOHw2L3ta9GtcpQ+9y+1aCe/dDI",
	"RUdlzSS4cBcCxUbbOUtzWBV7WA5mv5M2GCHNErC+2uRJhVGrsF5Ei9Z7OKsz6NWcsE0huP3AmGMFN8Nb",
	"X8iLeK+/OtoEBXXHFunMfHbr/EWBfVg8u+PXo8d+1ear1ywin+BLBsY+MsZrf0uA3nG+uVdI877qYQEt",
	"OT3fEJWsNBkOhQqZWCpjJ0+DIGgdZpWv6aZ3u8VnX9hokFNQArJJmy1BX9ioU8T4ISLGHSLq96lwgyWU",
	"if1yajM65LGU75+eahVloZ892/wrAAD//66+L2QfKwAA",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}

