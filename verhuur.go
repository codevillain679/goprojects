package main

// Verhuur vertegenwoordigt 1 specifieke verhuur.
type Verhuur struct {
	Verhuurnummer   int
	Verhuurdatum    string
	Bakfietsnummer  int
	Aantal_dagen    int
	Huurprijstotaal float64
	Klantnummer     int
	Verhuurder      int
}
