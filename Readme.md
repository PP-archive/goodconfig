# Good config
package goodconfig provides the basic *.ini files parsing options and supports the sections inheritance

## Example
The example of usage is located in the [/example](https://github.com/PavelPolyakov/goodconfig/tree/master/example) folder.

## Example in short
**application.ini:**
```
[production]
intValue = 33
boolValue = true
;comment here
stringValue = string me!
parent.child.subchild = hello
parent.child.subchild_2 = hello_2

[development : production]
boolValue = false
;comment here
stringValue = string me!
parent.child.subchild = bye
```
**part of the *.go file**
```Go
Config := goodconfig.NewConfig()
err := Config.Parse("./config/application.ini")

if err != nil {
    fmt.Println(err)
}

Config.SetDefaultSection("production")
fmt.Println(Config.Get("parent").Get("child").Get("subchild").ToString())

// or
fmt.Println(Config.Section("development").Get("parent").Get("child").Get("subchild").ToString())
```