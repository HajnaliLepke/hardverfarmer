package main

import (
	"fmt"
	"regexp"
	"slices"
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

type PhoneBrand struct {
	name        string
	phones      []Phone
	excel_color string
}

type PhoneRegExp struct {
	name     string
	p_regexp *regexp.Regexp
}

type PhoneRegExpCatalog struct {
	regexps []PhoneRegExp
}

// type PhoneCatalog struct {
// 	flip    []Phone
// 	oneplus []Phone
// 	sony    []Phone
// 	nothing []Phone
// 	xiaomi  []Phone
// 	huawei  []Phone
// 	pixel   []Phone
// 	honor   []Phone
// 	other   []Phone
// 	samsung []Phone
// 	apple   []Phone
// }

type PhoneCatalog struct {
	phonebrands []PhoneBrand
}

func main() {

	c_jofog := colly.NewCollector()

	c_hardapro := colly.NewCollector()

	const MAX_PRICE = 200000
	const MIN_PRICE = 50000

	const DEPTH = 50
	currentDepth := 0

	var foundPhones PhoneCatalog

	// INIT PHONEBRANDS
	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "flip",
			phones:      make([]Phone, 0),
			excel_color: "GOOD",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "fold",
			phones:      make([]Phone, 0),
			excel_color: "GOOD",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "samsung",
			phones:      make([]Phone, 0),
			excel_color: "GOOD",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "sony",
			phones:      make([]Phone, 0),
			excel_color: "GOOD",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "oneplus",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "pixel",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "honor",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "nothing",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "huawei",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "xiaomi",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "other",
			phones:      make([]Phone, 0),
			excel_color: "NEUTRAL",
		})

	foundPhones.phonebrands = append(foundPhones.phonebrands,
		PhoneBrand{
			name:        "apple",
			phones:      make([]Phone, 0),
			excel_color: "BAD",
		})

	var p_regexps PhoneRegExpCatalog

	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "flip",
		p_regexp: regexp.MustCompile("FLIP"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "fold",
		p_regexp: regexp.MustCompile("FOLD"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "samsung",
		p_regexp: regexp.MustCompile("SAMSUNG|GALAXY"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "apple",
		p_regexp: regexp.MustCompile("APPLE|IPHONE"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "huawei",
		p_regexp: regexp.MustCompile("HUAWEI|HAUWEI"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "xiaomi",
		p_regexp: regexp.MustCompile("XIAOMI|XAOMI"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "sony",
		p_regexp: regexp.MustCompile("SONY"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "nothing",
		p_regexp: regexp.MustCompile("NOTHING"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "oneplus",
		p_regexp: regexp.MustCompile("ONEPLUS|ONE PLUS"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "pixel",
		p_regexp: regexp.MustCompile("PIXEL"),
	})
	p_regexps.regexps = append(p_regexps.regexps, PhoneRegExp{
		name:     "honor",
		p_regexp: regexp.MustCompile("HONOR"),
	})

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

		var foundbrand = false

		for _, v := range p_regexps.regexps {
			if v.p_regexp.MatchString(strings.ToUpper(currentTitle)) {
				phone := slices.IndexFunc(foundPhones.phonebrands, func(p PhoneBrand) bool { return strings.ToUpper(p.name) == strings.ToUpper(v.name) })
				foundPhones.phonebrands[phone].phones = append(foundPhones.phonebrands[phone].phones,
					Phone{
						title: currentTitle,
						price: currentPriceInt,
						city:  currentCity,
						link:  currentLink,
						ram:   currentRam,
					})
				foundbrand = true
				break
			}
		}
		if foundbrand == false {
			phone := slices.IndexFunc(foundPhones.phonebrands, func(p PhoneBrand) bool { return strings.ToUpper(p.name) == strings.ToUpper("other") })
			foundPhones.phonebrands[phone].phones = append(foundPhones.phonebrands[phone].phones,
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
		for _, v := range foundPhones.phonebrands {
			phonesCount += len(v.phones)
		}
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
		currentCity := e.ChildText("div.uad-info div.uad-light")
		currentRam := getRamClassification(currentTitle)

		//currentRam = strconv.FormatBool(re512GB.MatchString(currentTitle)) + strconv.FormatBool(re256GB.MatchString(currentTitle)) + strconv.FormatBool(re128GB.MatchString(currentTitle)) + strconv.FormatBool(re512GBprob.MatchString(currentTitle)) + strconv.FormatBool(re256GBprob.MatchString(currentTitle)) + strconv.FormatBool(re128GBprob.MatchString(currentTitle))
		//"a[href].subject"

		//fmt.Println(currentTitle)
		// fmt.Printf("%v : %s \n", reSamsung.MatchString(currentTitle), currentTitle)
		currentPriceInt, err := strconv.Atoi(currentPrice)
		if err != nil {
			currentPriceInt = -1
		}

		var foundbrand = false

		for _, v := range p_regexps.regexps {
			if v.p_regexp.MatchString(strings.ToUpper(currentTitle)) {
				phone := slices.IndexFunc(foundPhones.phonebrands, func(p PhoneBrand) bool { return strings.ToUpper(p.name) == strings.ToUpper(v.name) })
				foundPhones.phonebrands[phone].phones = append(foundPhones.phonebrands[phone].phones,
					Phone{
						title: currentTitle,
						price: currentPriceInt,
						city:  currentCity,
						link:  currentLink,
						ram:   currentRam,
					})
				foundbrand = true
				break
			}
		}
		if foundbrand == false {
			phone := slices.IndexFunc(foundPhones.phonebrands, func(p PhoneBrand) bool { return strings.ToUpper(p.name) == strings.ToUpper("other") })
			foundPhones.phonebrands[phone].phones = append(foundPhones.phonebrands[phone].phones,
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
			for _, v := range foundPhones.phonebrands {
				phonesCount += len(v.phones)
			}
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

	fmt.Println("Done scraping\nStarting Sorting\n")

	sortPhones(foundPhones)
	fmt.Println("Done sorting\nStarting Excel\n")
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
		// LOOPING \\
	*/

	for i, p := range phones.phonebrands {
		//fmt.Println("Starting: " + p.name + "\n")
		title_cell1, err := excelize.CoordinatesToCellName(1+(i*5), 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		title_cell4, err := excelize.CoordinatesToCellName(4+(i*5), 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		subtitle_cell1, err := excelize.CoordinatesToCellName(1+(i*5), 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		subtitle_cell2, err := excelize.CoordinatesToCellName(2+(i*5), 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		subtitle_cell3, err := excelize.CoordinatesToCellName(3+(i*5), 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		subtitle_cell4, err := excelize.CoordinatesToCellName(4+(i*5), 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue("Phones", title_cell1, strings.ToUpper(p.name))

		f.MergeCell("Phones", title_cell1, title_cell4)
		styleColorofCol(f, "Phones", title_cell1[:len(title_cell1)-1]+":"+title_cell4[:len(title_cell4)-1], "C6EFCE")
		styleTitle(f, "Phones", 1+(i*5), 1, 4+(i*5), 1, p.excel_color)
		f.SetCellValue("Phones", subtitle_cell1, "Cím")
		styleTitle(f, "Phones", 1+(i*5), 2, 1+(i*5), 2, p.excel_color)
		f.SetCellValue("Phones", subtitle_cell2, "Ár")
		styleTitle(f, "Phones", 2+(i*5), 2, 2+(i*5), 2, p.excel_color)
		f.SetCellValue("Phones", subtitle_cell3, "Város")
		styleTitle(f, "Phones", 3+(i*5), 2, 3+(i*5), 2, p.excel_color)
		f.SetCellValue("Phones", subtitle_cell4, "Link")
		styleTitle(f, "Phones", 4+(i*5), 2, 4+(i*5), 2, p.excel_color)
		for j, v := range p.phones {
			//fmt.Println("Writing: " + v.title + "\n")
			cell, err := excelize.CoordinatesToCellName(1+(i*5), j+3)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.SetCellValue("Phones", cell, v.title)

			cell, err = excelize.CoordinatesToCellName(2+(i*5), j+3)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.SetCellValue("Phones", cell, v.price)

			cell, err = excelize.CoordinatesToCellName(3+(i*5), j+3)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.SetCellValue("Phones", cell, v.city)

			cell, err = excelize.CoordinatesToCellName(4+(i*5), j+3)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.SetCellValue("Phones", cell, v.link)
		}
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

func styleTitle(f *excelize.File, sheetName string, cellPositionx int, cellPositiony int, cellPositionx2 int, cellPositiony2 int, colorgoodorbad string) {
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
	if colorgoodorbad == "GOOD" { // GOOD
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
	} else if colorgoodorbad == "BAD" { // BAD
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

	for _, v := range phones.phonebrands {
		sort.Slice(v.phones, func(i, j int) bool {
			return v.phones[i].price > v.phones[j].price //&& v.phones[i].ram > v.phones[j].ram
		})
	}

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

/*

#
#	STOLENED !!!!!!
#

*/
