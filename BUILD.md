# How to build from Windows

### Install dependencies

##### C compiler

Download and install https://www.msys2.org/

Open `MSYS2 MinGW 64-bit` from the start menu and run:

```bash
pacman -Syu
pacman -S git mingw-w64-x86_64-toolchain
echo "export PATH=\$PATH:/d/Apps/Dev/Go/go1.19.5/bin:~/Go/bin" >> ~/.bashrc
```

##### Fyne Module

```bash
go get fyne.io/fyne/v2
```

##### Text to Speech

```bash
go get github.com/hajimehoshi/go-mp3
go get github.com/hajimehoshi/oto/v2
```

##### Cleanup

```bash
go mod tidy
```

### Packaging

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

```bash
fyne package -os windows -icon res\SyncTimer.icon.png
```
