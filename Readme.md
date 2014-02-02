# Good config

## Example
The example of usage is located in the http://github.com/PavelPolyakov/goodconfig/example/ folder.

## Example in short
**application.ini:**
```
[production]
hello = 1
hello_2 = 2
parent.child.hello = 3
parent.child.subchild = 3
```
**part of the *go file**
```Go
Config := goodconfig.NewConfig()
err := Config.Parse("./config/application.ini")

if err != nil {
    fmt.Println(err)
}

Config.SetDefaultSection("production")
fmt.Println(Config.Get("parent").Get("child").Get("subchild").ToInt())
```