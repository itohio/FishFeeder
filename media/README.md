This folder contains icons used to display various states

In order to generate binary representations of the icons you can run `convert2bin` utility in the project root like so:

```
go run ./media media/aquarium.png >icons/aquarium.go
```

Alternatively, from inside `icons` folder you can just run `go generate`

Note: if you are using environment configured for TinyGo, then you should run it in a separate console.

# Icons:
- Fish Food by Icons8: https://icons8.com/icon/jQmKQCKAiKna/fish-food
- Aquarium by Icons8: https://icons8.com/icon/1Wi51AkN6dYv/aquarium
- Water Filter by Icons8: https://icons8.com/icon/dGlyR3EdBrXx/water-filter
- Thermometer by Icons8: https://icons8.com/icon/poFZHQZ-CjsC/thermometer
