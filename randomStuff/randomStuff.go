package randomStuff

import (
	"math/rand"
)

var randomText = []string{
	"The Ultimate Wavetable Tool!",
	"Krul is a dancer!",
	"くるみ",
	"A tale about a land of lotuses",
	"Turning Konami SCC into YM2151",
	"Creeper? Aww man!",
	"When the impostor is sus",
	"Also try Furnace!",
	"YMF163",
	"https://en.touhouwiki.net/wiki/Kurumi",
	"Seraph of the Beginning",
	"Amogus",
	"E",
	"If you encounter a bug, Flandre Scarlet will destroy it.",
	"SEGA PC-Engine",
	"Sussy",
	"Turning furries into vampires...",
	"Will be probably rewritten in Rust...",
	"Here We Go(lang)!",
	"This program never segfaults",
	"/\\/\\/\\/\\/\\/\\/\\/\\",
	"16, 25, 30, 31, 30, 29, 26, 25, 25, 28, 31, 28, 18, 11, 10, 13, 17, 20, 22, 20, 15, 6, 0, 2, 6, 5, 3, 1, 0, 0, 1, 4",
	"LEGO Wavetable!",
	"Dance dance dance with my hands hands hands...",
	"If you're not careful and you noclip out of reality in the wrong areas, you'll end up in the Backrooms.",
	"Want to summon a demon? Program in php!",
	"Watch your neck!",
	"Gaussian interpolation = BEST INTERPOLATION",
	"WHAT THE F-",
	"Never gonna give you up",
	"Eating a Bad Apple!!",
}

func GetTitle() string {
	num := rand.Intn(len(randomText) - 1)
	return randomText[num]
}
