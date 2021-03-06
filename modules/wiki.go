package helper


import (
    // "fmt"
	"strings"
    "github.com/gocolly/colly"
)


func FetchData(url string) string {
	c := colly.NewCollector(
        colly.AllowedDomains("en.wikipedia.org"),
    )

	var text string
    c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
        raw_text := e.ChildText("p")
		raw_text = strings.ReplaceAll(raw_text, "\n", "")
		text += raw_text
    })

    c.Visit(url)
	// fmt.Println(text)
	return text
	// return "Synthetic diamond (also called laboratory-grown, laboratory-created, man-made, artisan created, or cultured diamond) is diamond that is produced in a controlled technological process (in contrast to naturally formed diamond, which is created through geological processes and obtained by mining). Unlike diamond simulants (imitations of diamond made of superficially-similar non-diamond materials), synthetic diamonds are composed of the same material as naturally formed diamonds – pure carbon crystallized in an isotropic 3D form – and share identical chemical and physical properties. Numerous claims of diamond synthesis were reported between 1879 and 1928; most of these attempts were carefully analyzed but none were confirmed. In the 1940s, systematic research of diamond creation began in the United States, Sweden and the Soviet Union, which culminated in the first reproducible synthesis in 1953. Further research activity yielded the discoveries of HPHT diamond and CVD diamond, named for their production method (high-pressure high-temperature and chemical vapor deposition, respectively). These two processes still dominate synthetic diamond production. A third method in which nanometer-sized diamond grains are created in a detonation of carbon-containing explosives, known as detonation synthesis, entered the market in the late 1990s. A fourth method, treating graphite with high-power ultrasound, has been demonstrated in the laboratory, but currently has no commercial application. The properties of synthetic diamond depend on the manufacturing process. Some synthetic diamonds have properties such as hardness, thermal conductivity and electron mobility that are superior to those of most naturally formed diamonds. Synthetic diamond is widely used in abrasives, in cutting and polishing tools and in heat sinks. Electronic applications of synthetic diamond are being developed, including high-power switches at power stations, high-frequency field-effect transistors and light-emitting diodes. Synthetic diamond detectors of ultraviolet (UV) light or high-energy particles are used at high-energy research facilities and are available commercially. Due to its unique combination of thermal and chemical stability, low thermal expansion and high optical transparency in a wide spectral range, synthetic diamond is becoming the most popular material for optical windows in high-power CO2 lasers and gyrotrons. It is estimated that 98% of industrial-grade diamond demand is supplied with synthetic diamonds.[1] Both CVD and HPHT diamonds can be cut into gems and various colors can be produced: clear white, yellow, brown, blue, green and orange. The advent of synthetic gems on the market created major concerns in the diamond trading business, as a result of which special spectroscopic devices and techniques have been developed to distinguish synthetic and natural diamonds."
}
