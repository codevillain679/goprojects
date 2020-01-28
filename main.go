// vanderbinckes
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jedib0t/go-pretty/table"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	login             = "root:pwfontys"
	db, err           = sql.Open("mysql", login+"@/vanderBinckesdb")
	medewerker        Medewerker
	klant             Klant
	verhuur           Verhuur
	bakfiets          Bakfiets
	accessoire        Accessoire
	verhuuraccessoire Verhuuraccessoire
	accessoires       = make([]Accessoire, 0)
)

func main() {
	login, medewerker := loginSuccesful(credentials())
	if login == true {
		fmt.Println("Inloggen succesvol!")
		fmt.Println("Welkom", medewerker)
		mainMenu(medewerker)
	} else {
		fmt.Println("Inloggen mislukt!")
		main()
	}
}

func credentials() (string, string) {
	//return "Bas", "2018-05-21"
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nLog in met gebruikersnaam [voornaam] en wachtwoord [datum in dienst]\nGebruikersnaam: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Wachtwoord: ")
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func loginSuccesful(username string, password string) (bool, Medewerker) {
	query := "select datum_in_dienst,achternaam,medewerkernummer,voornaam from medewerker where voornaam = ?"
	resultaat, err := db.Query(query, username)
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		err := resultaat.Scan(&medewerker.Datum_in_dienst, &medewerker.Achternaam, &medewerker.Medewerkernummer, &medewerker.Voornaam)
		if err != nil {
			panic(err)
		}
		return password == medewerker.Datum_in_dienst, medewerker
	}
	resultaat.Close()
	return false, medewerker
}

func mainMenu(medewerker Medewerker) {
	var action = scanInput("Hoofdmenu\n1.\tKlant toevoegen\n2.\tBakfiets verhuren\n3.\tKlanten inzien\n4.\tVerhuren inzien\n", false, 4)
	switch action {
	case "1":
		klant = addCustomer()
		break
	case "2":
		rentCargobike()
		break
	case "3":
		listCustomer()
		break
	case "4":
		listRental()
		listCargobike()
		break
	}
	mainMenu(medewerker)
}

