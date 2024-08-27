# NF Web Framework

### Usage

##### basic usage

- get param

  ```go
  func main() {
      app := nf.New()

      app.Get("/hello/:name", func(c *nf.Ctx) error {
          name := c.Param("name")
          return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
      })

      log.Fatal(app.Run("0.0.0.0:80"))
  }
  ```

- parse request query

  ```go
  func handleQuery(c *nf.Ctx) error {
      type Req struct {
          Name string   `query:"name"`
          Addr []string `query:"addr"`
      }

      var (
          err error
          req = Req{}
      )

      if err = c.QueryParser(&req); err != nil {
          return nf.NewNFError(400, err.Error())
      }

      return c.JSON(nf.Map{"query": req})
  }
  ```

- parse application/json body

  ```go
  func handlePost(c *nf.Ctx) error {
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
          return nf.NewNFError(400, err.Error())
      }

      // can parse body multi times
      if err = c.BodyParser(&reqMap); err != nil {
          return nf.NewNFError(400, err.Error())
      }

      return c.JSON(nf.Map{"struct": req, "map": reqMap})
  }
  ```

- pass local value

  ```go
  type User struct {
      Id int
      Username string
  }

  func main() {
      app := nf.New()
      app.Use(auth())

      app.Get("/item/list", list)
  }

  func auth() nf.HandlerFunc {
      return func(c *nf.Ctx) error {
          c.Locals("user", &User{Id: 1, Username:"user"})
          return c.Next()
      }
  }

  func list(c *nf.Ctx) error {
      user, ok := c.Locals("user").(*User)
      if !ok {
          return c.Status(401).SendString("login required")
      }

      ...
  }
  ```
