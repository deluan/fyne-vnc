module github.com/deluan/fyne-vnc

go 1.13

require (
	fyne.io/fyne v1.4.2
	github.com/amitbet/vnc2video v0.0.0-20190616012314-9d50b9dab1d9
	github.com/boz/go-throttle v0.0.0-20160922054636-fdc4eab740c1
	github.com/deluan/bring v0.0.7
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/amitbet/vnc2video => github.com/deluan/vnc2video v0.0.0-20210101045232-81fb4f50aef5

//replace github.com/amitbet/vnc2video => ../vnc2video