func scanInput(caption string, required bool, upper ...float64) string {
	fmt.Printf(caption)
	if required {
		fmt.Printf("(*) ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	var val, _ = strconv.ParseFloat(strings.Replace(input, "\n", "", -1), 64)
	if required && input == "" && len(upper) > 0 {
		return scanInput(caption, true, upper[0])
	}
	if required && input == "0" {
		return scanInput(caption, required)
	} else if len(upper) > 0 && val > upper[0] {
		return scanInput(caption, required, upper[0])
	} else if input == "exit" {
		os.Exit(0)
	}
	return input
}

func addCustomer() Klant {
	fmt.Println("Om een klant toe te voegen, voer de bijbehorende gegevens in. De velden met een (*) zijn verplicht. Om een invoer leeg te laten voer 0 in.")
	klant.Naam = scanInput("Naam: ", true)
	klant.Voornaam = scanInput("Voornaam ", true)
	klant.Postcode = scanInput("Postcode ", true)
	klant.Huisnummer, _ = strconv.Atoi(scanInput("Huisnummer ", true))
	klant.Huisnummer_toevoeging = scanInput("Huisnummer toevoeging: ", false)
	klant.Opmerkingen = scanInput("Opmerkingen: ", false)

	query := "INSERT INTO klant (klantnummer,naam,voornaam,postcode,huisnummer,huisnummer_toevoeging,opmerkingen) values ((select max(klantnummer)+1 from klant kl),?,?,?,?,?,?)"
	_, err := db.Query(query, klant.Naam, klant.Voornaam, klant.Postcode, klant.Huisnummer, klant.Huisnummer_toevoeging, klant.Opmerkingen)
	if err != nil {
		panic(err)
	}
	return klant
}

func rentCargobike() {
	var action = scanInput("Bakfiets verhuren\n1.\tNieuwe klant\n2.\tBestaande klant\n", false, 2)
	switch action {
	case "1":
		klant = addCustomer()
		klant.Klantnummer = int(listCustomer())
		break
	case "2":
		klant = scanCustomer()
		break
	}
	bakfiets = scanCargobike()
	fmt.Println(bakfiets)
	verhuur = scanRental()
	accessoires = make([]Accessoire, 0)
	accessoire = scanAccessoire()
	if confirmRental() {
		calculateRentPrice()
		updateCargobike()
		showReceipt()
		commitRental()
		listRental()
		listCargobike()
	}
}

func scanCustomer() Klant {
	fmt.Println("Om een bakfiets te huren/verhuren, voer de bijbehorende gegevens in. De velden met een (*) zijn verplicht. Om een invoer leeg te laten voer 0 in.")
	customerCount := listCustomer()
	input, err := strconv.Atoi(scanInput("Klantnummer: ", true, customerCount))
	if err != nil || input < 0 {
		fmt.Println("Invoer ongeldig! Voer een geldig klantnummer in.")
		return scanCustomer()
	}
	query := "Select klantnummer,naam,voornaam,postcode,huisnummer,ifnull(huisnummer_toevoeging,''),ifnull(opmerkingen,'') from klant where klantnummer = ?"
	resultaat, err := db.Query(query, input)
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		err := resultaat.Scan(&klant.Klantnummer, &klant.Naam, &klant.Voornaam, &klant.Postcode, &klant.Huisnummer, &klant.Huisnummer_toevoeging, &klant.Opmerkingen)
		if err != nil {
			panic(err)
		}
	}
	resultaat.Close()
	return klant
}

func scanCargobike() Bakfiets {
	fmt.Println("Om een bakfiets te huren/verhuren, voer de bijbehorende gegevens in. De velden met een (*) zijn verplicht. Om een invoer leeg te laten voer 0 in.")
	bakfietsCount := listCargobike()
	input, _ := strconv.Atoi(scanInput("Bakfietsnummer: ", true, bakfietsCount))
	query := "Select bakfietsnummer,naam,type,huurprijs,aantal,aantal_verhuurd from bakfiets where bakfietsnummer = ?"
	resultaat, err := db.Query(query, input)
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		err := resultaat.Scan(&bakfiets.Bakfietsnummer, &bakfiets.Naam, &bakfiets.Type, &bakfiets.Huurprijs, &bakfiets.Aantal, &bakfiets.Aantal_verhuurd)
		if err != nil {
			panic(err)
		}
	}
	resultaat.Close()
	return bakfiets
}

func scanRental() Verhuur {
	fmt.Println("Om een bakfiets te huren/verhuren, voer de bijbehorende gegevens in. De velden met een (*) zijn verplicht. Om een invoer leeg te laten voer 0 in.")
	input, err := strconv.Atoi(scanInput("Aantal dagen: ", true))
	if err != nil || input <= 0 {
		fmt.Printf("Invoer Ongeldig! ")
		return scanRental()
	}
	verhuur.Aantal_dagen = input

	query := "Insert into verhuur(verhuurnummer,verhuurdatum,bakfietsnummer,aantal_dagen,huurprijstotaal,klantnummer,verhuurder) values( (select count(*) + 1 from verhuur v), (select curdate()),?,?,?,?,?)"
	_, err = db.Query(query, bakfiets.Bakfietsnummer, verhuur.Aantal_dagen, 0, klant.Klantnummer, medewerker.Medewerkernummer)
	if err != nil {
		panic(err)
	}

	query = "select max(verhuurnummer) from verhuur"
	resultaat, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		err := resultaat.Scan(&verhuur.Verhuurnummer)
		if err != nil {
			panic(err)
		}
	}

	return verhuur
}

