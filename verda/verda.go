package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Car struct {
	name     string
	price    int
	location string
	link     string
	km       string
	age      string
}

type CarCatalog struct {
	nissanLeaf []Car
}

func main() {
	c_parking := colly.NewCollector()

	const DEPTH = 50
	currentDepth := 0

	var foundCars CarCatalog

	c_parking.OnHTML(".right-result-bloc", func(e *colly.HTMLElement) {
		// currentName := fmt.Sprintf("%s %s %s", e.ChildText(".tag_f_titre:nth-child(1)"), e.ChildText(".tag_f_titre:nth-child(2)"), e.ChildText(".tag_f_titre:nth-child(3)"))
		currentName := fmt.Sprintf("%s %s %s", e.ChildText(".title-block.brand"), strings.Replace(e.ChildText(".title-block.sub-title"), e.ChildText(".title-block.sub-title .nowrap"), "", -1), e.ChildText(".title-block.sub-title .nowrap"))
		currentPrice := strings.Replace(strings.Replace(e.ChildText(".price-block .prix"), " €", "", -1), ",", "", -1)
		currentLocation := e.ChildText(".location .upper")
		currentLink := "https://www.theparking.eu/" + e.ChildAttr(".external", "href")

		currentPriceInt, err := strconv.Atoi(currentPrice)
		if err != nil {
			currentPriceInt = -1
		}

		foundCars.nissanLeaf = append(foundCars.nissanLeaf,
			Car{
				name:     currentName,
				price:    currentPriceInt,
				location: currentLocation,
				link:     currentLink,
				age:      "1997",
				km:       "1000 km",
			})

		fmt.Printf("%s: %v€, %s \n", currentName, currentPriceInt, currentLocation)
	})

	c_parking.OnHTML(".pagin-top li.btn-next a", func(n *colly.HTMLElement) {
		nextLink := n.Attr("href")

		carsCount := len(foundCars.nissanLeaf)

		fmt.Printf("Parking Visiting:  %s\n Length Phones: %d\n", nextLink, carsCount)
		currentDepth++
		if currentDepth >= DEPTH || nextLink == "" {
			//letsExcelize(foundPhones)
			return
		}
		c_parking.Visit(nextLink)

	})

	c_parking.Visit(fmt.Sprintf("https://www.theparking.eu/used-cars/nissan-leaf.html#!/used-cars/nissan-leaf.html%%3Ftri%%3Dprix_croissant"))
	c_parking.Wait()

	//fmt.Printf("https://www.theparking.eu/used-cars/nissan-leaf.html#!/used-cars/nissan-leaf.html%%3Ftri%%3Dprix_croissant")
}
