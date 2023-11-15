package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

type Phone struct {
	title string
	price int
	city  string
	link  string
	ram   string
}

type PhoneCatalog struct {
	oneplus []Phone
	sony    []Phone
	nothing []Phone
	xiaomi  []Phone
	huawei  []Phone
	pixel   []Phone
	honor   []Phone
	other   []Phone
	samsung []Phone
	apple   []Phone
}

func main() {
	c_jofog := colly.NewCollector()

	c_hardapro := colly.NewCollector()

	const MAX_PRICE = 200000
	const MIN_PRICE = 70000

	const DEPTH = 50
	currentDepth := 0

	var foundPhones PhoneCatalog

	reSamsung := regexp.MustCompile("SAMSUNG|GALAXY")
	reApple := regexp.MustCompile("APPLE")
	reHuawei := regexp.MustCompile("HUAWEI")
	reXiaomi := regexp.MustCompile("XIAOMI|XAOMI")
	reSony := regexp.MustCompile("SONY")
	reNothing := regexp.MustCompile("NOTHING")
	reOnePlus := regexp.MustCompile("ONEPLUS|ONE PLUS")
	rePixel := regexp.MustCompile("PIXEL")
	reHonor := regexp.MustCompile("HONOR")

	/*
		----------------------------------------------------------
		----------------------------------------------------------
		||                       JÓFOGÁS                        ||
		----------------------------------------------------------
		----------------------------------------------------------
	*/
	c_jofog.OnHTML(".contentArea", func(e *colly.HTMLElement) {
		currentTitle := e.ChildText("section.subjectWrapper h3.item-title a.subject")
		currentLink := e.ChildAttr("section.subjectWrapper h3.item-title a.subject", "href")
		currentPrice := strings.Replace(e.ChildText("section.price div.priceBox h3.item-price span.price-value"), " ", "", -1)
		currentCity := strings.Replace(e.ChildText("section.cityname"), "  ,", ",", 1)
		currentRam := getRamClassification(currentTitle)

		//currentRam = strconv.FormatBool(re512GB.MatchString(currentTitle)) + strconv.FormatBool(re256GB.MatchString(currentTitle)) + strconv.FormatBool(re128GB.MatchString(currentTitle)) + strconv.FormatBool(re512GBprob.MatchString(currentTitle)) + strconv.FormatBool(re256GBprob.MatchString(currentTitle)) + strconv.FormatBool(re128GBprob.MatchString(currentTitle))
		//"a[href].subject"

		// fmt.Printf("%v : %s \n", reSamsung.MatchString(currentTitle), currentTitle)
		currentPriceInt, err := strconv.Atoi(currentPrice)
		if err != nil {
			currentPriceInt = -1
		}
		if reSamsung.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.samsung = append(foundPhones.samsung,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reApple.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.apple = append(foundPhones.apple,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reOnePlus.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.oneplus = append(foundPhones.oneplus,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reSony.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.sony = append(foundPhones.sony,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reNothing.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.nothing = append(foundPhones.nothing,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reHuawei.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.huawei = append(foundPhones.huawei,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reXiaomi.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.xiaomi = append(foundPhones.xiaomi,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if rePixel.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.pixel = append(foundPhones.pixel,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reHonor.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.honor = append(foundPhones.honor,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else {
			foundPhones.other = append(foundPhones.other,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		}
	})

	c_jofog.OnHTML(".ad-list-pager-item-next", func(n *colly.HTMLElement) {
		nextLink := n.Attr("href")
		phonesCount := 0
		phonesCount = phonesCount + len(foundPhones.oneplus)
		phonesCount = phonesCount + len(foundPhones.sony)
		phonesCount = phonesCount + len(foundPhones.honor)
		phonesCount = phonesCount + len(foundPhones.huawei)
		phonesCount = phonesCount + len(foundPhones.xiaomi)
		phonesCount = phonesCount + len(foundPhones.samsung)
		phonesCount = phonesCount + len(foundPhones.apple)
		phonesCount = phonesCount + len(foundPhones.pixel)
		phonesCount = phonesCount + len(foundPhones.nothing)
		phonesCount = phonesCount + len(foundPhones.other)
		fmt.Printf("Jófogás Visiting:  %s\n Length Phones: %d\n", nextLink, phonesCount)
		currentDepth++
		if currentDepth >= DEPTH || nextLink == "" {
			//letsExcelize(foundPhones)
			currentDepth = 0
			return
		}
		c_jofog.Visit(nextLink)
	})

	/*
		----------------------------------------------------------
		----------------------------------------------------------
		||                     HARDVERAPRÓ                      ||
		----------------------------------------------------------
		----------------------------------------------------------
	*/
	c_hardapro.OnHTML(".media-body", func(e *colly.HTMLElement) {
		currentTitle := e.ChildText("div.uad-title h1 a")
		currentLink := e.ChildAttr("div.uad-title h1 a", "href")
		currentPrice := strings.Replace(strings.Replace(e.ChildText("div.uad-info div.uad-price"), " ", "", -1), "Ft", "", -1)
		currentCity := e.ChildText("div.uad-info div.light")
		currentRam := getRamClassification(currentTitle)

		//currentRam = strconv.FormatBool(re512GB.MatchString(currentTitle)) + strconv.FormatBool(re256GB.MatchString(currentTitle)) + strconv.FormatBool(re128GB.MatchString(currentTitle)) + strconv.FormatBool(re512GBprob.MatchString(currentTitle)) + strconv.FormatBool(re256GBprob.MatchString(currentTitle)) + strconv.FormatBool(re128GBprob.MatchString(currentTitle))
		//"a[href].subject"

		//fmt.Println(currentTitle)
		// fmt.Printf("%v : %s \n", reSamsung.MatchString(currentTitle), currentTitle)
		currentPriceInt, err := strconv.Atoi(currentPrice)
		if err != nil {
			currentPriceInt = -1
		}
		if reSamsung.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.samsung = append(foundPhones.samsung,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reApple.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.apple = append(foundPhones.apple,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reOnePlus.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.oneplus = append(foundPhones.oneplus,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reSony.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.sony = append(foundPhones.sony,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reNothing.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.nothing = append(foundPhones.nothing,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reHuawei.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.huawei = append(foundPhones.huawei,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reXiaomi.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.xiaomi = append(foundPhones.xiaomi,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if rePixel.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.pixel = append(foundPhones.pixel,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else if reHonor.MatchString(strings.ToUpper(currentTitle)) {
			foundPhones.honor = append(foundPhones.honor,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		} else {
			foundPhones.other = append(foundPhones.other,
				Phone{
					title: currentTitle,
					price: currentPriceInt,
					city:  currentCity,
					link:  currentLink,
					ram:   currentRam,
				})

		}
	})

	c_hardapro.OnHTML("#forum-nav-top ~ ul.mr-md-auto > li.nav-arrow > a", func(n *colly.HTMLElement) {
		if n.Attr("rel") == "next" {
			nextLink := "https://hardverapro.hu" + n.Attr("href")
			phonesCount := 0
			phonesCount = phonesCount + len(foundPhones.oneplus)
			phonesCount = phonesCount + len(foundPhones.sony)
			phonesCount = phonesCount + len(foundPhones.honor)
			phonesCount = phonesCount + len(foundPhones.huawei)
			phonesCount = phonesCount + len(foundPhones.xiaomi)
			phonesCount = phonesCount + len(foundPhones.samsung)
			phonesCount = phonesCount + len(foundPhones.apple)
			phonesCount = phonesCount + len(foundPhones.pixel)
			phonesCount = phonesCount + len(foundPhones.nothing)
			phonesCount = phonesCount + len(foundPhones.other)
			fmt.Printf("Hardverapró Visiting:  %s\n Length Phones: %d\n", nextLink, phonesCount)
			currentDepth++
			if currentDepth >= DEPTH || nextLink == "" {
				//letsExcelize(foundPhones)
				return
			}
			c_hardapro.Visit(nextLink)
		}
	})

	c_jofog.Visit(fmt.Sprintf("https://www.jofogas.hu/magyarorszag/mobiltelefon?max_price=%d&min_price=%d&mobile_memory=3,4,5,6,7,8&mobile_os=1&sp=2", MAX_PRICE, MIN_PRICE))
	c_jofog.Wait()
	c_hardapro.Visit(fmt.Sprintf("https://hardverapro.hu/aprok/mobil/mobil/android/keres.php?stext=&stcid_text=&stcid=&stmid_text=&stmid=&minprice=%d&maxprice=%d&cmpid_text=&cmpid=&usrid_text=&usrid=&__buying=0&__buying=1&stext_none=", MIN_PRICE, MAX_PRICE))
	c_hardapro.Wait()

	fmt.Println("Done scraping\nStarting Excel")

	sortPhones(foundPhones)
	letsExcelize(foundPhones)

	fmt.Println("Done with Excel")
}

func letsExcelize(phones PhoneCatalog) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Phones")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	/*

		// OnePlus \\

	*/

	f.SetCellValue("Phones", "A1", "ONEPLUS")
	f.MergeCell("Phones", "A1", "D1")
	styleColorofCol(f, "Phones", "A:D", "C6EFCE")
	styleTitle(f, "Phones", 1, 1, 4, 1, 1)
	f.SetCellValue("Phones", "A2", "Cím")
	styleTitle(f, "Phones", 1, 2, 1, 2, 1)
	f.SetCellValue("Phones", "B2", "Ár")
	styleTitle(f, "Phones", 2, 2, 2, 2, 1)
	f.SetCellValue("Phones", "C2", "Város")
	styleTitle(f, "Phones", 3, 2, 3, 2, 1)
	f.SetCellValue("Phones", "D2", "Link")
	styleTitle(f, "Phones", 4, 2, 4, 2, 1)
	for i, v := range phones.oneplus {
		cell, err := excelize.CoordinatesToCellName(1, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(2, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(3, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(4, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Nothing \\

	*/

	f.SetCellValue("Phones", "F1", "NOTHING")
	f.MergeCell("Phones", "F1", "I1")
	styleTitle(f, "Phones", 6, 1, 9, 1, 0)
	f.SetCellValue("Phones", "F2", "Cím")
	styleTitle(f, "Phones", 6, 2, 6, 2, 0)
	f.SetCellValue("Phones", "G2", "Ár")
	styleTitle(f, "Phones", 7, 2, 7, 2, 0)
	f.SetCellValue("Phones", "H2", "Város")
	styleTitle(f, "Phones", 8, 2, 8, 2, 0)
	f.SetCellValue("Phones", "I2", "Link")
	styleTitle(f, "Phones", 9, 2, 9, 2, 0)
	for i, v := range phones.nothing {
		cell, err := excelize.CoordinatesToCellName(6, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(7, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(8, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(9, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Sony \\

	*/

	f.SetCellValue("Phones", "K1", "SONY")
	f.MergeCell("Phones", "K1", "N1")
	styleTitle(f, "Phones", 11, 1, 14, 1, 0)
	f.SetCellValue("Phones", "K2", "Cím")
	styleTitle(f, "Phones", 11, 2, 11, 2, 0)
	f.SetCellValue("Phones", "L2", "Ár")
	styleTitle(f, "Phones", 12, 2, 12, 2, 0)
	f.SetCellValue("Phones", "M2", "Város")
	styleTitle(f, "Phones", 13, 2, 13, 2, 0)
	f.SetCellValue("Phones", "N2", "Link")
	styleTitle(f, "Phones", 14, 2, 14, 2, 0)
	for i, v := range phones.sony {
		cell, err := excelize.CoordinatesToCellName(11, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(12, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(13, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(14, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Honor \\

	*/

	f.SetCellValue("Phones", "P1", "HONOR")
	f.MergeCell("Phones", "P1", "S1")
	styleTitle(f, "Phones", 16, 1, 19, 1, 0)
	f.SetCellValue("Phones", "P2", "Cím")
	styleTitle(f, "Phones", 16, 2, 16, 2, 0)
	f.SetCellValue("Phones", "Q2", "Ár")
	styleTitle(f, "Phones", 17, 2, 17, 2, 0)
	f.SetCellValue("Phones", "R2", "Város")
	styleTitle(f, "Phones", 18, 2, 18, 2, 0)
	f.SetCellValue("Phones", "S2", "Link")
	styleTitle(f, "Phones", 19, 2, 19, 2, 0)
	for i, v := range phones.honor {
		cell, err := excelize.CoordinatesToCellName(16, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(17, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(18, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(19, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Pixel \\

	*/

	f.SetCellValue("Phones", "U1", "PIXEL")
	f.MergeCell("Phones", "U1", "X1")
	styleTitle(f, "Phones", 21, 1, 24, 1, 0)
	f.SetCellValue("Phones", "U2", "Cím")
	styleTitle(f, "Phones", 21, 2, 21, 2, 0)
	f.SetCellValue("Phones", "V2", "Ár")
	styleTitle(f, "Phones", 22, 2, 22, 2, 0)
	f.SetCellValue("Phones", "W2", "Város")
	styleTitle(f, "Phones", 23, 2, 23, 2, 0)
	f.SetCellValue("Phones", "X2", "Link")
	styleTitle(f, "Phones", 24, 2, 24, 2, 0)
	for i, v := range phones.pixel {
		cell, err := excelize.CoordinatesToCellName(21, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(22, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(23, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(24, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Xiaomi \\

	*/

	f.SetCellValue("Phones", "Z1", "XIAOMI")
	f.MergeCell("Phones", "Z1", "AC1")
	styleTitle(f, "Phones", 26, 1, 29, 1, 0)
	f.SetCellValue("Phones", "Z2", "Cím")
	styleTitle(f, "Phones", 26, 2, 26, 2, 0)
	f.SetCellValue("Phones", "AA2", "Ár")
	styleTitle(f, "Phones", 27, 2, 27, 2, 0)
	f.SetCellValue("Phones", "AB2", "Város")
	styleTitle(f, "Phones", 28, 2, 28, 2, 0)
	f.SetCellValue("Phones", "AC2", "Link")
	styleTitle(f, "Phones", 29, 2, 29, 2, 0)
	for i, v := range phones.xiaomi {
		cell, err := excelize.CoordinatesToCellName(26, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(27, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(28, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(29, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Huawei \\

	*/

	f.SetCellValue("Phones", "AE1", "HUAWEI")
	f.MergeCell("Phones", "AE1", "AH1")
	styleTitle(f, "Phones", 31, 1, 34, 1, 0)
	f.SetCellValue("Phones", "AE2", "Cím")
	styleTitle(f, "Phones", 31, 2, 31, 2, 0)
	f.SetCellValue("Phones", "AF2", "Ár")
	styleTitle(f, "Phones", 32, 2, 32, 2, 0)
	f.SetCellValue("Phones", "AG2", "Város")
	styleTitle(f, "Phones", 33, 2, 33, 2, 0)
	f.SetCellValue("Phones", "AH2", "Link")
	styleTitle(f, "Phones", 34, 2, 34, 2, 0)
	for i, v := range phones.huawei {
		cell, err := excelize.CoordinatesToCellName(31, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(32, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(33, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(34, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Other \\

	*/

	f.SetCellValue("Phones", "AJ1", "OTHER")
	f.MergeCell("Phones", "AJ1", "AM1")
	styleTitle(f, "Phones", 36, 1, 39, 1, 0)
	f.SetCellValue("Phones", "AJ2", "Cím")
	styleTitle(f, "Phones", 36, 2, 36, 2, 0)
	f.SetCellValue("Phones", "AK2", "Ár")
	styleTitle(f, "Phones", 37, 2, 37, 2, 0)
	f.SetCellValue("Phones", "AL2", "Város")
	styleTitle(f, "Phones", 38, 2, 38, 2, 0)
	f.SetCellValue("Phones", "AM2", "Link")
	styleTitle(f, "Phones", 39, 2, 39, 2, 0)
	for i, v := range phones.other {
		cell, err := excelize.CoordinatesToCellName(36, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(37, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(38, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(39, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Samsung \\

	*/

	f.SetCellValue("Phones", "AO1", "SAMSUNG")
	f.MergeCell("Phones", "AO1", "AR1")
	styleColorofCol(f, "Phones", "AO:AR", "FFC7CE")
	styleTitle(f, "Phones", 41, 1, 44, 1, -1)
	f.SetCellValue("Phones", "AO2", "Cím")
	styleTitle(f, "Phones", 41, 2, 41, 2, -1)
	f.SetCellValue("Phones", "AP2", "Ár")
	styleTitle(f, "Phones", 42, 2, 42, 2, -1)
	f.SetCellValue("Phones", "AQ2", "Város")
	styleTitle(f, "Phones", 43, 2, 43, 2, -1)
	f.SetCellValue("Phones", "AR2", "Link")
	styleTitle(f, "Phones", 44, 2, 44, 2, -1)

	for i, v := range phones.samsung {
		cell, err := excelize.CoordinatesToCellName(41, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(42, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(43, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(44, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	/*

		// Apple \\

	*/

	f.SetCellValue("Phones", "AT1", "APPLE")
	f.MergeCell("Phones", "AT1", "AW1")
	styleColorofCol(f, "Phones", "AT:AW", "FFC7CE")
	styleTitle(f, "Phones", 46, 1, 49, 1, -1)
	f.SetCellValue("Phones", "AT2", "Cím")
	styleTitle(f, "Phones", 46, 2, 46, 2, -1)
	f.SetCellValue("Phones", "AU2", "Ár")
	styleTitle(f, "Phones", 47, 2, 47, 2, -1)
	f.SetCellValue("Phones", "AV2", "Város")
	styleTitle(f, "Phones", 48, 2, 48, 2, -1)
	f.SetCellValue("Phones", "AW2", "Link")
	styleTitle(f, "Phones", 49, 2, 49, 2, -1)

	for i, v := range phones.apple {
		cell, err := excelize.CoordinatesToCellName(46, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.title)

		cell, err = excelize.CoordinatesToCellName(47, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.price)

		cell, err = excelize.CoordinatesToCellName(48, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.city)

		cell, err = excelize.CoordinatesToCellName(49, i+3)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", cell, v.link)
	}

	autoFitCols(f, "Phones")

	// Save spreadsheet by the given path.
	if err := f.SaveAs("scrapings/phones.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func autoFitCols(f *excelize.File, sheetName string) {
	// Autofit all columns according to their text content
	cols, err := f.GetCols(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetColWidth(sheetName, name, name, float64(largestWidth))
	}
}

func styleTitle(f *excelize.File, sheetName string, cellPositionx int, cellPositiony int, cellPositionx2 int, cellPositiony2 int, colorgoodorbad int) {
	fontsize := 11.0
	if cellPositiony+cellPositiony2 == 2 {
		fontsize = 14.0
	}

	cell, err := excelize.CoordinatesToCellName(cellPositionx, cellPositiony)
	if err != nil {
		fmt.Println(err)
		return
	}
	cell2, err := excelize.CoordinatesToCellName(cellPositionx2, cellPositiony2)
	if err != nil {
		fmt.Println(err)
		return
	}

	var style int
	var styleerr error
	if colorgoodorbad == 1 { // GOOD
		style, styleerr = f.NewStyle(&excelize.Style{
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"C6EFCE"}, Pattern: 1},
			Font:      &excelize.Font{Bold: true, Size: fontsize},
			Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
			Border: []excelize.Border{
				{Type: "top", Style: 2, Color: "000000"},
				{Type: "bottom", Style: 2, Color: "000000"},
				{Type: "left", Style: 2, Color: "000000"},
				{Type: "right", Style: 2, Color: "000000"},
			},
		})
	} else if colorgoodorbad == -1 { // BAD
		style, styleerr = f.NewStyle(&excelize.Style{
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFC7CE"}, Pattern: 1},
			Font:      &excelize.Font{Bold: true, Size: fontsize},
			Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
			Border: []excelize.Border{
				{Type: "top", Style: 2, Color: "000000"},
				{Type: "bottom", Style: 2, Color: "000000"},
				{Type: "left", Style: 2, Color: "000000"},
				{Type: "right", Style: 2, Color: "000000"},
			},
		})
	} else { // NEUTRAL
		style, styleerr = f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: fontsize},
			Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
			Border: []excelize.Border{
				{Type: "top", Style: 2, Color: "000000"},
				{Type: "bottom", Style: 2, Color: "000000"},
				{Type: "left", Style: 2, Color: "000000"},
				{Type: "right", Style: 2, Color: "000000"},
			},
		})
	}
	if styleerr != nil {
		fmt.Println(styleerr)
		return
	}
	f.SetCellStyle(sheetName, cell, cell2, style)
}

func styleColorofCol(f *excelize.File, sheetName string, cols string, color string) {
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetColStyle(sheetName, cols, style)
}

func sortPhones(phones PhoneCatalog) {
	sort.Slice(phones.oneplus, func(i, j int) bool {
		return phones.oneplus[i].price > phones.oneplus[j].price //&& phones.oneplus[i].ram > phones.oneplus[j].ram
	})
	sort.Slice(phones.sony, func(i, j int) bool {
		return phones.sony[i].price > phones.sony[j].price //&& phones.sony[i].ram > phones.sony[j].ram
	})
	sort.Slice(phones.nothing, func(i, j int) bool {
		return phones.nothing[i].price > phones.nothing[j].price //&& phones.nothing[i].ram > phones.nothing[j].ram
	})
	sort.Slice(phones.xiaomi, func(i, j int) bool {
		return phones.xiaomi[i].price > phones.xiaomi[j].price //&& phones.xiaomi[i].ram > phones.xiaomi[j].ram
	})
	sort.Slice(phones.huawei, func(i, j int) bool {
		return phones.huawei[i].price > phones.huawei[j].price //&& phones.huawei[i].ram > phones.huawei[j].ram
	})
	sort.Slice(phones.pixel, func(i, j int) bool {
		return phones.pixel[i].price > phones.pixel[j].price //&& phones.pixel[i].ram > phones.pixel[j].ram
	})
	sort.Slice(phones.honor, func(i, j int) bool {
		return phones.honor[i].price > phones.honor[j].price //&& phones.honor[i].ram > phones.honor[j].ram
	})
	sort.Slice(phones.other, func(i, j int) bool {
		return phones.other[i].price > phones.other[j].price //&& phones.other[i].ram > phones.other[j].ram
	})
	sort.Slice(phones.samsung, func(i, j int) bool {
		return phones.samsung[i].price > phones.samsung[j].price //&& phones.samsung[i].ram > phones.samsung[j].ram
	})
	sort.Slice(phones.apple, func(i, j int) bool {
		return phones.apple[i].price > phones.apple[j].price //&& phones.apple[i].ram > phones.apple[j].ram
	})

}

func getRamClassification(title string) (ram string) {
	ram = ""

	re512GB := regexp.MustCompile("512GB|512 GB|512")
	re512GBprob := regexp.MustCompile("512")
	re256GB := regexp.MustCompile("256GB|256 GB|256")
	re256GBprob := regexp.MustCompile("256")
	re128GB := regexp.MustCompile("128GB|128 GB|128")
	re128GBprob := regexp.MustCompile("128")

	if re512GB.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	if re256GB.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	if re128GB.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	if re512GBprob.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	if re256GBprob.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	if re128GBprob.MatchString(title) {
		ram += "1"
	} else {
		ram += "0"
	}
	return
}
