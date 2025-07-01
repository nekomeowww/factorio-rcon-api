package interceptors

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/nekomeowww/factorio-rcon-api/v2/pkg/apierrors"
	"github.com/nekomeowww/fo"
	"google.golang.org/grpc/metadata"
)

func MetadataCookie() func(context.Context, *http.Request) metadata.MD {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		md := metadata.MD{}

		for _, cookie := range r.Cookies() {
			md.Append("header-cookie-"+cookie.Name, string(fo.May(json.Marshal(http.Cookie{
				Name:       cookie.Name,
				Value:      cookie.Value,
				Path:       cookie.Path,
				Domain:     cookie.Domain,
				Expires:    cookie.Expires,
				RawExpires: cookie.RawExpires,
				MaxAge:     cookie.MaxAge,
				Secure:     cookie.Secure,
				HttpOnly:   cookie.HttpOnly,
				SameSite:   cookie.SameSite,
				Raw:        cookie.Raw,
				Unparsed:   cookie.Unparsed,
			}))))
		}

		return md
	}
}

type Cookies []*http.Cookie

func (c Cookies) Cookie(name string) *http.Cookie {
	for _, cookie := range c {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

func CookiesFromContext(ctx context.Context) (Cookies, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, apierrors.NewErrInternal().WithError(errors.New("failed to get metadata from context")).WithCaller().AsStatus()
	}

	return CookiesFromMetadata(md)
}

func CookiesFromMetadata(md metadata.MD) (Cookies, error) {
	var cookies Cookies

	for k, v := range md {
		if len(v) == 0 {
			continue
		}

		if strings.HasPrefix(k, "header-cookie-") {
			var cookie http.Cookie

			err := json.Unmarshal([]byte(v[0]), &cookie)
			if err != nil {
				return nil, err
			}

			cookies = append(cookies, &cookie)
		}
	}

	return cookies, nil
}
