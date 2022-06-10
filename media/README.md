This folder contains icons used to display various states

In order to generate binary representations of the icons you can run `convert2bin` utility in the project root like so:

```
go run ./media media/aquarium.png >icons/aquarium.go
```

Alternatively, from inside `icons` folder you can just run `go generate`

Note: if you are using environment configured for TinyGo, then you should run it in a separate console.