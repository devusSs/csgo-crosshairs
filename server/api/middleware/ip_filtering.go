package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/system"
	"github.com/gin-gonic/gin"
)

var (
	AllowedDomain string

	privateIPBlocks []*net.IPNet
)

// https://gist.github.com/nanmu42/9c8139e15542b3c4a1709cb9e9ac61eb
func SetupPrivateIPBlock() error {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}

		privateIPBlocks = append(privateIPBlocks, block)
	}

	return nil
}

// https://gist.github.com/nanmu42/9c8139e15542b3c4a1709cb9e9ac61eb
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func CheckAllowedHostMiddleware(c *gin.Context) {
	if AllowedDomain == "*" {
		c.Next()
		return
	}

	allowedIP, err := system.GetIPForDynamicHost(AllowedDomain)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	remoteIP := net.ParseIP(c.RemoteIP())

	if c.RemoteIP() != allowedIP.String() && !isPrivateIP(remoteIP) && !remoteIP.IsPrivate() {
		// TODO: maybe add this to errors / database log and make it available for query for admins / engineers?
		log.Printf("%s Invalid IP / domain %s tried accessing ressource %s\n", logging.WarnSign, c.RemoteIP(), c.Request.RequestURI)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Your IP address is not allowed on this ressource."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	c.Next()
}
