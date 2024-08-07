package main

import "fmt"

/*
	- Kiválasztani a gyertya adatokat közlő forrást (1-2 évre visszamenőleg, napi gyertya adatok, 4 kategóriában: crypto, részvény, nyersanyag, fx)
	- Service, mely hívásra leszedi az utolsó időszak (paraméterezhető) adatait a kért típusra és az addícionálist hozzáadja az adatbázishoz. Goroutin használata
	- Adatbázis és hozzá repository, mely tartalmazza a napi gyertyaadatokat
	- Adatbázis és hozzá repository, mely az utolsó adatfrissítést idejét tartalmazza ill. típusonként az utolsó dátumot
	- Adatgyűjtő service-be olyan public fv, mely csak akkor indít adatfrissítést, ha aktuális nap még nem történt. Ez hívódna meg alkalmazás indításkor ill. egyéb bármilyen fv indításakor
	- gin-el controller készítése és induló index.html, mely https://github.com/go-echarts/go-echarts libbel kirajzolja a 4 db chart-ot
	- Dockerfile és fly.toml
	- Fő oldalon legyen lehetőség kategórián belül mást is kiválasztani
	- Ha egy chart-ra kattintunk, külön oldalon jelenjen csak ez meg és itt addícionális mutatók/görbék is legyenek (pl. mozgó átlag)
	- Bármilyen átlag/indikátor esetén unit tesztek
	- email-el legyen lehetőség belépni (don-task-es alapján). Ilyenkor 1 millió $ áll rendelkezésre, ezzel lehet:
		- limit áras adásvételt indítani kiválasztott típusra
		- piaci áras adásvételt indítani kiválasztott típusra
	- A megbízások mentődjenek el és új adatpont lementésekor hajtódjanak végre (ha lehet). Mindenhol unit tesztek
	- Minden user-nél időpecsétes könyveléssel látszódjon, hogy melyik típus hogy változott ($, kereskedési típusok) és hogy ezáltal aznap miből mennyi készlete volt. Mindenhol unit tesztek
	- Ügyfél vagyonának napi helyzetéről legyne chart, típus szerinti bontásban
	- Nyitott megbízásokról lista
	- Email fíók létrehozása
	- Végrehajtott megbízásról emailes értesítés legyen (ha kéri)
	- Minden nap emailes értesítés a portfólió aktuális állapotáról (ha kéri)
	- Legyen lehetőség short ügyletekre is
	- Legyen lehetőség tőkeáttételes ügyletekre is (long / short)
	- Automatikus/állandó eladási/vételi utasítások megadása generálisan ill. adott típusra (pl. X napos mozgó átlag magasabb, Y naposnál, vétel)
	- Értelmes indikátorok meghatározása / megismerése
	- Limit megadása automatikus műveleteknél (pl. short is, tőkeáttét is, őssz portfóli max x%-a)
*/

func main() {
	fmt.Println("Hello")
}