func scanAccessoire() Accessoire {
	fmt.Println("Om een accessoires toe te voegen, voer de bijbehorende gegevens in. Om te stoppen met accessoires toevoegen voer een 0 in.")
	accessoireCount := listAccessoire()
	input, _ := strconv.Atoi(scanInput("Accessoirenummer: ", false, accessoireCount))
	if input <= 0 {
		var accessoire Accessoire
		return accessoire
	}
	query := "Select accessoirenummer,naam,huurprijs from accessoire where accessoirenummer = ?"
	resultaat, err := db.Query(query, input)
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		err := resultaat.Scan(&accessoire.Accessoirenummer, &accessoire.Naam, &accessoire.Huurprijs)
		if err != nil {
			panic(err)
		}
	}
	verhuuraccessoire = scanVerhuuraccessoire()
	accessoires = append(accessoires, accessoire)
	resultaat.Close()
	return scanAccessoire()
}

func scanVerhuuraccessoire() Verhuuraccessoire {
	input, err := strconv.Atoi(scanInput("Aantal: ", true))
	if err != nil || input <= 0 {
		fmt.Printf("Invoer Ongeldig! ")
		return scanVerhuuraccessoire()
	}
	var verhuuraccessoire Verhuuraccessoire
	verhuuraccessoire.Aantal = input
	verhuuraccessoire.Accessoirenummer = accessoire.Accessoirenummer
	verhuuraccessoire.Verhuurnummer = verhuur.Verhuurnummer
	query := "Insert into verhuuraccessoire(verhuurnummer,accessoirenummer,aantal) values(?,?,?)"
	fmt.Println(verhuuraccessoire)
	_, err = db.Query(query, verhuuraccessoire.Verhuurnummer, verhuuraccessoire.Accessoirenummer, verhuuraccessoire.Aantal)
	if err != nil {
		panic(err)
	}
	accessoire.Huurprijs = accessoire.Huurprijs * float64(verhuuraccessoire.Aantal)
	return verhuuraccessoire
}

