// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package compression

import (
	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
)

// Supported compression levels
const (
	LevelNoCompression CompressLevel = iota - 1
	LevelDefaultCompression
	LevelBestSpeed
	LevelBestCompression
	LevelHuffmanOnly
)

type CompressLevel int

// Config ...
type Config struct {
	// Exclude defines a function to skip middleware.
	// Optional. Default: nil
	Exclude func(*fiber.Ctx) bool
	// Level of compression
	// Optional. Default value 0.
	Level CompressLevel
}

func FixConfig(configs []Config) Config {
	var cfg Config
	if len(configs) > 0 {
		cfg = configs[0]
	}
	return cfg
}

// Convert compress levels to correct int
// https://github.com/valyala/fasthttp/blob/master/compress.go#L17
func GetCompressHandlerLevel(level CompressLevel) int {
	switch level {
	case LevelNoCompression:
		return 0
	case LevelDefaultCompression:
		return 6
	case LevelBestSpeed:
		return 1
	case LevelBestCompression:
		return 9
	case LevelHuffmanOnly:
		return -2
	}
	return 6
}

// New ...
func New(configs ...Config) func(*fiber.Ctx) {
	// Init config
	cfg := FixConfig(configs)
	doNothing := func(c *fasthttp.RequestCtx) { return }
	level := GetCompressHandlerLevel(cfg.Level)
	compress := fasthttp.CompressHandlerLevel(doNothing, level)
	// Middleware function
	return func(c *fiber.Ctx) {
		c.Next()
		// Run only the request is not exclude
		if cfg.Exclude == nil || !cfg.Exclude(c) {
			compress(c.Fasthttp)
		}
	}
}
