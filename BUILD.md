# Linux Debian/Ubuntu
### JetBrains GoLand or IntelliJ IDEA Ultimate
_Source: https://wiki.debian.org/JetBrains_
```bash
curl -s https://s3.eu-central-1.amazonaws.com/jetbrains-ppa/0xA6E8698A.pub.asc | gpg --dearmor | sudo tee /usr/share/keyrings/jetbrains-ppa-archive-keyring.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/jetbrains-ppa-archive-keyring.gpg] http://jetbrains-ppa.s3-website.eu-central-1.amazonaws.com any main" | sudo tee /etc/apt/sources.list.d/jetbrains-ppa.list > /dev/null
sudo apt update
```
```bash
sudo apt install goland
```
```bash
sudo apt install intellij-idea-ultimate
```
### Required Packages
```bash
sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev
sudo apt-get install libasound2-dev
```
# Windows
##### C compiler

_Based on Fyne "Get Started" documentation: https://developer.fyne.io/started/_

Download and install https://www.msys2.org/

Start `MSYS2 MinGW 64-bit` from the start menu and run:

```bash
pacman -Syu
pacman -S git mingw-w64-x86_64-toolchain
GO_BIN_PATH=d/Projects/SDK/go1.21.4/bin
echo -e "export PATH=\$PATH:${GO_BIN_PATH}:~/Go/bin" >> ~/.bashrc
```

For the compiler to work you will need to set up the windows %PATH% variable to find these tools.
Add “C:\msys64\mingw64\bin” to the Path list.

# Dependancies
### Fyne Module
```bash
go get fyne.io/fyne/v2
```
### Text to Speech
```bash
go get github.com/hajimehoshi/go-mp3
go get github.com/hajimehoshi/oto/v2
```
### Get other requirements & cleanup
```bash
go get -u
go clean -modcache
go mod tidy
```
# Packaging
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```
```bash
fyne package -os windows -icon res\SyncTimer.icon.png
```