func listCustomer() float64 {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Naam", "Voornaam", "Postcode", "Huisnummer", "Toevoeging", "Opmerkingen"})
	query := "Select klantnummer,naam,voornaam,postcode,huisnummer,ifnull(huisnummer_toevoeging,''),ifnull(opmerkingen,'') from klant"
	resultaat, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	n := 0.0
	for resultaat.Next() {
		err := resultaat.Scan(&klant.Klantnummer, &klant.Naam, &klant.Voornaam, &klant.Postcode, &klant.Huisnummer, &klant.Huisnummer_toevoeging, &klant.Opmerkingen)
		if err != nil {
			panic(err)
		}
		t.AppendRow([]interface{}{klant.Klantnummer, klant.Naam, klant.Voornaam, klant.Postcode, klant.Huisnummer, klant.Huisnummer_toevoeging, klant.Opmerkingen})
		n = n + 1
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return n
}

func listCargobike() float64 {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Naam", "Type", "Huurprijs", "Aantal", "Aantal verhuurd"})
	query := "Select bakfietsnummer,naam,type,huurprijs,aantal,aantal_verhuurd from bakfiets order by bakfietsnummer"
	resultaat, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	n := 0.0
	for resultaat.Next() {
		err := resultaat.Scan(&bakfiets.Bakfietsnummer, &bakfiets.Naam, &bakfiets.Type, &bakfiets.Huurprijs, &bakfiets.Aantal, &bakfiets.Aantal_verhuurd)
		if err != nil {
			panic(err)
		}
		t.AppendRow([]interface{}{bakfiets.Bakfietsnummer, bakfiets.Naam, bakfiets.Type, bakfiets.Huurprijs, bakfiets.Aantal, bakfiets.Aantal_verhuurd})
		n = n + 1
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return n
}

func listRental() float64 {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Verhuurdatum", "Bakfietsnummer", "Aantal_dagen", "Huurprijstotaal", "Klantnummer", "Verhuurder"})
	query := "Select verhuurnummer,verhuurdatum,bakfietsnummer,aantal_dagen,huurprijstotaal,klantnummer,verhuurder from verhuur where verhuurdatum+aantal_dagen>= curdate() order by verhuurnummer"
	resultaat, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	n := 0.0
	for resultaat.Next() {
		resultaat.Scan(&verhuur.Verhuurnummer, &verhuur.Verhuurdatum, &verhuur.Bakfietsnummer, &verhuur.Aantal_dagen, &verhuur.Huurprijstotaal, &verhuur.Klantnummer, &verhuur.Verhuurder)
		t.AppendRow([]interface{}{verhuur.Verhuurnummer, verhuur.Verhuurdatum, verhuur.Bakfietsnummer, verhuur.Aantal_dagen, verhuur.Huurprijstotaal, verhuur.Klantnummer, verhuur.Verhuurder})
		n = n + 1
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return n
}

func listAccessoire() float64 {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Naam", "Huurprijs"})
	query := "Select accessoirenummer,naam,huurprijs from accessoire order by accessoirenummer"
	resultaat, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	n := 0.0
	var accessoire Accessoire
	for resultaat.Next() {
		resultaat.Scan(&accessoire.Accessoirenummer, &accessoire.Naam, &accessoire.Huurprijs)
		t.AppendRow([]interface{}{accessoire.Accessoirenummer, accessoire.Naam, accessoire.Huurprijs})
		n = n + 1
	}
	t.SetStyle(table.StyleLight)
	t.Render()
	return n
}

func calculateRentPrice() {
	verhuur.Huurprijstotaal = bakfiets.Huurprijs * float64(verhuur.Aantal_dagen)
	query := "Select huurprijs*aantal from accessoire,verhuuraccessoire where verhuurnummer=? and accessoire.accessoirenummer=verhuuraccessoire.accessoirenummer"
	resultaat, err := db.Query(query, verhuur.Verhuurnummer)
	var deelprijs float64
	if err != nil {
		panic(err)
	}
	for resultaat.Next() {
		resultaat.Scan(&deelprijs)
		verhuur.Huurprijstotaal += deelprijs
	}
}

func updateCargobike() {
	query := "update bakfiets set aantal_verhuurd = (select count(*) from verhuur where verhuur.bakfietsnummer=bakfiets.bakfietsnummer and verhuurdatum > (select curdate()-aantal_dagen))"
	_, err = db.Query(query)
	if err != nil {
		panic(err)
	}
}

func confirmRental() bool {
	var action = scanInput("Bestelling bevestigen?\nY/y.\tJa\nN/n.\tNee\n", false, 2)
	switch action {
	case "Y", "y":
		return true
	case "N", "n":
		return false
	}
	return false
}

func showReceipt() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow([]interface{}{bakfiets.Naam, bakfiets.Type, bakfiets.Huurprijs, float64(verhuur.Aantal_dagen) * bakfiets.Huurprijs})
	for _, accessoire := range accessoires {
		query := "Select huurprijs from accessoire where accessoirenummer=?"
		resultaat, _ := db.Query(query, accessoire.Accessoirenummer)
		var deelprijs float64
		for resultaat.Next() {
			resultaat.Scan(&deelprijs)
		}
		t.AppendRow([]interface{}{accessoire.Naam, accessoire.Huurprijs / deelprijs, deelprijs, accessoire.Huurprijs})
	}
	t.AppendRow([]interface{}{"", "", "", ""})
	t.AppendRow([]interface{}{"", "", "Totaal ", verhuur.Huurprijstotaal})
	t.AppendRow([]interface{}{klant.Klantnummer, klant.Voornaam, klant.Naam, ""})
	t.AppendRow([]interface{}{klant.Postcode, klant.Huisnummer, klant.Huisnummer_toevoeging, klant.Opmerkingen})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func commitRental() {
	query := "Update verhuur set huurprijstotaal = ? where verhuurnummer = ?"
	_, err := db.Query(query, verhuur.Huurprijstotaal, verhuur.Verhuurnummer)
	if err != nil {
		panic(err)
	}
}
