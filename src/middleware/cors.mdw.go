/*
This is a cors middleware for the gin http framwork and can be used to configure the cors behaviour of your application
*/
package middleware

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Defines the configuration of your cors middleware
*/
type Config struct {
	// All the allowed origins in an array. The default is "*"
	// The default cannot be used when AllowCredentials is true
	// [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	AllowedOrigins []string
	// All the allowed HTTP Methodes. The default is "*"
	// The default cannot be used when AllowCredentials is true
	// [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	AllowedMethods []string
	// All the allowed Headers that can be sent from the client. The default is "*"
	// The default cannot be used when AllowCredentials is true
	// Note that the Authorization header cannot be wildcarded and needs to be listed explicitly [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AllowedHeaders []string
	// The headers which should be readable by the client. The default is "*"
	// The default cannot be used when AllowCredentials is true
	// [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	ExposeHeaders []string
	// If you allow receiving cookies and Authorization headers. The default is false
	// [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	AllowCredentials bool
	// The maximum age of your preflight requests. The default is 1 day
	// [MDN Web Docs]
	//
	// [MDN Web Docs]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge time.Duration
}

/*
Adds wildcard to the array if no value was set.
Panics if the allowCredentials is true and the header is a wildcard
*/
func checkCredentials(header []string, allowCredentials bool) []string {
	if header == nil && !allowCredentials {
		header = []string{"*"}
		return header
	} else if (header == nil || slices.Contains(header, "*")) && allowCredentials {
		panic("The allowed origins must be set and cannot contain the \"*\" wildcard when AllowCredentials is true")
	} else {
		return header
	}
}

/*
Validates the config and sets empty values to their defaults if necessary
*/
func (c Config) validate() Config {
	c.AllowedOrigins = checkCredentials(c.AllowedOrigins, c.AllowCredentials)

	c.AllowedMethods = checkCredentials(c.AllowedMethods, c.AllowCredentials)

	c.AllowedHeaders = checkCredentials(c.AllowedHeaders, c.AllowCredentials)

	c.ExposeHeaders = checkCredentials(c.ExposeHeaders, c.AllowCredentials)

	if c.MaxAge == 0 {
		c.MaxAge = 24 * time.Hour
	}

	return c
}

func DefaultConfig() Config {
	return Config{}
}

func CorsMiddleware(config Config) gin.HandlerFunc {
	config = config.validate()
	return func(c *gin.Context) {
		currentOrigin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Vary", "Origin")

		if currentOrigin == "" {
			c.AbortWithStatus(http.StatusForbidden)
		}

		if !slices.Contains(config.AllowedOrigins, "*") && !slices.Contains(config.AllowedOrigins, currentOrigin) {
			c.AbortWithStatus(http.StatusForbidden)
		}

		preflight := strings.ToUpper(c.Request.Method) == "OPTIONS"
		if preflight {
			// Headers for preflight requests
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			c.Writer.Header().Set("Access-Control-Max-Age", config.MaxAge.String())
		}

		// Headers for all requests
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))

		if preflight {
			// If this is a preflight request we don't need to continue
			c.AbortWithStatus(204)
		}

		c.Next()
	}
}
