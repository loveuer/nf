# Ursa Web Framework

Ursa (🐻 Great Bear) is a fast, simple, and production-ready web framework for Go.

### Installation

```bash
go get github.com/loveuer/ursa
```

### Features

- ⚡ High performance with radix tree routing
- 🛡️ Built-in security middlewares (CORS, Secure headers)
- ⏱️ Configurable server timeouts
- 🔍 Request tracking with RequestID
- 📦 Zero-allocation optimizations
- 🧪 Comprehensive test coverage

### Usage

##### Basic usage

- Get route parameters

  ```go
  func main() {
      app := ursa.New()

      app.Get("/hello/:name", func(c *ursa.Ctx) error {
          name := c.Param("name")
          return c.JSON(ursa.Map{"status": 200, "data": "hello, " + name})
      })

      log.Fatal(app.Run("0.0.0.0:80"))
  }
  ```

- Parse request query

  ```go
  func handleQuery(c *ursa.Ctx) error {
      type Req struct {
          Name string   `query:"name"`
          Addr []string `query:"addr"`
      }

      var (
          err error
          req = Req{}
      )

      if err = c.QueryParser(&req); err != nil {
          return ursa.NewNFError(400, err.Error())
      }

      return c.JSON(ursa.Map{"query": req})
  }
  ```

- Parse application/json body

  ```go
  func handlePost(c *ursa.Ctx) error {
      type Req struct {
          Name string   `json:"name"`
          Addr []string `json:"addr"`
      }

      var (
          err error
          req = Req{}
          reqMap = make(map[string]interface{})
      )

      if err = c.BodyParser(&req); err != nil {
          return ursa.NewNFError(400, err.Error())
      }

      // can parse body multiple times
      if err = c.BodyParser(&reqMap); err != nil {
          return ursa.NewNFError(400, err.Error())
      }

      return c.JSON(ursa.Map{"struct": req, "map": reqMap})
  }
  ```

- Pass local values between middlewares

  ```go
  type User struct {
      Id int
      Username string
  }

  func main() {
      app := ursa.New()
      app.Use(auth())

      app.Get("/item/list", list)
  }

  func auth() ursa.HandlerFunc {
      return func(c *ursa.Ctx) error {
          c.Locals("user", &User{Id: 1, Username:"user"})
          return c.Next()
      }
  }

  func list(c *ursa.Ctx) error {
      user, ok := c.Locals("user").(*User)
      if !ok {
          return c.Status(401).SendString("login required")
      }

      // ...
      return nil
  }
  ```

### Middlewares

Ursa comes with built-in production-ready middlewares:

```go
app := ursa.New()

// CORS
app.Use(ursa.NewCORS())

// Security headers
app.Use(ursa.NewSecure())

// Request tracking
app.Use(ursa.NewRequestID())
```

### License

MIT
